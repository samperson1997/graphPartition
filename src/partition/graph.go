package partition

import (
	"math"
)

// Graph manage Nodes
type Graph struct {
	Nodes []*Node
	// vertex from 0~vertexSize
	vertexSize uint64
}

// NewGraph return a empty graph with vertexSize vertex
func NewGraph(vertexSize uint64) *Graph {
	g := Graph{
		Nodes: make([]*Node, vertexSize),
		vertexSize:vertexSize,
	}
	for i := uint64(0); i < vertexSize; i++ {
		g.Nodes[i] = NewNode(i)
	}
	return &g
}

// GetVertexSize return vertexsize 
func (g *Graph) GetVertexSize()uint64{
	return g.vertexSize
}

// AddEdge Add edge from src to dst
func (g *Graph) AddEdge(src, dst uint64) {
	g.Nodes[src].addNbr(dst)
}

// ChangeColor change node color
func (g *Graph) ChangeColor(id, color uint64) {
	g.Nodes[id].Color = color
}

// Node a vertex and Nbrlist
type Node struct {
	ID      uint64
	Nbrlist []uint64
	Color   uint64
}

// NewNode returns a node with id and empty list
func NewNode(id uint64) *Node {
	n := Node{
		ID:      id,
		Nbrlist: make([]uint64, 0),
		Color:   math.MaxUint64,
	}
	return &n
}

func (n *Node) addNbr(id uint64) {
	n.Nbrlist = append(n.Nbrlist, id)
}
