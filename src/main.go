package main

import (
	pt "gpartition/partition"
)

func main() {
	config := pt.LoadGraph("data.in", 5)

	shp := pt.NewSHPImpl(config)
	shp.InitBucket()
	maxIteration := 100
	for i := 0; i < maxIteration; i++ {
		pt.NextIteration(shp)
	}
	//shp.PrintResult()
}
