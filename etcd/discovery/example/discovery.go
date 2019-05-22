package main

import (
	"flag"
	"fmt"
	"github.com/KevinBaiSg/goSamples/etcd/discovery"
	"runtime"
	"time"
)

func main() {
	var role = flag.String("role", "", "master | worker")

	flag.Parse()

	endpoints := []string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:22379",
		"http://127.0.0.1:32379",
	}

	if *role == "master" {
		master, _ := discovery.NewMaster(endpoints, "services/")
		master.WatchNodes()
	} else if *role == "worker" {

		serverName := "localhost"
		ip := "127.0.0.1"

		workerInfo := discovery.WorkerInfo{
			Name: serverName,
			IP: ip,
			CPU: runtime.NumCPU(),
		}
		worker, _ := discovery.NewWorker(serverName, workerInfo, endpoints)

		go func() {
			time.Sleep(time.Second*20)
			worker.Stop()
		}()

		worker.Start()
	} else {
		fmt.Println("example -h for usage")
	}
}
