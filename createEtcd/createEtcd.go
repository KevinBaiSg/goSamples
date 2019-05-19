package main

import (
	"fmt"
	_ "net/http"

	"context"
	"github.com/coreos/etcd/client"
)

func main() {
	c, err := NewClient()
	if err != nil {
		// handle error
		fmt.Printf("NewClient error: %v", err)
		return
	}

	kAPI := client.NewKeysAPI(c)

	_, err = kAPI.Create(context.Background(), "/foo", "bar")
	if err != nil {
		fmt.Printf("Create error: %v", err)
	}

	// delete the newly created key only if the value is still "bar"
	response, err := kAPI.Delete(context.Background(), "/foo", &client.DeleteOptions{PrevValue: "bar"})
	if err != nil {
		fmt.Printf("Delete error: %v\n", err)
	}

	fmt.Printf("Delete success: %v\n", response)

}
