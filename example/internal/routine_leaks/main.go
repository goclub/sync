package main

import (
	"context"
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

		// 控制超时0.5秒
		// 而 some() 函数执行需要1秒
		// 这样 ch <- true 执行时就没有 case <-ch: 去接收通道
		// 因为已经 return ctx.Err() 退出了
		// 所以导致了routine 泄露
		log.Print(sleepOneSecondReturn(ctx))
	}
}
func sleepOneSecondReturn(ctx context.Context) error {
	// 将此行修改为 ch := make(chan bool, 1) 即可避免routine泄露
	// 可以在修改后观察 http://127.0.0.1:6060/debug/pprof/ 页面的数据
	ch := make(chan bool)
	go func() {
		time.Sleep(time.Second)
		ch <- true
		log.Print("routine return")
		return
	}()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
