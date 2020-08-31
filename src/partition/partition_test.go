package partition_test

import (
	"fmt"
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
func TestFanoutChange(t *testing.T) {
	config := pt.LoadGraph("test_data/youtube.in", 5)
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	fmt.Println(int(shp.CalcFanout()))
	for pt.NextIteration(shp) {
	}
	fmt.Println(int(shp.CalcFanout()))
}

//BenchmarkSHP a benchmark demo
func BenchmarkSHP(b *testing.B) {
	config := pt.LoadGraph("test_data/youtube.in", 5)
	b.Run(
		"social hash",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shp := pt.NewSHPImpl(config)
				b.StartTimer()
				shp.InitBucket()
				for pt.NextIteration(shp) {
				}
				b.StopTimer()
			}
		},
	)
}
func BenchmarkSHPParallel(b *testing.B) {
	config := pt.LoadGraph("test_data/youtube.in", 5)
	b.Run(
		"social hash Parallel",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shp := pt.NewSHPImpl(config)
				b.StartTimer()
				shp.InitBucket()
				for pt.NextIterationParallel(shp) {
				}
				b.StopTimer()
			}
		},
	)
}
func BenchmarkSHPWithBufferParallel(b *testing.B) {
	config := pt.LoadGraph("test_data/youtube.in", 5)
	b.Run(
		"social hash Parallel",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shp := pt.NewSHPImpl(config)
				b.StartTimer()
				shp.InitBucket()
				for pt.NextIterationWithBufferParallel(shp) {
				}
				b.StopTimer()
			}
		},
	)
}
