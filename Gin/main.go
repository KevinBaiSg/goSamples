package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化引擎
	engine := gin.Default()
	// 注册一个路由和处理函数
	v2 := engine.Group("/v2", middleware, middleware2)
	v2.Any("/login", WebRoot)
	// 绑定端口，然后启动应用
	engine.Run(":8080")

}

func WebRoot(context *gin.Context) {
	log.Println("WebRoot start")
	context.String(http.StatusOK, "hello, world")
	log.Println("WebRoot end")
}

func middleware(context *gin.Context) {
	log.Println("middleware start")
	log.Println("middleware before next")
	context.Next()
	log.Println("middleware after next")
	log.Println("middleware end")
}

func middleware2(context *gin.Context) {
	log.Println("middleware2 start")
	log.Println("middleware2 before next")
	context.Next()
	log.Println("middleware2 after next")
	log.Println("middleware2 end")
}