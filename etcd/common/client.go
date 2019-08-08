package common

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

func NewClient() (*clientv3.Client, error)  {
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	endpoints := viper.GetStringSlice("endpoints")

	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}
