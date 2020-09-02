package partition
import (
	"gpartition/common"
)

// PartitionType defines different partition
type PartitionType int8

const (
	//BdgPartitionType type
	BdgPartitionType PartitionType = iota
	//ShpPartitionType simple type
	ShpPartitionType
)


// Config all type of config
type Config struct {
	PartitionType
	BlockSize uint64
	VertexSize uint64
	BucketSize uint64
	Prob       float64
	Graph      *common.Graph
}

