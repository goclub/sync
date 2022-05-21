package xsync

import (
	"fmt"
	xerr "github.com/goclub/error"
	"runtime/debug"
)

type ErrPanic struct {
	Recover interface{}
	Stack   []byte
}

func (e *ErrPanic) Error() string {
	return fmt.Sprintf("%+v", e.Recover)
}
func AsErrPanic(err error) (as bool, errPanic ErrPanic) {
	var target *ErrPanic
	if xerr.As(err, &target) {
		return true, *target
	}
	return
}

// ChanError eg: err = <-errCh
type ChanError chan error

// 特意实现一个 Error() string 方法,作用是在 go vet 或者编辑器中出现 Unhandled error 的提示避免使用者忘记处理 ErrorChan
func (ChanError) Error() string { return "goclub/rabbitmq: Implemented to prompt Unhandled Error" }

func Go(routine func() (err error)) (errCh ChanError) {
	// 使用 1 缓存通道防止routine 泄露
	errCh = make(chan error, 1)
	go func() {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				err = xerr.WithStack(&ErrPanic{
					Recover: r,
					Stack:   debug.Stack(),
				})
			}
			errCh <- err
		}()
		routineErr := routine()
		if routineErr != nil {
			err = routineErr
		}
	}()
	return
}
