package service

import (
	"context"
	"errors"

	"github.com/1inkun/Gubill/internal/models"
	"gorm.io/gorm"
)

type MemberService struct {
	db *gorm.DB
}

func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{db: db}
}

var (
	ErrUniqueErr = models.NewBusinessError(400, "已添加过该数据")
)

func (s *MemberService) AddNewMemberPlan(ctx context.Context, name string, planType string, value int64, des string) (string, error) {
	var newData models.MemberPlan
	newData.Name = name
	newData.Type = planType
	newData.Value = value
	newData.Description = des
	// 直接利用Type的唯一性标签
	err := gorm.G[models.MemberPlan](s.db).Create(ctx, &newData)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", ErrUniqueErr
		}
		return "", ErrDatabaseErr
	}
	return newData.UUID, nil
}
