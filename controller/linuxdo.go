package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type LinuxdoUser struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	TrustLevel int    `json:"trust_level"`
	Silenced   bool   `json:"silenced"`
}

func CallAuthorize(c *gin.Context) {
	fmt.Println("111111111111111111", c.Request.URL.Query())
}
