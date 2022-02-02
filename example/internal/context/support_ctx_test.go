package context_test

import (
	"context"
	xerr "github.com/goclub/error"
	xsync "github.com/goclub/sync"
	"log"
	"testing"
	"time"
)

func TestSupportCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	data ,err := supportCtx(ctx) ; if err != nil {
	    xerr.PrintStack(err)
	} else {
		log.Print(data)
	}
}
// 固定1秒后返回字符串(支持 ctx)
func supportCtx(ctx context.Context) (data string, err error) {
	dataCh := make(chan string, 1)
	errCh := xsync.Go(func() (err error) {
		time.Sleep(time.Second)
		dataCh <- "abc"
		return
	})
	select {
	case data = <- dataCh:
	return
	case err = <-errCh:
	return
	case <- ctx.Done():
		err = ctx.Err()
	return
	}
}
