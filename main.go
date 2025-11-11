package main

import (
	api "fleetpilot/api"
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"fleetpilot/task"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	manager := api.GetHandlerManager()
	if err := manager.Register(&task.NmapClientParams{}); err != nil {
		logger.Error("注册nmap工具失败", err)
		os.Exit(1)
	}
}
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
	config.GlobalCfg = cfg
	logger.InitLogger(config.GlobalCfg.Log.Level, nil)
	fmt.Printf("Loaded config: %+v\n", config.GlobalCfg)

	router := gin.Default()

	// -- 受保护的路由，需要 Access Token
	auth := router.Group("/api").Use(api.AuthMiddleware())
	{
		// 校验用户信息，并生成token
		auth.GET("/userinfo", api.GetUserInfo)

		// 处理ws消息请求
		auth.GET("/ws", api.WsHandler)
	}

	// -- 公共 路由
	router.POST("/login", api.Login)
	router.POST("/registry", api.CreateUser)

	// -- 刷新 路由
	router.POST("/token/refresh", api.RefreshHanlder)

	// 启动 服务
	router.Run("0.0.0.0:8000")

}
