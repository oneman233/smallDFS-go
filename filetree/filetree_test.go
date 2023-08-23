package filetree

import (
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
