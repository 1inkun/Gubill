package models

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
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
