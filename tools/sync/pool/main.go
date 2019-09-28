package main

import (
	"log"
	"runtime"
	"sync"
)

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
