package main

import "github.com/coreos/etcd/client"

func NewClient() (client.Client, error)  {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}