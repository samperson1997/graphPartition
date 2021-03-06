package fshp_test

import (
	"fmt"
	"gpartition/common"
	"gpartition/fshp"
	"testing"
	"time"
)

func TestFShp(t *testing.T) {
	graph, err := common.LoadGraphFromPath("../test_data/youtube.in")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fshpConfig := fshp.FSHPConfig{
		Graph:      graph,
		VertexSize: graph.GetVertexSize(),
		BucketSize: 5,
		Prob:       0.5,
	}
	fshpImpl := fshp.NewFSHPImpl(fshpConfig)

	iter := 0
	will := true
	for will {
		t1 := time.Now().UnixNano()
		will = fshpImpl.NextIteration(iter)
		t2 := time.Now().UnixNano()
		fmt.Println(float64(t2-t1) / 1000000)
		iter++

	}
}
