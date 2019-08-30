# 信号量

信号量（英语：semaphore）又称为信号标，是一个同步对象，用于保持在0至指定最大值之间的一个计数值。
当线程完成一次对该semaphore对象的等待（wait）时，该计数值减一；当线程完成一次对semaphore对象的释放（release）时，
计数值加一。当计数值为0，则线程等待该semaphore对象不再能成功直至该semaphore对象变成signaled状态。
semaphore对象的计数值大于0，为signaled状态；计数值等于0，为nonsignaled状态.
semaphore对象适用于控制一个仅支持有限个用户的共享资源，是一种不需要使用忙碌等待（busy waiting）的方法。

## Golang 中 semaphore 数据结构

```go 
type waiter struct {
	n     int64
	ready chan<- struct{} 
}

func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}

type Weighted struct {
	size    int64
	cur     int64
	mu      sync.Mutex
	waiters list.List
}
```

## 主要方法

```go
func (s *Weighted) Acquire(ctx context.Context, n int64)
func (s *Weighted) TryAcquire(n int64)
func (s *Weighted) Release(n int64)
```

## 与其他类型对比
