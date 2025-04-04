package limiter

import (
	"sync"
	"time"
)

type Limiter interface {
	Do(unique string, f func(), timeout time.Duration)
}

type operation struct {
	timer *time.Timer
	f     func()
	mu    sync.Mutex
}

type limiter struct {
	ds map[string]*operation
	mu sync.Mutex
}

func New() Limiter {
	return &limiter{
		ds: make(map[string]*operation),
	}
}

func (o *operation) delay(f func(), timeout time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	// if func set call it and update time
	if o.f != nil {
		o.f()
		o.f = nil
		o.timer = time.AfterFunc(timeout, func() { o.delay(f, timeout) })
		return
	}
	// clear timer
	o.timer = nil
}

func (l *limiter) Do(unique string, f func(), timeout time.Duration) {
	// check if operation registered
	l.mu.Lock()
	defer l.mu.Unlock()
	o, ok := l.ds[unique]
	if !ok {
		o = &operation{}
		l.ds[unique] = o
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	// set timer if not set
	if o.timer == nil {
		// call single func
		f()
		o.timer = time.AfterFunc(timeout, func() { o.delay(f, timeout) })
		return
	}
	// rewrite func
	o.f = f
}
