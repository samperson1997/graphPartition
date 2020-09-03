package fshp

import (
	"gpartition/common"
	"math"
)

//
func (fshp *FSHPImpl) GetBucketFromId(id uint64) uint64 {
	if id > fshp.vertexSize {
		return math.MaxUint64
	}
	return fshp.vertex2Bucket[id]
}

//
func (fshp *FSHPImpl) GetGraph() *common.Graph {
	return fshp.graph
}

//
func (fshp *FSHPImpl) GetBucketSize() uint64 {
	return fshp.bucketSize
}

//
func (fshp *FSHPImpl) AfterCalc() {

}
