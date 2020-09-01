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
func TestFanoutChange(t *testing.T) {
	t.Logf("TestFanoutChange...")
	config := pt.LoadGraph("test_data/youtube.in", 5)
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	initFanout := shp.CalcFanout()
	iter := 0
	for pt.NextIteration(shp) && iter < 100 {
		iter++
	}
	resultFanout := shp.CalcFanout()
	if initFanout <= resultFanout {
		t.Fatalf("init fanout is better than resultFanout  with initFanout:%d  resultFanout:%d\n", int(initFanout), int(resultFanout))
	}
	if iter >= 100 {
		t.Fatalf("too much iteration, may be not convergence\n")
	}
}

func TestFanoutChangeParallel(t *testing.T) {
	t.Logf("TestFanoutChangeParallel...")
	config := pt.LoadGraph("test_data/youtube.in", 5)
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	initFanout := shp.CalcFanout()
	iter := 0
	for pt.NextIterationParallel(shp) && iter < 100 {
		iter++
	}
	resultFanout := shp.CalcFanout()
	if initFanout <= resultFanout {
		t.Fatalf("init fanout is better than resultFanout  with initFanout:%d  resultFanout:%d\n", int(initFanout), int(resultFanout))
	}
	if iter >= 100 {
		t.Fatalf("too much iteration, may be not convergence\n")
	}
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
