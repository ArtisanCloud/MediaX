#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${POWERX_SESSION_TOKEN_BASE_URL:-http://127.0.0.1:7070}
API_TOKEN=${SESSIONTOKEN_API_TOKEN:-${POWERX_SESSION_TOKEN_API_TOKEN:-dev-session-token}}
CURL_BIN=${CURL:-curl}

usage() {
  cat <<'EOF'
Usage:
  scripts/sessiontoken-debug.sh flow <flow_id>
      # 查看 Flow 状态/metadata
  scripts/sessiontoken-debug.sh followings "<session_token>"
      # 调用 GET /zhihu/v1/me/followings
  scripts/sessiontoken-debug.sh channels "<session_token>" <channel_id> [limit] [offset]
      # 调用 GET /zhihu/v1/channels/{channel_id}/articles
  scripts/sessiontoken-debug.sh article "<session_token>" <article_id>
      # 调用 GET /zhihu/v1/articles/{id}
  scripts/sessiontoken-debug.sh sanity "<session_token>"
      # 调用 POST /zhihu/v1/sanity/check

  所有命令都会默认读取环境变量：
    POWERX_SESSION_TOKEN_BASE_URL (默认 http://127.0.0.1:7070)
    SESSIONTOKEN_API_TOKEN 或 POWERX_SESSION_TOKEN_API_TOKEN
EOF
}

ensure_args() {
  local need=$1
  shift
  if [ "$#" -lt "$need" ]; then
    usage
    exit 1
  fi
}

pretty_print() {
  if command -v jq >/dev/null 2>&1; then
    jq '.'
  else
    cat
  fi
}

call_flow() {
  ensure_args 1 "$@"
  local flow_id=$1
  "$CURL_BIN" --noproxy '*' -sS -H "Authorization: Bearer ${API_TOKEN}" \
    "${BASE_URL}/session-token/flows/${flow_id}" | pretty_print
}

call_api() {
  local method=$1
  local path=$2
  local session_token=$3
  local body=${4:-}
  if [ -n "$body" ]; then
    printf '%s' "$body" | "$CURL_BIN" --noproxy '*' -sS -X "$method" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer ${API_TOKEN}" \
      -H "X-SessionToken: ${session_token}" \
      --data-binary @- \
      "${BASE_URL}${path}" | pretty_print
  else
    "$CURL_BIN" --noproxy '*' -sS -X "$method" \
      -H "Authorization: Bearer ${API_TOKEN}" \
      -H "X-SessionToken: ${session_token}" \
      "${BASE_URL}${path}" | pretty_print
  fi
}

cmd=${1:-}
shift || true

case "$cmd" in
  flow)
    call_flow "$@"
    ;;
  followings)
    ensure_args 1 "$@"
    call_api GET "/zhihu/v1/me/followings" "$1"
    ;;
  channels)
    ensure_args 2 "$@"
    session_token=$1
    channel_id=$2
    limit=${3:-20}
    offset=${4:-0}
    call_api GET "/zhihu/v1/channels/${channel_id}/articles?limit=${limit}&offset=${offset}" "$session_token"
    ;;
  article)
    ensure_args 2 "$@"
    call_api GET "/zhihu/v1/articles/${2}" "$1"
    ;;
  sanity)
    ensure_args 1 "$@"
    call_api POST "/zhihu/v1/sanity/check" "$1"
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
