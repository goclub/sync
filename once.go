package xsync

import "sync"

type Once struct {
	core sync.Once
}

func (o *Once) Do(f func () (err error)) (err error) {
	o.core.Do(func() {
		err = f()
	})
	return err
}
