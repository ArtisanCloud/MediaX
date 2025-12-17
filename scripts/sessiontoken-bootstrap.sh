#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
TEMPLATE_PATH="${TEMPLATE_PATH:-$ROOT_DIR/config.example.yaml}"
CONFIG_PATH="${SESSIONTOKEN_CONFIG:-$ROOT_DIR/config.yaml}"

if [ $# -ge 1 ]; then
  CONFIG_PATH=$1
fi

if [ ! -f "$TEMPLATE_PATH" ]; then
  printf "[sessiontoken-bootstrap] 模板不存在: %s\\n" "$TEMPLATE_PATH" >&2
  exit 1
fi

if [ ! -f "$CONFIG_PATH" ]; then
  cp "$TEMPLATE_PATH" "$CONFIG_PATH"
  printf "[sessiontoken-bootstrap] 已根据模板生成 %s\\n" "$CONFIG_PATH"
else
  printf "[sessiontoken-bootstrap] 已存在 %s，跳过复制\\n" "$CONFIG_PATH"
fi

if grep -E '\${[A-Z0-9_]+}' "$CONFIG_PATH" >/dev/null; then
  printf "[sessiontoken-bootstrap] 检测到未替换的占位符，请根据模板填充后再启动 SessionToken 服务\\n" >&2
  exit 1
fi

printf "[sessiontoken-bootstrap] 配置就绪：%s\\n" "$CONFIG_PATH"
