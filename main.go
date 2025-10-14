package main

import (
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"fleetpilot/usermanager"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	const configPath = "conf/config.yaml"

	// 1️. 检查或创建配置文件
	if err := config.EnsureConfigExists(configPath); err != nil {
		fmt.Println(err)
		return
	}

	// 2️. 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.InitLogger(cfg.Log.Level, nil)

	// 3. 配置路由

	router := gin.Default()
	user := router.Group("/user")
	{
		user.GET("/", usermanager.User)
		user.POST("/registry", usermanager.CreateUser)
	}

	// 启动服务
	router.Run("0.0.0.0:8000")
}
