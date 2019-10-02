package main

import (
	"fmt"
	"testing"
)

func TestChannel(t *testing.T)  {
	
}

func ExampleChannel() {
	c := make(chan int)
	done := make(chan struct{})

	//defer close(c)
	defer close(done)
	
	for i := 0; i < 10; i++ {
		go func(index int, ch chan int, d chan struct{}) {
			for {
				select {
				case <- d:
					return
				case value, ok := <- c:
					if ok == false {
						return
					}

					fmt.Println("index = ", index, "value = ", value)
				}
			}
		}(i, c, done)
	}

	for i := 0; i < 1000; i++ {
		c <- i
	}

	// Output:
	//
}
