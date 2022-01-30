package context_test

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestWithCancel (t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 可以尝试修改 AOption BOption 的参数后观察代码结果
	AOption :=  Option{
		Name: "a",
		ReturnError: true,
		Sleep: time.Second * 1,
	}
	BOption := Option{
		Name: "b",
		ReturnError: false,
		Sleep: time.Second * 2,
	}
	go func() {
		s, err := Call(ctx, AOption) ; if err != nil {
			cancel() // 调用 cancel() 时  <- ctx.Done() 会接收,并且通过 ctx.Err() 可以获取到 context.Canceled
			log.Printf("call a error: %v", err)
			return
		}
		log.Print("call a result: ", s)
	}()
	go func() {
		s, err := Call(ctx, BOption) ; if err != nil {
			cancel() // 调用 cancel() 时  <- ctx.Done() 会接收,并且通过 ctx.Err() 可以获取到 context.Canceled
			log.Printf("call b error: %v", err)
			return
		}
		log.Print("call b result: ", s)
	}()
	time.Sleep(time.Second*3)
	return
}
