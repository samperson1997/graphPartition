package partition

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadGraph load a graph with path
func LoadGraph(path string, BucketSize int) (c Config) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	str, err := reader.ReadString('\n')
	fmt.Fscanf(strings.NewReader(str), "%d", &c.VertexSize)
	c.Graph = NewGraph(int(c.VertexSize))
	var src, dst uint64
	edgeSize := 0
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		edgeSize++
		fmt.Fscanf(strings.NewReader(str), "%d %d", &src, &dst)
		if src < 0 || src >= c.VertexSize || dst < 0 || dst >= c.VertexSize {
			fmt.Println("err in edge with src: ", src, " dst: ", dst, "vertexSize ", c.VertexSize)
			return
		}
		c.Graph.AddEdge(src, dst)
		c.Graph.AddEdge(dst, src)
	}
	fmt.Println("load data from ", path, "vertex:", c.VertexSize, "edge:", edgeSize)
	c.BucketSize = uint64(BucketSize)
	c.Prob = 0.5
	return
}
func min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
