package main

import (
	"fmt"
	"math/rand"
	"time"
)

func ExampleOrDone()  {
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
	// Output:
	//
}

func ExampleTee()  {
	done    := make(chan interface{})
	in      := make(chan interface{}, 5)

	out1, out2 := tee(done, in)
	go func() {
		defer func() {
			fmt.Println("close done")
			close(done)
		}()
		// defer close(in)
		for i := 0; i < 5; i++ {
			in <- int(i)
		}
		//in <- int(1)
		time.Sleep(time.Second)
	}()

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
	// Output:
	//
}