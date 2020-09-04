package pshp

import (
	"runtime"
	"sync"
	"gpartition/common"
	"math"
	"sync/atomic"
	"fmt"
)
// GatherNeedTransVertices 
func (shp *SHPImpl) GatherNeedTransVertices() {
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			//shp.vertexTrans[bucketI][bucketJ] = 0
			shp.needTrans[bucketI][bucketJ] = make([]uint64, 0)
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
				if shp.vertex2Bucket[vertex] == shp.vertex2Target[vertex] {
					continue
				}
				//minGain, target := shp.calcSingleGain(shp.graph.Nodes[vertex])
				//if minGain < 0 {
				//	shp.vertex2Target[vertex] = target
				//	atomic.AddUint64(&shp.vertexTrans[shp.vertex2Bucket[vertex]][target], 1)
				//}
				shp.mutexBucket2Bucket[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]].Lock()
				shp.needTrans[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]] = append(shp.needTrans[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]], vertex)
				shp.mutexBucket2Bucket[shp.vertex2Bucket[vertex]][shp.vertex2Target[vertex]].Unlock()
			}
		} (beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
}

func (shp *SHPImpl) calcSingleGainSort(node *common.Node) (minGain float64, target uint64) {
	minGain = 10
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
	shp.gains[node.ID] = gains //is this okay???
	return
}

func (shp *SHPImpl) ComputMoveGainSortParallel() {
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
				minGain, target := shp.calcSingleGainSort(shp.graph.Nodes[vertex])
				if minGain < 0 {
					shp.vertex2Target[vertex] = target
					atomic.AddUint64(&shp.vertexTrans[shp.vertex2Bucket[vertex]][target], 1)
				}
			}
		} (beginvertex, min(beginvertex+segmentVertexSize, shp.vertexSize))
	}
	wg.Wait()
}

func (shp *SHPImpl) solve(arr []uint64, L int64, R int64, k int64, target int64) {
	//log.Println("solve", L, R, k)
    if L >= R {
        return
    }
    base := shp.gains[arr[L]][target]
    l := L - 1
    r := R + 1
    for {
    	//log.Println("for for for ", l, r)
        if l >= r {
            break
        }
        for {
            l += 1
            if l > R || shp.gains[arr[l]][target] >= base {
                break
            }
        }
        for {
            r -= 1
            if r < L || shp.gains[arr[r]][target] <= base {
                break
            }
        }
        if l < r {
            t := arr[l]
            arr[l] = arr[r]
            arr[r] = t
        }
    }
    //log.Println("enter next solve", L, r, k)
    //log.Println("enter next solve", r + 1, R, k - (r - L + 1))
    if (r - L + 1 >= k) {
        shp.solve(arr, L, r, k, target)
    } else {
        shp.solve(arr, r + 1, R, k - (r - L + 1), target)
    }
}

func (shp *SHPImpl) smallestK(arr []uint64, k uint64, target uint64) {
	//log.Println("smallestK", k, target)
	if k == uint64(0) {
		return
	}
    shp.solve(arr, int64(0), int64(len(arr) - 1), int64(k), int64(target))
}

func (shp *SHPImpl) SortNeedTransVertices() {

	var wg sync.WaitGroup
	//shp.bucketSize = 5
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {//could optimize original implementation!!!
		//log.Println("bucketI enter ", bucketI)
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			//log.Println(bucketI, bucketJ)
			//log.Println("why enter 5")
			if bucketI == bucketJ {
				continue
			}
			wg.Add(1)
			go func(bucketI, bucketJ uint64) {
				defer wg.Done()
				//log.Println("????")
				//log.Println(len(shp.vertexTrans), shp.bucketSize)
				//log.Println(len(shp.vertexTrans[bucketI]))
				//log.Println(bucketI, bucketJ)
				transNum := min(shp.vertexTrans[bucketI][bucketJ], shp.vertexTrans[bucketJ][bucketI])
				//fmt.Println(bucketI, bucketJ, transNum)
				if transNum != shp.vertexTrans[bucketI][bucketJ] {
					shp.smallestK(shp.needTrans[bucketI][bucketJ], transNum, bucketJ)//not checked yet??
				}
				//log.Println("before smallestK", len(shp.needTrans[bucketI][bucketJ]), transNum)
				for index := uint64(0); index < transNum; index++ {
					shp.toBeTransed[shp.needTrans[bucketI][bucketJ][index]] = true
				}
			} (bucketI, bucketJ)
		}
	}
	wg.Wait()
}

func (shp *SHPImpl) SetNewSortParallel() (ret bool) {
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
				if shp.vertex2Target[vertex] != shp.vertex2Bucket[vertex] && shp.toBeTransed[vertex] == true {
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

func (shp *SHPImpl) ResetNeedTransVertices() {

	var wg sync.WaitGroup
	for bucketI := uint64(0); bucketI < shp.bucketSize; bucketI++ {//could optimize original implementation!!!
		for bucketJ := uint64(0); bucketJ < shp.bucketSize; bucketJ++ {
			if bucketI == bucketJ {
				continue
			}
			wg.Add(1)
			go func(bucketI, bucketJ uint64) {
				defer wg.Done()
				transNum := min(shp.vertexTrans[bucketI][bucketJ], shp.vertexTrans[bucketJ][bucketI])
				// if transNum == shp.vertexTrans[bucketI][bucketJ] {
				// 	return
				// }
				//smallestK(shp.needTrans[bucketI][bucketJ], transNum, source, bucketJ)
				for index := uint64(0); index < transNum; index++ {
					shp.toBeTransed[shp.needTrans[bucketI][bucketJ][index]] = false
				}
			} (bucketI, bucketJ)
		}
	}
	wg.Wait()
}

func (shp *SHPImpl) PrintBucketNum() {
	nums := make([]int, shp.bucketSize)
	for beginvertex := uint64(0); beginvertex < shp.vertexSize; beginvertex += 1 {
		nums[shp.vertex2Bucket[beginvertex]] += 1
	}
	tot := 0
	for bucketI, num := range nums {
		fmt.Println(bucketI, num)
		tot += num
	}
	fmt.Println(tot)
}


func NextIterationWithSortParallel(shp *SHPImpl) bool {
	shp.PrintBucketNum()
	shp.PreComputeBucketParallel()
	//log.Println("PreCompute succeed")
	shp.ComputMoveGainSortParallel()
	//log.Println("ComputeMoveGainSort succeed")
	shp.GatherNeedTransVertices()
	//log.Println("GatherNeedTransVertices succeed")
	shp.SortNeedTransVertices()
	//log.Println("SortNeedTransVertices succeed")
	changed := shp.SetNewSortParallel()
	//log.Println("SetNewSortParallel succeed")
	shp.ResetNeedTransVertices()
	//log.Println("ResetNeedTransVertices succeed")
	return changed
}