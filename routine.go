package xsync

import (
	"sync"
)

type Routine struct {
	wg sync.WaitGroup
	err error
	recoverValue interface{}
}
func NewRoutine() *Routine {
	return new(Routine)
}
func (r *Routine) Go(routine func() error) {
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
func (r *Routine) Wait() (err error, recoverValue interface{}) {
	r.wg.Wait()
	return r.err, r.recoverValue
}