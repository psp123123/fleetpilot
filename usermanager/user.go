package usermanager

import (
	"fleetpilot/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建用户接口
func CreateUser(ctx *gin.Context) {

}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			logger.Error("get client token is null")
			c.JSON(401, gin.H{"msg": "missing token"})
			c.Abort()
			return
		}

		claims, err := VerifyAccessToken(token)
		if err != nil {
			logger.Error("接口守卫验证失败: %v", err)
			c.JSON(401, gin.H{"msg": "invalid or expired token"})
			c.Abort()
			return
		}

		// token合法，设置用户信息到context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// 获取用户信息，根据客户端传来的access token认证
func GetUserInfo(ctx *gin.Context) {
	accessToken := ctx.GetHeader("Authorization")
	logger.Debug("get client token is: %v", accessToken)
	claims, err := VerifyAccessToken(accessToken)
	if err != nil {
		logger.Error("verify token error: ", err)
		ctx.JSON(401, gin.H{"msg": "invalid or expired token"})
		ctx.Abort()
		return
	}

	// 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"user": claims.Username,
	})
}
