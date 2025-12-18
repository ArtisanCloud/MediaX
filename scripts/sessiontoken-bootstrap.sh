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

session_block="$(python3 - "$CONFIG_PATH" <<'PY'
import sys
from pathlib import Path

path = Path(sys.argv[1])
lines = path.read_text().splitlines()
start = None
for idx, line in enumerate(lines):
    if line.startswith("session_token_providers:"):
        start = idx
        break

if start is None:
    sys.exit(1)

end = len(lines)
for idx in range(start + 1, len(lines)):
    text = lines[idx]
    stripped = text.strip()
    if stripped == "":
        continue
    if not text.startswith(" "):
        end = idx
        break

print("\n".join(lines[start:end]))
PY
)" || {
  printf "[sessiontoken-bootstrap] 未在 %s 中找到 session_token_providers 段落，请确认配置文件\\n" "$CONFIG_PATH" >&2
  exit 1
}

placeholder_lines="$(printf '%s\\n' "$session_block" | grep -En '\$\{[A-Z0-9_:-]+\}' || true)"
if [ -n "$placeholder_lines" ]; then
  printf "[sessiontoken-bootstrap] SessionToken 配置仍含占位符，请补齐以下字段：\\n%s\\n" "$placeholder_lines" >&2
  exit 1
fi

printf "[sessiontoken-bootstrap] SessionToken 配置就绪：%s\\n" "$CONFIG_PATH"
