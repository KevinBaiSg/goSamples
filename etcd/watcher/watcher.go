package main

import (
	"context"
	"log"
	"path/filepath"

	. "github.com/KevinBaiSg/goSamples/etcd/common"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

func main() {
	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("filepath directory error ", err)
		return
	}
	viper.AddConfigPath(dir)

	c, err := NewClient()
	if err != nil {
		log.Fatal("Error: NewClient Failed:", err)
		return
	}

	watcher := clientv3.NewWatcher(c)
	ctx := context.Background()

	log.Printf("start watch")

	for {
		watchChans := watcher.Watch(ctx, "/events", clientv3.WithPrefix())
		for watch := range watchChans {
			for _, event := range watch.Events {
				log.Printf("Watch: %s %q: %q \n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}
}

