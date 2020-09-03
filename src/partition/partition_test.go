package partition_test

import (
	"fmt"
	"gpartition/common"
	"gpartition/partition"
	"testing"
)

func LoadConfigFromGraphPath() {

}

func TestFanout(t *testing.T) {
	graph, err := common.LoadGraphFromPath("../test_data/youtube.in")
	if err != nil {
		t.Fatalf(err.Error())
	}
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
	shpConfig := bdgConfig
	shpConfig.PartitionType = partition.ShpPartitionType
	bdgImpl, err := partition.NewPartition(bdgConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	bdgFanout := partition.CalcFanout(bdgImpl)
	shpImpl, err := partition.NewPartition(shpConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	shpFanout := partition.CalcFanout(shpImpl)
	fmt.Println(graph.GetVertexSize())
	fmt.Printf("result fanout shpfanout : %f\n, bdg fanout: %f\n",
		float64(shpFanout)/float64(graph.GetVertexSize()),
		float64(bdgFanout)/float64(graph.GetVertexSize()))
}

func TestBucketBalance(t *testing.T) {
	fmt.Println("step1")
	graph, err := common.LoadGraphFromPath("../test_data/youtube.in")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("step2")
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

	fmt.Println("step3")
	shpConfig := bdgConfig
	shpConfig.PartitionType = partition.ShpPartitionType

	fmt.Println("step4")
	// bdgImpl, err := partition.NewPartition(bdgConfig)

	fmt.Println("step5")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("step6")
	// partition.GetEachBucketVolumn(bdgImpl)

	fmt.Println("step7")
	shpImpl, err := partition.NewPartition(shpConfig)

	fmt.Println("step8")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("step9")
	partition.GetEachBucketVolumn(shpImpl)

}
