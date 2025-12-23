#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${CLIENTTOKEN_BASE_URL:-http://127.0.0.1:7072}
API_TOKEN=${CLIENTTOKEN_API_TOKEN:-dev-clienttoken}
PROVIDER_CODE=${DOUYIN_PROVIDER_CODE:-byte_dance_douyin_clienttoken}
APP_CODE=${DOUYIN_APP_CODE:-douyin_service}
AUTH_MODE=${DOUYIN_AUTH_MODE:-default}
CONFIG_PATH=${CLIENTTOKEN_CONFIG_PATH:-config.yaml}
CURL_BIN=${CURL:-curl}

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required for JSON escaping" >&2
  exit 1
fi

usage() {
  cat <<'USAGE'
Usage:
  scripts/clienttoken-douyin.sh refresh [provider] [app] [mode] [config]
      # 刷新 DouYin client_token 并回显 TTL/缓存 key
  scripts/clienttoken-douyin.sh cache [show|clear] [provider] [app] [mode] [config]
      # 查看或清除 DouYin client_token 缓存
  scripts/clienttoken-douyin.sh call <action> [method] [query] [body|@file] [provider] [app] [mode] [config]
      # 通过 /client-token/call 触发 DouYin API

环境变量：
  CLIENTTOKEN_BASE_URL       默认 http://127.0.0.1:7072
  CLIENTTOKEN_API_TOKEN      默认 dev-clienttoken
  DOUYIN_PROVIDER_CODE       默认 byte_dance_douyin_clienttoken
  DOUYIN_APP_CODE            默认 douyin_service
  DOUYIN_AUTH_MODE           默认 default
  CLIENTTOKEN_CONFIG_PATH    默认 config.yaml
  CURL                       覆盖 curl 可执行文件
USAGE
}

pretty_print() {
  if command -v jq >/dev/null 2>&1; then
    jq '.'
  else
    cat
  fi
}

json_escape() {
  python3 - <<'PY' "$1"
import json, sys
print(json.dumps(sys.argv[1]))
PY
}

payload_base() {
  local provider=${1:-$PROVIDER_CODE}
  local app=${2:-$APP_CODE}
  local mode=${3:-$AUTH_MODE}
  local config=${4:-$CONFIG_PATH}
  cat <<JSON
{
  "provider_code": $(json_escape "${provider}"),
  "provider_app": $(json_escape "${app}"),
  "provider_auth_mode": $(json_escape "${mode}"),
  "config_path": $(json_escape "${config}")
}
JSON
}

post_json() {
  local path=$1
  local body=$2
  printf '%s' "$body" | "$CURL_BIN" --noproxy '*' -sS -X POST \
    -H "Authorization: Bearer ${API_TOKEN}" \
    -H 'Content-Type: application/json' \
    --data-binary @- \
    "${BASE_URL}${path}" | pretty_print
}

get_with_query() {
  local path=$1
  shift
  "$CURL_BIN" --noproxy '*' -sS -G \
    -H "Authorization: Bearer ${API_TOKEN}" \
    "$@" \
    "${BASE_URL}${path}" | pretty_print
}

cmd_refresh() {
  local body=$(payload_base "$@")
  post_json "/client-token/token" "$body"
}

cmd_cache() {
  local action=${1:-show}
  shift || true
  local provider=${1:-$PROVIDER_CODE}
  local app=${2:-$APP_CODE}
  local mode=${3:-$AUTH_MODE}
  local config=${4:-$CONFIG_PATH}

  case "$action" in
    show|list)
      get_with_query "/client-token/cache" \
        --data-urlencode "provider_code=${provider}" \
        --data-urlencode "provider_app=${app}" \
        --data-urlencode "provider_auth_mode=${mode}" \
        --data-urlencode "config_path=${config}"
      ;;
    clear|delete)
      local payload=$(payload_base "$provider" "$app" "$mode" "$config")
      "$CURL_BIN" --noproxy '*' -sS -X DELETE \
        -H "Authorization: Bearer ${API_TOKEN}" \
        -H 'Content-Type: application/json' \
        --data-binary "${payload}" \
        "${BASE_URL}/client-token/cache" | pretty_print
      ;;
    *)
      echo "Unknown cache action: ${action}" >&2
      usage
      exit 1
      ;;
  esac
}

cmd_call() {
  if [ $# -lt 1 ]; then
    usage
    exit 1
  fi
  local action=$1
  local method=${2:-POST}
  local query=${3:-}
  local bodyPayload=${4:-}
  local provider=${5:-$PROVIDER_CODE}
  local app=${6:-$APP_CODE}
  local mode=${7:-$AUTH_MODE}
  local config=${8:-$CONFIG_PATH}

  local body_field
  if [ -n "$bodyPayload" ] && [[ "$bodyPayload" == @* ]]; then
    local file=${bodyPayload#@}
    body_field=$(json_escape "$(cat "$file")")
  else
    body_field=$(json_escape "${bodyPayload}")
  fi

  local merged=$(cat <<JSON
{
  "provider_code": $(json_escape "${provider}"),
  "provider_app": $(json_escape "${app}"),
  "provider_auth_mode": $(json_escape "${mode}"),
  "config_path": $(json_escape "${config}"),
  "action": $(json_escape "${action}"),
  "method": $(json_escape "${method}"),
  "query": $(json_escape "${query}"),
  "body": ${body_field}
}
JSON
)
  post_json "/client-token/call" "$merged"
}

cmd=${1:-}
shift || true

case "$cmd" in
  refresh)
    cmd_refresh "$@"
    ;;
  cache)
    cmd_cache "$@"
    ;;
  call)
    cmd_call "$@"
    ;;
  ""|-h|--help|help)
    usage
    ;;
  *)
    echo "Unknown command: $cmd" >&2
    usage
    exit 1
    ;;
esac
