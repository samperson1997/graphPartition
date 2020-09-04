package pshp

import (
	"fmt"
	"gpartition/common"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

// LoadGraph load a graph with path
func LoadGraph(path string, BucketSize int, prob float64) (c SHPConfig, err error) {
	c.Graph, err = common.LoadGraphFromPath(path)
	if err != nil {
		return SHPConfig{}, err
	}
	c.VertexSize = c.Graph.GetVertexSize()
	c.Prob = defaultProb
	c.BucketSize = defaultBucketSize
	if BucketSize >= 1 {
		c.BucketSize = uint64(BucketSize)
	}
	if prob <= 1 && prob >= 0 {
		c.Prob = prob
	}
	return c, nil
}

type bucket struct {
}

type transfer struct {
	from uint64
	to   uint64
}

type transferNeed struct {
	buffer     []uint64
	bufferSize int64
}

// SHPConfig SHPImpl config
type SHPConfig struct {
	VertexSize uint64
	BucketSize uint64
	Prob       float64
	Graph      *common.Graph
}

// SHPImpl calc SHP partition
type SHPImpl struct {
	bucketSize uint64
	vertexSize uint64
	prob       float64

	//vertex2Bucket
	vertex2Bucket []uint64
	vertex2Target []uint64

	//pre calculation Bucket
	nbrBucket [][]int

	probability [][]float64
	vertexTrans [][]uint64

	tf transferNeed
	// graph manage all graph data
	graph *common.Graph

	//for sort version
	needTrans          [][][]uint64
	mutexBucket2Bucket [][]sync.Mutex
	gains              [][]float64
	toBeTransed        []bool
}

// NewSHPImpl a new shpimpl with Config
func NewSHPImpl(c SHPConfig) *SHPImpl {
	shp := SHPImpl{
		graph:         c.Graph,
		vertex2Bucket: make([]uint64, c.VertexSize),
		vertex2Target: make([]uint64, c.VertexSize),
		bucketSize:    c.BucketSize,
		vertexSize:    c.VertexSize,
		prob:          c.Prob,
	}
	shp.tf.buffer = make([]uint64, c.VertexSize)
	shp.tf.bufferSize = 0
	shp.probability = make([][]float64, c.BucketSize)
	shp.vertexTrans = make([][]uint64, c.BucketSize)
	b := c.BucketSize
	arena := make([]float64, b*b)
	for i := range shp.probability {
		shp.probability[i] = arena[i*int(b) : (i+1)*int(b)]
	}
	arena1 := make([]uint64, b*b)
	for i := range shp.vertexTrans {
		shp.vertexTrans[i] = arena1[i*int(b) : (i+1)*int(b)]
	}
	shp.nbrBucket = make([][]int, c.VertexSize)

	//for sort
	shp.needTrans = make([][][]uint64, c.BucketSize)
	for i := range shp.needTrans {
		shp.needTrans[i] = make([][]uint64, c.BucketSize)
	}
	shp.mutexBucket2Bucket = make([][]sync.Mutex, c.BucketSize)
	for i := range shp.mutexBucket2Bucket {
		shp.mutexBucket2Bucket[i] = make([]sync.Mutex, c.BucketSize)
		fmt.Println(len(shp.mutexBucket2Bucket))
		//need extra shp.mutexBucket2Bucket[i] = append(shp.mutexBucket2Bucket[i], mu)?????
	}
	shp.gains = make([][]float64, c.VertexSize)
	for i := range shp.gains {
		shp.gains[i] = make([]float64, c.BucketSize)
	}
	shp.toBeTransed = make([]bool, c.VertexSize)
	return &shp
}

func (shp *SHPImpl) GetVertexSize() uint64 {
	return shp.vertexSize
}

func (shp *SHPImpl) calcSingleGain(node *common.Node) (minGain float64, target uint64) {
	minGain = 0.1
	preBucket := shp.vertex2Bucket[node.ID]
	shp.vertex2Target[node.ID] = preBucket
	// init target is not change
	target = shp.vertex2Bucket[node.ID]
	gains := make([]float64, shp.bucketSize)
	vertex := node.ID
	for _, nbrNode := range node.Nbrlist {
		uBucket := shp.vertex2Bucket[node.ID]
		nb := shp.nbrBucket[nbrNode]
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			if bucketJ != shp.vertex2Bucket[vertex] {
				gains[bucketJ] += math.Pow(1-shp.prob, float64(nb[bucketJ])) - math.Pow(1-shp.prob, float64(nb[uBucket]-1))
			}
		}
	}
	for bucketI, gain := range gains {
		if gain < minGain {
			minGain = gain
			target = uint64(bucketI)
		}
	}
	return
}

// ComputMoveGainParallel parallel compute maxgain of each vertex
func (shp *SHPImpl) ComputMoveGainParallel() {
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			shp.vertexTrans[bucketI][bucketJ] = 0
		}
	}
	parallel := uint64(runtime.NumCPU())

	segmentVertexSize := (shp.vertexSize + parallel - 1) / parallel
	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < shp.vertexSize; beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for vertex := begin; vertex != end; vertex++ {
				minGain, target := shp.calcSingleGain(shp.graph.Nodes[vertex])
				if minGain < 0 {
					shp.vertex2Target[vertex] = target
					atomic.AddUint64(&shp.vertexTrans[shp.vertex2Bucket[vertex]][target], 1)
				}
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
}

// InitBucket set every vertex a init bucket
func (shp *SHPImpl) InitBucket() {
	for i := uint64(0); i < shp.vertexSize; i++ {
		shp.vertex2Bucket[i] = rand.Uint64() % shp.bucketSize
	}
}

// ComputMoveGain compute maxgain of each vertex
func (shp *SHPImpl) ComputMoveGain() {
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			shp.vertexTrans[bucketI][bucketJ] = 0
		}
	}
	changeNumber := 0
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		minGain, target := shp.calcSingleGain(shp.graph.Nodes[vertex])
		if minGain < 0 {
			shp.vertex2Target[vertex] = target
			shp.vertexTrans[shp.vertex2Bucket[vertex]][target]++
			changeNumber++
		}
	}
	log.Println("vertex can process :", changeNumber)
}

// ComputMoveProb compute very probility from bucket to target
func (shp *SHPImpl) ComputMoveProb() {

	for bucketI := 0; uint64(bucketI) < shp.bucketSize; bucketI++ {
		for bucketJ := 0; uint64(bucketJ) < shp.bucketSize; bucketJ++ {
			if shp.vertexTrans[bucketI][bucketJ] != 0 {
				shp.probability[bucketI][bucketJ] = float64(min(shp.vertexTrans[bucketI][bucketJ], shp.vertexTrans[bucketJ][bucketI])) / float64(shp.vertexTrans[bucketI][bucketJ])
			} else {
				shp.probability[bucketI][bucketJ] = 0
			}
		}
	}
}

func (shp *SHPImpl) setNewSegment(begin, end uint64) {
	for vertex := begin; vertex != end; vertex++ {
		if shp.vertex2Target[vertex] != shp.vertex2Bucket[vertex] &&
			rand.Float64() < shp.probability[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] {
			shp.vertex2Bucket[vertex] = shp.vertex2Target[vertex]
		}
	}
}

// SetNewParallel parallel check bucket to set
func (shp *SHPImpl) SetNewParallel() (ret bool) {
	parallel := uint64(runtime.NumCPU())
	// fmt.Println("parallel with ", parallel, "cpu")
	var isSet atomic.Value
	isSet.Store(false)
	segmentVertexSize := (shp.vertexSize + parallel - 1) / parallel
	var wg sync.WaitGroup
	number := int64(0)
	for beginvertex := uint64(0); beginvertex < shp.vertexSize; beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for vertex := begin; vertex != end; vertex++ {
				if shp.vertex2Target[vertex] != shp.vertex2Bucket[vertex] &&
					rand.Float64() < shp.probability[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] {
					shp.vertex2Bucket[vertex] = shp.vertex2Target[vertex]
					isSet.Store(true)
					atomic.AddInt64(&number, 1)
				}
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
	return isSet.Load().(bool)
}

// SetNew check bucket to set
func (shp *SHPImpl) SetNew() (ret bool) {
	ret = false
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		if shp.vertex2Target[vertex] != shp.vertex2Bucket[vertex] &&
			rand.Float64() < shp.probability[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] {
			shp.vertex2Bucket[vertex] = shp.vertex2Target[vertex]
			ret = true
		}
	}
	return
}

// PrintResult print all result
func (shp *SHPImpl) PrintResult() {
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		fmt.Println("vertex:", vertex, " bucket:", shp.vertex2Bucket[vertex])
	}

}
func (shp *SHPImpl) calcSingleFanout(vertex uint64) (fanout float64) {
	ns := make([]uint64, shp.bucketSize)
	fanout = 0
	for _, nbrNode := range shp.graph.Nodes[vertex].Nbrlist {
		uBucket := shp.vertex2Bucket[nbrNode]
		ns[uBucket]++
	}
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		if ns[bucketI] > 0 {
			fanout++
		}
	}
	return
}
func (shp *SHPImpl) calcSingleFanout2(vertex uint64) (fanout float64) {
	ns := make([]uint64, shp.bucketSize)
	fanout = 0
	for _, nbrNode := range shp.graph.Nodes[vertex].Nbrlist {
		uBucket := shp.vertex2Bucket[nbrNode]
		ns[uBucket]++
	}
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		if ns[bucketI] > 0 {
			fanout++
		}
	}
	if uint64(fanout) == 0 && uint64(fanout) != uint64(len(shp.graph.Nodes[vertex].Nbrlist)) {
		fmt.Println("error shows!!!")
	}
	return
}
func (shp *SHPImpl) calcSinglepFanout(vertex uint64) (fanout float64) {
	ns := make([]uint64, shp.bucketSize)
	fanout = 0
	for _, nbrNode := range shp.graph.Nodes[vertex].Nbrlist {
		uBucket := shp.vertex2Bucket[nbrNode]
		ns[uBucket]++
	}
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		fanout += 1 - math.Pow(1-shp.prob, float64(ns[bucketI]))

	}
	return
}

// CalcFanout for test
func (shp *SHPImpl) CalcFanout() (fanout float64) {
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		fanout += shp.calcSingleFanout(vertex)
	}
	return
}

func (shp *SHPImpl) CalcFanout2() (fanout float64) {
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		fanout += shp.calcSingleFanout2(vertex)
	}
	return
}

func (shp *SHPImpl) computeBucketSingle(node *common.Node) []int {
	nb := make([]int, shp.bucketSize)
	for _, nbrNode := range node.Nbrlist {
		uBucket := shp.vertex2Bucket[nbrNode]
		nb[uBucket]++
	}
	//TODO
	//shp.nbrBucket[node.ID] = nb
	return nb
}

// PreComputeBucket calc bucket of every node
func (shp *SHPImpl) PreComputeBucket() {
	//TODO make it parallel
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		shp.nbrBucket[vertex] = shp.computeBucketSingle(shp.graph.Nodes[vertex])
	}
}

// PreComputeBucketParallel calc bucket of every node
func (shp *SHPImpl) PreComputeBucketParallel() {
	//TODO make it parallel
	parallel := uint64(runtime.NumCPU())
	segmentVertexSize := (shp.vertexSize + parallel - 1) / parallel
	var wg sync.WaitGroup

	for beginvertex := uint64(0); beginvertex < shp.vertexSize; beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for vertex := begin; vertex != end; vertex++ {
				shp.nbrBucket[vertex] = shp.computeBucketSingle(shp.graph.Nodes[vertex])
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
}

// NextIteration calc a iteration
func NextIteration(shp *SHPImpl) bool {
	log.Println("[process new NextIteration]")
	shp.PreComputeBucket()
	shp.ComputMoveGain()
	shp.ComputMoveProb()
	return shp.SetNew()
}

// NextIterationParallel process a iteration with a iteration
func NextIterationParallel(shp *SHPImpl) (ret bool) {
	shp.PreComputeBucket()
	shp.ComputMoveGainParallel()
	shp.ComputMoveProb()
	ret = shp.SetNewParallel()
	return
}

// Calc calc all
func (shp *SHPImpl) Calc() {
	shp.InitBucket()
	iter := 0
	for NextIterationWithSortParallel(shp) && iter < 100 {
		iter++
	}
}

func (shp *SHPImpl) GetBucketFromId(id uint64) uint64 {
	if id > shp.vertexSize {
		return math.MaxUint64
	}
	return shp.vertex2Bucket[id]
}
func (shp *SHPImpl) GetGraph() *common.Graph {
	return shp.graph
}
func (shp *SHPImpl) GetBucketSize() uint64 {
	return shp.bucketSize
}
func (shp *SHPImpl) AfterCalc() {

}
