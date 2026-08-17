# Gubill API 文档

## 通用约定

- 基础路径：用户端 `http://localhost:8080/api/v1`，管理端 `http://localhost:8081/api/v1`
- 响应体统一为 `{code, status, msg, data}`
- 认证：`Authorization` 请求头，兼容 `Bearer <token>` 与裸 token 两种格式
- 登录有效期：JWT 6 小时
- 分页：`page`（默认 1）、`page_size`（默认 10，上限 100）

## 用户端接口

### 注册

```http
POST /api/v1/user/register
Content-Type: application/json

{"username":"tom","nickname":"汤姆","password":"pass1234","email":"tom@example.com"}
```

密码要求：6-30 位，且同时包含小写字母与数字。

### 登录

```http
POST /api/v1/user/login
Content-Type: application/json

{"username":"tom","password":"pass1234"}
```

返回 `data.token`，后续请求放入 `Authorization` 头。

### 签到

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/sign` | 开始签到，返回 `signId` |
| GET | `/sign?status=0` | 查询自己的签到记录（按状态过滤） |
| GET | `/sign/:sign_id` | 签到详情 |
| PUT | `/sign/:sign_id` | 结算，生成 pays 并返回 `payId` |

### 会员

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/member_plan` | 全部会员计划 |
| GET | `/member_plan/:plan_id` | 计划详情 |
| POST | `/member_plan/:plan_id` | 生成会员订单 |
| GET | `/member_list` | 我的会员状态 |
| GET | `/member_order?page=1` | 我的会员订单 |
| GET | `/member_order/:order_id` | 订单详情 |
| DELETE | `/member_order/:order_id` | 取消订单 |
| POST | `/member_order/:order_id` | 结算订单 |

### 支付

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/pay?page=1` | 我的支付记录 |
| GET | `/pay/:pay_id` | 支付单详情 |

### 支付流程示例

```bash
# 1. 登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"pass1234"}' | jq -r .data.token)

# 2. 开始签到
SIGN_ID=$(curl -s -X POST http://localhost:8080/api/v1/sign \
  -H "Authorization: Bearer $TOKEN" | jq -r .data.signId)

# 3. 结算（TODO(支付接入)：接入真实渠道后返回支付二维码/链接）
PAY=$(curl -s -X PUT http://localhost:8080/api/v1/sign/$SIGN_ID \
  -H "Authorization: Bearer $TOKEN")
PAY_ID=$(echo $PAY | jq -r .data.payId)
echo "支付单: $PAY_ID"
```

### 真实支付渠道接入（TODO）

当前版本未接入支付渠道，支付动作（确认/退款）由支付模块统一操作 pays。接入方式：

1. 实现 `internal/payment.Gateway` 接口（微信 Native / 支付宝当面付，接口内已附示例注释）。
2. 在 `cmd/main/main.go` 的 `TODO(支付接入)` 处创建实例并注入。
3. 新增渠道异步回调接口：验签后调用 `PaymentService.ConfirmPaid(payId)`，状态机自动完成业务联动。

## 管理端接口

全部接口要求 `Role=Admin` 的 JWT。管理员创建方式见 README。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/member_plan` | 新增会员计划 |
| PUT | `/member_plan/:plan_id` | 修改会员计划 |
| GET | `/member_list?page=1` | 全部会员情况 |
| POST | `/member_list` | 添加会员 |
| PUT | `/member_list/:member_id` | 修改会员 |
| GET | `/member_list/:member_id` | 会员详情 |
| GET | `/member_order?page=1` | 全部会员订单 |
| GET | `/member_order/:order_id` | 订单详情 |
| PUT | `/member_order/:order_id` | 修改订单状态/金额 |
| GET | `/sign?page=1` | 全部签到记录 |
| GET | `/sign/:sign_id` | 签到详情 |
| PUT | `/sign/:sign_id` | 修改签到状态/金额 |
| GET | `/pay?page=1` | 全部支付记录 |
| GET | `/pay/:pay_id` | 支付单详情 |
| POST | `/pay/:pay_id/confirm` | 确认收款（兜底） |
| POST | `/pay/:pay_id/refund` | 全额退款 |

## 状态码与错误

- 成功：HTTP 200，`code=200, status="success"`
- 业务错误：HTTP 400/401/403/429，`code` 与 HTTP 状态一致，`status="fail"`
- 内部错误：HTTP 500，统一返回 `"内部错误"` 或数据库错误文案
- 限流：每 IP 10 req/s（burst 20），超出返回 429
