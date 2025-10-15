// 关于所有的token相关内容
package usermanager

import (
	"errors"
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义Claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成token，包含access Token & refresh Token
// accessToken
func GenerateAccessToken(userID string, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.GlobalCfg.Jwt.AccessExp)),
			Issuer:    config.GlobalCfg.Jwt.IssuedAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(config.GlobalCfg.Jwt.AccessSecret))
}

// refreshToken
func GenerateRefreshoken(userID string, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.GlobalCfg.Jwt.RefreshExp)),
			Issuer:    config.GlobalCfg.Jwt.IssuedAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(config.GlobalCfg.Jwt.RefreshSecret))
}

// 验证token
// 验证accessToken
func VerifyAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.GlobalCfg.Jwt.AccessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	logger.Error("invalid token")
	return nil, errors.New("invalid token")
}

// 验证refreshToken
func VerifyRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return config.GlobalCfg.Jwt.RefreshSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	logger.Error("refresh token invalid")
	return nil, errors.New("refresh token invalid")
}

// 刷新接口生成新的accessToken
func RefreshTokenHandler(refreshToken string) (string, error) {
	// 验证refreshToken
	claims, err := VerifyRefreshToken(refreshToken)
	if err != nil {
		logger.Error("refreshToken has problem")
		return "", err
	}

	return GenerateAccessToken(claims.UserID, claims.Username)
}
