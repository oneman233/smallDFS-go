package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// HashFunc 哈希函数定义，将 byte 数组映射到 uint32
type HashFunc func(data []byte) uint32

// Map 一致性哈希核心数据结构
type Map struct {
	hashFunc    HashFunc       // 哈希函数
	replicas    int            // 每个节点生成几个虚拟节点
	fakeNodes   []int          // 有序的虚拟节点位置
	fakeNodeMap map[int]string // 虚拟节点到实际节点的映射
}

// New 新建一致性哈希，默认的哈希函数为 crc32.ChecksumIEEE
func New(replicas int, hashFunc HashFunc) *Map {
	tmp := &Map{
		hashFunc:    hashFunc,
		replicas:    replicas,
		fakeNodeMap: make(map[int]string),
	}
	if tmp.hashFunc == nil {
		tmp.hashFunc = crc32.ChecksumIEEE
	}
	return tmp
}

// Add 向一致性哈希中插入节点，自动生成虚拟节点并排序
func (mp *Map) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < mp.replicas; i++ {
			fakeNode := int(mp.hashFunc([]byte(node + strconv.Itoa(i))))
			mp.fakeNodes = append(mp.fakeNodes, fakeNode)
			mp.fakeNodeMap[fakeNode] = node
		}
	}
	sort.Ints(mp.fakeNodes)
}

// Get 获取 key 对应的节点
func (mp *Map) Get(key string) string {
	// 如果当前没有节点，返回空值
	if len(mp.fakeNodes) == 0 {
		return ""
	}

	// 计算 key 对应的哈希值
	hashVal := int(mp.hashFunc([]byte(key)))
	// 找到第一个大于等于 key 的节点
	idx := sort.Search(len(mp.fakeNodes), func(i int) bool {
		return mp.fakeNodes[i] >= hashVal
	})
	// 如果 key 大雨所有节点，那么把它映射到第一个节点（一致性哈希是一个数环）
	idx %= len(mp.fakeNodes)

	// 从虚拟节点映射到实际节点，并将实际节点作为返回值返回
	return mp.fakeNodeMap[mp.fakeNodes[idx]]
}
