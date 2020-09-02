package bdg

import (
	"bufio"
	"fmt"
	. "gpartition/common"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFanout(t *testing.T) {
	config := LoadGraphForBDG("../data.in", 3, 6)

	bdg := NewBDGImpl(config)
	bdg.bfs()
	bdg.deterministicGreedy()

	for i := range bdg.buckets {
		fmt.Print(i, ": ")
		for block := bdg.buckets[i].Front(); block != nil; block = block.Next() {
			fmt.Print(block.Value.(uint64), ",")
		}
		fmt.Println()
	}
}

// LoadGraph load a graph with path
func LoadGraphForBDG(path string, bucketSize int, blockSize int) (c BDGConfig) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	str, err := reader.ReadString('\n')
	fmt.Fscanf(strings.NewReader(str), "%d", &c.VertexSize)
	c.Graph = NewGraph(c.VertexSize)
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
	c.BucketSize = uint64(bucketSize)
	c.BlockSize = uint64(blockSize)
	return
}
