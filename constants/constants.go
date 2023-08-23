package constants

import "os"

// DefaultFileMode 打开文件的默认模式
var DefaultFileMode = os.FileMode(0777)

// DefaultJSONName 序列化文件树的默认 json 名
var DefaultJSONName = "FileTree.json"
