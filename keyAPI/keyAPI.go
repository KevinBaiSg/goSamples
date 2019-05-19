package main

import (
	"context"
	"github.com/KevinBaiSg/etcdSample/common"
	"github.com/coreos/etcd/client"
	"log"
	_ "net/http"
)

func main() {
	c, err := common.NewClient()
	if err != nil {
		log.Fatal("Error: NewClient Failed:", err)
		return
	}

	kAPI := client.NewKeysAPI(c)
	ctx := context.Background()

	response, err := kAPI.Set(ctx, "/foo", "bar", nil)
	if err != nil {
		log.Fatal("Error: Set Failed:", err)
	}
	log.Print("Success: Set Ok; Response:", response)

	response, err = kAPI.Get(ctx, "/foo", nil)
	if err != nil {
		log.Fatal("Error: Get Failed:", err)
	}
	log.Print("Success: Get Ok; Value:", response.Node.Value)

	response, err = kAPI.Update(ctx, "/foo", "bar2")
	if err != nil {
		log.Fatal("Error: Update Failed:", err)
	}
	log.Print("Success: Update Ok; Response:", response)

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

	response, err = kAPI.Get(ctx, "/foo", nil)
	if err != nil {
		log.Fatal("Error: Get Failed:", err)
	}
	log.Print("Success: Get Ok; Value:", response.Node.Value)

	response, err = kAPI.Delete(ctx, "/foo", &client.DeleteOptions{PrevValue: "bar2"})
	if err != nil {
		log.Fatal("Error: Delete Failed:", err)
	}

	log.Print("Success: Delete Ok; Response:", response)

	response, err = kAPI.Get(ctx, "/foo", nil)
	if err != nil {
		log.Fatal("Error: Get Failed:", err)
	}
	log.Print("Success: Get Ok; Value:", response.Node.Value)
}
