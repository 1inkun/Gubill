package models

import "time"

type NewOrderRequest struct {
	UserId int64 `json:"user_id" binding:"required" example:"100001"`
}

// 用于向数据库中插入的数据
type NewOrder struct {
	UserId  int64 `json:"user_id"`
	StartAt int64 `json:"start_at"`
	Status  int64 `json:"status"`
}

// 在响应中的数据
type OrderResponse struct {
	Id      int64     `json:"id"`
	UserId  int64     `json:"user_id"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Status  int64     `json:"status"`
	Price   int64     `json:"price"`
}

// 结算时的响应
type FinishOrderResponse struct {
	PayOrderId int64 `json:"pay_order_id"`
}

// 查询方法所用的参数
type SearchOrderParams struct {
	Id       int64 `form:"id"`
	UserId   int64 `form:"user_id"`
	Status   int64 `form:"status"`
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
