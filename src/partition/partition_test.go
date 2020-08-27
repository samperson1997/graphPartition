package partition_test

import (
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
