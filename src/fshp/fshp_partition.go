package fshp

import (
	"fmt"
	"gpartition/common"
	"math"
	"runtime"
	"sync"
)

// FSHPConfig SHPImpl config
type FSHPConfig struct {
	VertexSize uint64
	BucketSize uint64
	Prob       float64
	Graph      *common.Graph
}

// FSHPImpl calc SHP partition
type FSHPImpl struct {
	bucketSize uint64
	vertexSize uint64
	prob       float64

	//vertex2Bucket
	vertex2Bucket []uint64
	vertex2Target []uint64

	//pre calculation Bucket
	nbrBucket [][]int32

	probability [][]float64
	vertexTrans [][]uint64

	buffer     []uint64
	buffer2    []uint64
	bufferSize int32
	// graph manage all graph data
	graph *common.Graph

	mayChangeSize int32
	mayChange     []uint64

	notChangeSize int32
	// not changed
	notChanged    []uint64
	notChangedVis []bool

	// use and clear
	vis  []int32
	vis1 []int32
}

// NewFSHPImpl a new shpimpl with Config
func NewFSHPImpl(c FSHPConfig) *FSHPImpl {
	fshp := FSHPImpl{
		graph:         c.Graph,
		vertex2Bucket: make([]uint64, c.VertexSize),
		vertex2Target: make([]uint64, c.VertexSize),
		bucketSize:    c.BucketSize,
		vertexSize:    c.VertexSize,
		prob:          c.Prob,
		vis:           make([]int32, c.VertexSize),
		vis1:          make([]int32, c.VertexSize),

		buffer:  make([]uint64, c.VertexSize),
		buffer2: make([]uint64, c.VertexSize),
	}
	fmt.Println(c.BucketSize)
	fshp.bufferSize = 0
	fshp.notChangedVis = make([]bool, c.VertexSize)
	fshp.notChanged = make([]uint64, c.VertexSize)
	fshp.probability = make([][]float64, c.BucketSize)
	fshp.vertexTrans = make([][]uint64, c.BucketSize)
	b := c.BucketSize
	arena := make([]float64, b*b)
	for i := range fshp.probability {
		fshp.probability[i] = arena[i*int(b) : (i+1)*int(b)]
	}
	arena1 := make([]uint64, b*b)
	for i := range fshp.vertexTrans {
		fshp.vertexTrans[i] = arena1[i*int(b) : (i+1)*int(b)]
	}
	fshp.nbrBucket = make([][]int32, c.VertexSize)
	return &fshp
}

func (fshp *FSHPImpl) ComputMoveProb() {
	for bucketI := 0; uint64(bucketI) < fshp.bucketSize; bucketI++ {
		for bucketJ := 0; uint64(bucketJ) < fshp.bucketSize; bucketJ++ {
			if fshp.vertexTrans[bucketI][bucketJ] != 0 {
				fshp.probability[bucketI][bucketJ] = float64(min(fshp.vertexTrans[bucketI][bucketJ], fshp.vertexTrans[bucketJ][bucketI])) / float64(fshp.vertexTrans[bucketI][bucketJ])
			} else {
				fshp.probability[bucketI][bucketJ] = 0
			}
		}
	}
}

func (fshp *FSHPImpl) calcSingleGain(node *common.Node) (minGain float64, target uint64) {
	minGain = 0.1
	// init target is not change
	target = fshp.vertex2Bucket[node.ID]
	gains := make([]float64, fshp.bucketSize)
	vertex := node.ID
	for _, nbrNode := range node.Nbrlist {
		uBucket := fshp.vertex2Bucket[node.ID]
		nb := fshp.nbrBucket[nbrNode]
		for bucketJ := uint64(0); bucketJ < fshp.bucketSize; bucketJ++ {
			if bucketJ != fshp.vertex2Bucket[vertex] {
				gains[bucketJ] += math.Pow(1-fshp.prob, float64(nb[bucketJ])) - math.Pow(1-fshp.prob, float64(nb[uBucket]-1))
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

// PrintResult print all result
func (fshp *FSHPImpl) PrintResult() {
	for vertex := uint64(0); vertex < fshp.vertexSize; vertex++ {
		fmt.Println("vertex:", vertex, " bucket:", fshp.vertex2Bucket[vertex])
	}

}

func (fshp *FSHPImpl) computeBucketSingle(node *common.Node) []int32 {
	nb := make([]int32, fshp.bucketSize)
	for _, nbrNode := range node.Nbrlist {
		uBucket := fshp.vertex2Bucket[nbrNode]
		nb[uBucket]++
	}
	//TODO
	//fshp.nbrBucket[node.ID] = nb
	return nb
}

// PreComputeBucket calc bucket of every node
func (fshp *FSHPImpl) PreComputeBucket() {
	//TODO make it parallel
	for vertex := uint64(0); vertex < fshp.vertexSize; vertex++ {
		fshp.nbrBucket[vertex] = fshp.computeBucketSingle(fshp.graph.Nodes[vertex])
	}
}

// PreComputeBucketParallel calc bucket of every node
func (fshp *FSHPImpl) PreComputeBucketParallel() {
	//TODO make it parallel
	parallel := uint64(runtime.NumCPU())
	segmentVertexSize := (fshp.vertexSize + parallel - 1) / parallel
	var wg sync.WaitGroup

	for beginvertex := uint64(0); beginvertex < fshp.vertexSize; beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for vertex := begin; vertex != end; vertex++ {
				fshp.nbrBucket[vertex] = fshp.computeBucketSingle(fshp.graph.Nodes[vertex])
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, fshp.vertexSize))
	}
	wg.Wait()
}

func (fshp *FSHPImpl) calcSingleFanout(vertex uint64) (fanout float64) {
	ns := make([]uint64, fshp.bucketSize)
	fanout = 0
	for _, nbrNode := range fshp.graph.Nodes[vertex].Nbrlist {
		uBucket := fshp.vertex2Bucket[nbrNode]
		ns[uBucket]++
	}
	for bucketI := uint64(0); bucketI < fshp.bucketSize; bucketI++ {
		if ns[bucketI] > 0 {
			fanout++
		}
	}
	return
}
func (fshp *FSHPImpl) calcSinglepFanout(vertex uint64) (fanout float64) {
	ns := make([]uint64, fshp.bucketSize)
	fanout = 0
	for _, nbrNode := range fshp.graph.Nodes[vertex].Nbrlist {
		uBucket := fshp.vertex2Bucket[nbrNode]
		ns[uBucket]++
	}
	for bucketI := uint64(0); bucketI < fshp.bucketSize; bucketI++ {
		fanout += 1 - math.Pow(1-fshp.prob, float64(ns[bucketI]))

	}
	return
}

// CalcFanout for test
func (shp *FSHPImpl) CalcFanout() (fanout float64) {
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		fanout += shp.calcSingleFanout(vertex)
	}
	return
}

// GetEachBucketVolumn get each bucket volumn
func (shp *FSHPImpl) GetEachBucketVolumn() {
	g := shp.graph
	ns := make([]int, shp.bucketSize)

	for vertex := uint64(0); vertex < g.GetVertexSize(); vertex++ {
		uBucket := shp.vertex2Bucket[vertex]
		ns[uBucket]++
	}

	for bucketI, size := range ns {
		fmt.Println(bucketI, size)
	}
}
