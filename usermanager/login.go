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
	logger.Debug("get client info user:%v", cond["username"])

	// 获取数据库信息
	retUsername, retErr := backend.GetMysqlOneData("user", cond)
	logger.Debug("get info from mysql user info:%v", retUsername)

	if retErr != nil || len(retUsername.Username) == 0 {
		logger.Error("查询用户失败")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user invalid",
			"token":   "",
		})
	} else {
		// 生成token
		accessToken, err := GenerateAccessToken(retUsername.UserID, retUsername.Username)
		if err != nil {
			logger.Error("genera access token error:%v", err)
		}
		refreshToken, err := GenerateRefreshoken(retUsername.UserID, retUsername.Username)
		if err != nil {
			logger.Error("genera refresh token error:%v", err)
		}
		logger.Debug("get data:%v", retUsername)

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,

			"data": gin.H{
				"accessToken":  accessToken,
				"refreshToken": refreshToken,
			},
			"message": "ok",
		})
	}
}
