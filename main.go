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

	// 初始化日志
	bootstrap.InitializeLog()
	global.App.Log.Info("initialize log success")

	// 初始化数据库
	bootstrap.InitializeDB()
	// 程序结束前关闭数据库链接
	defer func() {
		if global.App.DB != nil {
			sqlDB, _ := global.App.DB.DB()
			sqlDB.Close()
		}
	}()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(fmt.Sprintf(":%s", global.App.Config.App.Port))
}
