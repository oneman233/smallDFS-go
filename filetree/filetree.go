package filetree

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	value  string
	isFile bool
	sons   map[string]*TreeNode
}

type FileTree struct {
	Root *TreeNode
}

func New() *FileTree {
	return &FileTree{
		Root: &TreeNode{
			value:  "/",
			isFile: false,
			sons:   make(map[string]*TreeNode),
		},
	}
}

// Find 返回指定的路径是否存在，如果存在返回其是否是文件
func (t *FileTree) Find(path string) (bool, bool) {
	folders := strings.Split(path, "/")
	node := t.Root
	for _, folder := range folders {
		if node.sons[folder] != nil {
			node = node.sons[folder]
		} else {
			return false, false
		}
	}
	return true, node.isFile
}

// Insert 向文件树中插入路径或文件，如果路径或文件存在则返回错误
func (t *FileTree) Insert(path string, isFile bool) error {
	if isPath, _ := t.Find(path); isPath {
		return fmt.Errorf("path %s already existed", path)
	}

	folders := strings.Split(path, "/")
	node := t.Root
	for _, folder := range folders {
		if node.sons[folder] == nil {
			newNode := &TreeNode{
				value:  folder,
				isFile: false,
				sons:   make(map[string]*TreeNode),
			}
			node.sons[folder] = newNode
		}
		node = node.sons[folder]
	}
	node.isFile = isFile
	return nil
}

// List 输出文件树某个节点的所有子节点
func List(node *TreeNode) {
	for _, son := range node.sons {
		fmt.Printf("%s ", son.value)
	}
	fmt.Println()
}

// Tree 输出文件树某个节点的所有子树
func Tree(node *TreeNode, counter int) {
	for i := 0; i < counter; i++ {
		fmt.Printf("-")
	}
	fmt.Println(node.value)
	for _, son := range node.sons {
		Tree(son, counter+1)
	}
}
