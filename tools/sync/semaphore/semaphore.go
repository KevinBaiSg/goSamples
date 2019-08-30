package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func main()  {
	var wg sync.WaitGroup

	maxWorkers := runtime.GOMAXPROCS(0)
	log.Printf("maxWorkers: %d ", maxWorkers)

	sem := semaphore.NewWeighted(int64(maxWorkers))
	ctx := context.TODO()

	for i := 0; i < 64; i++ {
		wg.Add(1)

		go func(i int) {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("task: %d Failed to acquire semaphore: %v", i, err)
				return
			}
			time.Sleep(time.Duration(2) * time.Second)

			log.Printf("task: %d successful", i)
			sem.Release(1)
			wg.Done()
		}(i)
	}

	wg.Wait()
	log.Printf("task all finish",)
}
