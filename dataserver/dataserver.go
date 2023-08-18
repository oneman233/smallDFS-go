package dataserver

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net/http"
	"os"
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
	_ = os.Mkdir(ds.folderName, 0777)
	_ = os.Chdir(ds.folderName)
	// 设置 http handler
	http.HandleFunc("/upload", ds.uploadHandler)
	http.HandleFunc("/download", ds.downloadHandler)
	ds.log("listening")
	// 监听端口
	log.Fatal(http.ListenAndServe(":"+ds.port, nil))
}

func (ds *DataServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	ds.log("upload")
	data, _ := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	pbReq := &pb.UploadFileRequest{}
	_ = proto.Unmarshal(data, pbReq)

	file := pbReq.GetFile()
	path := pbReq.GetPath()

	if strings.Contains(path, "/") { // 需要新建路径
		err := os.MkdirAll(getPath(path), 0777)
		if err != nil {
			panic(err)
		}
	}

	createFile(path, []byte(file))
}

func getPath(pathAndFilename string) string {
	for i := len(pathAndFilename) - 1; i >= 0; i-- {
		if pathAndFilename[i] == '/' {
			return pathAndFilename[0:i]
		}
	}
	return ""
}

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

}
