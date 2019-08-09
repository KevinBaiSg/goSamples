package common

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

func NewClient() (*clientv3.Client, error)  {
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	endpoints := viper.GetStringSlice("endpoints")
	user := viper.GetString("user")
	password := viper.GetString("password")

	cfg := clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: time.Second * 30,
		Username: user,
		Password: password,
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}
