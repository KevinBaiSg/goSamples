package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		value = int8(0)
		m sync.Mutex
		c *sync.Cond
	)

	c = sync.NewCond(&m)

	go func() {
		m.Lock()
		for value == int8(1) {
			c.Wait()
			continue
		}
		fmt.Println("value == int8(1)")
		value = 0
		m.Unlock()
	}()

	fmt.Println("m.Lock()")
	m.Lock()
	value = int8(1)
	c.Signal()
	m.Unlock()
}