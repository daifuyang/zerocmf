// 验证码
package code

import (
	"math/rand"
	"time"
)

// 生成随机验证码
func GenerateRandomCode(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	charSet := "0123456789" // 使用数字字符集
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		code[i] = charSet[randomIndex]
	}
	return string(code)
}
