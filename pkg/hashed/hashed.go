package hashed

import (
	"golang.org/x/crypto/bcrypt"
)

// 加密密码
func Password(password string, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
