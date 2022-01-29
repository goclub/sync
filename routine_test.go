package xsync_test

import (
	"errors"
	xsync "github.com/goclub/sync"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoutine_Go(t *testing.T) {
	// no err no panic
	{
		exec := false
		routine := xsync.NewRoutine()
		routine.Go(func() error {
			exec = true
			return nil
		})
		err, recoverValue := routine.Wait()
		assert.Equal(t, exec, true)
		assert.NoError(t, err)
		assert.Nil(t, recoverValue)
	}
	// has err no panic
	{
		exec := false
		routine := xsync.NewRoutine()
		routine.Go(func() error {
			exec = true
			return errors.New("abc")
		})
		err, recoverValue := routine.Wait()
		assert.Equal(t, exec, true)
		assert.EqualError(t, err, "abc")
		assert.Nil(t, recoverValue)
	}
	// no err has panic
	{
		exec := false
		routine := xsync.NewRoutine()
		routine.Go(func() error {
			exec = true
			panic("xyz")
			return nil
		})
		err, recoverValue := routine.Wait()
		assert.Equal(t, exec, true)
		assert.NoError(t, err)
		assert.Equal(t, recoverValue, "xyz")
	}
}