#!/usr/bin/env bash

set -euo pipefail

DEFAULT_BASE_URL=${BASE_URL:-http://127.0.0.1:7071}
DEFAULT_ITERATIONS=${ITERATIONS:-10}
API_TOKEN=${ACCESSTOKEN_API_TOKEN:-dev-accesstoken}
PAYLOAD_FILE=${PAYLOAD_FILE:-}
PAYLOAD_TEMPLATE='{"provider_code":"byte_dance_douyin","provider_app":"douyin","action":"douyin.video.list","payload":{"cursor":0,"count":1}}'

BASE_URL=${1:-$DEFAULT_BASE_URL}
ITERATIONS=${2:-$DEFAULT_ITERATIONS}

if [ -n "$PAYLOAD_FILE" ]; then
  if [ ! -f "$PAYLOAD_FILE" ]; then
    echo "指定的 payload 文件不存在：$PAYLOAD_FILE" >&2
    exit 1
  fi
  PAYLOAD=$(cat "$PAYLOAD_FILE")
else
  PAYLOAD=${PAYLOAD:-$PAYLOAD_TEMPLATE}
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "需要 curl 命令，请先安装。" >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "需要 python3 以计算 p95。" >&2
  exit 1
fi

durations=()
failures=0
fail_messages=()
tmp_resp=$(mktemp)

for i in $(seq 1 "$ITERATIONS"); do
  meta=$(curl -sS -o "$tmp_resp" -w "%{http_code} %{time_total}" \
    -H "Authorization: Bearer $API_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD" \
    "$BASE_URL/accesstoken/call" || true)

  code=$(printf '%s\n' "$meta" | awk '{print $1}')
  elapsed=$(printf '%s\n' "$meta" | awk '{print $2}')
  durations+=("$elapsed")

  if [ "${code:-000}" -ge 400 ]; then
    failures=$((failures + 1))
    body=$(cat "$tmp_resp")
    fail_messages+=("第 ${i} 次返回 ${code}: ${body}")
  fi
done

rm -f "$tmp_resp"

if [ "${#durations[@]}" -eq 0 ]; then
  echo "没有可用的请求数据，Benchmark 终止。" >&2
  exit 1
fi

p95=$(printf '%s\n' "${durations[@]}" | python3 - <<'PY'
import sys, math
values = [float(line.strip()) for line in sys.stdin if line.strip()]
if not values:
    print("0.0")
    sys.exit(0)
values.sort()
idx = max(math.ceil(0.95 * len(values)) - 1, 0)
print(f"{values[idx]:.4f}")
PY
)

avg=$(printf '%s\n' "${durations[@]}" | python3 - <<'PY'
import sys
values = [float(line.strip()) for line in sys.stdin if line.strip()]
if not values:
    print("0.0")
    sys.exit(0)
print(f"{sum(values)/len(values):.4f}")
PY
)

echo "调用基准完成：迭代 ${#durations[@]} 次，失败 ${failures} 次。"
echo "平均耗时 ${avg}s，p95=${p95}s"

if [ "$failures" -gt 0 ]; then
  echo "⚠️ 以下请求失败："
  for msg in "${fail_messages[@]}"; do
    echo "  - $msg"
  done
fi

threshold_flag=$(P95="$p95" python3 - <<'PY'
import os
p95 = float(os.environ["P95"])
print(1 if p95 <= 2.0 else 0)
PY
)

if [ "$threshold_flag" -eq 1 ]; then
  echo "✅ p95 ≤ 2 秒，符合 SC-003 要求。"
else
  echo "⚠️ p95 超过 2 秒，请检查 DouYin API 状态或本地网络。"
  exit 2
fi
