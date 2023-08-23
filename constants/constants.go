package constants

import "os"

// DefaultFileMode 打开文件的默认模式
var DefaultFileMode = os.FileMode(0777)

// DefaultJSONName 序列化文件树的默认 json 文件名
var DefaultJSONName = "FileTree.json"

// DefaultUploadPath 默认的 http 上传路径
var DefaultUploadPath = "/upload"

// DefaultDownloadPath 默认的 http 下载路径
var DefaultDownloadPath = "/download"
