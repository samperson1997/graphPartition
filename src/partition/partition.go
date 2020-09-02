package partition

import (
	"gpartition/bdg"
	"gpartition/common"
	"gpartition/pshp"
	"fmt"
)

// Partition can partion graph
type Partition interface {
	GetBucketFromId(uint64) uint64
	GetGraph() *common.Graph
	GetBucketSize() uint64
	AfterCalc()
	Calc()
}

// NewPartition returns a partition with type
func NewPartition(c Config)(Partition,error){
	switch c.PartitionType{
		case BdgPartitionType:
		{
			return bdg.NewBDGImpl(bdg.BDGConfig{
				VertexSize: c.VertexSize,
				BlockSize :c.BlockSize,
				BucketSize:c.BucketSize ,
				Graph : c.Graph,
			}),nil
		}
		case ShpPartitionType:{
			return pshp.NewSHPImpl(pshp.SHPConfig{
				VertexSize: c.VertexSize,
				Prob :c.Prob,
				BucketSize:c.BucketSize,
				Graph : c.Graph,
			}),nil
		}	
	}
	return nil,fmt.Errorf("no such type")
}


func calcSingleFanout(vertex uint64, graph *common.Graph, p Partition) (fanout int) {
	fanout = 0
	ns := make(map[uint64]bool, 0)
	for _, nbrNode := range graph.Nodes[vertex].Nbrlist {
		uBucket := p.GetBucketFromId(nbrNode)
		ns[uBucket] = true
	}
	fanout = len(ns)

	return
}

// CalcFanout for test
// partition is not calced
func CalcFanout(p Partition) (fanout int) {
	p.Calc()
	p.AfterCalc()
	g := p.GetGraph()
	for vertex := uint64(0); vertex < g.GetVertexSize(); vertex++ {
		fanout += calcSingleFanout(vertex, g, p)
	}
	return
}

func GetEachBucketVolumn(p Partition){
	p.Calc()
	p.AfterCalc()
	g := p.GetGraph()
	ns := make([]int, p.GetBucketSize())

	for vertex := uint64(0); vertex < g.GetVertexSize(); vertex++ {
		uBucket := p.GetBucketFromId(vertex)
		ns[uBucket]++
	}

	for bucketI,size := range ns{
		fmt.Println(bucketI,size)
	}
}
