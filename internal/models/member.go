package models

import "time"

type SearchMemberParams struct {
	UserId int64 `json:"user_id" example:"100001"`
}

type MemberData struct {
	Id      int64     `json:"id" example:"100001"`
	UserId  int64     `json:"user_id" example:"100001"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Status  int64     `json:"status" example:"1"`
}
