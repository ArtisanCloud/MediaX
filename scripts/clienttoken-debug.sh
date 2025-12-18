#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${CLIENTTOKEN_BASE_URL:-http://127.0.0.1:7072}
API_TOKEN=${CLIENTTOKEN_API_TOKEN:-dev-clienttoken}
PROVIDER_CODE=${CLIENTTOKEN_PROVIDER_CODE:-wechat}
APP_CODE=${CLIENTTOKEN_APP_CODE:-official_account}
AUTH_MODE=${CLIENTTOKEN_AUTH_MODE:-default}
CONFIG_PATH=${CLIENTTOKEN_CONFIG_PATH:-config.yaml}
CURL_BIN=${CURL:-curl}

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required for JSON escaping" >&2
  exit 1
fi

usage() {
  cat <<'USAGE'
Usage:
  scripts/clienttoken-debug.sh token [provider] [app] [mode] [config]
      # 主动刷新 access_token 并查看 TTL
  scripts/clienttoken-debug.sh cache [provider] [app] [mode] [config]
      # 查看缓存命中与剩余 TTL
  scripts/clienttoken-debug.sh call <action> [method] [query] [body] [provider] [app] [mode] [config]
      # 调用公众号 API（action 形如 cgi-bin/getcallbackip）
  scripts/clienttoken-debug.sh validate <signature> <timestamp> <nonce> [echostr]
      # 触发 /client-token/message/validate，确认签名是否正确
  scripts/clienttoken-debug.sh callbacks [list|clear]
      # 查看或清空回调日志（/client-token/message/callbacks）

环境变量：
  CLIENTTOKEN_BASE_URL        默认 http://127.0.0.1:7072
  CLIENTTOKEN_API_TOKEN       默认 dev-clienttoken
  CLIENTTOKEN_PROVIDER_CODE   默认 wechat
  CLIENTTOKEN_APP_CODE        默认 official_account
  CLIENTTOKEN_AUTH_MODE       默认 default
  CLIENTTOKEN_CONFIG_PATH     默认 config.yaml
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

cmd_token() {
  local body=$(payload_base "$@")
  post_json "/client-token/token" "$body"
}

cmd_cache() {
  local provider=${1:-$PROVIDER_CODE}
  local app=${2:-$APP_CODE}
  local mode=${3:-$AUTH_MODE}
  local config=${4:-$CONFIG_PATH}
  get_with_query "/client-token/cache" \
    --data-urlencode "provider_code=${provider}" \
    --data-urlencode "provider_app=${app}" \
    --data-urlencode "provider_auth_mode=${mode}" \
    --data-urlencode "config_path=${config}"
}

cmd_call() {
  if [ $# -lt 1 ]; then
    usage
    exit 1
  fi
  local action=$1
  local method=${2:-GET}
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

cmd_validate() {
  if [ $# -lt 3 ]; then
    usage
    exit 1
  fi
  local signature=$1
  local timestamp=$2
  local nonce=$3
  local echostr=${4:-}
  local body=$(cat <<JSON
{
  "provider_code": $(json_escape "${PROVIDER_CODE}"),
  "provider_app": $(json_escape "${APP_CODE}"),
  "provider_auth_mode": $(json_escape "${AUTH_MODE}"),
  "config_path": $(json_escape "${CONFIG_PATH}"),
  "signature": $(json_escape "${signature}"),
  "timestamp": $(json_escape "${timestamp}"),
  "nonce": $(json_escape "${nonce}"),
  "echostr": $(json_escape "${echostr}")
}
JSON
)
  post_json "/client-token/message/validate" "$body"
}

cmd_callbacks() {
  local sub=${1:-list}
  case "$sub" in
    list)
      get_with_query "/client-token/message/callbacks"
      ;;
    clear)
      "$CURL_BIN" --noproxy '*' -sS -X DELETE \
        -H "Authorization: Bearer ${API_TOKEN}" \
        "${BASE_URL}/client-token/message/callbacks" | pretty_print
      ;;
    *)
      echo "Unknown callbacks action: ${sub}" >&2
      exit 1
      ;;
  esac
}

cmd=${1:-}
shift || true

case "$cmd" in
  token)
    cmd_token "$@"
    ;;
  cache)
    cmd_cache "$@"
    ;;
  call)
    cmd_call "$@"
    ;;
  validate)
    cmd_validate "$@"
    ;;
  callbacks)
    cmd_callbacks "$@"
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
