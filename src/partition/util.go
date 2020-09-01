package partition

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)
// LoadGraphFromPath load a graph with path
func LoadGraphFromPath(path string) (g *Graph,err error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(err.Error())
		return nil,err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	str, err := reader.ReadString('\n')
	var vertexSize uint64
	fmt.Fscanf(strings.NewReader(str), "%d", &vertexSize)
	
	g = NewGraph(vertexSize)
	var src, dst uint64
	edgeSize := 0
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		edgeSize++
		fmt.Fscanf(strings.NewReader(str), "%d %d", &src, &dst)
		if src < 0 || src >= vertexSize || dst < 0 || dst >= vertexSize {
			return nil,fmt.Errorf("err in edge with src: %d dst: %d vertexSize: %d",src,dst,vertexSize)
		}
		g.AddEdge(src, dst)
		g.AddEdge(dst, src)
	}
	return g,nil
}
