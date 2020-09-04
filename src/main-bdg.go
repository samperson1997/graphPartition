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
	graph, _ := common.LoadGraphFromPath("test_data/youtube.in")
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
	bdgImpl, err := partition.NewPartition(bdgConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	bdgFanout := partition.CalcFanout(bdgImpl)
	fmt.Printf("result bdg fanout: %f\n", float64(bdgFanout)/float64(graph.GetVertexSize()))

	shpConfig := bdgConfig
	shpConfig.PartitionType = partition.ShpPartitionType
	shpImpl, err := partition.NewPartition(shpConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	shpFanout := partition.CalcFanout(shpImpl)
	fmt.Printf("result shp fanout: %f\n", float64(shpFanout)/float64(graph.GetVertexSize()))

	brutefconfig := shpConfig
	bfimp, err := partition.NewPartition(brutefconfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	bfimp.(*pshp.SHPImpl).InitBucket()
	bfFanout := bfimp.(*pshp.SHPImpl).CalcFanout()
	fmt.Printf("result bf fanout: %f\n", bfFanout/float64(graph.GetVertexSize()))
}
