package models

type PayOrderData struct {
	Id            int64  `json:"id"`
	UserId        int64  `json:"user_id"`
	Type          string `json:"type"`
	OrderNo       int64  `json:"order_no"`
	TransactionId string `json:"transaction_id"`
	Provider      string `json:"privider"`
	ClientIP      string `json:"client_ip"`
	Status        int64  `json:"status"`
	TotalFee      int64  `json:"total_fee"`
	Description   string `json:"description"`
}

type SearchPayOrderDataParams struct {
	Id       int64 `form:"id"`
	UserId   int64 `form:"user_id"`
	Status   int64 `form:"status"`
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
