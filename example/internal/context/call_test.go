package context_test

import (
	"context"
	xerr "github.com/goclub/error"
	"time"
)

type Option struct {
	Name string
	ReturnError bool
	Sleep time.Duration
}
// 通过 Option 控制是否返回错误和运行时间
func Call(ctx context.Context, opt Option) (string, error) {
	resultCh := make(chan string)
	errCh := make(chan error)
	go func() {
		time.Sleep(opt.Sleep)
		if opt.ReturnError {
			errCh <- xerr.New(opt.Name + " error")
			return
		}
		resultCh <- opt.Name + " success"
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
