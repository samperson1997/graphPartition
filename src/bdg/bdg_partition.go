package bdg

import (
	"container/list"
	"fmt"
	. "gpartition/partition"
	"math"
	"math/rand"
	"sort"
)

type BDGConfig struct {
	VertexSize uint64
	BlockSize  uint64
	BucketSize uint64
	Graph      *Graph
}

type Block struct {
	id      uint64
	nodes   *list.List
	nbrlist map[uint64]bool
}

// BDGImpl calc BDG partition
type BDGImpl struct {
	blockSize  uint64
	vertexSize uint64
	bucketSize uint64

	blocks  map[uint64]*Block
	buckets []*list.List

	// graph manage all graph data
	graph *Graph
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

// BDGImpl a new bdgimpl with Config
func NewBDGImpl(c BDGConfig) *BDGImpl {
	bdg := BDGImpl{
		graph:      c.Graph,
		blockSize:  c.BlockSize,
		vertexSize: c.VertexSize,
		bucketSize: c.BucketSize,
		blocks:     make(map[uint64]*Block, c.BlockSize),
		buckets:    make([]*list.List, c.BucketSize),
	}
	for i := range bdg.buckets {
		bdg.buckets[i] = list.New()
		bdg.buckets[i].Init()
	}

	return &bdg
}

// bfs cut an input graph into fine-grained blocks
func (bdg *BDGImpl) bfs() {
	var queue = list.New()
	queue.Init()

	chosenSrc := make(map[uint64]bool, bdg.blockSize)

	// add random source nodes to blocks and change color
	for i := uint64(0); i < bdg.blockSize; i++ {
		srcId := uint64(rand.Intn(int(bdg.vertexSize)))
		_, ok := chosenSrc[srcId]
		if ok {
			i--
			continue
		}
		chosenSrc[srcId] = true
		queue.PushBack(srcId)
		bdg.graph.ChangeColor(srcId, i)
		bdg.blocks[i] = NewBlock(i)
		bdg.blocks[i].nodes.PushBack(srcId)
	}

	for queue.Len() > 0 {
		var tmpQueue = list.New()
		tmpQueue.Init()
		g := bdg.graph.Nodes
		for node := queue.Front(); node != nil; node = node.Next() {
			color := g[node.Value.(uint64)].Color
			for _, nbrNode := range g[node.Value.(uint64)].Nbrlist {
				if g[nbrNode].Color == math.MaxUint64 {
					g[nbrNode].Color = color
					// add neighbor node into block
					bdg.blocks[color].nodes.PushBack(nbrNode)
					tmpQueue.PushBack(nbrNode)
				} else if g[nbrNode].Color != color {
					bdg.blocks[color].addBlockNbr(g[nbrNode].Color)
				}
			}
		}
		queue = tmpQueue
	}

	// running CC finding algorithm (like Hash-Min) on uncolored vertices
}

// deterministicGreedy deterministic greedy strategy
func (bdg *BDGImpl) deterministicGreedy() {
	// ===== print blocks info ======
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
	// ====== end print ======

	// sort buckets by nodes size
	sortedBlocks := sortBlocksByNodesNum(bdg.blocks)

	// add first block into first bucket
	bdg.buckets[0].PushBack(sortedBlocks[0].id)
	for i := 1; i < len(sortedBlocks); i++ {
		block := sortedBlocks[i]
		var bset = make(map[uint64]bool)
		var pset = make(map[uint64]bool)
		for nbr := range block.nbrlist {
			bset[nbr] = true
		}
		j := 0.0
		for i := 0; i < len(bdg.buckets); i++ {
			blocksInWorker := bdg.buckets[0]
			for blockInWorker := blocksInWorker.Front(); blockInWorker != nil; blockInWorker = blockInWorker.Next() {
				pset[blockInWorker.Value.(uint64)] = true
			}
			retainSize := 0
			for key := range bset {
				_, ok := pset[key]
				if ok {
					retainSize++
				}
			}
			j = math.Max(j, float64(retainSize*(1-len(pset)/int(bdg.vertexSize))))
		}
		bdg.buckets[int(j)].PushBack(block.id)
	}
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type BlockList []Block

func (p BlockList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p BlockList) Len() int {
	return len(p)
}

func (p BlockList) Less(i, j int) bool {
	return p[i].nodes.Len() > p[j].nodes.Len()
}

// A function to turn a map into a BlockList, then sort and return it.
func sortBlocksByNodesNum(m map[uint64]*Block) BlockList {
	blockList := make(BlockList, len(m))
	i := 0
	for _, v := range m {
		blockList[i] = *v
		i++
	}
	sort.Sort(blockList)
	return blockList
}
