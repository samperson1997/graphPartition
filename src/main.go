package main

import (
	"bufio"
	"fmt"
	pt "gpartition/partition"
	"io"
	"os"
	"strings"
)

func loadGraph(path string) (c pt.Config) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("err=%v", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	str, err := reader.ReadString('\n')
	fmt.Fscanf(strings.NewReader(str), "%d", &c.VertexSize)
	c.Graph = pt.NewGraph(int(c.VertexSize))
	var src, dst uint64
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Fscanf(strings.NewReader(str), "%d %d", &src, &dst)
		c.Graph.AddEdge(src, dst)
		c.Graph.AddEdge(dst, src)
	}
	c.BucketSize = 3
	c.Prob = 0.5
	return
}

func main() {
	config := loadGraph("data.in")

	shp := pt.NewSHPImpl(config)
	shp.InitBucket()
	maxIteration := 100
	for i := 0; i < maxIteration; i++ {
		pt.NextIterationParallel(shp)
	}
	shp.PrintResult()
}
