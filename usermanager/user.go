package usermanager

import "github.com/gin-gonic/gin"

// 获取用户的接口
func User(ctx *gin.Context) {

}

// 创建用户接口
func CreateUser(ctx *gin.Context) {

}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"msg": "missing token"})
			c.Abort()
			return
		}

		claims, err := VerifyAccessToken(token)
		if err != nil {
			c.JSON(401, gin.H{"msg": "invalid or expired token"})
			c.Abort()
			return
		}

		// token合法，设置用户信息到context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
