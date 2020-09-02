package bdg

import (
	"container/list"
	"fmt"
	. "gpartition/common"
	"math"
	"math/rand"
	"sort"
	"time"
)

type BDGConfig struct {
	VertexSize  uint64
	SrcNodesNum uint64
	StepNum     uint64
	BucketSize  uint64
	Graph       *Graph
}

type Block struct {
	id      uint64
	nodes   *list.List
	nbrlist map[uint64]bool
}

// BDGImpl calc BDG partition
type BDGImpl struct {
	srcNodesNum uint64
	stepNum     uint64
	vertexSize  uint64
	bucketSize  uint64

	blocks  []*Block
	buckets []*list.List

	// graph manage all graph data
	graph         *Graph
	vertex2Bucket []uint64
}

func NewBlock(id uint64) *Block {
	block := Block{
		id:      id,
		nodes:   list.New(),
		nbrlist: make(map[uint64]bool),
	}
	block.nodes.Init()
	return &block
}

func (block *Block) addBlockNbr(id uint64) {
	block.nbrlist[id] = true
}

// NewBDGImpl a new bdgimpl with Config
func NewBDGImpl(c BDGConfig) *BDGImpl {
	bdg := BDGImpl{
		graph:         c.Graph,
		srcNodesNum:   c.SrcNodesNum,
		stepNum:       c.StepNum,
		vertexSize:    c.VertexSize,
		bucketSize:    c.BucketSize,
		buckets:       make([]*list.List, c.BucketSize),
		vertex2Bucket: make([]uint64, c.VertexSize),
	}
	for i := range bdg.buckets {
		bdg.buckets[i] = list.New()
		bdg.buckets[i].Init()
	}

	return &bdg
}

// bfs cut an input graph into fine-grained blocks
func (bdg *BDGImpl) bfs() {
	g := bdg.graph.Nodes
	nodesToVisit := bdg.vertexSize
	itr := 0

	for nodesToVisit > 0 {
		var queue = list.New()
		queue.Init()
		chosenSrc := make(map[uint64]bool, bdg.srcNodesNum)

		// add random source nodes to blocks and change color
		for i := uint64(0); i < bdg.srcNodesNum; i++ {
			rand.Seed(time.Now().UnixNano())
			srcId := uint64(rand.Intn(int(bdg.vertexSize)))
			_, ok := chosenSrc[srcId]
			if ok || g[srcId].Color != math.MaxUint64 {
				i--
				continue
			}
			chosenSrc[srcId] = true
			queue.PushBack(srcId)
			bdg.graph.ChangeColor(srcId, i)
			blockId := uint64(itr)*bdg.srcNodesNum + i
			bdg.blocks = append(bdg.blocks, NewBlock(blockId))
			bdg.blocks[blockId].nodes.PushBack(srcId)
			nodesToVisit--
		}

		step := uint64(0)
		for node := queue.Front(); node != nil && queue.Len() > 0 && step < bdg.stepNum; node = node.Next() {
			color := g[node.Value.(uint64)].Color
			for _, nbrNode := range g[node.Value.(uint64)].Nbrlist {
				if g[nbrNode].Color == math.MaxUint64 {
					g[nbrNode].Color = color
					// add neighbor node into block
					bdg.blocks[color].nodes.PushBack(nbrNode)
					queue.PushBack(nbrNode)
					nodesToVisit--
				} else if g[nbrNode].Color != color {
					bdg.blocks[color].addBlockNbr(g[nbrNode].Color)
				}
			}
			step++
		}
		itr++
	}
}

// deterministicGreedy deterministic greedy strategy
func (bdg *BDGImpl) deterministicGreedy() {
	capacity := float64(bdg.vertexSize) / float64(bdg.bucketSize)

	// sort buckets by nodes size
	sort.Slice(bdg.blocks, func(i, j int) bool {
		return bdg.blocks[i].nodes.Len() > bdg.blocks[j].nodes.Len()
	})

	// ===== print blocks info ======
	bdg.printBlocksInfo()

	// add bucket num of blocks into buckets firstly
	for i := 0; i < int(bdg.bucketSize); i++ {
		bdg.buckets[i].PushBack(bdg.blocks[i].id)
	}
	for k := int(bdg.bucketSize); k < len(bdg.blocks); k++ {
		block := bdg.blocks[k]
		var bset = make(map[uint64]bool)
		for nbr := range block.nbrlist {
			// add nodes in nbr block into bset
			for node := bdg.blocks[nbr].nodes.Front(); node != nil; node = node.Next() {
				bset[node.Value.(uint64)] = true
			}
		}
		j := 0
		tmp := 0.0
		for i := 0; i < len(bdg.buckets); i++ {
			var pset = make(map[uint64]bool)
			for blockInWorker := bdg.buckets[i].Front(); blockInWorker != nil; blockInWorker = blockInWorker.Next() {
				// add nodes in block in worker into pset
				for node := bdg.blocks[blockInWorker.Value.(uint64)].nodes.Front(); node != nil; node = node.Next() {
					pset[node.Value.(uint64)] = true
				}
			}
			retainSize := 0
			for key := range pset {
				_, ok := bset[key]
				if ok {
					retainSize++
				}
			}
			if float64(retainSize)*(1-float64(len(pset))/capacity) > tmp {
				j = i
			}
		}
		bdg.buckets[j].PushBack(block.id)
	}
}

func (bdg *BDGImpl) printBlocksInfo() {
	for k, v := range bdg.blocks {
		fmt.Print("Key:", k, ",nodes:", v.nodes.Len(), "(")
		for node := v.nodes.Front(); node != nil; node = node.Next() {
			fmt.Print(node.Value, ",")
		}
		fmt.Println(")")
		fmt.Print("Neighbors:")
		for k1 := range v.nbrlist {
			fmt.Print(k1, ",")
		}
		fmt.Println()
		fmt.Println("=======")
	}
}

func (bdg *BDGImpl) Calc() {
	bdg.bfs()
	bdg.deterministicGreedy()
}

func (bdg *BDGImpl) AfterCalc() {
	for i := range bdg.buckets {
		for blockInWorker := bdg.buckets[i].Front(); blockInWorker != nil; blockInWorker = blockInWorker.Next() {
			for node := bdg.blocks[blockInWorker.Value.(uint64)].nodes.Front(); node != nil; node = node.Next() {
				bdg.vertex2Bucket[node.Value.(uint64)] = uint64(i)
			}
		}
	}
}

func (bdg *BDGImpl) GetBucketFromId(id uint64) uint64 {
	if id > bdg.vertexSize {
		return math.MaxUint64
	}
	return bdg.vertex2Bucket[id]
}

func (bdg *BDGImpl) GetGraph() *Graph {
	return bdg.graph
}

func (bdg *BDGImpl) GetBucketSize() uint64 {
	return bdg.bucketSize
}
