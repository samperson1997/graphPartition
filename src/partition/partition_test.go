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
	fmt.Printf("result fanout shpfanout : %d bdg fanout: %d\n", shpFanout, bdgFanout)
}

func TestBucketBalance(t *testing.T) {
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
	partition.GetEachBucketVolumn(bdgImpl)
	shpImpl, err := partition.NewPartition(shpConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	partition.GetEachBucketVolumn(shpImpl)

}
