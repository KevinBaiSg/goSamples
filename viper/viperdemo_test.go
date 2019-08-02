package redis

import (
	"fmt"

	"github.com/spf13/viper"
)

func ExampleNewClient() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found; ignore error if desired"))
		} else {
			panic(fmt.Errorf("Config file was found but another error was produced, error: %s \n", err))
		}
	}
	fmt.Println(err)


	dbType := viper.GetString("type")
	fmt.Println(dbType)

	url := viper.GetString("environments.url")
	fmt.Println(url)

	port := viper.GetInt("environments.port")
	fmt.Println(port)

	password := viper.GetString("environments.password")
	fmt.Println(password)

	// Output: <nil>
	// redis
	// 127.0.0.1
	// 6379
	// 123456

}