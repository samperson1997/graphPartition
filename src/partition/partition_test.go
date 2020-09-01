package partition_test

import (
	"fmt"
	pt "gpartition/partition"
	"runtime"
	"testing"
	"time"
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

func BenchmarkSHPEachIteration(b *testing.B) {
	config := pt.LoadGraph("test_data/youtube.in", 5)
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()
	beginTime := time.Now().UnixNano() / 1000
	for iter := 0; iter < 1000; iter++ {

		fmt.Printf("social hash Parallel iter: %d\n", iter)
		time.Now().UnixNano()
		time1 := time.Now().UnixNano() //纳秒
		shp.PreComputeBucket()
		time2 := time.Now().UnixNano() //纳秒
		shp.ComputMoveGainParallel()
		time3 := time.Now().UnixNano() //纳秒
		shp.ComputMoveProb()
		time4 := time.Now().UnixNano() //纳秒
		shp.SetNewParallel()
		time5 := time.Now().UnixNano() //纳秒
		fmt.Println(time2-time1, time3-time2, time4-time3, time5-time4, time5-time1)
		runtime.GC()
		endTime := time.Now().UnixNano() / 1000
		fmt.Println("process minisecond", endTime-beginTime)
	}
}
