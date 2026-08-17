# 贡献指南

感谢你对 Gubill 感兴趣！这是一个学习性质的按时计费场馆计费项目，任何形式的贡献（代码、文档、Issue 反馈）都欢迎。

## 工作流程

1. **先提 Issue**：Bug、新功能或疑问请先创建 Issue 讨论，避免重复劳动与方向偏差。
2. **Fork 并创建分支**：从 `main` 切出功能分支，命名建议 `feat/xxx`、`fix/xxx`、`docs/xxx`。
3. **开发与自测**：本地跑通 `go vet ./...` 与 `go test ./...`，新逻辑补充对应测试。
4. **提交 Pull Request**：描述改动动机、实现方式与测试情况；维护者会尽快 review。

## 本地开发环境

### 依赖

- Go 1.26+
- 无需 C 编译器（SQLite 使用纯 Go 驱动）

### 启动

```bash
cd cmd/main
go run main.go
```

首次使用需创建管理员：

```bash
cd cmd/main
go run ../cli create-admin --username admin --password 你的密码 --email admin@example.com
```

### 代码规范

- 统一使用 `gofmt` 格式化，提交前运行 `gofmt -l .` 确认无输出。
- 通过 `go vet ./...` 与 `go test ./...`。
- 业务状态流转优先复用 `internal/models` 中的状态常量，避免魔法数字。
- 支付相关改动必须经过 `PaymentService`，不直接操作支付表。
- 注释与提交信息建议中文，保持简洁明确。

### 提交信息约定

建议遵循以下前缀：

- `feat:` 新功能
- `fix:` 缺陷修复
- `docs:` 文档
- `refactor:` 重构
- `test:` 测试
- `chore:` 构建/依赖/杂项

示例：`fix: 修正支付单过期后无法重新结算的问题`

## Pull Request 检查清单

- [ ] 分支基于最新的 `main`
- [ ] `gofmt -l .` 无输出
- [ ] `go vet ./...` 通过
- [ ] `go test ./...` 全绿（新增/修改逻辑有对应测试）
- [ ] README / docs 在涉及行为变化时已同步更新

## 行为准则

请阅读并遵守 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。
