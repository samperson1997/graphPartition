package main

import (
	"fmt"

	pt "gpartition/pshp"


)

func main() {

	config,err := pt.LoadGraph("test_data/lj.in", 5, 0.5)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	fmt.Println(int(shp.CalcFanout()))
	iter := 0
	for pt.NextIterationParallel(shp) || iter < 500 {
		fmt.Println("CalcFanout", int(shp.CalcFanout()))
		iter++
	}

}
