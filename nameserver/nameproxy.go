package nameserver

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
	"smallDFS/pb"
)

type NameProxy struct {
	uploadPath   string
	downloadPath string
}

// UploadFile 向 addr 发送 byte 数组形式的文件 file，存储路径为 path，返回上传结果
func (np *NameProxy) UploadFile(file []byte, path string, addr string) string {
	pbReq := &pb.UploadFileRequest{
		File: string(file),
		Path: path,
	}
	httpReq, _ := proto.Marshal(pbReq)
	httpRes, err := http.Post(addr+np.uploadPath, "multipart/form-data", bytes.NewReader(httpReq))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(httpRes.Body)

	body, _ := io.ReadAll(httpRes.Body)
	pbRes := &pb.UploadFileResponse{}
	_ = proto.Unmarshal(body, pbRes)
	return pbRes.Message
}

// DownloadFile 从 addr 下载路径为 path 的文件，返回值是 byte 数组形式
func (np *NameProxy) DownloadFile(path string, addr string) []byte {
	return nil
}
