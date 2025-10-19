// 验证用户名和密码的代码
package usermanager

import (
	"fleetpilot/common/logger"

	"golang.org/x/crypto/bcrypt"
)

// 将密码字段生成bcrypt格式
func EncodeBcrypt(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密码转换失败")
	}
	return string(hash)
}
