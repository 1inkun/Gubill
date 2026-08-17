// Package payment 定义支付网关抽象。
//
// 模拟支付已回退移除；接入微信/支付宝等真实渠道时，
// 只需实现本接口并在 cmd/main 的 TODO 标注处注入，业务层无需改动。
package payment

import "context"

// Order 描述一笔待创建的下单请求。
type Order struct {
	PayId        string // 本地支付单 UUID
	UserId       string // 支付用户
	BusinessType string // 业务类型：sign / member
	BusinessId   string // 业务单 UUID
	Value        int64  // 金额（单位：分）
	ExpireTime   int64  // 支付单过期时间（Unix 秒）
}

// CreateResult 是网关下单成功后的返回信息。
type CreateResult struct {
	PayUrl         string // 支付链接 / 二维码内容（微信 Native、支付宝当面付等）
	ChannelOrderNo string // 渠道侧订单号
}

// RefundRequest 描述退款请求。
type RefundRequest struct {
	PayId  string // 本地支付单 UUID
	Value  int64  // 原支付金额（分）
	Refund int64  // 退款金额（分），当前仅支持全额退款
}

// Gateway 是支付网关的统一接口。
//
// TODO(支付接入)：参考各渠道官方文档实现本接口，例如：
//   - 微信支付 Native：CreateOrder 调用统一下单换取 code_url（二维码内容）；
//     VerifyNotifySignature 按官方文档校验回调签名与金额；Refund 调用退款接口。
//   - 支付宝当面付：CreateOrder 调用预下单获取二维码；回调验签使用支付宝公钥。
//
// 实现完成后在 cmd/main 的 TODO 标注处创建实例并注入 PaymentService。
type Gateway interface {
	// Name 返回渠道标识，如 "wechat" / "alipay"，会写入支付单 Channel 字段。
	Name() string
	// CreateOrder 调用渠道下单接口，返回支付链接/二维码与渠道订单号。
	CreateOrder(ctx context.Context, order *Order) (*CreateResult, error)
	// VerifyNotifySignature 校验渠道异步通知的签名（验签失败应拒绝回调）。
	VerifyNotifySignature(payload []byte, signature string) bool
	// Refund 调用渠道退款接口（失败时返回错误，业务侧将中止退款）。
	Refund(ctx context.Context, req *RefundRequest) error
	// NewTransactionId 生成渠道交易号（支付成功回调携带的真实交易号优先）。
	NewTransactionId() string
}
