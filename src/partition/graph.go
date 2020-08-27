package partition

import "container/list"

// Graph manage nodes
type Graph struct {
	nodes []*Node
}

// NewGraph return a empty graph with vertexSize vertex
func NewGraph(vertexSize int) *Graph {
	g := Graph{
		nodes: make([]*Node, vertexSize),
	}
	for i := 0; i < vertexSize; i++ {
		g.nodes[i] = NewNode(uint64(i))
	}
	return &g
}

// AddEdge Add edge from src to dst
func (g *Graph) AddEdge(src, dst uint64) {
	g.nodes[src].addNbr(dst)
}

type nbrNode struct {
	id uint64
}

// Node a vertex and nbrlist
type Node struct {
	id      uint64
	nbrlist *list.List
}

// NewNode returns a node with id and empty list
func NewNode(id uint64) *Node {
	n := Node{
		id:      id,
		nbrlist: list.New(),
	}
	n.nbrlist.Init()
	return &n
}
func (n *Node) addNbr(id uint64) {
	n.nbrlist.PushBack(nbrNode{id})
}
