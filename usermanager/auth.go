// 验证用户名和密码的代码
package usermanager

import (
	"fleetpilot/common/logger"

	"golang.org/x/crypto/bcrypt"
)

// 比对验证密码
func ComparePass(sourcepass, verifypass string) bool {
	logger.Debug("sourcePass is %v,verifypass is %v", sourcepass, verifypass)
	err := bcrypt.CompareHashAndPassword([]byte(sourcepass), []byte(verifypass))
	if err != nil {
		logger.Error("compare password error:%v", err)
	}
	return err == nil
}
