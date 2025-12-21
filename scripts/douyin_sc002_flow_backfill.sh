#!/usr/bin/env bash

set -euo pipefail

DEFAULT_BASE_URL=${BASE_URL:-http://127.0.0.1:7071}
DEFAULT_LIMIT=${LIMIT:-20}
API_TOKEN=${ACCESSTOKEN_API_TOKEN:-dev-accesstoken}
EXPECTED_BACKEND=${EXPECTED_BACKEND:-}

BASE_URL=${1:-$DEFAULT_BASE_URL}
LIMIT=${2:-$DEFAULT_LIMIT}

if ! command -v curl >/dev/null 2>&1; then
  echo "需要 curl 命令，请先安装。" >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "需要 python3 以解析 JSON。" >&2
  exit 1
fi

flows_json=$(curl -fsS -H "Authorization: Bearer $API_TOKEN" \
  "$BASE_URL/api/oauth/tokens?provider_code=byte_dance_douyin&limit=$LIMIT")

tmp=$(mktemp)
printf '%s' "$flows_json" >"$tmp"
mapfile -t FLOW_ROWS < <(python3 <<'PY' "$tmp" "$LIMIT"
import json, sys
path = sys.argv[1]
limit = int(sys.argv[2])
with open(path, "r", encoding="utf-8") as fh:
    data = json.load(fh)
flows = data.get("flows", [])[:limit]
for flow in flows:
    identifier = flow.get("flow_id") or ""
    backend = flow.get("storage_backend") or "-"
    status = flow.get("status") or "-"
    print(f"{identifier}|{backend}|{status}")
PY
)
rm -f "$tmp"

total=${#FLOW_ROWS[@]}
if [ "$total" -eq 0 ]; then
  echo "未在 $BASE_URL/api/oauth/tokens 取得 DouYin Flow，请先通过 /debug 生成授权。" >&2
  exit 1
fi

success=0
fail_list=()
backend_mismatch=0
declare -A backend_count=()

for row in "${FLOW_ROWS[@]}"; do
  IFS='|' read -r flow_id backend status <<<"$row"
  if [ -z "$flow_id" ]; then
    continue
  fi
  response_tmp=$(mktemp)
  http_meta=$(curl -sS -o "$response_tmp" -w "%{http_code}" \
    -H "Authorization: Bearer $API_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"flow_id\":\"$flow_id\"}" \
    "$BASE_URL/accesstoken/flow/replay" || true)

  if [ "$http_meta" = "200" ]; then
    success=$((success + 1))
  else
    fail_body=$(cat "$response_tmp")
    fail_list+=("$flow_id ($http_meta): $fail_body")
  fi
  rm -f "$response_tmp"

  backend_count["$backend"]=$(( ${backend_count["$backend"]:-0} + 1 ))
  if [ -n "$EXPECTED_BACKEND" ] && [ "$backend" != "$EXPECTED_BACKEND" ]; then
    backend_mismatch=$((backend_mismatch + 1))
  fi
done

success_rate=$(python3 - <<PY
success = $success
total = $total
print(f"{(success/total)*100:.2f}")
PY
)

echo "Flow 回填成功 ${success}/${total}，成功率 ${success_rate}%（基准 URL：$BASE_URL）"
echo "storage_backend 分布："
for key in "${!backend_count[@]}"; do
  echo "  - ${key:-unknown}: ${backend_count[$key]} 条"
done

if [ -n "$EXPECTED_BACKEND" ] && [ "$backend_mismatch" -gt 0 ]; then
  echo "⚠️ 有 ${backend_mismatch} 条 Flow 的 storage_backend ≠ ${EXPECTED_BACKEND}，请确认正在验证的模式。"
fi

if [ "${#fail_list[@]}" -gt 0 ]; then
  echo "以下 Flow 回填失败："
  for item in "${fail_list[@]}"; do
    echo "  - $item"
  done
  exit 2
fi

threshold_pass=$(RATE="$success_rate" python3 - <<'PY'
import os
rate = float(os.environ["RATE"])
print(1 if rate >= 95 else 0)
PY
)

if [ "$threshold_pass" -eq 1 ]; then
  echo "✅ 满足 SC-002（≥95% 命中率）。"
else
  echo "⚠️ 低于 95% 命中率，请查看失败列表并确认 Redis/内存配置。"
  exit 3
fi
