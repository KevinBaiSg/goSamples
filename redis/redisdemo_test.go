package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var redisdb *redis.Client

func init() {
	redisdb = redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func ExampleNewClient() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: %s \n", err))
	}

	url := viper.GetString("environments.url")
	port := viper.GetInt("environments.port")
	password := viper.GetString("environments.password")

	redisdb = redis.NewClient(&redis.Options{
		Addr:     url + ":" + strconv.Itoa(port),
		Password: password,
		DB:       0,
	})
	pong, err := redisdb.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	err := redisdb.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisdb.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := redisdb.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}