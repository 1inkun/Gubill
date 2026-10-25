package models

import (
	"fmt"
)

type BusinessErr struct {
	HttpCode int    `json:"code"`
	ErrCode  int    `json:"errCode"`
	Msg      string `json:"msg"`
	Detail   string `json:"detail"`
}

func (e *BusinessErr) Error() string {
	return fmt.Sprintf("ErrCode=%d, ErrMsg=%s", e.ErrCode, e.Msg)
}

func NewBusinessErr(code int, errCode int, msg string, detail string) *BusinessErr {
	return &BusinessErr{
		HttpCode: code,
		ErrCode:  errCode,
		Msg:      msg,
		Detail:   detail,
	}
}

// 业务错误
var (
	// http相关
	ErrBindingData = NewBusinessErr(400, 100001, "数据绑定失败", "请检查输入的数据是否有误")
	// 用户相关
	ErrWrongPasswd = NewBusinessErr(401, 101001, "密码错误", "请检查输入的密码是否正确")
	ErrLoginStatus = NewBusinessErr(401, 101002, "登录状态有误", "登录状态有误，请尝试重新登录")
	// 订单相关
	ErrExistUnFinishedOrder = NewBusinessErr(400, 102001, "签到失败", "存在尚未完成的订单")
	ErrIllegalOperation     = NewBusinessErr(401, 102002, "非法操作", "涉及非法操作，请联系管理员")
	ErrBadParams            = NewBusinessErr(400, 102003, "参数错误", "查询参数有误，请检查请求的参数")
	ErrNoDataExist          = NewBusinessErr(400, 102004, "找不到数据", "找不到对应的数据")
	// 商品相关
	ErrNoSuchGood = NewBusinessErr(400, 103001, "商品不存在", "数据库中找不到对应的商品信息")
	// 支付订单相关
	ErrNoSuchPayOrder          = NewBusinessErr(400, 103001, "订单不存在", "数据库中找不到对应的支付订单")
	ErrExistUnFinishedPayOrder = NewBusinessErr(400, 103002, "请结清付款", "存在尚未完成的支付订单")
)
