package skiplist

import (
	"fmt"
	"math/rand"
	"sync"
)

type SkipList struct {
	mtx          sync.RWMutex // 更新数据的读写锁
	maxLevel     int          // 跳表最大高度
	currentLevel int          // 跳表当前高度
	elementCount int          // 跳表中元素个数
	head         *Node        // 跳表头节点指针
}

// 随机生成节点的高度值
func (l *SkipList) getRandomLevel() int {
	// 初始高度为 1
	level := 1
	// 每次有 1/2 的概率将高度增加一层
	for rand.Int()%2 == 1 {
		level++
	}
	// 如果大于最大高度，则将其设置为最大高度
	if level > l.maxLevel {
		level = l.maxLevel
	}
	// 返回结果
	return level
}

// NewSkipList 新建一个跳表
func NewSkipList(maxLevel int) *SkipList {
	return &SkipList{
		mtx:          sync.RWMutex{},
		maxLevel:     maxLevel,
		currentLevel: 0,
		elementCount: 0,
		head:         NewNode("", "", maxLevel),
	}
}

// Put 更新跳表元素，不存在则新建
func (l *SkipList) Put(key string, value interface{}) {
	// 更新数据需要获取写锁
	l.mtx.Lock()
	defer l.mtx.Unlock()

	// 当前节点
	currentNode := l.head
	// 需要更新 forward 的节点数组
	updateNode := make([]*Node, l.maxLevel+1)

	// 遍历每一层，找到这一层中最后一个小于待插入节点的位置
	for i := l.currentLevel; i >= 0; i-- {
		for currentNode.forward[i] != nil && currentNode.forward[i].key < key {
			currentNode = currentNode.forward[i]
		}
		// 记录待更新节点位置
		updateNode[i] = currentNode
	}

	// 目前处于第 0 层，也即全体数据的有序链表，因此需要转移到下一个节点，便于判断元素是否已经存在
	currentNode = currentNode.forward[0]

	// 元素已经存在
	if currentNode != nil && currentNode.key == key {
		// 从底层向上遍历
		for i := 0; i <= l.currentLevel; i++ {
			// 如果发现某一层待更新节点的下一个节点不是更新节点，那么直接退出循环，说明在更高层中也没有更新节点了
			if updateNode[i].forward[i].key != key {
				break
			}
			// 更新节点对应的 value
			updateNode[i].forward[i].value = value
		}

		return
	}

	// 获取新元素高度
	randomLevel := l.getRandomLevel()

	// 如果新元素高度大于跳表高度
	if randomLevel > l.currentLevel {
		// 对于那些超过当前高度的层，待更新节点应当为头节点
		for i := l.currentLevel + 1; i <= randomLevel; i++ {
			updateNode[i] = l.head
		}
		// 更新跳表当前高度
		l.currentLevel = randomLevel
	}

	// 新建节点
	insertNode := NewNode(key, value, randomLevel)

	for i := 0; i <= randomLevel; i++ {
		// 让新节点的下一个节点指向待更新节点的下一个节点
		insertNode.forward[i] = updateNode[i].forward[i]
		// 在待更新节点后面插入新节点
		updateNode[i].forward[i] = insertNode
	}

	// 更新跳表元素个数
	l.elementCount++
}

// 找到小于等于 key 的最后一个节点
func (l *SkipList) search(key string) *Node {
	current := l.head

	for i := l.currentLevel; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key <= key {
			current = current.forward[i]
		}
	}

	return current
}

// Get 获取指定 key 对应的 value
func (l *SkipList) Get(key string) (interface{}, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	node := l.search(key)

	if node.key == key {
		return node.value, nil
	} else {
		return nil, fmt.Errorf("key %s existed", key)
	}
}

// LowerBound 找到小于等于 key 的最后一个节点对应的 key，找不到则返回空字符串
func (l *SkipList) LowerBound(key string) string {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	node := l.search(key)

	return node.key
}

// Delete 删除指定 key
func (l *SkipList) Delete(key string) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	// 当前节点
	currentNode := l.head
	// 需要更新 forward 的节点数组
	updateNode := make([]*Node, l.maxLevel+1)

	// 遍历每一层，找到这一层中最后一个小于待插入节点的位置
	for i := l.currentLevel; i >= 0; i-- {
		for currentNode.forward[i] != nil && currentNode.forward[i].key < key {
			currentNode = currentNode.forward[i]
		}
		// 记录待更新节点位置
		updateNode[i] = currentNode
	}

	currentNode = currentNode.forward[0]
	if currentNode != nil && currentNode.key == key {
		// 从底层向上遍历
		for i := 0; i <= l.currentLevel; i++ {
			// 如果发现某一层待更新节点的下一个节点不是删除节点，那么直接退出循环，说明在更高层中也没有删除节点了
			if updateNode[i].forward[i].key != key {
				break
			}
			// 删除节点
			updateNode[i].forward[i] = currentNode.forward[i]
		}

		// 删除掉没有节点的层
		for l.currentLevel > 0 && l.head.forward[l.currentLevel] == nil {
			l.currentLevel--
		}

		// 更新元素数量
		l.elementCount--
	}
}

// Show 输出跳表到控制台
func (l *SkipList) Show() {
	fmt.Println("*****skiplist*****")
	for i := 0; i <= l.currentLevel; i++ {
		// 找到每一层的第一个节点
		node := l.head.forward[i]

		fmt.Printf("Level %d:\n", i)
		for node != nil {
			fmt.Printf("%s:%d; ", node.key, node.value.(int))
			node = node.forward[i]
		}
		fmt.Println()
	}
}
