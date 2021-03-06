package partition_test

import (
	"fmt"
	"gpartition/common"
	"gpartition/partition"
	"gpartition/pshp"
	"testing"
)

func TestFanout(t *testing.T) {
	fmt.Println("TestFanout.....")
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
	tshpConfig := shpConfig
	tshpConfig.PartitionType = partition.TShpPartitionType
	fshpImpl, err := partition.NewPartition(tshpConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	tshpFanout := partition.CalcFanout(fshpImpl)

	burteConfig := shpConfig
	tshpConfig.PartitionType = partition.ShpPartitionType
	bruteforceImpl, err := partition.NewPartition(burteConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	bruteforceImpl.(*pshp.SHPImpl).InitBucket()
	bruteFanout := bruteforceImpl.(*pshp.SHPImpl).CalcFanout()

	fmt.Println(graph.GetVertexSize())
	fmt.Printf("result fanout shpfanout : %f\n, fshp fanout: %f, bdg fanout: %f brute force : %f\n",
		float64(shpFanout)/float64(graph.GetVertexSize()),
		float64(bdgFanout)/float64(graph.GetVertexSize()),
		float64(tshpFanout)/float64(graph.GetVertexSize()),
		float64(bruteFanout)/float64(graph.GetVertexSize()),
	)
}

func TestBucketBalance(t *testing.T) {
	fmt.Println("TestBucketBalance.....")

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

	tshpConfig := shpConfig
	tshpConfig.PartitionType = partition.TShpPartitionType
	fshpImpl, err := partition.NewPartition(tshpConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	partition.GetEachBucketVolumn(fshpImpl)

}
