package limiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/subliker/logger/zap"
)

func TestLimiter(t *testing.T) {
	assert := assert.New(t)

	limiter := New()
	type data struct {
		num  int
		time time.Time
	}
	c := make(chan data, 12)

	var wg sync.WaitGroup
	nums := 5
	interval := 500
	timeout := time.Second
	for i := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter.Do("test", func() {
				c <- data{
					num:  i,
					time: time.Now(),
				}
			}, timeout)
		}()
		time.Sleep(time.Millisecond * time.Duration(interval))
	}

	wg.Wait()
	l := time.Time{}
	time.Sleep(time.Second * 2)
	mn := timeout + time.Second
	for {
		select {
		case d := <-c:
			if (l != time.Time{}) {
				mn = min(mn, d.time.Sub(l))
				zap.Logger.Infof("interval: %s", d.time.Sub(l))
			}
			l = d.time
			continue
		default:
			close(c)
		}
		break
	}
	assert.Equal(true, mn >= timeout)
}
