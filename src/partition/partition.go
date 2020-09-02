package partition

import (
	"gpartition/pshp"
	"gpartition/bdg"
	"gpartition/common"
)

type partition interface{
	GetBucketFromId(uint64)uint64
	GetGraph()*common.Graph
	Calc()
}

func NewPartition(c Config)partition{
	switch c.PartitionType{
	case shpPartitionType:
		{
			return bdg.NewBDGImpl(bdg.BDGConfig{
				VertexSize: c.VertexSize,
				BlockSize :c.BlockSize,
				BucketSize:c.BucketSize ,
				Graph : c.Graph,
			})
		}
		case bdgPartitionType:{
			return pshp.NewSHPImpl(pshp.SHPConfig{
				VertexSize: c.VertexSize,
				Prob :c.Prob,
				BucketSize:c.BucketSize,
				Graph : c.Graph,
			})
		}
	default:

	}
	
}


func  calcSingleFanout(vertex uint64,graph *common.Graph,p partition) (fanout int) {
	fanout = 0
	ns := make (map[uint64]bool,0)
	for _, nbrNode := range graph.Nodes[vertex].Nbrlist {
		uBucket := p.GetBucketFromId(nbrNode)
		ns[uBucket] = true
	}
	fanout = len(ns)

	return
}

// CalcFanout for test
func  CalcFanout(p partition) (fanout int) {
	g := p.GetGraph()
	for vertex := uint64(0); vertex < g.GetVertexSize(); vertex++ {
		fanout += calcSingleFanout(vertex,g,p)
	}
	return
}
