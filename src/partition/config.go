package partition
import (
	"gpartition/common"
)

type PartitionType int8

const (
	//shpPartitionType type
	shpPartitionType PartitionType = iota
	//bdgPartitionType simple type
	bdgPartitionType
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

