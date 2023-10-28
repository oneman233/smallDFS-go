# intro

* 参考自[smallDFS](https://github.com/oneman233/smallDFS)
    * 用 go 语言重构了 cpp 实现的文件树
    * 将原项目的负载均衡策略直接替换为一致性哈希
* 引入基于 http 和 protobuf 的通信机制，真正实现了分布式
* 引入一致性哈希存储的文件名到节点的映射，支持虚拟节点
* 引入布隆过滤器降低文件树访问压力
* 支持文件树序列化和反序列化，同步对布隆过滤器进行更新

# usage

1. 在 main 函数中启动 nameserver
2. 在 main 函数中启动 dataserver
3. 在 nameserver 中执行命令 put a.txt b/a.txt
4. 项目目录下会新建目录 ds-1/b，同时在该目录下生成 a.txt 的三份拷贝
5. 在 nameserver 中执行命令 read b/a.txt 可以下载 a.txt 到本地
6. help 命令可以展示各个命令的使用方法
7. dump/undump 命令用于序列化和反序列化文件树（正常退出时自动保存序列化文件）
8. tree 命令可以用于展示文件树结构

# 项目结构
```text
.
│  go.mod
│  go.sum
│  main.go
│  README.md
│
├─bloomfilter
│      bloomfilter.go
│
├─consistenthash
│      consistenthash.go
│      consistenthash_test.go
│
├─constants
│      constants.go
│
├─dataserver
│      dataserver.go
│
├─filetree
│      filetree.go
│      filetree_test.go
│
├─nameserver
│      nameproxy.go
│      nameserver.go
│
├─pb
│      complie.sh
│      pb.go
│      pb.proto
│
└─skiplist
        node.go
        skiplist.go
        skiplist_test.go

```

# todos

- [x] 上传文件返回消息
- [x] 下载文件实现
- [x] 布隆过滤器实现
- [x] 实现文件树序列化为 json
- [x] 接入文件树序列化
- [x] 接入布隆过滤器
- [x] help 命令实现
- [x] 布隆过滤器在树反序列化后重置
- [x] 跳表实现及测试
- [x] 文件树支持 delete
- [ ] 抽离出单独的日志逻辑
- [ ] delete 实现
