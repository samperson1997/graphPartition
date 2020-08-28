package partition_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	pt "gpartition/partition"
)

// TODO
func TestFanout(t *testing.T) {
	config := pt.Config{}
	shp := pt.NewSHPImpl(config)
	pt.NextIteration(shp)
	shp.CalcFanout()
}

func loadGraph(path string) (c pt.Config) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(err)
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
		if src < 0 || src >= c.VertexSize || dst < 0 || dst >= c.VertexSize {
			fmt.Println("err in edge with src: ", src, " dst: ", dst, "vertexSize ", c.VertexSize)
			return
		}
		c.Graph.AddEdge(src, dst)
		c.Graph.AddEdge(dst, src)
	}
	c.BucketSize = 3
	c.Prob = 0.5
	return
}

//BenchmarkSHP a benchmark demo
func BenchmarkSHP(b *testing.B) {
	config := loadGraph("data.in")
	shp := pt.NewSHPImpl(config)
	b.Run(
		"social hash",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for pt.NextIteration(shp) {

				}
			}
		},
	)
}
