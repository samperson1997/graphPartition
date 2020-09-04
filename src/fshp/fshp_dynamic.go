package fshp

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

func (fshp *FSHPImpl) Calc() {

	will := true
	iter := 0
	for will {
		will = fshp.NextIteration(iter)
		iter++

	}
}

// InitBucket set every vertex a init bucket
func (fshp *FSHPImpl) InitBucket() {
	for i := uint64(0); i < fshp.vertexSize; i++ {
		fshp.vertex2Bucket[i] = rand.Uint64() % fshp.bucketSize
		fshp.vertex2Target[i] = fshp.vertex2Bucket[i]
		fshp.buffer[i] = i
		//fshp.notChangedVis[i] = false
	}
	fshp.notChangeSize = 0
	fshp.bufferSize = int32(fshp.vertexSize)
	for i := uint64(0); i < fshp.bucketSize; i++ {
		for j := uint64(0); j < fshp.bucketSize; j++ {
			fshp.vertexTrans[i][j] = 0
		}
	}
}

// NextIteration
func (fshp *FSHPImpl) NextIteration(iter int) (isChange bool) {
	isChange = true
	if iter == 0 {
		fshp.InitBucket()
		fshp.PreComputeBucketParallel()
	} else {
		// now buffer stores all the vertex that may changed.
		// we will recalc all the
		isChange = fshp.GetMayChangeVertexToBuffer()
		if isChange == false {
			return
		}
		// now buffer save all the vertex that need changed
		// vis and vis1 is empty
		// we will recalc all the nbrBucket in [buffer vertex] nbrs
		fshp.GetNbrBucketOfMayChangeVertexAndSetNewToBuffer()

		// now buffer save all the vertex which is [need changed 's 1-hop nbr]
		// we calc 2-hop nbr of need changed 's gain
		fshp.GetNewChangeGainVertexsToBuffer()

	}
	// now buffer save all the  2-hop nbr of need changed 's
	// and set vertex

	fshp.ComputNewChangeMoveGainSubsetVertex()

	fshp.ComputMoveProb()

	fshp.MergeBufferAndNotChangedPreIteration()
	return
}
func (fshp *FSHPImpl) testf() bool {
	var tmpVertexTrans [][]int64
	tmpVertexTrans = make([][]int64, fshp.bucketSize)
	arena1 := make([]int64, fshp.bucketSize*fshp.bucketSize)

	for i := range tmpVertexTrans {
		tmpVertexTrans[i] = arena1[i*int(fshp.bucketSize) : (i+1)*int(fshp.bucketSize)]
	}
	for i := uint64(0); i < fshp.bucketSize; i++ {
		for j := uint64(0); j < fshp.bucketSize; j++ {
			tmpVertexTrans[i][j] = 0
		}
	}
	for i := uint64(0); i < fshp.vertexSize; i++ {
		tmpVertexTrans[fshp.vertex2Bucket[i]][fshp.vertex2Target[i]]++
	}
	for i := uint64(0); i < fshp.bucketSize; i++ {
		for j := uint64(0); j < fshp.bucketSize; j++ {
			if i != j && tmpVertexTrans[i][j] != fshp.vertexTrans[i][j] {
				return false
			}
		}
	}
	return true
}
func (fshp *FSHPImpl) testBucket() bool {
	var tmpVertexTrans [][]int64
	tmpVertexTrans = make([][]int64, fshp.bucketSize)
	arena1 := make([]int64, fshp.bucketSize*fshp.bucketSize)

	for i := range tmpVertexTrans {
		tmpVertexTrans[i] = arena1[i*int(fshp.bucketSize) : (i+1)*int(fshp.bucketSize)]
	}
	for i := uint64(0); i < fshp.bucketSize; i++ {
		for j := uint64(0); j < fshp.bucketSize; j++ {
			tmpVertexTrans[i][j] = 0
		}
	}
	for i := uint64(0); i < fshp.vertexSize; i++ {
		tmpVertexTrans[fshp.vertex2Bucket[i]][fshp.vertex2Target[i]]++
	}
	for i := uint64(0); i < fshp.bucketSize; i++ {
		for j := uint64(0); j < fshp.bucketSize; j++ {
			if i != j && tmpVertexTrans[i][j] != fshp.vertexTrans[i][j] {
				return false
			}
		}
	}
	return true
}

func (fshp *FSHPImpl) testIfAll() {

	viss := make([]bool, fshp.vertexSize)
	for j := 0; j < int(fshp.bufferSize); j++ {
		viss[fshp.buffer[j]] = true
	}

	fmt.Println()
	cnt := 0
	x := 0
	for i := 0; i < int(fshp.vertexSize); i++ {
		if fshp.vertex2Bucket[i] != fshp.vertex2Target[i] {
			x++
		}
		if fshp.vertex2Bucket[i] != fshp.vertex2Target[i] && viss[i] == false {
			cnt++

		}
	}
	if cnt != -1 {
		fmt.Println("cnt", fshp.bufferSize, fshp.notChangeSize, x, cnt)
	}
}

// GetMayChangeVertexToBuffer check every vertex
// and set them to buffer
func (fshp *FSHPImpl) GetMayChangeVertexToBuffer() bool {
	parallel := uint64(runtime.NumCPU())

	var ret atomic.Value
	ret.Store(false)
	tp := int32(0)
	fshp.notChangeSize = 0
	segmentVertexSize := (uint64(fshp.bufferSize) + parallel - 1) / parallel

	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < uint64(fshp.bufferSize); beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for i := begin; i != end; i++ {
				vertex := fshp.buffer[i]
				if fshp.vertex2Target[vertex] != fshp.vertex2Bucket[vertex] &&
					rand.Float64() < fshp.probability[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]] {
					stats := atomic.AddInt32(&tp, 1)
					fshp.buffer2[stats-1] = vertex
					fshp.notChangedVis[vertex] = false
					ret.Store(true)
				} else if fshp.vertex2Target[vertex] != fshp.vertex2Bucket[vertex] {
					nstats := atomic.AddInt32(&fshp.notChangeSize, 1)
					fshp.notChanged[nstats-1] = vertex
					fshp.notChangedVis[vertex] = true
				}
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, fshp.vertexSize))
	}
	wg.Wait()
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp
	return ret.Load().(bool)
}

// GetNbrBucketOfMayChangeVertexAndSetNewToBuffer defines
// buffer shows the vertex which buckets **should** set to target
func (fshp *FSHPImpl) GetNbrBucketOfMayChangeVertexAndSetNewToBuffer() {
	tp := int32(0)
	parallel := uint64(runtime.NumCPU())

	segmentVertexSize := (uint64(fshp.bufferSize) + parallel - 1) / parallel

	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < uint64(fshp.bufferSize); beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for i := begin; i != end; i++ {
				vertex := fshp.buffer[i]
				for _, nbrNode := range fshp.graph.Nodes[vertex].Nbrlist {
					if atomic.CompareAndSwapInt32(&fshp.vis[nbrNode], 0, 1) {
						stats := atomic.AddInt32(&tp, 1)
						fshp.buffer2[stats-1] = nbrNode
					}
					atomic.AddInt32(&fshp.nbrBucket[nbrNode][fshp.vertex2Bucket[vertex]], -1)
					atomic.AddInt32(&fshp.nbrBucket[nbrNode][fshp.vertex2Target[vertex]], 1)
				}
				atomic.AddInt64(&fshp.vertexTrans[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]], -1)
				fshp.vertex2Bucket[vertex] = fshp.vertex2Target[vertex]
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, fshp.vertexSize))
	}
	wg.Wait()
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp
}

//
func (fshp *FSHPImpl) GetNewChangeGainVertexsToBuffer() {
	tp := int32(0)

	parallel := uint64(runtime.NumCPU())

	segmentVertexSize := (uint64(fshp.bufferSize) + parallel - 1) / parallel

	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < uint64(fshp.bufferSize); beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for i := begin; i != end; i++ {

				vertex := fshp.buffer[i]
				for _, nbrNode := range fshp.graph.Nodes[vertex].Nbrlist {
					if atomic.CompareAndSwapInt32(&fshp.vis1[nbrNode], 0, 1) {
						stats := atomic.AddInt32(&tp, 1)
						fshp.buffer2[stats-1] = nbrNode
					}
				}
				fshp.vis[vertex] = 0
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, fshp.vertexSize))
	}
	wg.Wait()
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp

}

// ComputNewChangeMoveGainSubsetVertex c
func (fshp *FSHPImpl) ComputNewChangeMoveGainSubsetVertex() {
	tp := int32(0)
	parallel := uint64(runtime.NumCPU())

	segmentVertexSize := (uint64(fshp.bufferSize) + parallel - 1) / parallel

	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < uint64(fshp.bufferSize); beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for i := begin; i != end; i++ {

				vertex := fshp.buffer[i]
				minGain, target := fshp.calcSingleGain(fshp.graph.Nodes[vertex])
				if minGain < 0 {
					fshp.notChangedVis[vertex] = false
					stats := atomic.AddInt32(&tp, 1)
					fshp.buffer2[stats-1] = vertex
					atomic.AddInt64(&fshp.vertexTrans[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]], -1)
					atomic.AddInt64(&fshp.vertexTrans[fshp.vertex2Bucket[vertex]][target], 1)
					fshp.vertex2Target[vertex] = target
				} else {
					if fshp.notChangedVis[vertex] == true {
						fshp.notChangedVis[vertex] = false
					}
					atomic.AddInt64(&fshp.vertexTrans[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]], -1)
					fshp.vertex2Target[vertex] = fshp.vertex2Bucket[vertex]

				}
				fshp.vis1[vertex] = 0
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, fshp.vertexSize))
	}
	wg.Wait()
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp
}

// MergeBufferAndNotChangedPreIteration
func (fshp *FSHPImpl) MergeBufferAndNotChangedPreIteration() {
	for i := int32(0); i < fshp.notChangeSize; i++ {
		vertex := fshp.notChanged[i]
		if fshp.notChangedVis[vertex] {
			fshp.buffer[fshp.bufferSize] = vertex
			fshp.bufferSize++
		}
		fshp.notChangedVis[vertex] = false
	}
}
