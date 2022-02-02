package xsync

import (
	"fmt"
	xerr "github.com/goclub/error"
	"runtime/debug"
)

type ErrPanic struct {
	Recover interface{}
	Stack []byte
}
func (e *ErrPanic) Error() string {
	return fmt.Sprintf("%+v", e.Recover)
}
func IsErrPanic(err error) (is bool, errPanic ErrPanic) {
	var target *ErrPanic
	if xerr.As(err, &target) {
		return true, *target
	}
	return
}
func Go(routine func() (err error)) (errCh chan error) {
	// 使用 1 缓存通道防止routine 泄露
	errCh = make(chan error, 1)
	go func() {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				err = xerr.WithStack(&ErrPanic{
					Recover: r,
					Stack: debug.Stack(),
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