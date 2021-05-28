package xsync

import "sync"

type ErrorRecover struct {
	Err error
	Recover interface{}
}
func Go(errRecoverCh chan ErrorRecover, routine func() (err error)) {
	// @TODO 是否要判断 errRecoverCh == nil
	errRecover := ErrorRecover{}
	go func() {
		defer func() {
			r := recover()
			// r != nil 是可以省略的，但是在防御编程角度还是加上比较好
			if r != nil {
				errRecover.Recover = r
				errRecoverCh <- errRecover
			}
		}()
		errRecover.Err = routine()
		errRecoverCh <- errRecover
	}()
}

type WaitGroupRoutine struct {
	wg sync.WaitGroup
	err error
	recoverValue interface{}
}
func (r *WaitGroupRoutine) AddGo(routine func() (err error)) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		defer func() {
			recoverValue := recover() ; if recoverValue != nil {
				r.recoverValue = recoverValue
			}
		}()
		err := routine() ; if err != nil {
			r.err = err
		}
	}()
}
func (r *WaitGroupRoutine) Wait() (err error, recoverValue interface{}) {
	r.wg.Wait()
	return r.err, r.recoverValue
}