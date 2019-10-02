# select-channel 模板

说起 Go 的语言特性必然会提到其高并发的的特性，以及 channel 和 select，下面就总结几个《Concurrency in Go》中提到的模板。

## 约束 channel 的所有权

《Concurrency in Go》书中定义了 channel 的所有权为实例化、写入和关闭channel的goroutine。下面这个模板很好的约束了 channel 的所有权，并且把他以只读的方式传给其他channel

```go
func changeOwner() <-chan int {
    results := make(chan int, 5)
    go func() {
       defer close(results) // 所有权关闭 channel
       
       for i:= 0; i < 5; i++ {
           results <- i 
       }
    }()
    return results
}

func consumer(results <-chan int) {
    for result := range results {
        fmt.Printf("Received: %d\n", result)
    }
    
    fmt.Printf("Done receiving!\n")
}

results := changeOwner() 
consumer(results)
```

## 配合超时

```go
var c chan int
select {
    case <- c:  // 永远不会解锁, 因为 c 为 nil
    case <- time.After(3 * time.Second):
        fmt.Printf("Timed out\n")
}
```

## for-select

for-select 应该是最常见的模型了

```go
for _, s := range []string{"a", "b", "c"} {
    select {
        case <-done:
            // wait done
        case stringStream <- s:
    }
}
```

## for-select-default
```go
for {
    select {
    case <-done:
        // wait done
    default: 
        // 非抢占式的任务
    }
    // 或者这里执行非抢占式的任务
}
```

## cancel 

```go
newRandStream := func(done <- chan interface{}) <- chan int {
    randStream := make(chan int)
    
    go func() {
        defer fmt.Println("newRandStream closure exited.")
        defer close(randStream) // 结束关闭 randStream
        
        for {
            select {
                case randStream <- rand.Int():
                case done: // 等待外部 close
                    return 
            }
        }
    }()
    
    return randStream
}

done := make(chan interface{})
randStream := newRandStream(done)

fmt.Println("3 random ints:")
for i := 1; i <=3; i++ {
    fmt.Printf("%d,: %d\n", i, <-randStream)
}
close(done) // 关闭 done 会结束 newRandStream

// 等待

如果 goroutine 负责创建 goroutine， 它也负责确保它可以停止 goroutine

```

## or-done-channel

我们经常会从一个channel中连续读取数据，但是我们又没办法判断channel的状态是关闭还是继续存在，但是为了防止goroutine泄露，还需要做一些其他的操作。

比如在下面代码下我们没法确定是否能够正常退出，因为 myChan 可能不会不会关闭

```go
for val := range myChan {
    // 使用 val
}
```

仿照上例增加 done 代码会复杂不少：

```go
loop：
for {
    select {
        case <- done:
            break loop
        case maybeVal, ok := <- myChan:
            if ok == false {
                return 
            }
            // 使用 val
    }
}
```

下面看书中提到的模板

```go
orDone := func(done, c <- chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <- c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

// 调用
done := make(chan interface{})
	myChan := make(chan interface{})
	go func() {
		timer := time.After(time.Second)
		for {
			select {
			case myChan <- rand.Int():
			case <- timer:
				close(done); return
			}
		}
	}()
	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
```

orDone 实现可能会显得复杂，但是使用确实很方便
适用场景：需要从 channel 连续读取数据，但是又不能依靠这个 channel 来结束流程，需要增加一个 done 来控制流程

## tee-channel

适用场景：数据从一个channel进入，流向两个独立的 channel 

```go
tee := func(
		done    <-chan interface{},
		in      <-chan interface{},
	)(_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func(){
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	// 适用
	done    := make(chan interface{})
	in      := make(chan interface{})

	out1, out2 := tee(done, in)
	go func() {
		in <- int(1)
		time.Sleep(5000)
		defer close(done)
		// defer close(in)
	}()

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
```

其实这就是一分二

## 生成器

将一组离散值转换为一个 channel 上的数据流

```go
func generator(done <-chan interface{}, integers ...int) <-chan int {
    intStream := make(chan int)
    go func() {
        defer close(intStream)
        
        for _, integer := range integers {
            select {
            case <-done:
                return 
            case intStream <- integer:
            }
        }
    }()
    return intStream 
}
```

示例代码见 [goSamples](https://github.com/KevinBaiSg/goSamples/tree/master/tools/channel-select)