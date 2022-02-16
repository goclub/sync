package xsync_test

import (
	xerr "github.com/goclub/error"
	xsync "github.com/goclub/sync"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoutine_Go(t *testing.T) {
	// no err no panic
	{
		errCh, err := xsync.Go(func() (err error) {
			return nil
		}) ; assert.NoError(t, err)
		err = <-errCh
		assert.Equal(t,err, nil)
	}
	// has err no panic
	{
		errCh, err := xsync.Go(func() (err error) {
			return xerr.New("abc")
		}) ; assert.NoError(t, err)
		err = <-errCh
		assert.Error(t,err, "abc")
	}
	// no err has panic
	{
		errCh, err := xsync.Go(func() (err error) {
			panic(1)
			return nil
		}) ; assert.NoError(t, err)
		err = <-errCh
		assert.Error(t,err, "1")
		as, errPanic := xsync.AsErrPanic(err)
		assert.Equal(t, as, true)
		assert.Equal(t,errPanic.Recover, 1)
		// xerr.PrintStack(err)
		// log.Print(string(errPanic.Stack))
	}
}