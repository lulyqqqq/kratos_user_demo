package util

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt 加密密码
func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// ComparePasswords 验证密码
func ComparePasswords(password string, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(inputPassword))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
