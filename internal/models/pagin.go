package models

type Pagin struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}
