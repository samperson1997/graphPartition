package main

import (
	"fmt"
	pt "gpartition/partition"
)

func main() {
	config := pt.LoadGraph("partition/test_data/youtube.in", 5)
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	fmt.Println(int(shp.CalcFanout()))
	iter := 0
	for pt.NextIteration(shp) && iter < 500 {
		fmt.Println("CalcFanout", int(shp.CalcFanout()))
		iter++
	}

}
