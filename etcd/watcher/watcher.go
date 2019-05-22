package main

import (
	"context"
	"github.com/KevinBaiSg/goSamples/etcd/common"
	"github.com/coreos/etcd/client"
	"log"
)

func main() {
	c, err := common.NewClient()
	if err != nil {
		log.Fatal("Error: NewClient Failed:", err)
		return
	}

	kAPI := client.NewKeysAPI(c)
	ctx := context.Background()

	watcher := kAPI.Watcher("/foo", &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(ctx)
		if err != nil {
			log.Println("Error watch workers:", err)
			break
		}
		switch res.Action {
		case "expire":
			log.Println("Watcher Next Action: ", res.Action)
		case "get", "set", "delete", "update", "create", "compareAndSwap", "compareAndDelete":
			log.Printf("Watcher Next Action: %s; Value: %v", res.Action, res.Node.Value)
		}
	}
}

