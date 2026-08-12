package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenNewPasswdHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CheckPassword(passwordHash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, err
}
