package main

import (
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"fleetpilot/usermanager"

	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	password := "admin123"
	hash := "$2a$10$V1QkUuFHzGwlzP4XoQ3H1O5GqDg8Tb7RQOkBxMZfJdUqv7Wf5v9bi"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("密码不匹配:", err)
	} else {
		fmt.Println("登录成功")
	}

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
	auth := router.Group("/api").Use(usermanager.AuthMiddleware())
	{
		// 校验用户信息，并生成token
		auth.POST("/", usermanager.User)
	}

	// -- 公共路由
	router.POST("/login", usermanager.Login)
	router.POST("/registry", usermanager.CreateUser)

	// -- 刷新路由
	router.POST("/token/refresh", usermanager.RefreshHanlder)

	// 启动服务
	router.Run("0.0.0.0:8000")

}
