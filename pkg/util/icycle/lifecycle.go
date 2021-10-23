package icycle

import (
	"sync"
	"sync/atomic"
)

// Cycle ...
type Cycle struct {
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	done    chan struct{}
	quit    chan error
	closing uint32
	waiting uint32
}

// NewCycle 新增循环周期
func NewCycle() *Cycle {
	return &Cycle{
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
		quit:    make(chan error),
		closing: 0,
		waiting: 0,
	}
}

// Run 启动一个goroutine
func (c *Cycle) Run(fn func() error) {
	c.mu.Lock()
	//todo add check options panic before waiting
	defer c.mu.Unlock()
	c.wg.Add(1)
	go func(c *Cycle) {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.quit <- err
		}
	}(c)
}

// Done 阻塞并返回一个chan错误
func (c *Cycle) Done() <-chan struct{} {
	if atomic.CompareAndSwapUint32(&c.waiting, 0, 1) {
		go func(c *Cycle) {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.wg.Wait()
			close(c.done)
		}(c)
	}
	return c.done
}

// DoneAndClose 阻塞并关闭
func (c *Cycle) DoneAndClose() {
	<-c.Done()
	c.Close()
}

// Close 关闭
func (c *Cycle) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if atomic.CompareAndSwapUint32(&c.closing, 0, 1) {
		close(c.quit)
	}
}

// Wait 阻塞生命周期
func (c *Cycle) Wait() <-chan error {
	return c.quit
}
