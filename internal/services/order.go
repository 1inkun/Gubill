package services

import (
	"context"
	"time"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/repositories"
	"github.com/1inkun/billApp/internal/utils"
)

type OrderService struct {
	orderQuery *repositories.OrderQuery
}

func NewOrderService(orderQuery *repositories.OrderQuery) *OrderService {
	return &OrderService{
		orderQuery: orderQuery,
	}
}

func (s *OrderService) NewOrder(ctx context.Context, userId int64) (models.OrderResponse, error) {
	// 检查是否存在尚未完成的订单
	ok, err := s.orderQuery.CheckBeforeNew(ctx, userId)
	if err != nil {
		return models.OrderResponse{}, err
	}
	if !ok {
		return models.OrderResponse{}, models.ErrExistUnFinishedOrder
	}
	// 检查是否存在尚未完成的支付订单
	ok, err = s.orderQuery.CheckPayBeforeNew(ctx, userId)
	if err != nil {
		return models.OrderResponse{}, err
	}
	if !ok {
		return models.OrderResponse{}, models.ErrExistUnFinishedPayOrder
	}
	// 新建订单
	newOrder := models.NewOrder{
		UserId:  userId,
		StartAt: time.Now().Unix(),
		Status:  0,
	}
	// log.Println(newOrder)
	orderId, err := s.orderQuery.NewOrder(ctx, newOrder)
	if err != nil {
		return models.OrderResponse{}, nil
	}
	// 用得到的orderId去获取全部的orderData
	data, err := s.orderQuery.GetOrderDataById(ctx, orderId)
	if err != nil {
		return models.OrderResponse{}, err
	}
	return data, nil
}

func (s *OrderService) GetOrderDataById(ctx context.Context, orderId int64, userId int64) (models.OrderResponse, error) {
	if orderId < 0 {
		return models.OrderResponse{}, models.ErrBadParams
	}
	// 先从redis缓存中查找
	data, err := s.orderQuery.GetOrderDataFromCache(ctx, orderId)
	if err != nil {
		return models.OrderResponse{}, err
	}
	if data.Id <= 0 {
		// 表示还没找到数据,去数据库中继续查找
		data, err = s.orderQuery.GetOrderDataById(ctx, orderId)
		if err != nil {
			return models.OrderResponse{}, err
		}
	}
	// 验证查询是否合法
	if data.UserId != userId {
		return models.OrderResponse{}, models.ErrIllegalOperation
	}
	return data, nil
}

func (s *OrderService) GetOrderDataByUserId(ctx context.Context, params models.SearchOrderParams) (any, error) {
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
	datas, err := s.orderQuery.GetOrderDataByUserId(ctx, params.UserId, params.Status, limit, offset)
	if err != nil {
		return nil, err
	}
	paginationData := utils.PaginationData{
		Data:     datas,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	return paginationData, nil
}

func (s *OrderService) FinishOrder(ctx context.Context, orderId int64, userId int64) (models.FinishOrderResponse, error) {
	// 检查订单信息
	orderData, err := s.orderQuery.GetOrderDataById(ctx, orderId)
	if err != nil {
		return models.FinishOrderResponse{}, err
	}
	if orderData.UserId != userId || orderData.Status != 0 {
		return models.FinishOrderResponse{}, models.ErrIllegalOperation
	}
	// 确认无误,修改订单状态并生成支付订单
	res, err := s.orderQuery.FinishOrder(ctx, orderId)
	if err != nil {
		return models.FinishOrderResponse{}, err
	}

	if res > 0 {
		// 拼接返回的内容
		resp := models.FinishOrderResponse{
			PayOrderId: res,
		}
		return resp, nil
	} else {
		return models.FinishOrderResponse{}, nil
	}
}
