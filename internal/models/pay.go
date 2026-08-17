package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pay struct {
	Basic
	UserId       string
	BusinessType string
	// BusinessId 不做数据库唯一约束：作废/过期的历史支付单允许被新单替换，
	// "同一业务单同一时间最多一张有效支付单"由服务层保证。
	BusinessId    string
	Channel       string
	Value         int64
	Status        int64
	PayMethods    int64
	TransactionId string
	ExpireTime    int64
	PayAt         int64
	RefundValue   int64
	RefundAt      int64
}

func (p *Pay) BeforeCreate(tx *gorm.DB) (err error) {
	p.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}

type PayRes struct {
	UUID          string `json:"payId"`
	UserId        string `json:"userId"`
	BusinessType  string `json:"businessType"`
	BusinessId    string `json:"businessId"`
	Channel       string `json:"channel"`
	Value         int64  `json:"value"`
	Status        int64  `json:"status"`
	TransactionId string `json:"transactionId"`
	ExpireTime    int64  `json:"expireTime"`
	PayAt         int64  `json:"payAt"`
	RefundValue   int64  `json:"refundValue"`
	RefundAt      int64  `json:"refundAt"`
	CreatedAt     int    `json:"createdAt"`
}

// ToPayRes 将 Pay 转换为响应结构
func (p *Pay) ToPayRes() PayRes {
	return PayRes{
		UUID:          p.UUID,
		UserId:        p.UserId,
		BusinessType:  p.BusinessType,
		BusinessId:    p.BusinessId,
		Channel:       p.Channel,
		Value:         p.Value,
		Status:        p.Status,
		TransactionId: p.TransactionId,
		ExpireTime:    p.ExpireTime,
		PayAt:         p.PayAt,
		RefundValue:   p.RefundValue,
		RefundAt:      p.RefundAt,
		CreatedAt:     p.CreatedAt,
	}
}
