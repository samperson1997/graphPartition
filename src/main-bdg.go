package main

import (
	"fmt"
	"gpartition/common"
	"gpartition/partition"
	"gpartition/pshp"

	_ "net/http/pprof"
)

func main() {
	fmt.Println("start")
	graph, _ := common.LoadGraphFromPath("test_data/lj.in")
	bdgConfig := partition.Config{
		PartitionType: partition.BdgPartitionType,
		Graph:         graph,
		VertexSize:    graph.GetVertexSize(),
		BucketSize:    10,
		// for shp
		Prob: 0.5,
		// for bdg
		SrcNodesNum: 4000,
		StepNum:     100000,
	}
	bdgImpl, _ := partition.NewPartition(bdgConfig)
	bdgFanout := partition.CalcFanout(bdgImpl)
	fmt.Printf("result bdg fanout: %f\n", float64(bdgFanout)/float64(graph.GetVertexSize()))

	shpConfig := bdgConfig
	shpConfig.PartitionType = partition.ShpPartitionType
	shpImpl, _ := partition.NewPartition(shpConfig)
	shpFanout := partition.CalcFanout(shpImpl)
	fmt.Printf("result shp fanout: %f\n", float64(shpFanout)/float64(graph.GetVertexSize()))

	brutefconfig := shpConfig
	bfimp, _ := partition.NewPartition(brutefconfig)
	bfimp.(*pshp.SHPImpl).InitBucket()
	bfFanout := bfimp.(*pshp.SHPImpl).CalcFanout()
	fmt.Printf("result bf fanout: %f\n", bfFanout/float64(graph.GetVertexSize()))
}
