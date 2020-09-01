package partition

type partition interface{
	GetBucketFromId(uint64)uint64
	GetGraph()*Graph
}


func  calcSingleFanout(vertex uint64,graph *Graph,p partition) (fanout int) {
	fanout = 0
	ns := make (map[uint64]bool,0)
	for _, nbrNode := range graph.Nodes[vertex].Nbrlist {
		uBucket := p.GetBucketFromId(nbrNode)
		ns[uBucket] = true
	}
	fanout = len(ns)

	return
}

// CalcFanout for test
func  CalcFanout(p partition) (fanout int) {
	g := p.GetGraph()

	for vertex := uint64(0); vertex < g.GetVertexSize(); vertex++ {
		fanout += calcSingleFanout(vertex,g,p)
	}
	return
}
