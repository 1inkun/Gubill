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
	UserName      string `gorm:"unique;column:username"`
	NickName      string `gorm:"column:nickname"`
	AvatarFile    string
	Email         string `gorm:"unique"`
	PasswordHash  string
	Role          string
	LastLoginIP   string
	LastLoginDate int64
	RegisterIP    string
	RegisterDate  int64
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.NewString()
	// log.Println(u.UUID)
	var userData []User
	db := tx.Where("uuid = ?", u.UUID).Find(&userData)
	if db.Error != nil {
		return db.Error
	}
	if len(userData) != 0 {
		return errors.New("生成了无效的UUID")
	}
	return nil
}
