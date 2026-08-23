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

type LoginData struct {
	Username string `json:"username" binding:"required,alphanum,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=30,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}

type RegisterData struct {
	UserName string `json:"username" binding:"required,alphanum,min=4,max=20"`
	NickName string `json:"nickname" binding:"min=0,max=16"`
	Password string `json:"password" binding:"required,min=6,max=30,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	Email    string `json:"email" binding:"required,email"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}
