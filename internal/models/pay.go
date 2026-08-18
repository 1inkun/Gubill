package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pay struct {
	Basic
	BusinessType  string
	BusinessId    string `gorm:"unique"`
	UserId        string
	Value         int64
	Status        int64
	PayMethods    int64
	TransactionId string
	ExpireTime    int64
	PayAt         int64
}

func (p *Pay) BeforeCreate(tx *gorm.DB) (err error) {
	p.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}
