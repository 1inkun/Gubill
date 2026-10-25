package models

type Good struct {
	Id          int64  `json:"id" example:"1"`
	Name        string `json:"name" example:"周卡会员"`
	Type        string `json:"type" example:"weekly"`
	Description string `json:"description" example:"在会员有效期内免费使用场地"`
	Status      int    `json:"status" example:"1"`
	Price       int    `json:"price" example:"3000"`
}

type BuyGoodResponse struct {
	PayOrderId int64 `json:"pay_order_id"`
}
