package partition_test

import (
	"gpartition/common"
	"testing"
	"gpartition/partition"
)

func LoadConfigFromGraphPath(){
}

func TestFanout(t *testing.T)  {
	graph,err := common.LoadGraphFromPath("../test_data/youtube.in")
	if err != nil{
		t.Fatalf(err.Error())
	}
	bdgConfig := partition.Config{
		PartitionType: partition.BdgPartitionType,
		Graph:graph,
		VertexSize:graph.GetVertexSize(),
		BucketSize:5,
		// for shp
		Prob:0.5,
		// for bdg
		BlockSize:10,
	}
	shpConfig := bdgConfig
	shpConfig.PartitionType = partition.ShpPartitionType
	bdgImpl,err:= partition.NewPartition(bdgConfig)
	if err != nil{
		t.Fatalf(err.Error())
	}
	bdgFanout := partition.CalcFanout(bdgImpl)
	shpImpl,err:= partition.NewPartition(shpConfig)
	if err != nil{
		t.Fatalf(err.Error())
	}
	shpFanout :=  partition.CalcFanout(shpImpl)

	t.Logf("result fanout shpfanout : %d bdg fanout: %d",bdgFanout,shpFanout )
}