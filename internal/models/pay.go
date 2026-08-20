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

type PayRes struct {
	UUID         string `json:"payId"`
	BusinessType string `json:"businessType"`
	BusinessId   string `json:"businessId"`
	UserId       string `json:"userId"`
	Value        int64  `json:"value"`
	Status       int64  `json:"status"`
	ExpireTime   int64  `json:"expire_time"`
	PayAt        int64  `json:"pay_at"`
}

func (p *Pay) BeforeCreate(tx *gorm.DB) (err error) {
	p.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}
