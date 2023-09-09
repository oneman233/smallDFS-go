package filetree

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"smallDFS/constants"
	"strings"
)

type TreeNode struct {
	Value  string               `json:"value"`
	IsFile bool                 `json:"isFile"`
	Sons   map[string]*TreeNode `json:"sons"`
}

type FileTree struct {
	Root *TreeNode
}

func New() *FileTree {
	return &FileTree{
		Root: &TreeNode{
			Value:  "/",
			IsFile: false,
			Sons:   make(map[string]*TreeNode),
		},
	}
}

// Find 返回指定的路径是否存在，如果存在返回其是否是文件
func (t *FileTree) Find(path string) (bool, bool) {
	folders := strings.Split(path, "/")
	node := t.Root
	for _, folder := range folders {
		if node.Sons[folder] != nil {
			node = node.Sons[folder]
		} else {
			return false, false
		}
	}
	return true, node.IsFile
}

// Insert 向文件树中插入路径或文件，如果路径或文件存在则返回错误
func (t *FileTree) Insert(path string, isFile bool) error {
	if isPath, _ := t.Find(path); isPath {
		return fmt.Errorf("path %s already existed", path)
	}

	folders := strings.Split(path, "/")
	node := t.Root
	for _, folder := range folders {
		if node.Sons[folder] == nil {
			newNode := &TreeNode{
				Value:  folder,
				IsFile: false,
				Sons:   make(map[string]*TreeNode),
			}
			node.Sons[folder] = newNode
		}
		node = node.Sons[folder]
	}
	node.IsFile = isFile
	return nil
}

// Delete 删除文件树中的指定文件
func (t *FileTree) Delete(path string) error {
	// 如果路径不是文件，那么返回错误
	if _, isFile := t.Find(path); !isFile {
		return fmt.Errorf("path %s is not a file", path)
	}

	folders := strings.Split(path, "/")
	node := t.Root
	// 最后一个文件节点的下标
	lastFolderIdx := len(folders) - 1

	for i, folder := range folders {
		// 找到文件节点的父目录节点
		if i == lastFolderIdx {
			break
		}
		node = node.Sons[folder]
	}

	// 删除父目录节点中的文件指针
	delete(node.Sons, folders[lastFolderIdx])
	return nil
}

// Tree 输出文件树某个节点的所有子树
func Tree(node *TreeNode, counter int) {
	for i := 0; i < counter; i++ {
		fmt.Printf("-")
	}
	fmt.Println(node.Value)
	for _, son := range node.Sons {
		Tree(son, counter+1)
	}
}

// Dump 序列化文件树为 json 并保存在本地
func (t *FileTree) Dump(jsonName string) {
	// 序列化为 json
	j, _ := json.Marshal(t)
	// 打开文件
	file, _ := os.OpenFile(jsonName, os.O_CREATE|os.O_RDWR, constants.DefaultFileMode)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	// 写入 json
	_, _ = file.Write(j)
}

// UnDump 读取本地 json 文件并转换为 FileTree
func UnDump(jsonName string) (*FileTree, error) {
	// 打开 json 文件
	file, err := os.Open(jsonName)
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	j, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// json 解码并赋值给 FileTree
	t := &FileTree{}
	err = json.Unmarshal(j, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetAllPath 获取文件树指定节点包含的全部路径
func GetAllPath(node *TreeNode, currentPath string) []string {
	// 如果不为根节点则拼接当前节点的路径
	if node.Value != "/" {
		currentPath += "/" + node.Value
	}
	// 生成结果数组
	var ans []string
	// 当前路径不为空则将其添加进结果数组
	if currentPath != "" {
		ans = append(ans, currentPath)
	}
	// 添加每个子节点对应的路径
	for _, son := range node.Sons {
		ans = append(ans, GetAllPath(son, currentPath)...)
	}
	return ans
}
