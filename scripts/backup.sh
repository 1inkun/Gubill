#!/bin/sh
# Gubill SQLite 在线备份脚本
#
# 用法：
#   本地：  ./scripts/backup.sh
#   Docker：docker compose exec -T gubill /app/scripts/backup.sh
#   Cron：  0 2 * * * docker compose exec -T gubill /app/scripts/backup.sh >> /var/log/gubill-backup.log 2>&1
set -eu

DB_PATH="${DB_PATH:-/app/data/data.db}"
BACKUP_DIR="${BACKUP_DIR:-/app/backups}"

if [ ! -f "$DB_PATH" ]; then
  echo "数据库不存在: $DB_PATH" >&2
  exit 1
fi

mkdir -p "$BACKUP_DIR"
STAMP=$(date +%Y%m%d-%H%M%S)
OUT="$BACKUP_DIR/gubill-$STAMP.db"

# 优先使用 sqlite3 的在线备份接口，不可用时退化为文件拷贝
if command -v sqlite3 >/dev/null 2>&1; then
  sqlite3 "$DB_PATH" ".backup '$OUT'"
else
  cp "$DB_PATH" "$OUT"
fi

# 清理 30 天前的备份
find "$BACKUP_DIR" -name 'gubill-*.db' -mtime +30 -delete 2>/dev/null || true

echo "备份完成: $OUT"
