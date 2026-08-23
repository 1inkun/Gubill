package repository

import "gorm.io/gorm"

type MemberQuery struct {
	db *gorm.DB
}

func NewMemberQuery(db *gorm.DB) *MemberQuery {
	return &MemberQuery{db: db}
}
