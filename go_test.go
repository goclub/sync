package xsync_test

import (
	"context"
	xsync "github.com/goclub/sync"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	errRecoverCh := xsync.Go(ctx, func() (err error) {
		time.Sleep(time.Millisecond*100)
		return nil
	})
	errRecover := <-errRecoverCh
	assert.ErrorIs(t, errRecover.Err, context.DeadlineExceeded)
}