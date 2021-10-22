package ibuffer

import "sync"

// Unbounded Unbounded是一个不使用额外goroutine实现的无界缓冲区，这通常用于在gRPC中将更新从一个实体传递到另一个实体。
// 该类型上的所有方法都是线程安全的，除了用于同步的底层互斥锁外，不会阻塞任何东西。
// Unbounded支持使用' interface{} '通道存储任何类型的值。
// 这意味着对Put()的调用会导致额外的内存分配，而且用户在读取时还需要进行类型断言。
// 对于性能关键的代码路径，强烈反对使用Unbounded，最好定义该缓冲区的新类型特定实现。
type Unbounded struct {
	c       chan interface{}
	mu      sync.Mutex
	backlog []interface{}
}

// NewUnbounded 新建Unbounded对象
func NewUnbounded() *Unbounded {
	return &Unbounded{c: make(chan interface{}, 1)}
}

// Put 添加t到unbounded缓冲区
func (b *Unbounded) Put(t interface{}) {
	b.mu.Lock()
	if len(b.backlog) == 0 {
		select {
		case b.c <- t:
			b.mu.Unlock()
			return
		default:
		}
	}
	b.backlog = append(b.backlog, t)
	b.mu.Unlock()
}

// Load 加载
func (b *Unbounded) Load() {
	b.mu.Lock()
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = nil
			b.backlog = b.backlog[1:]
		default:
		}
	}
	b.mu.Unlock()
}

// Get 获取
func (b *Unbounded) Get() <-chan interface{} {
	return b.c
}
