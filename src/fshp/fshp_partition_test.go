package fshp_test

import (
	"fmt"
	"gpartition/common"
	"gpartition/fshp"
	"testing"
)

func TestFuckShp(t *testing.T) {
	graph, err := common.LoadGraphFromPath("../test_data/youtube.in")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fshpConfig := fshp.FSHPConfig{
		Graph:      graph,
		VertexSize: graph.GetVertexSize(),
		BucketSize: 5,
		// for shp
		Prob: 0.5,
		// for bdg
	}
	fshpImpl := fshp.NewFSHPImpl(fshpConfig)
	for iter := 0; iter < 50; iter++ {
		fshpImpl.NextIteration(iter)
		fmt.Println("fanout :", fshpImpl.CalcFanout())
		fshpImpl.GetEachBucketVolumn()
	}

	//	fmt.Printf("result fanout fshpfanout : %d\n", fshpFanout)

	if err != nil {
		t.Fatalf(err.Error())
	}
	//	fshpFanout = partition.CalcFanout(fshpImpl)
	//	fmt.Printf("result fanout fshpfanout : %d\n", fshpFanout)
}
