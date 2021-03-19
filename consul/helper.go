package main

import (
	"fmt"
	"log"

	consulApi "github.com/hashicorp/consul/api"
)

func Register(id string) {
	fmt.Println("test begin .")
	config := consulApi.DefaultConfig()
	// config.Address = "localhost"
	fmt.Println("defautl config : ", config)
	client, err := consulApi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	// 创建一个新服务。
	registration := new(consulApi.AgentServiceRegistration)
	registration.ID = id
	registration.Name = "user-tomcat"
	registration.Port = 8080
	registration.Tags = []string{"user-tomcat"}
	registration.Address = "127.0.0.1"

	// 增加check。
	check := new(consulApi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	// 设置超时 5s。
	check.Timeout = "5s"
	// 设置间隔 5s。
	check.Interval = "5s"
	// 注册check服务。
	registration.Check = check
	log.Println("get check.HTTP:",check)

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatal("register server error : ", err)
	}
}

func deregister(id string)  {
	fmt.Println("test begin .")
	config := consulApi.DefaultConfig()
	// config.Address = "localhost"
	fmt.Println("defautl config : ", config)
	client, err := consulApi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	err = client.Agent().ServiceDeregister(id)
	if err != nil {
		log.Fatal("register server error : ", err)
	}
}