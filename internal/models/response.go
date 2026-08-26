package models

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

type PaginData struct {
	Results  any `json:"results"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func NewResponse(code int, status string, msg string, data any) Response {
	var Response = Response{
		code,
		status,
		msg,
		data,
	}
	return Response
}
