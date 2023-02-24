package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hildxd/backend-template/bootstrap"
	"github.com/hildxd/backend-template/global"
)

func main() {
	// 初始化配置
	bootstrap.InitializeConfig()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(fmt.Sprintf(":%s", global.App.Config.App.Port))
}
