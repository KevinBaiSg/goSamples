# sync.Pool

在高并发的情况下任何东西的生成和销毁都会影响到最终的性能，由于 Go 使用了 GC 的方式进行内存管理，用过的对象会被销毁掉。为了能够重用对象，Go 为我们提供了 sync.Pool 这个工具。

## 使用方法
```go
func main()  {
	myPool := &sync.Pool{New: func() interface{} {
		log.Printf("Creating new instance")
		return int(0)
	}}

	instance := myPool.Get().(int)
	log.Println("instance = ", instance)
	instance = 10
	myPool.Put(instance)
	instance = myPool.Get().(int)
	log.Println("instance = ", instance)
	myPool.Put(instance)
	runtime.GC() // 不一定每次会成功
	instance = myPool.Get().(int)
	log.Println("instance = ", instance)
}

```

上边这段代码比较简单的说明了 sync.Pool 怎么使用，我们只需要构造一个 sync.Pool 结构体，并且把 New 这一项赋值。我在最后一部分增加了一次 runtime.GC()，多数情况下你就能看到如下的输出

```
2019/09/28 16:30:12 Creating new instance
2019/09/28 16:30:12 instance =  0
2019/09/28 16:30:12 instance =  10
2019/09/28 16:30:12 Creating new instance
2019/09/28 16:30:12 instance =  0
```
也就是说我们 put 进去的 instance 被销毁了，是调用我们定义 Pool 时的 New 方法重新产生的。在平时我们不用调用 runtime.GC()，go 的 GC 也会做同样的事情，所以 sync.Pool 是不能够当线程池来使用的。

既然是介绍 sync.Pool 当然要先说说它的好处
```go
func TestPool(t *testing.T)  {
	var numCalcsCreated int
	calcPool := &sync.Pool{New: func() interface{} {
		numCalcsCreated++
		mem := make([]byte, 1024)
		return &mem
	}}

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			//(*mem)[0] = 'a'
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	t.Log(numCalcsCreated, "calculators were created")
}
```
这段代码是 《Concurrency in Go》 中的一个示例，运行这段代码可以发现最后的计算结果还是很感动的, 值有8次 create。
```
=== RUN   TestPool
--- PASS: TestPool (0.28s)
    pool_test.go:35: 8 calculators were created
PASS

Process finished with exit code 0
```

这两个示例代码也可以在 [Github](https://github.com/KevinBaiSg/goSamples/tree/master/tools/sync/pool) 看到。

## 使用场景
sync.Pool是可伸缩的，并发安全的。其大小仅受限于内存的大小，可以被看作是一个存放可重用对象的值的容器。

任何存放区其中的值可以在任何时候被删除而不通知，在高负载下可以动态的扩容，在不活跃时对象池会收缩。

## 源码分析-v1.13

- 结构

```go
// Local per-P Pool appendix.
type poolLocalInternal struct {
	private interface{} // Can be used only by the respective P.
	shared  poolChain   // Local P can pushHead/popHead; any P can popTail.
}

type poolLocal struct {
	poolLocalInternal

	// Prevents false sharing on widespread platforms with
	// 128 mod (cache line size) = 0 .
	pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}
```

> poolLocalInternal 定义了两个值:
>> *private* 用于保存 P 私有的值
>> *shared* 用于保存 shared 的值，Local P 可以通过 pushHead/popHead 获取，其他 P 可以通过 popTail获取，但是这两种获取都需要加锁。

- 两个接口

```go
// Put adds x to the pool.
func (p *Pool) Put(x interface{}) {
	if x == nil { // 如果 x 是 nil，直接返回
		return
	}
	
	....
	
	l, _ := p.pin() // 绑定 goroutine 到 P 上，返回 P 的 poolLocal 和 Pid
	if l.private == nil { // 判断 private 是否为 nil
		l.private = x      // 保存 x 到 private
		x = nil
	}
	if x != nil {   // 如果 x 不为空，保存到 shared 中
		l.shared.pushHead(x)
	}
	runtime_procUnpin() // 调用 p.pin() 后必须调用该函数
	
	...
}

func (p *Pool) Get() interface{} {
	...
	
	l, pid := p.pin()   // 绑定 goroutine 到 P 上，返回 P 的 poolLocal 和 Pid
	x := l.private      // 读取 private
	l.private = nil
	
	/*
	1. 判断 x 是否为空，如果为空，尝试从 shared 中获取，如果还是失败，调用 getSlow() 从其他 P 中获取
	*/
	if x == nil {       
		// Try to pop the head of the local shard. We prefer
		// the head over the tail for temporal locality of
		// reuse.
		x, _ = l.shared.popHead()
		if x == nil {
			x = p.getSlow(pid)
		}
	}
	runtime_procUnpin()
	
	...
	
	if x == nil && p.New != nil {
		x = p.New()
	}
	return x
}
```

基本的注释我已经写在上边了，这里我们再把整个过程走一下。
当调用 sync.Pool 的 Get 或者 Put 方法时，都会先把 goroutine 固定到某个 P 上，就是代码中的 p.pin() 这个操作，同时会返回 p 的 poolLocal，这个结构中会有私有和共享两种对象，私有对象只有对应的 P 能够访问，因为一个P同一时间只能执行一个goroutine，因此对私有对象存取操作是不需要加锁的。共享列表是和其他P分享的，因此操作共享列表是需要加锁的。


## 建议
就当给自己之后使用 sync.Pool 总结一些建议

* sync.Pool 只是对对象的复用，不可以当做对象保存（比如连接池），因为 Pool 中的对象随时可能被 GC 了。
* Get() 操作最少0次加锁(直接读取private)，最坏的情况可能会有多次加锁，包括全局锁，所以这块要特别注意。
* Put() 操作最少0次加锁（直接写入private），最坏1次加锁。
* sync.Pool 的开销也是不可以忽略的。

## 注意的问题

曹大在 [《几个 Go 系统可能遇到的锁问题》](http://xargin.com/lock-contention-in-go/) 提到一个问题，就是在使用 **[fasttemplate](https://github.com/valyala/fasttemplate)** 时，由于内部使用了 sync.Pool, 但是开发者并没有复用 Template 对象，从而造成了 sync.Pool 频繁的进入 p.pinSlow 流程，以及频繁的触发 allPoolsMu（全局锁）。这种情况平时并不会有问题，但是全局锁在高并发的时候，绝对是性能杀手。
从文中提到的案例可见，sync.Pool 的目的是为了性能提升，但是错误的使用会造成更坏的情况发生。而且案例中的情况更加特殊，开发者使用的是 Template，只是使用的时候触发了副作用，所以提早 **压测**。

## 参考
1. [sync.Pool 源码](https://golang.org/src/sync/pool.go)   
2. Concurrency in Go
3. [几个 Go 系统可能遇到的锁问题](http://xargin.com/lock-contention-in-go/) 
4. [深入Golang之sync.Pool详解](https://www.cnblogs.com/sunsky303/p/9706210.html)    
5. [go语言的官方包sync.Pool的实现原理和适用场景](https://blog.csdn.net/yongjian_lian/article/details/42058893) 

