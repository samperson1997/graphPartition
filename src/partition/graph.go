package partition

import (
	"math"
)

// Graph manage Nodes
type Graph struct {
	Nodes []*Node
}

// NewGraph return a empty graph with vertexSize vertex
func NewGraph(vertexSize int) *Graph {
	g := Graph{
		Nodes: make([]*Node, vertexSize),
	}
	for i := 0; i < vertexSize; i++ {
		g.Nodes[i] = NewNode(uint64(i))
	}
	return &g
}

// AddEdge Add edge from src to dst
func (g *Graph) AddEdge(src, dst uint64) {
	g.Nodes[src].addNbr(dst)
}

func (g *Graph) ChangeColor(id, color uint64) {
	g.Nodes[id].Color = color
}

type NbrNode struct {
	Id uint64
}

// Node a vertex and Nbrlist
type Node struct {
	id      uint64
	Nbrlist []uint64
	Color   uint64
}

// NewNode returns a node with id and empty list
func NewNode(id uint64) *Node {
	n := Node{
		id:      id,
		Nbrlist: make([]uint64, 0),
		Color:   math.MaxUint64,
	}
	return &n
}

func (n *Node) addNbr(id uint64) {
	n.Nbrlist = append(n.Nbrlist, id)
}
