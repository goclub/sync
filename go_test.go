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
		errRecoverCh := xsync.Go(func() (err error) {
			return nil
		})
		errRecover := <-errRecoverCh
		assert.Equal(t,errRecover.Err, nil)
		assert.Equal(t,errRecover.Recover, nil)
	}
	// has err no panic
	{
		errRecoverCh := xsync.Go(func() (err error) {
			return xerr.New("abc")
		})
		errRecover := <-errRecoverCh
		assert.Error(t,errRecover.Err, "abc")
		assert.Equal(t,errRecover.Recover, nil)
	}
	// no err has panic
	{
		errRecoverCh := xsync.Go(func() (err error) {
			panic(1)
			return nil
		})
		errRecover := <-errRecoverCh
		assert.Equal(t,errRecover.Err, nil)
		assert.Equal(t,errRecover.Recover, 1)
	}
}