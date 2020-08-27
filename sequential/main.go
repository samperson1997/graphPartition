package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

func Min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

type bucket struct {
}
type nbrNode struct {
	id uint64
}
type Node struct {
	id      uint64
	nbrlist *list.List
}
type Graph struct {
	nodes []*Node
}

func (g *Graph) addEdge(u, v uint64) {
	g.nodes[u].addNbr(v)
}

type Config struct {
	vertexSize uint64
	bucketSize uint64
	prob       float64
	graph      *Graph
}

type SHPImpl struct {
	graph           *Graph
	buckets         []bucket
	vertex2Bucket   []uint64
	vertex2Target   []uint64
	bucketSize      uint64
	vertexSize      uint64
	probability     [][]float64
	S_bucket_target [][]uint64
	prob            float64
}

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

func NewGraph(vertexSize int) *Graph {
	g := Graph{
		nodes: make([]*Node, vertexSize),
	}
	for i := 0; i < vertexSize; i++ {
		g.nodes[i] = NewNode(uint64(i))
	}
	return &g
}

func NewSHPImpl(c Config) *SHPImpl {
	shp := SHPImpl{
		graph:         c.graph,
		vertex2Bucket: make([]uint64, c.vertexSize),
		vertex2Target: make([]uint64, c.vertexSize),
		bucketSize:    c.bucketSize,
		vertexSize:    c.vertexSize,
		prob:          c.prob,
	}
	b := c.bucketSize
	shp.probability = make([][]float64, c.bucketSize)
	shp.S_bucket_target = make([][]uint64, c.bucketSize)
	arena := make([]float64, b*b)
	for i := range shp.probability {
		shp.probability[i] = arena[i*int(b) : (i+1)*int(b)]
	}
	arena1 := make([]uint64, b*b)
	for i := range shp.S_bucket_target {
		shp.S_bucket_target[i] = arena1[i*int(b) : (i+1)*int(b)]
	}
	return &shp
}
func (shp *SHPImpl) calc_gain(vertex uint64) (gains []float64) {
	ns := make([]uint64, shp.bucketSize)
	gains = make([]float64, shp.bucketSize)

	total := 0
	for nbr := shp.graph.nodes[vertex].nbrlist.Front(); nbr != nil; nbr = nbr.Next() {
		u := nbr.Value.(nbrNode).id
		uBucket := shp.vertex2Bucket[u]
		ns[uBucket]++
		total++
	}
	for bucket_i := uint64(0); bucket_i < shp.bucketSize; bucket_i++ {
		gains[bucket_i] = float64(ns[bucket_i]) * (-shp.prob)
		gains[bucket_i] += float64(ns[shp.vertex2Bucket[vertex]]) * (shp.prob)
	}
	return
}

func (shp *SHPImpl) initBucket() {
	for i := uint64(0); i < shp.vertexSize; i++ {
		shp.vertex2Bucket[i] = rand.Uint64() % shp.bucketSize
	}
}
func (shp *SHPImpl) ComputMoveGain() {
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		max_gain := float64(-1)
		preBucket := shp.vertex2Bucket[vertex]
		target := shp.vertex2Bucket[vertex]
		gains := shp.calc_gain(vertex)
		for bucket_i := uint64(0); bucket_i < shp.bucketSize; bucket_i++ {
			gain := gains[bucket_i]
			if gain > max_gain {
				max_gain = gain
				target = bucket_i
			}
		}
		shp.vertex2Target[vertex] = target
		if max_gain > 0 {
			shp.S_bucket_target[preBucket][target]++
		}

	}
}

func (shp *SHPImpl) ComputMoveProb() {

	for bucket_i := 0; uint64(bucket_i) < shp.bucketSize; bucket_i++ {
		for bucket_j := 0; uint64(bucket_j) < shp.bucketSize; bucket_j++ {
			if shp.S_bucket_target[bucket_i][bucket_j] != 0 {
				shp.probability[bucket_i][bucket_j] = float64(Min(shp.S_bucket_target[bucket_i][bucket_j], shp.S_bucket_target[bucket_j][bucket_i])) / float64(shp.S_bucket_target[bucket_i][bucket_j])
			} else {
				shp.probability[bucket_i][bucket_j] = 0
			}
		}
	}
}

func (shp *SHPImpl) SetNew() {

	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		if shp.vertex2Target[vertex] != shp.vertex2Bucket[vertex] &&
			rand.Float64() < shp.probability[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] {
			shp.vertex2Bucket[vertex] = shp.vertex2Target[vertex]
		}
	}
	for bucket_i := uint64(0); bucket_i < shp.bucketSize; bucket_i++ {
		for bucket_j := uint64(0); bucket_j < shp.bucketSize; bucket_j++ {
			shp.S_bucket_target[bucket_i][bucket_j] = 0
		}
	}

}

func (shp *SHPImpl) PrintResult() {

	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		fmt.Println("vertex:", vertex, " bucket:", shp.vertex2Bucket[vertex])
	}

}

//LoadGraph :
//<VertexNumber>
//<edge1Src> <edge1Dst>
func LoadGraph(path string) (c Config) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("err=%v", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	str, err := reader.ReadString('\n')
	fmt.Fscanf(strings.NewReader(str), "%d", &c.vertexSize)
	c.graph = NewGraph(int(c.vertexSize))
	var src, dst uint64
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Fscanf(strings.NewReader(str), "%d %d", &src, &dst)
		c.graph.addEdge(src, dst)
		c.graph.addEdge(dst, src)
	}
	c.bucketSize = 5
	c.prob = 0.5
	return
}
func NextIteration(shp *SHPImpl) {
	shp.ComputMoveGain()
	shp.ComputMoveProb()
	shp.SetNew()
}
func main() {
	config := LoadGraph("data.in")

	shp := NewSHPImpl(config)
	shp.initBucket()
	maxIteration := 100
	for i := 0; i < maxIteration; i++ {
		NextIteration(shp)
	}
	shp.PrintResult()
}
