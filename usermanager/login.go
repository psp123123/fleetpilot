package usermanager

import (
	"fleetpilot/backend"
	"fleetpilot/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义登陆信息结构体
type LoginInfo struct {
	Username string `json:"user"`
	Password string `json:"passwd"`
}

func Login(ctx *gin.Context) {
	// 获取登陆信息
	var logininfo LoginInfo
	err := ctx.ShouldBind(&logininfo)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误" + err.Error(),
		})
	}

	cond := map[string]interface{}{
		"username": logininfo.Username,
	}
	logger.Debug("get client info user:", cond["username"])
	// 获取数据库信息
	ret, retErr := backend.GetMysqlOneData("user", cond)
	if retErr != nil || len(ret) == 0 {
		logger.Error("查询用户失败")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"token":   "",
			"message": "user invalid",
		})
	}

	logger.Info("get data", ret)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,

		"token": "123",

		"message": "ok",
	})
}
