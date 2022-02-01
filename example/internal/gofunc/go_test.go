package go_test

import (
	xerr "github.com/goclub/error"
	xsync "github.com/goclub/sync"
	"log"
	"testing"
)

// 等待 routine 完成
func TestWaitRoutineDone(t *testing.T) {
	errRecoverCh := xsync.Go(func() (err error) {
		return xerr.New("some error")
	})
	errRecover := <-errRecoverCh
	if errRecover.Err != nil {
		log.Print("err: ", errRecover.Err)
	}
	if errRecover.Recover != nil {
		log.Print("Recover: ", errRecover.Recover)
	}
}
// routine 通过channel 返回字符串或者 ErrorRecover
func TestGetStringOrError(t *testing.T) {
	nameCh := make(chan string)
	// 修改 sendError 或 sendPanic 为 true 来观察运行结果
	sendPanic := false
	sendError := false
	errRecoverCh := xsync.Go(func() (err error) {
		if sendPanic {
			panic("some panic")
		}
		if sendError {
			return xerr.New("some error")
		}
		nameCh <- "goclub"
		return
	})
	// 不使用 select 会导致死锁
	select {
	case errRecover := <-errRecoverCh:
		if errRecover.Err != nil {
			log.Print("err: ", errRecover.Err)
		}
		if errRecover.Recover != nil {
			log.Print("Recover: ", errRecover.Recover)
		}
	case name := <-nameCh:
		log.Print("name: ", name)
	}
}
