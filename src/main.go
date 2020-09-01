package main

import (
	"fmt"
	pt "gpartition/pshp"

	_ "net/http/pprof" 
	"net/http"

	"runtime"
	"os"
	"log"

)

func main() {

    runtime.GOMAXPROCS(1) // 限制 CPU 使用数，避免过载
    runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
    runtime.SetBlockProfileRate(1) // 开启对阻塞操作的跟踪

    go func() {
        // 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
        if err := http.ListenAndServe(":6060", nil); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()

	config,err := pt.LoadGraph("partition/test_data/youtube.in", 5, 0.5)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	shp := pt.NewSHPImpl(config)
	shp.InitBucket()

	fmt.Println(int(shp.CalcFanout()))
	iter := 0
	for pt.NextIterationParallel(shp) || iter < 500 {
		fmt.Println("CalcFanout", int(shp.CalcFanout()))
		iter++
	}

}
