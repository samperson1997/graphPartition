package fshp

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

func (fshp *FSHPImpl) Calc() {

}

// InitBucket set every vertex a init bucket
func (fshp *FSHPImpl) InitBucket() {
	for i := uint64(0); i < fshp.vertexSize; i++ {
		fshp.vertex2Bucket[i] = rand.Uint64() % fshp.bucketSize
		fshp.vertex2Target[i] = fshp.vertex2Bucket[i]
		fshp.buffer[i] = i
	}
	fshp.bufferSize = int32(fshp.vertexSize)
	for i := uint64(0); i < fshp.bucketSize; i++ {
		for j := uint64(0); j < fshp.bucketSize; j++ {
			fshp.vertexTrans[i][j] = 0
		}
	}
}
func (fshp *FSHPImpl) NextIteration(iter int) {
	if iter == 0 {
		fshp.InitBucket()
		fshp.PreComputeBucketParallel()

	} else {
		// now buffer stores all the vertex that may changed.
		// we will recalc all the
		fshp.GetMayChangeVertexToBuffer()
		fmt.Println("GetMayChangeVertexToBuffer", fshp.bufferSize)

		// now buffer save all the vertex that need changed
		// vis and vis1 is empty
		// we will recalc all the nbrBucket in [buffer vertex] nbrs
		fshp.GetNbrBucketOfMayChangeVertexAndSetNewToBuffer()
		fmt.Println("GetNbrBucketOfMayChangeVertexAndSetNewToBuffer", fshp.bufferSize)

		// now buffer save all the vertex which is [need changed 's 1-hop nbr]
		// we calc 2-hop nbr of need changed 's gain
		fshp.GetNewChangeGainVertexsToBuffer()
		fmt.Println("GetNewChangeGainVertexsToBuffer", fshp.bufferSize)

	}
	// now buffer save all the  2-hop nbr of need changed 's
	// and set vertex

	fshp.ComputNewChangeMoveGainSubsetVertex()
	fmt.Println("ComputNewChangeMoveGainSubsetVertex", fshp.bufferSize)

	fshp.ComputMoveProb()
	fmt.Println("ComputMoveProb", fshp.bufferSize)

	//
	fshp.MergeBufferAndNotChangedPreIteration()
	fmt.Println("MergeBufferAndNotChangedPreIteration", fshp.bufferSize)
}
func (fshp *FSHPImpl) testf() bool {
	var tmpVertexTrans [][]uint64
	tmpVertexTrans = make([][]uint64, fshp.bucketSize)
	arena1 := make([]uint64, fshp.bucketSize*fshp.bucketSize)

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
				fmt.Println("not match", i, j, tmpVertexTrans[i][j], fshp.vertexTrans[i][j])
				return false
			}
		}
	}
	return true
}

// ComputMoveGainWithBufferParallelAll parallel compute maxgain of each vertex
// then buffer store the all target change vertex
func (fshp *FSHPImpl) ComputMoveGainWithBufferParallelAll() {
	for bucketI := uint64(0); bucketI < fshp.bucketSize; bucketI++ {
		for bucketJ := uint64(0); bucketJ < fshp.bucketSize; bucketJ++ {
			fshp.vertexTrans[bucketI][bucketJ] = 0
		}
	}
	parallel := uint64(runtime.NumCPU())
	atomic.StoreInt32(&fshp.bufferSize, 0)
	// TODO
	// make parallel seperate stable
	segmentVertexSize := (fshp.vertexSize + parallel - 1) / parallel
	var wg sync.WaitGroup
	for beginvertex := uint64(0); beginvertex < fshp.vertexSize; beginvertex += segmentVertexSize {
		wg.Add(1)
		go func(begin, end uint64) {
			defer wg.Done()
			for vertex := begin; vertex != end; vertex++ {
				minGain, target := fshp.calcSingleGain(fshp.graph.Nodes[vertex])
				if minGain < 0 {
					fshp.vertex2Target[vertex] = target
					atomic.AddUint64(&fshp.vertexTrans[fshp.vertex2Bucket[vertex]][target], 1)
					tp := atomic.AddInt32(&fshp.bufferSize, 1)
					fshp.buffer[tp-1] = vertex
				}
			}
		}(beginvertex, min(beginvertex+segmentVertexSize, fshp.vertexSize))
	}
	wg.Wait()
	fmt.Println(fshp.bufferSize)
}

// GetMayChangeVertexToBuffer check every vertex
// and set them to buffer
func (fshp *FSHPImpl) GetMayChangeVertexToBuffer() (ret bool) {
	ret = false
	tp := int32(0)
	fshp.notChangeSize = 0
	for i := 0; i < int(fshp.bufferSize); i++ {
		vertex := fshp.buffer[i]
		if fshp.vertex2Target[vertex] != fshp.vertex2Bucket[vertex] &&
			rand.Float64() < fshp.probability[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]] {
			stats := atomic.AddInt32(&tp, 1)
			fshp.buffer2[stats-1] = vertex
			fshp.notChangedVis[vertex] = false
			ret = true
		} else if fshp.vertex2Target[vertex] != fshp.vertex2Bucket[vertex] {
			stats := atomic.AddInt32(&fshp.notChangeSize, 1)
			fshp.notChanged[stats-1] = vertex
			fshp.notChangedVis[vertex] = true
		}
	}
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp
	return
}

//buffer shows the vertex which buckets **should** set to target
func (fshp *FSHPImpl) GetNbrBucketOfMayChangeVertexAndSetNewToBuffer() {
	tp := int32(0)
	for i := int32(0); i < fshp.bufferSize; i++ {
		vertex := fshp.buffer[i]
		for _, nbrNode := range fshp.graph.Nodes[vertex].Nbrlist {
			if atomic.CompareAndSwapInt32(&fshp.vis[nbrNode], 0, 1) {
				stats := atomic.AddInt32(&tp, 1)
				fshp.buffer2[stats-1] = nbrNode
			}
			ns := make([]int32, 5)
			for _, nbn := range fshp.graph.Nodes[nbrNode].Nbrlist {
				ns[fshp.vertex2Bucket[nbn]]++
			}
			for k, v := range ns {
				if v != fshp.nbrBucket[nbrNode][k] {
					fmt.Println("qwq", nbrNode, v, fshp.nbrBucket[nbrNode][k], fshp.vertex2Bucket[vertex], fshp.vertex2Target[vertex])
				}
			}
			atomic.AddInt32(&fshp.nbrBucket[nbrNode][fshp.vertex2Bucket[vertex]], -1)
			atomic.AddInt32(&fshp.nbrBucket[nbrNode][fshp.vertex2Target[vertex]], 1)
		}
		fshp.vertexTrans[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]]--
		fshp.vertex2Bucket[vertex] = fshp.vertex2Target[vertex]
	}
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp
}

func (fshp *FSHPImpl) GetNewChangeGainVertexsToBuffer() {
	tp := int32(0)
	for i := int32(0); i < fshp.bufferSize; i++ {
		vertex := fshp.buffer[i]
		for _, nbrNode := range fshp.graph.Nodes[vertex].Nbrlist {
			if atomic.CompareAndSwapInt32(&fshp.vis1[nbrNode], 0, 1) {
				stats := atomic.AddInt32(&tp, 1)
				fshp.buffer2[stats-1] = nbrNode
			}
		}
		fshp.vis[vertex] = 0
	}
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp

}

func (fshp *FSHPImpl) ComputNewChangeMoveGainSubsetVertex() {
	tp := int32(0)
	for i := int32(0); i < fshp.bufferSize; i++ {
		vertex := fshp.buffer[i]
		minGain, target := fshp.calcSingleGain(fshp.graph.Nodes[vertex])
		if minGain < 0 {
			if fshp.notChangedVis[vertex] == false {
				fshp.notChangedVis[vertex] = true
				stats := atomic.AddInt32(&tp, 1)
				fshp.buffer2[stats-1] = vertex
			}
		} else {
			if fshp.notChangedVis[vertex] == true {
				fshp.notChangedVis[vertex] = false
			}
		}
		fshp.vertexTrans[fshp.vertex2Bucket[vertex]][fshp.vertex2Target[vertex]]--
		fshp.vertexTrans[fshp.vertex2Bucket[vertex]][target]++
		fshp.vertex2Target[vertex] = target
		fshp.vis1[vertex] = 0
	}
	fshp.buffer, fshp.buffer2 = fshp.buffer2, fshp.buffer
	fshp.bufferSize = tp
}

//
func (fshp *FSHPImpl) MergeBufferAndNotChangedPreIteration() {
	for i := int32(0); i < fshp.notChangeSize; i++ {
		vertex := fshp.notChanged[i]
		if fshp.notChangedVis[vertex] {
			fshp.buffer[fshp.bufferSize] = vertex
			fshp.bufferSize++
		}
		fshp.notChangedVis[vertex] = false
	}
	fshp.notChangeSize = 0
}
