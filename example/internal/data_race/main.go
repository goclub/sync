package main

import (
	"log"
	"sync"
)
// 运行 go run --race main.go  可以检测代码中的并发问题 (go test 也可以加上 --race)
// 但是 --race 需要消耗一定的性能,所以不建议在正式环境中一直运行 --race
// 删除 /**/// 代码开启互斥锁,即可避免数据竞争
func main () {
	count := 0
	/**/// var lock sync.Mutex
	// WaitGroup 的作用是等待 routine 执行结束
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			/**/// lock.Lock()
			count++
			/**/// lock.Unlock()
		}()
	}
	wg.Wait()
	log.Print("count: ", count)
	// 打印结果 10000 因为routine并发操作数据出现了数据竞争
}
