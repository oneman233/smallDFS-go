package nameserver

import (
	"bufio"
	"fmt"
	"os"
	"smallDFS/consistenthash"
	"smallDFS/filetree"
	"strconv"
	"strings"
)

type NameServer struct {
	consistentHash *consistenthash.Map // 一致性哈希用于保存某个文件对应的分块被存储到了哪个 DataServer
	fileTree       *filetree.FileTree  // 文件树用于保存当前存储的文件及目录
	fileReplicas   int                 // 上传的文件创建几份副本
	stdinReader    *bufio.Reader       // 从 stdin 中读取命令的 reader
	proxy          *NameProxy          // 网络通信的代理
}

func New(fileReplicas int, fakeNodes int) *NameServer {
	return &NameServer{
		consistentHash: consistenthash.New(fakeNodes, nil),
		fileTree:       filetree.New(),
		fileReplicas:   fileReplicas,
		stdinReader:    bufio.NewReader(os.Stdin),
		proxy: &NameProxy{
			uploadPath:   "/upload",
			downloadPath: "/download",
		},
	}
}

func (ns *NameServer) parseCmd() []string {
	fmt.Print("smallDFS> ")
	// 从标准输入流读取一整行
	line, err := ns.stdinReader.ReadString('\n')
	// 去除行尾的换行符
	line = strings.TrimSpace(line)
	if err != nil {
		panic(err)
	}
	return strings.Split(line, " ")
}

func (ns *NameServer) Run() {
	for {
		params := ns.parseCmd()

		if len(params) == 0 {
			fmt.Println("blank line")
			continue
		}

		switch params[0] {
		case "quit":
			os.Exit(0)
		case "put":
			_ = ns.put(params[1], params[2])
		case "read":
			ns.read(params[1])
		case "tree":
			ns.tree()
		default:
			fmt.Println("wrong command")
		}
	}
}

func (ns *NameServer) put(localPath string, remotePath string) error {
	// 查看是否存在同名文件
	_, isFile := ns.fileTree.Find(remotePath)
	if isFile {
		return fmt.Errorf("remote path %s existed", remotePath)
	}

	// 读取本地文件
	localFile, err := os.ReadFile("./" + localPath)
	if err != nil {
		panic(err)
	}

	// 插入远程路径
	_ = ns.fileTree.Insert(remotePath, true)

	// 创建文件副本
	for i := 0; i < ns.fileReplicas; i++ {
		// 拼接文件副本名
		replicaName := remotePath + "-" + strconv.Itoa(i)
		// 获取实际存储节点地址
		nodeAddr := ns.consistentHash.Get(replicaName)
		// 调用网络通信代理上传文件
		ns.proxy.UploadFile(localFile, replicaName, nodeAddr)
	}

	return nil
}

func (ns *NameServer) read(remotePath string) {

}

func (ns *NameServer) tree() {
	filetree.Tree(ns.fileTree.Root, 1)
}

// Add 注册新 DataServer
func (ns *NameServer) Add(remoteAddr string) {
	ns.consistentHash.Add(remoteAddr)
}
