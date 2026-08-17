package models

// 支付单状态
const (
	PayStatusPending  = 0  // 待支付
	PayStatusPaid     = 1  // 已支付
	PayStatusVoid     = -1 // 已作废（过期或取消）
	PayStatusRefunded = 3  // 已退款
)

// 支付单业务类型
const (
	PayBusinessSign   = "sign"
	PayBusinessMember = "member"
)

// 签到状态
const (
	SignStatusActive     = 0 // 进行中
	SignStatusFinished   = 1 // 已完成
	SignStatusPendingPay = 2 // 待支付
)

// 会员订单状态
const (
	MemberOrderStatusCreated    = 0  // 已创建/待结算
	MemberOrderStatusPaid       = 1  // 已支付
	MemberOrderStatusPendingPay = 2  // 待支付
	MemberOrderStatusCanceled   = -1 // 已取消
)

// 会员状态
const (
	MemberStatusInvalid = 0 // 已失效
	MemberStatusValid   = 1 // 有效
)

// PayStatusText 返回支付单状态的中文描述
func PayStatusText(status int64) string {
	switch status {
	case PayStatusPending:
		return "待支付"
	case PayStatusPaid:
		return "已支付"
	case PayStatusVoid:
		return "已作废"
	case PayStatusRefunded:
		return "已退款"
	default:
		return "未知状态"
	}
}

// SignStatusText 返回签到状态的中文描述
func SignStatusText(status int64) string {
	switch status {
	case SignStatusActive:
		return "进行中"
	case SignStatusFinished:
		return "已完成"
	case SignStatusPendingPay:
		return "待支付"
	default:
		return "未知状态"
	}
}

// MemberOrderStatusText 返回会员订单状态的中文描述
func MemberOrderStatusText(status int64) string {
	switch status {
	case MemberOrderStatusCreated:
		return "待结算"
	case MemberOrderStatusPaid:
		return "已支付"
	case MemberOrderStatusPendingPay:
		return "待支付"
	case MemberOrderStatusCanceled:
		return "已取消"
	default:
		return "未知状态"
	}
}
