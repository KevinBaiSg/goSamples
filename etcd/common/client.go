package common

import (
	"github.com/coreos/etcd/clientv3"
)

func NewClient() (*clientv3.Client, error)  {
	endpoints := []string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:22379",
		"http://127.0.0.1:32379",
	}

	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}