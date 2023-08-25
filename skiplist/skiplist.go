package skiplist

import (
	"fmt"
	"sync"
)

type SkipList struct {
	mtx          sync.RWMutex
	maxLevel     int
	currentLevel int
	elementCount int
	head         *Node
}

func (l *SkipList) getRandomLevel() int {
	return 0
}

func (l *SkipList) Insert(key string, value interface{}) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	currentNode := l.head
	updateNode := make([]*Node, l.maxLevel+1)

	for i := l.currentLevel; i >= 0; i-- {
		for currentNode.forward[i] != nil && currentNode.forward[i].key < key {
			currentNode = currentNode.forward[i]
		}
		updateNode[i] = currentNode
	}

	currentNode = currentNode.forward[0]

	if currentNode != nil && currentNode.key == key {
		return fmt.Errorf("key %s existed", key)
	}

	return nil
}
