package context_test

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s, err := SleepAndReturn(ctx, time.Second*2)
	log.Print("s: ", s)
	log.Print("err: ", err)
	// 故意sleep足够的时间来观察 ctx.Done() 触发后 SleepAndReturn 中的 routine 是否还在执行
	time.Sleep(time.Second*3)
}

func SleepAndReturn(ctx context.Context, sleep time.Duration) (string, error) {
	resultCh := make(chan string)
	go func() {
		log.Print("start sleep")
		time.Sleep(sleep)
		log.Print("end sleep")
		resultCh <- "some"
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case s := <-resultCh:
		return s, nil
	}
}
