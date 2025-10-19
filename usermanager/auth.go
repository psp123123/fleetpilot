// 验证用户名和密码的代码
package usermanager

import "golang.org/x/crypto/bcrypt"

// 比对验证密码
func ComparePass(sourcepass, verifypass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(sourcepass), []byte(verifypass))
	return err == nil
}
