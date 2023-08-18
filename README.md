# intro

* 参考自[smallDFS](https://github.com/oneman233/smallDFS)
    * 用 go 语言重构了 cpp 实现的文件树
    * 将原项目的负载均衡策略直接替换为一致性哈希
* 引入基于 http 和 protobuf 的通信机制，真正实现了分布式
* 引入一致性哈希存储的文件名到节点的映射，支持虚拟节点
* 引入布隆过滤器降低文件树访问压力

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
- [ ] 文件拆分为 chunk
