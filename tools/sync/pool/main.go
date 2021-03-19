package main

import (
	"fmt"
	"time"
)

//func main()  {
//	myPool := &sync.Pool{New: func() interface{} {
//		log.Printf("Creating new instance")
//		return int(0)
//	}}
//
//	instance := myPool.Get().(int)
//	log.Println("instance = ", instance)
//	instance = 10
//	myPool.Put(instance)
//	instance = myPool.Get().(int)
//	log.Println("instance = ", instance)
//	myPool.Put(instance)
//	runtime.GC() // 不一定每次会成功
//	instance = myPool.Get().(int)
//	log.Println("instance = ", instance)
//}

func main()  {
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
					case out2<-val:
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
		time.Sleep(1000)
		defer close(done)
		// defer close(in)
	}()

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}