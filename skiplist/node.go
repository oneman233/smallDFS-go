package skiplist

type Node struct {
	key       string
	value     interface{}
	nodeLevel int
	forward   []*Node
}

func NewNode(key string, value interface{}, nodeLevel int) *Node {
	return &Node{
		key:       key,
		value:     value,
		nodeLevel: nodeLevel,
		forward:   make([]*Node, nodeLevel),
	}
}
