package main

import (
	"fmt"
)

func orDone(done, c <- chan interface{}) <-chan interface{} {
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
					fmt.Println("Timeout")
				}
			}
		}
	}()
	return valStream
}

func tee(
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

func bridge(
	done 		<-chan interface{},
	chanStream 	<-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)

		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
					stream = maybeStream
				}
			case <-done:
				return
			}
			for val := range orDone(done, stream){
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()

	return valStream
}