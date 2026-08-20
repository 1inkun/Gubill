package service

import (
	"context"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/repository"
)

type PayService struct {
	payQuery *repository.PayQuery
}

func NewPayService(payQuery *repository.PayQuery) *PayService {
	return &PayService{payQuery: payQuery}
}

func (s *PayService) GetUserPayOrders(ctx context.Context, userId string, page int, pageSize int) ([]models.PayRes, error) {
	// 交由仓储层处理
	res, err := s.payQuery.GetUserPayOrders(ctx, userId, page, pageSize)
	if err != nil {
		return nil, err
	}
	// 获取表中的数据量
	nums, err := s.payQuery.GetUserPayOrderDataNums(ctx, userId)
	if err != nil {
		return nil, err
	}
	// log.Println(nums)
	_ = nums
	return res, nil
}
