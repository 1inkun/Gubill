package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Basic struct {
	UUID      string `gorm:"primaryKey;unique;column:uuid"`
	CreatedAt int
	UpdatedAt int
	DeletedAt gorm.DeletedAt
}

type User struct {
	Basic
	UserName      string `gorm:"unique;column:username;index"`
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
	return nil
}
