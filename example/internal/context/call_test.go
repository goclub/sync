package context_test

import (
	"context"
	xerr "github.com/goclub/error"
	"time"
)

type Option struct {
	Output string
	ReturnError bool
	Sleep time.Duration
}
// 通过 Option 控制是否返回错误和运行时间
func Call(ctx context.Context, opt Option) (string, error) {
	// 使用 1 缓冲通道防止内存泄露
	resultCh := make(chan string, 1)
	errCh := make(chan error)
	go func() {
		time.Sleep(opt.Sleep)
		if opt.ReturnError {
			errCh <- xerr.New(opt.Output + " error")
			return
		}
		resultCh <- opt.Output + " success"
	}()
	select {
	case <- ctx.Done():
		return "", ctx.Err()
	case err := <- errCh:
		return "", err
	case result := <- resultCh:
		return result, nil
	}
}