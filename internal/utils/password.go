package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswdHash(passwd string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(passwd), 10)
	if err != nil {
		return "", err
	}
	passwdHash := string(hashByte)
	return passwdHash, nil
}

func ComparePasswdHash(passwdHash string, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwdHash), []byte(passwd))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
