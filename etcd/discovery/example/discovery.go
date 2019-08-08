package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/KevinBaiSg/goSamples/etcd/discovery"
	"github.com/spf13/viper"
)

func main() {
	var role = flag.String("role", "", "master | worker")

	flag.Parse()

	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("filepath directory error ", err)
		return
	}
	viper.AddConfigPath(dir)

	if *role == "master" {
		master, _ := discovery.NewMaster("services/")
		master.WatchNodes()
	} else if *role == "worker" {

		serverName := "localhost"
		ip := "127.0.0.1"

		workerInfo := discovery.WorkerInfo{
			Name: serverName,
			IP: ip,
			CPU: runtime.NumCPU(),
		}
		worker, _ := discovery.NewWorker(serverName, workerInfo)

		go func() {
			time.Sleep(time.Second*20)
			worker.Stop()
		}()

		worker.Start()
	} else {
		fmt.Println("example -h for usage")
	}
}
