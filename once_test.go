package xsync_test

import (
	"errors"
	xsync "github.com/goclub/sync"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnce(t *testing.T) {
	{
		var once xsync.Once
		err := once.Do(func() (err error) {
			return nil
		})
		assert.NoError(t, err)
		err = once.Do(func() (err error) {
			return nil
		})
		assert.NoError(t, err)
	}
	{
		var once xsync.Once
		err1 := once.Do(func() (err error) {
			return errors.New("abc")
		})
		assert.Equal(t,err1.Error(), "abc")
		err2 := once.Do(func() (err error) {
			return errors.New("efg")
		})
		assert.NoError(t, err2)
	}
}