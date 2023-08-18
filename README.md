# intro

* 参考自[smallDFS](https://github.com/oneman233/smallDFS)
    * 用 go 语言重构了 cpp 实现的文件树
    * 将原项目的负载均衡策略直接替换为一致性哈希
* 引入基于 http 和 protobuf 的通信机制，真正实现了分布式
* 引入一致性哈希存储的文件名到节点的映射，支持虚拟节点
* 引入布隆过滤器降低文件树访问压力

# usage

1. 在 main 函数中启动 nameserver
2. 在 main 函数中启动 dataserver
3. 在 nameserver 中执行命令 put a.txt b/a.txt
4. 项目目录下会新建目录 ds-1/b，同时在该目录下生成 a.txt 的三份拷贝

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

- [x] 上传文件返回消息
- [ ] 下载文件实现
- [x] 布隆过滤器实现
- [ ] 文件拆分为 chunk
- [ ] 持久化
- [x] help 命令实现