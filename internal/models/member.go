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

type MemberPlan struct {
	Basic
	Name        string
	Type        string `gorm:"unique"`
	Value       int64
	Description string
}

type MemberOrders struct {
	Basic
	PlanId string
	UserId string
	Value  int64
	Status int64
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
