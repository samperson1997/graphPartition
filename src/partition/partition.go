package partition

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

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

// Config SHPImpl config
type Config struct {
	VertexSize uint64
	BucketSize uint64
	Prob       float64
	Graph      *Graph
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
	graph *Graph
}

// NewSHPImpl a new shpimpl with Config
func NewSHPImpl(c Config) *SHPImpl {
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
	return &shp
}

func (shp *SHPImpl) calcSingleGain(node *Node) (gains []float64) {
	gains = make([]float64, shp.bucketSize)
	vertex := node.id
	for _, nbrNode := range node.Nbrlist {
		uBucket := shp.vertex2Bucket[node.id]
		nb := shp.nbrBucket[nbrNode]
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			if bucketJ != shp.vertex2Bucket[vertex] {
				gains[bucketJ] += math.Pow(1-shp.prob, float64(nb[bucketJ])) - math.Pow(1-shp.prob, float64(nb[uBucket]-1))
			}
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
				minGain := math.MaxFloat64
				preBucket := shp.vertex2Bucket[vertex]
				shp.vertex2Target[vertex] = preBucket
				target := preBucket
				gains := shp.calcSingleGain(shp.graph.Nodes[vertex])
				for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
					gain := gains[bucketI]
					if gain < minGain {
						minGain = gain
						target = bucketI
					}
				}
				if minGain < 0 {
					shp.vertex2Target[vertex] = target
					atomic.AddUint64(&shp.vertexTrans[preBucket][target], 1)
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
		minGain := math.MaxFloat64
		preBucket := shp.vertex2Bucket[vertex]
		shp.vertex2Target[vertex] = preBucket
		target := preBucket
		gains := shp.calcSingleGain(shp.graph.Nodes[vertex])
		for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
			gain := gains[bucketI]
			if gain < minGain {
				minGain = gain
				target = bucketI
			}
		}
		if minGain < 0 {
			shp.vertex2Target[vertex] = target
			shp.vertexTrans[preBucket][target]++
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
	fmt.Println("setnew number :", number)
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

func (shp *SHPImpl) computeBucketSingle(node *Node) []int {
	nb := make([]int, shp.bucketSize)
	for _, nbrNode := range node.Nbrlist {
		uBucket := shp.vertex2Bucket[nbrNode]
		nb[uBucket]++
	}
	shp.nbrBucket[node.id] = nb
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
func NextIterationParallel(shp *SHPImpl) bool {
	log.Println("[process new NextIterationParallel]")
	shp.PreComputeBucket()
	shp.ComputMoveGainParallel()
	shp.ComputMoveProb()
	return shp.SetNewParallel()
}
