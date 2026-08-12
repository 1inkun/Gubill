package service

import "gorm.io/gorm"

type UserService struct {
	db *gorm.DB
}

func NewUserServer(db *gorm.DB) *UserService {
	return &UserService{db: db}
}
