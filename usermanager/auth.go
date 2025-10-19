// 验证用户名和密码的代码
package usermanager

import (
	"fleetpilot/common/logger"

	"golang.org/x/crypto/bcrypt"
)

// 比对验证密码
func ComparePass(hashedPassword, inputPassword string) bool {
	logger.Debug("hashedPassword: %v, inputPassword: %v", hashedPassword, inputPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		logger.Error("compare password error: %v", err)
		return false
	}
	return true
}
