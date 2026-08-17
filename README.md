# Gubill

使用 Go 实现的按时计费场馆计费应用（学习项目，已预留真实支付渠道接入位）。

![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-支付接入预留-yellow)

> 当前未接入任何支付渠道；接入微信/支付宝只需实现 `internal/payment.Gateway` 接口，并在 `cmd/main` 的 `TODO(支付接入)` 标注处注入。

## 文档导航

- [架构文档](docs/architecture.md) - 系统分层、数据模型、状态机与核心流程
- [API 文档](docs/api.md) - 全部接口约定与调用示例
- [贡献指南](CONTRIBUTING.md) - 如何参与开发
- [安全策略](SECURITY.md) - 漏洞报告方式
- [行为准则](CODE_OF_CONDUCT.md)
- [变更日志](CHANGELOG.md)

## 功能特性

- **用户端 API**（`:8080`）：注册/登录（JWT）、签到计时、按时长阶梯计费结算、会员计划浏览与购买、支付记录查询。
- **管理端 API**（`:8081`）：会员计划管理、会员订单管理（确认收款/取消/退款）、会员管理、签到记录、支付记录查询。
- **统一支付网关抽象**（`internal/payment.Gateway`）：支付状态机、过期作废、取消联动、管理端确认收款/退款均已就绪，真实渠道接口预留（TODO 标注）。
- **管理员引导**：`gubill-cli create-admin` 命令创建首个管理员。

## 技术栈

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.26+ |
| Web 框架 | Gin |
| 数据库 | SQLite（`gorm.io/driver/sqlite`，依赖 CGO） |
| ORM | GORM v2 |
| 认证 | JWT（HS256，6 小时有效期） |
| 密码 | bcrypt |
| 限流 | 每 IP 10 req/s（burst 20），空闲 IP 自动清理 |

## 快速开始（本地）

```powershell
# 进入配置目录并启动
cd cmd/main
go run main.go
```

启动后：

- 用户端 API：http://localhost:8080
- 管理端 API：http://localhost:8081
- 数据库文件：`cmd/main/data/data.db`（首次启动自动建表）

### 创建管理员

```powershell
cd cmd/main
go run ../cli create-admin --username admin --password 你的密码 --email admin@example.com
```

然后用该账号调用管理端 API（`Authorization: Bearer <token>`）。

## 配置说明

遵循上游约定：**`cmd/main/config.yaml` 为唯一配置源**，启动时由 `config` 包读取并写入环境变量（`os.Setenv`），不可通过环境变量覆盖。

| 配置项（config.yaml 键） | 说明 | 默认值 |
|--------------------------|------|--------|
| `JWT.Salt` | JWT 签名密钥（生产环境务必修改） | 我是马牛逼 |
| `Billing.SinglePrice` | 签到单价（分） | 500 |
| `Pay.ExpireMinutes` | 支付单有效期（分钟） | 30 |
| `Server.PublicBaseUrl` | 服务公开地址（接入真实支付后用于构造回调地址） | http://localhost:8080 |

仅 `CONFIG_PATH` 支持通过环境变量指定其他配置文件路径（默认 `./config.yaml`），参考 [.env.example](.env.example)。

## 支付流程与真实渠道接入

```
结算（签到/会员订单）——业务终点：生成 pays status=0（30 分钟有效），返回 payId
支付动作——统一操作 pays：确认收款/渠道回调 → status=1 → 业务单联动
过期：定时任务 + 惰性检查自动作废（status=-1）
退款：管理端确认收款后全额退款（status=3）
```

业务服务（Sign/Member）与支付服务（PaymentService）无依赖：业务只负责生成 pays，支付所需数据全部从 pays 读取。

支付单状态：`0 待支付 / 1 已支付 / -1 已作废 / 3 已退款`。

**接入步骤**：

1. 在 `internal/payment/gateway.go` 中按 `Gateway` 接口注释（微信 Native / 支付宝当面付示例）实现真实渠道。
2. 在 `cmd/main/main.go` 的 `TODO(支付接入)` 标注处创建实例并注入 `PaymentService`。
3. 编写渠道回调接口（验签后调用 `PaymentService.ConfirmPaid`），业务状态机无需改动。

> 未接入支付渠道时，支付动作由管理端手工确认收款/退款完成。

## API 摘要

统一前缀 `/api/v1`，响应体 `{code, status, msg, data}`，JWT 放 `Authorization` 头（兼容 `Bearer ` 前缀与裸 token）。

用户端：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /user/register | 注册 |
| POST | /user/login | 登录 |
| POST | /sign | 开始签到 |
| GET | /sign?status= | 我的签到 |
| GET | /sign/:sign_id | 签到详情 |
| PUT | /sign/:sign_id | 结算（生成支付单） |
| GET | /member_plan | 全部会员计划 |
| GET | /member_plan/:plan_id | 计划详情 |
| POST | /member_plan/:plan_id | 生成会员订单 |
| GET | /member_list | 我的会员状态 |
| GET | /member_order?page=&page_size= | 我的会员订单 |
| GET | /member_order/:order_id | 订单详情 |
| DELETE | /member_order/:order_id | 取消订单 |
| POST | /member_order/:order_id | 结算订单 |
| GET | /pay?page=&page_size= | 我的支付记录 |
| GET | /pay/:pay_id | 支付单详情 |

管理端（需 Admin 角色）：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST/PUT | /member_plan[/:plan_id] | 新增/修改会员计划 |
| GET/POST/PUT | /member_list[/:member_id] | 会员管理 |
| GET | /member_order[/:order_id] | 会员订单查询 |
| PUT | /member_order/:order_id | 修改订单 |
| GET/PUT | /sign[/:sign_id] | 签到管理 |
| GET | /pay?page=&page_size= | 全部支付记录 |
| GET | /pay/:pay_id | 支付单详情 |
| POST | /pay/:pay_id/confirm | 确认收款 |
| POST | /pay/:pay_id/refund | 全额退款 |

## Docker 部署

```bash
# 1. 生产环境请先修改 cmd/main/config.yaml 中的 JWT.Salt 等配置
# 2. 构建并启动
docker compose up -d --build

# 3. 创建管理员（在容器内执行）
docker compose exec gubill ./gubill-cli create-admin --username admin --password 你的密码 --email admin@example.com
```

数据保存在 Docker 卷 `gubill-data`（容器内 `/app/data/data.db`）。

### 备份与恢复

```bash
# 手动备份
docker compose exec -T gubill /app/scripts/backup.sh

# 定时备份（宿主机 crontab）
0 2 * * * docker compose exec -T gubill /app/scripts/backup.sh >> /var/log/gubill-backup.log 2>&1

# 恢复：将备份文件复制回卷内数据库路径后重启
docker compose cp ./backups/gubill-xxx.db gubill:/app/data/data.db
docker compose restart gubill
```

## 从旧版本升级

> 适用于此前版本创建的数据库文件（`data.db`）。

- 本次版本移除了 `pays.business_id` 的**数据库唯一约束**（改为服务层保证"同一业务单同时最多一张有效支付单"）。SQLite 的 `AutoMigrate` 不会删除已存在的唯一索引，因此旧库若曾对同一业务单生成过两张支付单（过期作废后重新结算），仍可能触发 `UNIQUE constraint failed: pays.business_id`。
- 处理方式（二选一）：
  1. 测试/演示环境可直接删除旧 `data.db` 重新初始化；
  2. 保留数据则需重建 `pays` 表（复制旧数据到新表后替换），或自行执行表重建迁移脚本。
- 支付单过期后，签到/会员订单支持**直接重新结算**（自动作废旧单并生成新单），无需再手动重置状态。

## 测试

```bash
go vet ./...
go test ./...
```

覆盖：计费阶梯边界、JWT、支付状态机（幂等/过期/作废/退款）、签到/会员结算流程、接口集成（可见性/权限）。

## 项目结构

```
cmd/main       服务入口（用户端 + 管理端双端口）
cmd/cli        运维命令（create-admin）
internal/
  config       配置加载（YAML + 环境变量）
  models       数据模型与状态常量
  repository   SQLite 连接与建表
  payment      Gateway 接口（待接入真实支付渠道）
  service      业务层（user/sign/member/pay）
  handler      API 处理
  middlewares  限流 / JWT / 统一错误 / 管理页鉴权
  utils        计费 / 密码 / JWT / 分页
  router       路由注册
  testutil     测试工具（内存 SQLite + 测试 Fake 网关）
```

## 许可

[MIT](LICENSE)

参与本项目即表示你同意遵守 [行为准则](CODE_OF_CONDUCT.md)；提交前请阅读 [贡献指南](CONTRIBUTING.md)。
