#!/usr/bin/env bash
# 后端契约回归：注册 -> 登录 -> 解码 JWT 取 user_id -> 签到 -> 查订单 -> 结算 -> 在场名单
# 用法：BASE=http://host.docker.internal:8080/api/v1 bash scripts/e2e.sh
set -euo pipefail
BASE=${BASE:-http://host.docker.internal:8080/api/v1}
TS=$(date +%s)
U="e2e$TS"          # 用户名 4-20 位，纯数字后缀
P="zxc123"

step() { printf '\n=== %s ===\n' "$1"; }

step "register"
curl -sS -X POST "$BASE/register" -H 'content-type: application/json' \
  -d "{\"username\":\"$U\",\"nickname\":\"E2E测试\",\"password\":\"$P\",\"email\":\"$U@example.com\"}"; echo

step "login"
LOGIN=$(curl -sS -X POST "$BASE/login" -H 'content-type: application/json' \
  -d "{\"username\":\"$U\",\"password\":\"$P\"}")
echo "$LOGIN"
TOKEN=$(printf '%s' "$LOGIN" | sed -n 's/.*"token" *: *"\([^"]*\)".*/\1/p')
[ -n "$TOKEN" ] || { echo "FAIL: 未拿到 token"; exit 1; }
echo "token长度=${#TOKEN}（内容不在日志中显示）"

step "decode jwt user_id"
USER_ID=$(node -e 'const p=process.argv[1].split(".")[1];console.log(JSON.parse(Buffer.from(p,"base64url")).user_id)' "$TOKEN")
echo "user_id=$USER_ID"

step "checkin (POST /order)"
curl -sS -X POST "$BASE/order" -H "Authorization: Bearer $TOKEN" -H 'content-type: application/json' \
  -d "{\"user_id\":$USER_ID}"; echo

step "list running orders (status=0)"
curl -sS "$BASE/order?user_id=$USER_ID&status=0&page=1&page_size=10" -H "Authorization: Bearer $TOKEN"; echo

step "list all orders (status=-2)"
ALL=$(curl -sS "$BASE/order?user_id=$USER_ID&status=-2&page=1&page_size=10" -H "Authorization: Bearer $TOKEN")
echo "$ALL"
ORDER_ID=$(printf '%s' "$ALL" | sed -n 's/.*"id":\([0-9]*\).*/\1/p')
echo "order_id=$ORDER_ID"

step "settle (PUT /order/:id)"
[ -n "$ORDER_ID" ] || { echo "FAIL: 未拿到 order_id"; exit 1; }
curl -sS -X PUT "$BASE/order/$ORDER_ID" -H "Authorization: Bearer $TOKEN"; echo

step "order detail by id"
curl -sS "$BASE/order?id=$ORDER_ID" -H "Authorization: Bearer $TOKEN"; echo

step "status (在场名单)"
curl -sS "$BASE/status"; echo

echo
echo "OK: 全流程跑通"
