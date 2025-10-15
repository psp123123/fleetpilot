package usermanager

import (
	"fleetpilot/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义登陆信息结构体
type LoginInfo struct {
	Username string
	Password string
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
	logger.Debug("获取的登陆信息：", logininfo)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"token":   "123",
		"message": "ok",
	})
}
