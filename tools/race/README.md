
这个例子看起来没任何问题，但是实际上，time.AfterFunc是会另外启动一个goroutine来进行计时和执行func()。
由于func中有对t(Timer)进行操作(t.Reset)，而主goroutine也有对t进行操作(t=time.AfterFunc)。
这个时候，其实有可能会造成两个goroutine对同一个变量进行竞争的情况。

golang在1.1之后引入了竞争检测的概念。我们可以使用go run -race 或者 go build -race 来进行竞争检测。
golang语言内部大概的实现就是同时开启多个goroutine执行同一个命令，并且纪录每个变量的状态。
```bash
go run -race race.go
```

当然这个参数会引发CPU和内存的使用增加，所以基本是在测试环境使用，不是在正式环境开启。


