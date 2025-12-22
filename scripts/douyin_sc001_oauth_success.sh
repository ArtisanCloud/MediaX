#!/usr/bin/env bash

set -euo pipefail

LOG_FILE=${1:-logs/accesstoken-server-info.log}
SAMPLE_COUNT=${2:-5}

if [ ! -f "$LOG_FILE" ]; then
  echo "日志文件不存在：$LOG_FILE" >&2
  exit 1
fi

if ! command -v rg >/dev/null 2>&1; then
  echo "需要 ripgrep (rg) 来解析日志，请先安装。" >&2
  exit 1
fi

all_events=$(rg --no-line-number --text 'accesstoken-server: event=oauth\.(complete|start\.error)' "$LOG_FILE" 2>/dev/null | rg 'provider_code=byte_dance_douyin' || true)
if [ -z "$all_events" ]; then
  echo "未在 $LOG_FILE 中找到 DouYin OAuth 相关日志。" >&2
  exit 1
fi

recent_events=$(printf '%s\n' "$all_events" | tail -n "$SAMPLE_COUNT")
attempts=0
success=0
fail_lines=()
while IFS= read -r line; do
  [ -z "$line" ] && continue
  attempts=$((attempts + 1))
  if [[ "$line" == *"event=oauth.complete"* ]]; then
    success=$((success + 1))
  else
    fail_lines+=("$line")
  fi
done <<<"$recent_events"

if [ "$attempts" -eq 0 ]; then
  echo "没有可统计的 OAuth 样本，请确认日志格式。" >&2
  exit 1
fi

rate=$(python3 - <<PY
success = $success
attempts = $attempts
print(f"{(success/attempts)*100:.2f}")
PY
)

echo "最近 ${attempts} 次 DouYin OAuth：成功 ${success} 次，成功率 ${rate}%（日志：$LOG_FILE）"

needs_help=$(RATE="$rate" python3 - <<'PY'
import os
rate = float(os.environ["RATE"])
print(1 if rate < 80 else 0)
PY
)

if [ "$needs_help" -eq 1 ]; then
  echo "⚠️ 成功率低于 80%，建议检查以下项："
  echo "   1) DouYin client_id/client_secret/scope 是否配置正确"
  echo "   2) redirect_url 与调试服务监听地址是否匹配"
  echo "   3) Redis/内存是否保留了最近的 Flow（可用 /api/oauth/tokens 排查）"
  echo "最近失败日志："
  for entry in "${fail_lines[@]}"; do
    echo "  - $entry"
  done
  exit 2
else
  echo "✅ 成功率满足 SC-001 要求（≥80%）。"
fi
