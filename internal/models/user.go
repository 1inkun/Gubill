package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Basic struct {
	UUID      string `gorm:"primaryKey;column:uuid"`
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

type User struct {
	Basic
	UserName      string
	NickName      string
	AvatarFile    string
	Email         string
	PasswordHash  string
	Role          string
	LastLoginIP   string
	LastLoginDate int64
	RegisterIP    string
	RegisterDate  int64
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.NewString()
	var userData []User
	db := tx.Where("uuid = ?", u.UUID).Find(userData)
	if db.Error != nil {
		return db.Error
	}
	if len(userData) != 0 {
		return errors.New("生成了无效的UUID")
	}
	return nil
}
