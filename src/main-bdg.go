package main

import (
	"fmt"
	"gpartition/common"
	"gpartition/partition"

	_ "net/http/pprof"
)

func main() {
	fmt.Println("start")
	graph, _ := common.LoadGraphFromPath("test_data/youtube.in")
	bdgConfig := partition.Config{
		PartitionType: partition.BdgPartitionType,
		Graph:         graph,
		VertexSize:    graph.GetVertexSize(),
		BucketSize:    5,
		// for shp
		Prob: 0.5,
		// for bdg
		SrcNodesNum: 1000,
		StepNum:     10000,
	}
	bdgImpl, _ := partition.NewPartition(bdgConfig)
	bdgFanout := partition.CalcFanout(bdgImpl)
	fmt.Printf("result bdg fanout: %d\n", bdgFanout)
}
