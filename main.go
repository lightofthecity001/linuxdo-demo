package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 创建默认的路由引擎
	r := gin.Default()

	// 2. 定义路由
	r.GET("/authorize", func(c *gin.Context) {
		queryParam := c.Request.URL.Query()
		fmt.Println("111111111111", queryParam)
		// 返回 JSON 格式的数据
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 3. 启动服务 (默认在 0.0.0.0:8080)
	r.Run(":10110") // 也可以指定端口，如 r.Run(":8081")
}
