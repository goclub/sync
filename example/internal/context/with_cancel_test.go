package context_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

func TestWithCancel (t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 可以尝试修改 AOption BOption 的参数后观察代码结果
	AOption :=  Option{
		Output: "a",
		ReturnError: true,
		Sleep: time.Second * 1,
	}
	BOption := Option{
		Output: "b",
		ReturnError: false,
		Sleep: time.Second * 2,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s, err := Call(ctx, AOption) ; if err != nil {
			cancel() // 调用 cancel() 时  Call(ctx, BOption) 函数中的 <- ctx.Done() 会接收,并且通过 ctx.Err() 可以获取到 context.Canceled
			log.Printf("call a error: %v", err)
			return
		}
		log.Print("call a result: ", s)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s, err := Call(ctx, BOption) ; if err != nil {
			cancel() // 调用 cancel() 时  Call(ctx, AOption) 函数中的 <- ctx.Done() 会接收,并且通过 ctx.Err() 可以获取到 context.Canceled
			log.Printf("call b error: %v", err)
			return
		}
		log.Print("call b result: ", s)
	}()
	wg.Wait()
	return
}
