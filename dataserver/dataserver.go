package dataserver

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net/http"
	"os"
	"smallDFS/constants"
	"smallDFS/pb"
	"strings"
)

type DataServer struct {
	port       string
	folderName string
}

func New(port string, folderName string) *DataServer {
	return &DataServer{
		port:       port,
		folderName: folderName,
	}
}

func (ds *DataServer) log(format string, v ...interface{}) {
	log.Printf("[DataServer %s] %s", ds.folderName, fmt.Sprintf(format, v...))
}

func (ds *DataServer) Run() {
	ds.log("run")
	// 创建并更改工作目录
	_ = os.Mkdir(ds.folderName, constants.DefaultFileMode)
	_ = os.Chdir(ds.folderName)
	// 设置 http handler
	http.HandleFunc(constants.DefaultUploadPath, ds.uploadHandler)
	http.HandleFunc(constants.DefaultDownloadPath, ds.downloadHandler)
	ds.log("listening")
	// 监听端口
	log.Fatal(http.ListenAndServe(":"+ds.port, nil))
}

// 上传文件 handler
func (ds *DataServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	ds.log("upload")
	// 读取 http 请求内容
	data, _ := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	// 解码 http 请求
	pbReq := &pb.UploadFileRequest{}
	_ = proto.Unmarshal(data, pbReq)

	// 获取上传的文件
	file := pbReq.GetFile()
	// 获取文件存储路径
	path := pbReq.GetPath()

	// 如果是形如 a/b.txt 这样的路径，则需要新建路径
	if strings.Contains(path, "/") {
		err := os.MkdirAll(getPath(path), constants.DefaultFileMode)
		if err != nil {
			panic(err)
		}
	}

	// 在 path 创建文件
	createFile(path, []byte(file))

	// 返回成功消息
	pbRes := &pb.UploadFileResponse{Message: "success"}
	httpRes, _ := proto.Marshal(pbRes)
	_, _ = w.Write(httpRes)
}

// getPath 从文件名中获取路径，例如 a/b/c/d.txt -> a/b/c
func getPath(pathAndFilename string) string {
	for i := len(pathAndFilename) - 1; i >= 0; i-- {
		if pathAndFilename[i] == '/' {
			return pathAndFilename[0:i]
		}
	}
	return ""
}

// createFile 在指定路径新建文件并写入内容，必须确保路径存在
func createFile(name string, data []byte) {
	f, _ := os.Create(name)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	_, _ = f.Write(data)
}

func (ds *DataServer) downloadHandler(w http.ResponseWriter, r *http.Request) {
	ds.log("download")
	// 读取 http 请求内容
	data, _ := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	// 解码 http 请求
	pbReq := &pb.DownloadFileRequest{}
	_ = proto.Unmarshal(data, pbReq)

	// 获取文件路径并读取文件
	path := pbReq.GetPath()
	file, _ := os.ReadFile(path)

	// 封装下载文件 response
	pbRes := &pb.DownloadFileResponse{
		Message: "",
		File:    string(file),
	}
	httpRes, _ := proto.Marshal(pbRes)
	_, _ = w.Write(httpRes)
}
