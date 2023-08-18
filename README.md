# 项目结构
```text
.
├── README.md
├── consistenthash
│        ├── consistenthash.go
│        └── consistenthash_test.go
├── dataserver
│        └── dataserver.go
├── filetree
│        └── filetree.go
├── go.mod
├── go.sum
├── main.go
├── nameserver
│        ├── nameproxy.go
│        └── nameserver.go
└── pb
    ├── pb.pb.go
    └── pb.proto
```

# todos

- [ ] 上传文件返回 protobuf
- [ ] 下载文件实现
- [ ] 布隆过滤器实现