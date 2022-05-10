package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	// 开启 pprof 分析
	_ "net/http/pprof"
	"time"
)
func main () {
	go func() {
		// 打开 http://127.0.0.1:6060/debug/pprof/
		// 不断刷新页面观察 goroutine 左侧数字,可以发现一直在增长
		addr := ":6060"
		log.Print("http://127.0.0.1"+addr+"/debug/pprof/")
		log.Print(http.ListenAndServe(addr, nil))
	}()
	for {
		// 每 0.1 秒执行一次
		time.Sleep(time.Millisecond*100)


		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()
		func1(ctx)
	}
}

func func1(ctx context.Context) {
	resp := make(chan int)
	go func() {
		time.Sleep(time.Second * 5)     // 模拟处理逻辑
		resp <- 1
	}()
	// 超时机制
	select {
	case <-ctx.Done():
		fmt.Println("ctx timeout")
		fmt.Println(ctx.Err())
	case <-resp:
		fmt.Println("done")
	}
	return
}
