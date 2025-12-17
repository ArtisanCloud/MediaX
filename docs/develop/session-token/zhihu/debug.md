# Zhihu SessionToken Debug 指南

> 目标：在 MediaX SDK 仓库内独立启动 SessionToken 服务（:7070），创建/查询 Flow，调试 harvester 回调链路，定位 `metadata.session_token is empty` 等问题，无需依赖插件仓库。

## 1. 依赖与环境准备

1. Go 1.18+、Redis（默认 `127.0.0.1:6379`）。
2. 保证当前目录为 `/private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX`。
3. 统一环境变量（可写入 shell profile）：
   ```bash
   export POWERX_SESSION_TOKEN_BASE_URL="http://127.0.0.1:7070"
   export POWERX_SESSION_TOKEN_API_TOKEN="dev-session-token"
   export POWERX_SESSION_TOKEN_CALLBACK_URL="https://plugin.local/api/v1/admin/platforms/session-token/callback"
   export SESSIONTOKEN_REDIS_ADDR="127.0.0.1:6379"
   export SESSIONTOKEN_ZHIHU_API_VERSION="v4" # 可选，覆盖 config.yaml 里的 service.api_version
   ```
4. 如无 `config.yaml`，执行 `make sessiontoken-bootstrap`（会从 `config.example.yaml` 复制模板并提示缺失字段）。

## 2. 启动 SessionToken 服务

1. 确认 Redis 已启动；如需临时实例，可执行 `docker run --rm -p 6379:6379 redis:7-alpine`。
2. 在 MediaX 根目录运行：
   ```bash
   make sessiontoken
   # 或
   go run ./cmd/sessiontoken -config config.yaml
   ```
3. 日志中看到 `sessiontoken: server listening addr=:7070 api_token=dev-session-token` 即表示成功。

## 3. 内置调试页面（/debug）

1. 默认情况下，服务会挂载 `http://127.0.0.1:7070/debug` 页面（若设置 `SESSIONTOKEN_DISABLE_DEBUG_PAGE=1` 可关闭）。
2. 启动 `make sessiontoken` 后，浏览器打开上述地址即可使用“Provider → App”级联表单创建 Flow。页面特性：
   - **默认值**：首次进入会自动填充 `API Token` 与 `Callback URL`（默认 `http://127.0.0.1:7070/debug/callback`），后续保存在 `localStorage`。
   - **模板**：切换 Provider/App 或点击“应用模板”将覆盖 metadata、Token、Callback，并为 `state` 生成随机值，便于多 Provider 覆盖式调试。
   - **会话复用**：勾选“浏览器已登录则直接复用 Cookie”时，页面会在 metadata 中写入 `reuse_session=true`，SessionToken 服务会优先尝试使用同租户/账号最近一次成功的 Cookie；若缓存存在且未过期，则 Flow 会立刻 `succeeded`，无需再次手动登录。
   - **API 版本选择**：若 Provider 提供多个 API 版本（当前 Zhihu 支持 `v4`），会显示额外下拉框；创建 Flow 时，页面会自动将 `metadata.api_version` 写入请求体，并与 `SESSIONTOKEN_ZHIHU_API_VERSION`/`config.yaml` 保持一致，便于回溯。
   - **全局函数**：`window.createFlow`/`pollFlow`/`copyAuthorizeURL` 暴露到全局，方便在浏览器控制台或自动化脚本里重复这些动作。
   - **Mock 回调**：自带 `/debug/callback` 的 ring buffer（最近 20 条），支持刷新/清空按钮，也可直接 `POST` payload 进行重放。
3. Flow 操作流程：
   - 点击“创建 Flow”后，输出区会打印完整响应，并自动填充 Flow ID、authorize_url；若要更换账号，点击“复制 authorize_url” 在新的浏览器 profile 中登录即可。
   - 如果勾选了“复用 Cookie”且同租户/账号最近一次成功的 Flow 仍然有效，则本次 Flow 会在几百毫秒内直接 `succeeded` 并触发回调，输出中可看到复用的 session_token；若缓存不存在或已过期，则需要继续下一步人工登录。
   - **Metadata 等待**：未勾选复用时，Flow 也不会立即失败。Orchestrator 会在后台每 2 秒查询一次 Redis，最多等待 2 分钟，直到浏览器脚本把 `metadata.session_token` 写回。只有超时或脚本写入空值时才会触发 `metadata session_token not ready before timeout`，因此可以放心先登录再执行写入操作。
   - “查询 Flow”会读取最新状态，输出 `metadata.session_token` 与 `code/message` 便于定位。
4. 页面不会向服务端持久化任何口令，所有 Token/Callback 仅存储在浏览器 `localStorage`，请仅在本地环境使用。

### 3.1 Provider/App 模板设计

- `provider` 目前内置 Zhihu / RedBook / Google，均可扩展，只需在 `providerCatalog` 增加条目。
- 每个 provider 可以挂多个 `apps`，配置项包含：
  - `code`: 最终写入 `provider_app_code` 的值（如 `zhihu_article`、`redbook_web`、`youtube_web`）。
  - `label`: 页面展示文案。
  - `metadata`: 建议写入的 JSON 片段（例如 Zhihu Web 默认 `{"login_mode":"pc"}`，小红书扫码场景默认 `{"login_mode":"qr"}`）。
- 若内置列表满足不了场景，可切换到「自定义」，输入任意 `provider_app_code` 并覆盖 metadata。模板仍会复用 Provider 对应的默认 `token/callback`。
- “应用模板”按钮区分「强制覆盖」与「智能保留」：若手动改过 metadata，输入框 `data-autofill=0` 后再点击模板不会被重置；若希望彻底替换，则点击一次“应用模板”按钮即可。

### 3.2 浏览器模拟/登录策略

为了在 MediaX SDK 范畴内复现插件的行为，推荐以下步骤：

1. **授权 URL 注入**：通过调试页创建 Flow，复制 `authorize_url`，在同一台机器打开新 Chrome/Edge 标签页，并使用干净的 profile（或 `chrome --user-data-dir=/tmp/mediax-debug-profile`）以免旧 Cookie 干扰。
2. **脚本执行策略**：
   - 当前授权页不会自动回写 Cookie，需要在登录成功后手动操作：打开 DevTools Console，执行后文脚本读取 `document.cookie` 并通过 `POST /debug/flows/<flow_id>/metadata` 回写。
   - 推荐使用 Playwright/Puppeteer 驱动的 CLI（`tools/sessiontoken-browser/`）。首次运行前需在仓库根目录执行 `pnpm install`（或 `npm install`）及 `npx playwright install chromium` 安装依赖与浏览器。随后运行 `pnpm sessiontoken:browser --flow <id> [--base http://127.0.0.1:7070] [--api-token ...]`，CLI 会自动拉起 Chromium、导航至 authorize_url、等待你完成登录，然后捕获 Cookie 并调用 `/debug/flows/<flow_id>/metadata` 写回。脚本结束时会在终端输出 `/debug/callback` 的结果，避免手动复制 JSON。默认会读取 `SESSIONTOKEN_API_TOKEN`/`POWERX_SESSION_TOKEN_API_TOKEN` 环境变量，也可通过 `--api-token` 显式指定。
3. **Cookie/Storage 监控**：Zhihu 场景需要抓取浏览器实际写入的 Cookie，例如 `SESSIONID`、`JOID`、`osd`、`q_c1`、`d_c0`、`unlock_ticket`、`z_c0` 等。调试流程建议：
   - 打开 DevTools → Application → Cookies，确认上述字段存在；留意 `q_c1`、`d_c0`、`z_c0` 中的 `|timestamp|`、`2|1:0|...` 结构需完整保留。
   - 创建一个虚拟字段 `session_token`，其值就是 `document.cookie` 的原始串（示例：`SESSIONID=7zKo...; JOID=...; ...; z_c0=2|1:0|...`），后续所有封装接口通过 `X-SessionToken: <原始串>` 回传。
   - 可在 Console 粘贴以下脚本快速生成 metadata：
     ```js
     (() => {
       const cookie = document.cookie;
       const pick = name => (new RegExp(`${name}=([^;]+)`)).exec(cookie)?.[1] || '';
       return {
         session_token: cookie,
         cookie_sessionid: pick('SESSIONID'),
         cookie_joid: pick('JOID'),
         cookie_osd: pick('osd'),
         cookie_q_c1: pick('q_c1'),
         cookie_d_c0: pick('d_c0'),
         cookie_unlock_ticket: pick('unlock_ticket'),
         cookie_z_c0: pick('z_c0'),
         credentials_note: navigator.userAgent,
         credentials_expires_hint: ''
       };
     })();
     ```
   - **更新 Flow**：使用浏览器 Console 直接调用，或运行 Playwright CLI 自动写回：
     ```js
     fetch('http://127.0.0.1:7070/debug/flows/<flow_id>/metadata', {
       method: 'POST',
       headers: {'Content-Type': 'application/json'},
       body: JSON.stringify(/* 上面的 JSON */)
     })
     ```
     Playwright CLI 内部也会发起同样的请求；服务端收到后会立刻触发 harvester 并在 `/debug/callback` 中展示成功 payload；只有在极端场景下才需要回退到直接操作 Redis。
4. **API 鉴权验证**：
   - 在浏览器执行 `fetch('https://www.zhihu.com/api/v4/members/me', {credentials:'include'})` 等自检请求，确认返回 200，证明 Cookie 可用。
   - 也可使用 `curl` 加 `-H "Cookie: $session_token"` 的方式验证抓到的 Token 是否能访问文章 API。
5. **可视化/回放**：若流程失败，可在 `/debug/callback` 每条记录里看到 `flow_id`、`provider_code`、payload 以及 header（含签名），方便复制重放；还可以直接调用 `POST /debug/callback` 注入自定义 payload 验证插件兼容性。

在长期规划中，我们会补充一个“浏览器模拟器”面板，允许：

- 选择 Provider/App 后自动打开相应的 Playwright 脚本（例如 `scripts/simulators/zhihu_web.ts`），并通过 WebSocket 把 Cookie/Storage 回写到 Flow metadata。
- 针对需要手机环境的场景（如小红书 `login_mode=qr`），提供二维码扫描/协议登录的沙箱，确保 harvester 行为与真实终端保持一致。
- 支付宝/Google 这类需要 OAuth 的平台，会在调试页内展示“授权线程”状态，指引测试者在外部窗口完成 MFA。

上述策略可确保你在 SDK 仓库内就能模拟“浏览器启动 → 登录 → Cookie 抓取 → 回调 → Flow 结束”的整套流程，而无需依赖插件端。

### 3.3 调试脚本

当你想在终端里快速查询 Flow 或调用封装好的 Zhihu API，可使用 `scripts/sessiontoken-debug.sh`：

```bash
scripts/sessiontoken-debug.sh flow stf_xxx
scripts/sessiontoken-debug.sh followings "$SESSION_TOKEN"
scripts/sessiontoken-debug.sh channels "$SESSION_TOKEN" zhihu_column_id 10 0
scripts/sessiontoken-debug.sh sanity "$SESSION_TOKEN"
```

脚本会自动读取 `POWERX_SESSION_TOKEN_BASE_URL` 与 `SESSIONTOKEN_API_TOKEN`/`POWERX_SESSION_TOKEN_API_TOKEN`，并在系统安装 `jq` 时自动美化输出。通过 `/debug` 创建 Flow 后，将回调中的 `metadata.session_token` 复制给脚本，即可在本地验证 API 是否仍可访问，或触发 `ZH_COOKIE_EXPIRED` 回调。

> 若想复用现有 Flow，可在本地直接操作 Redis：`redis-cli --raw keys 'sessionToken:flow:*'` 列出所有 key，然后 `redis-cli --raw GET "sessionToken:flow:<id>" | jq '.'` 查看对应 metadata。确认 Flow 仍在有效期后，把 ID 粘到 `/debug` → “Flow ID” 输入框，再点“从 Flow 填充 SessionToken”即可调试。

### 3.4 API 调试面板（Beta）

为了减少“抓 Cookie → 切终端 → curl”的上下文切换，调试页新增了“API 调试”组件，默认与 Provider 选择联动，并具备以下能力：

1. **接口目录**：自动列出当前 Provider + 版本（例如 Zhihu v4）的可调试 API（followings/channels/articles/sanity...），展示 `METHOD PATH` + 说明，方便直接点选。切换顶部 API 版本下拉时，接口列表会同步刷新，保证与 `pkg/client/zhihu/web/sessionTokenClient/<version>` 的实现保持一致。
2. **路径/Query 解析**：输入 `channel_id=xxx&limit=10` 之类的键值对时，调试页会自动把 `{channel_id}` 写入路径、并将剩余字段拼成 query string，无需手写 URL。
3. **SessionToken 自动注入**：点击“从 Flow 填充 SessionToken”会调用 `GET /session-token/flows/<flow_id>`，优先读取 `metadata.session_token`；如为空，会用 `cookie_sessionid/cookie_joid/...` 拼出临时串写入 `X-SessionToken`。
4. **一键发起请求**：填写 Body（JSON，可选）后点击“发送请求”，页面会自动携带 `Authorization: Bearer <API Token>` 与 `X-SessionToken`，直接请求当前 SessionToken 服务的 `/zhihu/v1/*` handler，并在下方 `pre` 区块显示 HTTP 状态码与格式化响应。
5. **失效复现**：若接口需要 SessionToken 而输入框为空，页面会提示补充；因此可以先让 Flow 成功，再通过该面板重放 API、DIY 过期 Token、观察 `sessiontoken_callback`。

通过该面板可以实现“创建 Flow → Playwright 抓 Cookie → 一键拉接口 → Mock 失效”的闭环，全程无需离开 `/debug`。

## 4. 监控日志与指标

- 一般排查使用：
  ```bash
  tail -f logs/sessiontoken-info.log | rg 'sessiontoken_(metric|callback)'
  ```
- 若需要完整堆栈/错误：
  ```bash
  tail -f logs/sessiontoken-error.log
  ```
- 典型日志：
  - `sessiontoken_metric: action=create_flow ... status=pending ...`
  - `sessiontoken_callback: failed ... metadata.session_token is empty`

## 5. 手动创建/查询 Flow（不依赖插件）

1. 创建 Flow：
   ```bash
   FLOW_RESP=$(curl --noproxy "*" -s -X POST "$POWERX_SESSION_TOKEN_BASE_URL/session-token/flows" \
     -H "Authorization: Bearer $POWERX_SESSION_TOKEN_API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "provider_code": "zhihu",
           "provider_app_code": "zhihu_article",
           "account_id": "acct_debug",
           "tenant_uuid": "tenant_debug",
           "state": "ui-debug-001",
           "callback_url": "'$POWERX_SESSION_TOKEN_CALLBACK_URL'",
           "metadata": {"login_mode": "pc"}
         }')
   echo "$FLOW_RESP" | jq '.'
   FLOW_ID=$(echo "$FLOW_RESP" | jq -r '.flow.flow_id')
   ```
2. 查询 Flow：
   ```bash
   curl --noproxy "*" -s -H "Authorization: Bearer $POWERX_SESSION_TOKEN_API_TOKEN" \
     "$POWERX_SESSION_TOKEN_BASE_URL/session-token/flows/$FLOW_ID" | jq '.'
   ```
3. 观察日志里 `sessiontoken_metric: action=get_flow ...` 与 `sessiontoken_callback` 是否出现。

## 6. 调试 Harvester：定位 `metadata.session_token is empty`

### 5.1 检查 Flow metadata（Redis）

1. 找到 Redis key：`sessionToken:flow:<flow_id>`。
2. 导出 JSON：
   ```bash
   redis-cli --raw GET "sessionToken:flow:$FLOW_ID" | jq '.' > /tmp/flow_$FLOW_ID.json
   ```
3. 查看 `metadata` 字段是否含有 `session_token`（完整 `document.cookie` 串）以及 `cookie_sessionid/cookie_joid/.../cookie_z_c0` 等拆分字段。若为空，表示浏览器脚本尚未把 Cookie 写回 Flow。
4. 若需要短期修改 metadata 进行实验，可编辑 `/tmp/flow_*.json` 并写回：
   ```bash
   cat /tmp/flow_$FLOW_ID.json | redis-cli -x SET "sessionToken:flow:$FLOW_ID"
   redis-cli EXPIRE "sessionToken:flow:$FLOW_ID" 3600
   ```
   > **注意**：写回需要保持 JSON 结构正确，并手动设置 TTL（默认留存 6h）。

### 5.2 在 SDK 内模拟 Harvester

1. 运行单测，验证 harvester 对 metadata 的期望：
   ```bash
   go test ./server/zhihu/sessionToken/harvester -run TestHarvester -v
   ```
2. 使用 FlowOrchestrator 单测验证闭环：
   ```bash
   go test ./pkg/client/sessionToken -run FlowOrchestrator -v
   ```
3. 若要单步调试，可用 Delve：
   ```bash
   dlv test ./server/zhihu/sessionToken/harvester -- -test.run TestHarvester
   ```

### 5.3 浏览器脚本联调（本地）

1. 创建 Flow 后，复制响应中的 `authorize_url` 到浏览器，按脚本提示完成登录并等待“凭证已写入”提示。
2. 在 DevTools → Application → Cookies 中确认 `SESSIONID`、`JOID`、`osd`、`q_c1`、`d_c0`、`unlock_ticket`、`z_c0` 等字段存在，并记录 `document.cookie` 的完整内容。
3. 使用上一节提供的脚本或 Postman，将完整 `session_token` 与拆分字段写入 Flow metadata。调试页可新增“导入 Cookie JSON”按钮（待实现），目前可在 Console 中执行 `fetch('/debug/flows/<flow_id>/metadata', ...)`。
4. 如需复用旧会话，可保存 metadata 并设置 `session_token_reuse_existing_session=true`（可选扩展位），避免重复扫码。

## 7. 验证回调

1. 若你临时将 `POWERX_SESSION_TOKEN_CALLBACK_URL` 指向 httpbin（如 `https://httpbin.org/post`），可在响应体看到 MediaX 发送的 payload；headers 将包含 `X-MediaX-Signature`、`X-MediaX-Timestamp` 等。
2. 要手动校验签名，可保存请求体为 `payload.json`，再执行：
   ```bash
   openssl dgst -sha256 -mac hmac -macopt hexkey:<callback_secret_hex> payload.json
   ```

## 8. 常见问题速查

| 症状 | 排查思路 |
| --- | --- |
| `metadata.session_token is empty` | 检查浏览器脚本是否执行完、metadata 是否写入、`WatchCookies/WatchHeaders` 是否覆盖新的 cookie 名称。 |
| `sessiontoken_callback: failed ... error=EOF` | Callback URL 不通或证书问题，使用 `curl -v $POWERX_SESSION_TOKEN_CALLBACK_URL` 先验证。 |
| Flow 一直 pending | 关注 orchestrator 日志，确认 harvester 是否 panic；必要时清理 Redis 旧 flow 重新创建。 |
| 回调签名不通过 | 确认 `callback.callback_secret` 与插件一致，或在插件里同步更新 secret。 |

## 9. 清理

- 停止服务：`Ctrl+C`（orchestrator 会优雅关闭）。
- 清理 Redis：
  ```bash
  redis-cli KEYS 'sessionToken:flow:*' | xargs -r redis-cli DEL
  redis-cli KEYS 'sessionToken:flow:state:*' | xargs -r redis-cli DEL
  ```
- 删除临时日志：`rm logs/sessiontoken-*.log`（可选）。

## 10. 附录：metadata 快速查看脚本

```bash
#!/usr/bin/env bash
set -euo pipefail
FLOW_ID=${1:-}
if [ -z "$FLOW_ID" ]; then
  echo "Usage: scripts/debug-flow.sh <flow_id>" && exit 1
fi
redis-cli --raw GET "sessionToken:flow:$FLOW_ID" | jq '.metadata'
```

放入 `scripts/debug-flow.sh` 后 `chmod +x`，即可快速打印 metadata。后续可扩展为写入 mock `session_token`、重放回调等功能。
