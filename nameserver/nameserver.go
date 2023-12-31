package nameserver

import (
	"bufio"
	"fmt"
	"os"
	"smallDFS/bloomfilter"
	"smallDFS/consistenthash"
	"smallDFS/constants"
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
	bloom          *bloomfilter.Bloom  // 布隆过滤器
}

// New 新建 NameServer，三个参数分别是文件备份数、虚拟节点数、布隆过滤器哈希次数
func New(fileReplicas int, fakeNodes int, bloomHashCount int) *NameServer {
	return &NameServer{
		consistentHash: consistenthash.New(fakeNodes, nil),
		fileTree:       filetree.New(),
		fileReplicas:   fileReplicas,
		stdinReader:    bufio.NewReader(os.Stdin),
		proxy: &NameProxy{
			uploadPath:   constants.DefaultUploadPath,
			downloadPath: constants.DefaultDownloadPath,
		},
		bloom: bloomfilter.New(bloomHashCount),
	}
}

// 转换用户输入的命令为参数列表
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

// Run 启动 NameServer
func (ns *NameServer) Run() {
	for {
		params := ns.parseCmd()

		if len(params) == 0 {
			fmt.Println("blank line")
			continue
		}

		switch params[0] {
		case "quit":
			func() {
				// 在退出前自动 dump 文件树
				ns.dump()
				os.Exit(0)
			}()
		case "put":
			_ = ns.put(params[1], params[2])
		case "read":
			_ = ns.read(params[1])
		case "tree":
			ns.tree()
		case "help":
			ns.help()
		case "dump":
			ns.dump()
		case "undump":
			_ = ns.unDump()
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

	// 插入布隆过滤器
	ns.bloom.Insert(remotePath)

	// 创建文件副本
	for i := 0; i < ns.fileReplicas; i++ {
		// 拼接文件副本名
		replicaName := remotePath + "-" + strconv.Itoa(i)
		// 获取实际存储节点地址
		nodeAddr := ns.consistentHash.Get(replicaName)
		// 调用网络通信代理上传文件
		msg := ns.proxy.UploadFile(localFile, replicaName, nodeAddr)
		fmt.Println(msg)
	}

	return nil
}

func (ns *NameServer) read(remotePath string) error {
	// 利用布隆过滤器检查远程路径是否存在
	if !ns.bloom.Check(remotePath) {
		return fmt.Errorf("remote path %s not existed", remotePath)
	}

	// 检查远程路径是否存在
	isPath, isFile := ns.fileTree.Find(remotePath)
	if !isPath || !isFile {
		return fmt.Errorf("remote path %s not existed", remotePath)
	}

	success := false
	var file []byte
	var err error

	// 尝试下载可用副本，下载成功则退出循环
	for i := 0; i < ns.fileReplicas; i++ {
		replicaName := remotePath + "-" + strconv.Itoa(i)
		nodeAddr := ns.consistentHash.Get(replicaName)
		// 下载的文件名要加上 "-0"
		file, err = ns.proxy.DownloadFile(replicaName, nodeAddr)
		if err == nil {
			success = true
			break
		}
	}

	if !success {
		return fmt.Errorf("no avaliable copys of file %s", remotePath)
	} else {
		// 将文件内容写入当前路径，文件名不加 "-0"
		folders := strings.Split(remotePath, "/")
		fileName := folders[len(folders)-1]
		err := os.WriteFile(fileName, file, constants.DefaultFileMode)
		if err != nil {
			return err
		}
	}

	return nil
}

// 打印整个文件树
func (ns *NameServer) tree() {
	filetree.Tree(ns.fileTree.Root, 1)
}

// 序列化文件树
func (ns *NameServer) dump() {
	ns.fileTree.Dump(constants.DefaultJSONName)
}

// 反序列化文件树并赋值给 NameServer
func (ns *NameServer) unDump() error {
	t, err := filetree.UnDump(constants.DefaultJSONName)
	if err != nil {
		return err
	}
	ns.fileTree = t
	// 清空布隆过滤器
	ns.bloom.Purge()
	// 获取当前文件树的全部路径
	paths := filetree.GetAllPath(ns.fileTree.Root, "")
	// 将全部路径插入布隆过滤器
	for _, path := range paths {
		ns.bloom.Insert(path)
	}
	return nil
}

// Add 注册新 DataServer
func (ns *NameServer) Add(remoteAddr string) {
	ns.consistentHash.Add(remoteAddr)
}

// help 命令输出使用指南
func (ns *NameServer) help() {
	fmt.Println("upload file: put <local file name> <remote path>")
	fmt.Println("download file: read <remote path>")
	fmt.Println("check file tree: tree")
	fmt.Println("dump file tree: dump")
	fmt.Println("undump file tree: undump")
	fmt.Println("close server: exit")
}
