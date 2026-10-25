### GuBill
基于 `Go` 编写的一个场馆按时计价应用,包含 Go 后端和小程序前端

## 需求
- Go >= 1.26
- Make
- Node.js >= 22

## 后端
后端基于 `Gin` 框架,在终端中输入以下命令开启开发服务器:
```bash
make dev
```

## 前端
前端为使用 `uni-app` 实现的小程序前端,使用以下命令启动开发服务器:
```bash
cd wxapp
npm run dev:mp-weixin
```
运行成功后使用 `微信开发者工具` 打开 `dist` 目录

## 许可证
本项目遵循 `MIT` 协议开源,详情参考 `LICENSE` 文件
