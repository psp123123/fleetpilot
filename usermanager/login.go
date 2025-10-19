package usermanager

import (
	"fleetpilot/backend"
	"fleetpilot/common/logger"
	"net/http"
	"time"

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

	//
	cond := map[string]interface{}{
		"username": logininfo.Username,
		//"password_hash": logininfo.Password,
	}
	logger.Debug("get client info user:%v", cond["username"])

	// 获取数据库信息
	retUsername, retErr := backend.GetMysqlOneData("user", cond)
	logger.Debug("get info from mysql users password length:%v,value is %v", len(retUsername.PasswordHash), retUsername.PasswordHash)

	if retErr != nil || len(retUsername.Username) == 0 || !ComparePass(retUsername.PasswordHash, logininfo.Password) {
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

		// 设置HttpOnly Cookie存放refresh token
		ctx.SetCookie(
			"refreshToken",
			refreshToken,
			int(7*24*time.Hour.Seconds()),
			"/token/refresh",
			"",
			false,
			true,
		)

		logger.Debug("get data:%v", retUsername.UserID)

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,

			"data": gin.H{
				"accessToken": accessToken,
				"user":        retUsername.UserID,
				"username":    retUsername.Username,
				"userID":      retUsername.UserID,
				"nickname":    retUsername.Nickname,
			},
			"message": "ok",
		})
	}
}

// 刷新接口
func RefreshHanlder(ctx *gin.Context) {
	// 从cookie中读取refresh token
	refreshtoken, err := ctx.Cookie("refreshToken")
	if err != nil {
		logger.Error("get refreshtoken from client cookie error:%v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorized error"})
		return
	}

	// 验证refreshtoken
	claims, err := VerifyRefreshToken(refreshtoken)
	if err != nil {
		logger.Error("get refreshtoken from client cookie error:%v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorized error"})
		return
	}

	// 验证refresh无问题，重新生成
	newRefreshToken, err := GenerateRefreshoken(claims.UserID, claims.Username)
	if err != nil {
		logger.Error("general new refresh token failed")
	}
	newAccessToken, err := GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		logger.Error("genera access token error:%v", err)
	}

	// 下发新的refreshtoken
	ctx.SetCookie(
		"refreshToken",
		newRefreshToken,
		int(7*24*time.Hour.Seconds()),
		"/auth/refresh",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": newAccessToken,
	})

}
