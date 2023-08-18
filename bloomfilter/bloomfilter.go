package bloomfilter

import (
	"crypto/sha256"
	"strconv"
)

type Bloom struct {
	filter      map[string]bool
	hashCounter int
}

// New 新建布隆过滤器
func New(hashCounter int) *Bloom {
	return &Bloom{
		filter:      make(map[string]bool),
		hashCounter: hashCounter,
	}
}

// Insert 向布隆过滤器中插入元素
func (b *Bloom) Insert(key string) {
	for i := 0; i < b.hashCounter; i++ {
		hashBytes := sha256.Sum256([]byte(key + strconv.Itoa(i)))
		hashString := string(hashBytes[:])
		b.filter[hashString] = true
	}
}

// Check 检查布隆过滤器中是否有该元素
func (b *Bloom) Check(key string) bool {
	for i := 0; i < b.hashCounter; i++ {
		hashBytes := sha256.Sum256([]byte(key + strconv.Itoa(i)))
		hashString := string(hashBytes[:])
		if b.filter[hashString] == false { // 只要有一个哈希值不存在，就说明该元素不存在
			return false
		}
	}
	return true // 所有哈希值都存在则需要进一步判断
}

// Purge 清空布隆过滤器
func (b *Bloom) Purge() {
	b.filter = make(map[string]bool)
}
