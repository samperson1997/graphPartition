package partition
import (
	"gpartition/common"
)

type PartitionType int8

const (
	//BdgPartitionType type
	BdgPartitionType PartitionType = iota
	//ShpPartitionType simple type
	ShpPartitionType
)

// Config 
type Config struct {
	PartitionType
	BlockSize uint64
	VertexSize uint64
	BucketSize uint64
	Prob       float64
	Graph      *common.Graph
}

