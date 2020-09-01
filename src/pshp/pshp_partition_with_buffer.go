package pshp

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

// SetNewWithBufferParallel parallel check bucket to set
func (shp *SHPImpl) SetNewWithBufferParallel() (ret bool) {
	parallel := uint64(runtime.NumCPU())
	// fmt.Println("parallel with ", parallel, "cpu")
	var isSet atomic.Value
	isSet.Store(false)
	segmentVertexSize := (uint64(shp.tf.bufferSize) + parallel - 1) / parallel
	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < uint64(shp.tf.bufferSize); beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for i := begin; i != end; i++ {
				vertex := shp.tf.buffer[i]
				if rand.Float64() < shp.probability[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] {
					shp.vertex2Bucket[vertex] = shp.vertex2Target[vertex]
					isSet.Store(true)
				}
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
	return isSet.Load().(bool)
}

// SetNewWithBuffer check bucket to set
func (shp *SHPImpl) SetNewWithBuffer() (ret bool) {
	ret = false
	for i := 0; i < int(shp.tf.bufferSize); i++ {
		vertex := shp.tf.buffer[i]
		if shp.vertex2Target[vertex] != shp.vertex2Bucket[vertex] &&
			rand.Float64() < shp.probability[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] {
			shp.vertex2Bucket[vertex] = shp.vertex2Target[vertex]
			ret = true
		}
	}
	return
}

// ComputMoveGainWithBufferParallel parallel compute maxgain of each vertex
func (shp *SHPImpl) ComputMoveGainWithBufferParallel() {
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			shp.vertexTrans[bucketI][bucketJ] = 0
		}
	}
	parallel := uint64(runtime.NumCPU())
	atomic.StoreInt64(&shp.tf.bufferSize, 0)
	// TODO
	// make parallel seperate stable
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
					tp := atomic.AddInt64(&shp.tf.bufferSize, 1)
					shp.tf.buffer[tp-1] = vertex
				}
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
	fmt.Println(shp.tf.bufferSize)
}

// ComputMoveGainWithBuffer compute maxgain of each vertex
func (shp *SHPImpl) ComputMoveGainWithBuffer() {
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			shp.vertexTrans[bucketI][bucketJ] = 0
		}
	}
	for vertex := uint64(0); vertex < shp.vertexSize; vertex++ {
		minGain, target := shp.calcSingleGain(shp.graph.Nodes[vertex])

		if minGain < 0 {
			shp.vertex2Target[vertex] = target
			shp.vertexTrans[shp.vertex2Bucket[vertex]][target]++
			shp.tf.buffer[shp.tf.bufferSize] = vertex
			shp.tf.bufferSize++
		}
	}
}

// NextIterationWithBufferParallel process a iteration with a iteration
func NextIterationWithBufferParallel(shp *SHPImpl) bool {
	shp.PreComputeBucketParallel()
	shp.ComputMoveGainWithBufferParallel()
	shp.ComputMoveProb()
	return shp.SetNewWithBufferParallel()
}
