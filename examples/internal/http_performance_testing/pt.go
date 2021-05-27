package pt

import (
	"fmt"
	xsync "github.com/goclub/sync"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
	"time"
)

type PT struct {
	Steps func() (steps []Step, err error)
	client *http.Client
}

type Step struct {
	Request *http.Request
	Check func(t *testing.T, step Step, resp *http.Response)
}
type RunOpt struct {
	Times uint64
	Interval time.Duration
}
func (pt PT) Run(t *testing.T, opt RunOpt) (err error) {

	counter := struct{
		residueTimes uint64
		sync.Mutex
	}{
		residueTimes: opt.Times,
	}
	if pt.client == nil {
		pt.client = &http.Client{}
	}
	routine := xsync.Routine{}
	ticker := time.NewTicker(opt.Interval)
	// 开启 routine 处理 <-ticker.C
	routine.Go(func() error {
		// 使用 for 配合 <-ticker.C 堵塞运行，直到从 ticker.C 接收到值
		for {
			<-ticker.C
			// 数据竞争
			counter.Lock()
			counter.residueTimes--
			counter.Unlock()
			// 千万不要写 < 0 会发生数据越界 18446744073709551615
			if counter.residueTimes <= 0 {
				break
			}
			// 开启 routine 异步发送请求，不同步等待请求响应
			routine.Go(func() (err error) {
				steps, err := pt.Steps() ; if err != nil {
				    return err
				}
				for _, step := range steps {
					resp, err := pt.client.Do(step.Request)
					assert.NoError(t, err)
					defer resp.Body.Close()
					if step.Check == nil {
						assert.Equal(t, resp.StatusCode, 200)
					} else {
						step.Check(t, step, resp)
					}
				}
				return nil
			})
		}
		return nil
	})
	err , recoverValue := routine.Wait() ; if err != nil {
	    return
	} ; if recoverValue != nil {
	    return fmt.Errorf("go routine panic: %v", recoverValue)
	}
	return
}