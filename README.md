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
5. 在 nameserver 中执行命令 read b/a.txt 可以下载 a.txt 到本地
6. help 命令可以展示各个命令的使用方法
7. tree 命令可以用于展示文件树结构

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
- [x] 下载文件实现
- [x] 布隆过滤器实现
- [x] 实现文件树序列化为 json
- [ ] 文件拆分为 chunk
- [ ] 接入文件树序列化
- [ ] 接入布隆过滤器
- [x] help 命令实现
