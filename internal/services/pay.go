package services

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/repositories"
	"github.com/1inkun/billApp/internal/utils"
)

type PayService struct {
	payQuery *repositories.PayQuery
	mu       *sync.Mutex
}

func NewPayService(payQuery *repositories.PayQuery) *PayService {
	var mu sync.Mutex
	return &PayService{
		payQuery: payQuery,
		mu:       &mu,
	}
}

func (s *PayService) GetPayOrderDataById(ctx context.Context, orderId int64) (models.PayOrderData, error) {
	key := "payOrder:" + strconv.Itoa(int(orderId))
	// 先去缓存中查找
	var data models.PayOrderData
	data, err := s.payQuery.GetPayOrderDataByIdFromCache(ctx, key)
	if err != nil {
		// redis出错或数据不存在
		data = models.PayOrderData{}
	}
	if data.Id > 0 {
		// 数据获取成功
		return data, nil
	}
	// 再去数据库中查找
	data, err = s.payQuery.GetPayOrderDataById(ctx, orderId)
	if err != nil {
		return models.PayOrderData{}, err
	}
	// 存入缓存中
	dataStr, err := json.Marshal(data)
	if err != nil {
		log.Printf("在GetPayOrderDataById中序列化json时出错:%s", err)
		return models.PayOrderData{}, err
	}
	s.payQuery.SetCache(ctx, key, string(dataStr))
	return data, nil
}

func (s *PayService) GetPayOrderDataByUserId(ctx context.Context, params models.SearchPayOrderDataParams) (utils.PaginationData, error) {
	limit := params.PageSize
	if params.PageSize <= 0 || 20 < params.PageSize {
		params.PageSize = 10
		limit = 10
	}
	offset := limit * (params.Page - 1)
	if params.Page <= 0 || 200 < params.Page {
		params.Page = 1
		offset = limit * 0
	}
	data, err := s.payQuery.GetPayOrderDataByUserId(ctx, params.UserId, params.Status, offset, limit)
	if err != nil {
		return utils.PaginationData{}, nil
	}
	paginationData := utils.PaginationData{
		Data:     data,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	return paginationData, nil
}
