package filetree

import (
	"fmt"
	"smallDFS/constants"
	"testing"
)

// 文件树序列化测试
func TestJSON(t *testing.T) {
	tree := New()
	_ = tree.Insert("b", false)
	_ = tree.Insert("c/a.txt", true)
	_ = tree.Insert("c", false)
	tree.Dump(constants.DefaultJSONName)
	treeRead, err := UnDump(constants.DefaultJSONName)
	if err != nil {
		panic(err)
	}
	Tree(treeRead.Root, 1)
}

// 获取全部路径函数测试
func TestGetAllPath(t *testing.T) {
	tree := New()
	_ = tree.Insert("b", false)
	_ = tree.Insert("c/a.txt", true)
	_ = tree.Insert("c", false)
	_ = tree.Insert("d.txt", true)
	paths := GetAllPath(tree.Root, "")
	fmt.Println(len(paths))
	for _, path := range paths {
		fmt.Println(path)
	}
}

func TestFileTree_Delete(t *testing.T) {
	tree := New()
	_ = tree.Insert("b", false)
	_ = tree.Insert("c/a.txt", true)
	_ = tree.Insert("c", false)
	_ = tree.Insert("d.txt", true)
	Tree(tree.Root, 1)
	_ = tree.Delete("c")
	_ = tree.Delete("c/d.txt")
	_ = tree.Delete("d.txt")
	_ = tree.Delete("c/a.txt")
	Tree(tree.Root, 1)
}
