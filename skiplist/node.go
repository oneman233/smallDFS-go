package skiplist

type Node struct {
	key       string      // 节点的 key 值
	value     interface{} // 节点的 value，可以为任意类型
	nodeLevel int         // 节点高度
	forward   []*Node     // forward[i] 代表该节点在第 i 层对应的 key 相同的节点的下一个节点
}

func NewNode(key string, value interface{}, nodeLevel int) *Node {
	return &Node{
		key:       key,
		value:     value,
		nodeLevel: nodeLevel,
		forward:   make([]*Node, nodeLevel+1),
	}
}
