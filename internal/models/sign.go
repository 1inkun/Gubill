package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sign struct {
	Basic
	UserId  string
	StartAt int64
	EndAt   int64
	Status  int64
	Value   int64
}

type SignRes struct {
	UUID    string `json:"signId"`
	UserId  string `json:"userId"`
	StartAt int64  `json:"start_at"`
	EndAt   int64  `json:"end_at"`
	Status  int64  `json:"status"`
	Value   int64  `json:"value"`
}

func (s *Sign) BeforeCreate(tx *gorm.DB) (err error) {
	s.UUID = uuid.NewString()
	// log.Println(u.UUID)
	return nil
}
