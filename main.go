package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"

	"linuxdo-demo/controller"
)

func main() {
	controller.LoadEnv()
	controller.StartBlacklistSync()
	// // 1. 创建默认的路由引擎
	r := gin.Default()

	// OIDC兼容LinuxDO
	// r.GET("/oauth2/authorize", controller.CallAuthorize)

	// r.GET("/oauth2/callback", controller.Callback)

	// r.POST("/oauth2/token", controller.GetToken)

	r.GET("/api/user", controller.GetUser)

	// // 3. 启动服务 (默认在 0.0.0.0:8080)
	r.Run(":10110") // 也可以指定端口，如 r.Run(":8081")
}
