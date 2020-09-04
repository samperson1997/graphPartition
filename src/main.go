package main

import (
	"fmt"
	pt "gpartition/pshp"
	"time"
)

func main() {

	config, err := pt.LoadGraph("test_data/youtube.in", 5, 0.5)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	fmt.Println(int(shp.CalcFanout()))
	iter := 0
	will := true
	for will {
		t1 := time.Now().UnixNano()
		will = pt.NextIterationParallel(shp)
		t2 := time.Now().UnixNano()
		fmt.Println(float64(t2-t1) / 1000000)
		iter++
	}

}
