package main

import (
	"context"
	"log"
	_ "net/http"
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
		log.Fatal("Error: NewClient Failed: ", err)
		return
	}

	ctx := context.Background()

	kv := clientv3.NewKV(c)

	putResponse, err := kv.Put(ctx, "/foo", "bar")
	if err != nil {
		log.Fatal("Error: Put Failed:", err)
	}
	if putResponse.PrevKv == nil {
		log.Print("Success: Put Ok; PrevKv KeyValue empty")
	} else {
		log.Print("Success: Put Ok; PrevKv KeyValue:", string(putResponse.PrevKv.Value))
	}

	getResponse, err := kv.Get(ctx, "/foo")
	if err != nil {
		log.Fatal("Error: Get Failed:", err)
	}
	if len(getResponse.Kvs) == 0 {
		log.Print("Success: Get Ok; Value is empty")
	} else {
		log.Print("Success: Get Ok; first Value:", string(getResponse.Kvs[0].Value))
	}

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()
	//response, err = kAPI.Set(ctx, "/ping", "pong", nil)
	//if err != nil {
	//	if err == context.DeadlineExceeded {
	//		// request took longer than 5s
	//	} else {
	//		// handle error
	//	}
	//}

	deleteResponse, err := kv.Delete(ctx, "/foo", clientv3.WithPrevKV()) //
	if err != nil {
		log.Fatal("Error: Delete Failed:", err)
	}

	log.Print("Success: Delete Ok; deleteResponse:", string(deleteResponse.PrevKvs[0].Value))

	getResponse, err = kv.Get(ctx, "/foo")
	if err != nil {
		log.Fatal("Error: Get Failed:", err)
	}
	if len(getResponse.Kvs) == 0 {
		log.Print("Success: Get Ok; Value is empty")
	} else {
		log.Print("Success: Get Ok; first Value:", string(getResponse.Kvs[0].Value))
	}
}
