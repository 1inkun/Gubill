package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberList struct {
	Basic
	UserId string `gorm:"unique"`
	Status int64
	EndAt  int64
}

type MemberListRes struct {
	UUID   string `json:"memberId"`
	UserId string `json:"userId"`
	Status int64  `json:"status"`
	EndAt  int64  `json:"end_at"`
}

type MemberPlan struct {
	Basic
	Name        string
	Type        string `gorm:"unique"`
	Value       int64
	Description string
}

type MemberPlanRes struct {
	UUID        string `json:"planId"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       int64  `json:"value"`
	Description string `json:"des"`
}

type MemberOrders struct {
	Basic
	PlanId string
	UserId string
	Value  int64
	Status int64
}

type MemberOrderRes struct {
	UUID   string `json:"orderId"`
	PlanId string `json:"planId"`
	UserId string `json:"userId"`
	Value  int64  `json:"value"`
	Status int64  `json:"status"`
}

func (ml *MemberList) BeforeCreate(tx *gorm.DB) (err error) {
	ml.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}

func (mp *MemberPlan) BeforeCreate(tx *gorm.DB) (err error) {
	mp.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}

func (mo *MemberOrders) BeforeCreate(tx *gorm.DB) (err error) {
	mo.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}
