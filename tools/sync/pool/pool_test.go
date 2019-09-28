package main

import (
	"sync"
	"testing"
)

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
