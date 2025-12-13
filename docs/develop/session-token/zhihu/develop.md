# Zhihu SessionToken PRD（Browser Cookie Orchestrator）

> 目标：为外部应用提供一套可重复获取知乎登录 Cookie、封装常用接口并可回调告警的模拟登录闭环，以便订阅频道、抓取仅登录可见的文章，以及模拟发表内容。

## 1. 背景与目标

1. 当前 session-token 调试页只能手动观察 Flow；缺少“抓取 Cookie → 存储 → 接口复用”的完整链路。
2. 外部应用需要定期拉取**订阅频道列表**、**频道最新文章**、**登录后可见但公开的文章**，以及**发布自有文章**。
3. 需要标准化的回调协议：包含 Cookie、有效期、失效告警（`message + code`），并与 MediaX 现有 SessionToken Flow 管理规范兼容。

**成果要求**

- 浏览器模拟登录策略可复用：生成 Flow → 打开 `authorize_url` → 采集 Cookie → 写回 metadata → 调用回调。
- Cookie 生命周期清晰：外部保存 `session_token`，请求任何封装 API 必须带回该 token。
- 当 Cookie 失效或命中风控时，SessionToken 服务需要回调外部并给出明确 `code/message`，同时标记 Flow `failed`。
- 本地需要 Playbook：如何调试登录、查看 metadata、用封装 API 联调。

## 2. 用户故事与范围

| 编号 | 用户 | 场景 | 价值 |
| --- | --- | --- | --- |
| US-01 | 平台运营 | 定期刷新知乎账号会话，确保抓站脚本持续可用 | 减少手工登录成本 |
| US-02 | 订阅服务 | 获取“我关注的频道/话题”列表 | 驱动内容推荐 |
| US-03 | 内容采集 | 拉取频道/关注/推荐 feed、仅登录可见文章 | 保障抓取覆盖 |
| US-04 | 内容分发 | 使用同账号发布文章 | 支撑“代发”能力 |

**不在范围**：知乎官方 OAuth、移动端扫码（本期仅 web-PC）、付费专栏、视频上传。

## 3. 方案概述

### 3.1 浏览器模拟登录策略

1. `make sessiontoken` 启动服务 → `/debug` 创建 Flow（Provider=Zhihu，App=zhihu_article）。
2. 调试页提供“打开模拟器”按钮（或提示手动复制 `authorize_url`）；浏览器完成登录。
3. 登陆成功后，通过注入脚本读取 `document.cookie`，并**完整原样**拼接为一个虚拟字段 `session_token`（示例：`SESSIONID=7zKo...; JOID=VV0W...; osd=VF4T...; q_c1=f5693a25...; d_c0=kpYUP2...; unlock_ticket=AEBC...; z_c0=2|1:0|...`）。注意浏览器里并没有名为 `session_token` 的 Cookie，这只是我们在 metadata 与回调中定义的聚合字符串，便于外部系统直接存取。调用任何接口时必须在请求头 `X-SessionToken` 回传这一整串值。
4. 同时拆分关键 Cookie 写入 metadata 中的 `cookie_<name>`，留给内部风控/补发用途。
5. 记录 `captured_at`、`credentials_note`（UA/设备信息）、`credentials_expires_hint`（知乎响应头或默认 7d）。
4. metadata 写回 Redis 后，Harvester 读取 `session_token` 并执行 `CompleteFlowSuccess`，触发回调。

> **Metadata 等待机制**：Flow 创建后会一直处于 `pending`，Orchestrator 每 2 秒轮询一次 Redis，直到 `metadata.session_token` 不为空才真正启动 Harvester。最长等待 2 分钟；若超时会在日志/回调中看到 `metadata session_token not ready before timeout`，提示开发者重新登录或排查脚本。这样即使浏览器窗口已登录，也可以直接复用现有 Cookie，而不会在 Flow 创建后立刻失败。

#### 3.1.1 Cookie 字段定义

| Zehihu Cookie | metadata key | 说明 | 备注 |
| --- | --- | --- | --- |
| `SESSIONID` | `cookie_sessionid` | Web 登录核心会话 ID，用于绝大多数 API | 失效后 401 |
| `JOID` / `osd` | `cookie_joid` / `cookie_osd` | 登录设备标识，常与 `SESSIONID` 搭配验证 | 两个字段成对出现 |
| `q_c1` | `cookie_q_c1` | 访问轨迹 ID，有助于模拟浏览器行为 | 需要携带 `|timestamp|` 部分 |
| `d_c0` | `cookie_d_c0` | 设备持久 ID，部分接口校验 | 需完整值（含 `=` 与 `|timestamp`） |
| `unlock_ticket` | `cookie_unlock_ticket` | 登录凭证刷新所需 | 可能用于重发内容 |
| `z_c0` | `cookie_z_c0` | 旧版本脚本依赖的 Token，部分接口仍校验 | 需要完整 `2|1:0|...` 格式 |
| 其他（`z_c0`、`tgw_l7_route` 等） | `cookie_<lower_name>` | 如果页面出现新的 cookie，保留扩展能力 | watch 列表允许自定义 |

`session_token` = `SESSIONID + JOID + osd + q_c1 + d_c0 + unlock_ticket + ...` 的原始拼接字符串，顺序无需强制，但必须在内部/回调保持一致。

> 额外 metadata：
> - `reuse_session`: `"true"` 时表示允许直接复用缓存的会话；默认为 `"false"`。

### 3.2 Cookie 失效与重试

- SessionToken 服务在以下事件将 Flow 标记失败并回调：
  1. Harvester 未获取到 `session_token` => `code=ZH_COOKIE_EMPTY`
  2. 调用知乎接口返回 401/403 => `code=ZH_COOKIE_EXPIRED`
  3. 风控提示（接口返回 `errcode`/`r=true` 等）=> `code=ZH_COOKIE_RISK`
- 回调 payload 中要携带 **完整 `session_token` 字符串** 以及拆分后的 `cookie_*` 字段，便于外部持久化；失败场景还需要附带最后一次访问的 API。
- 失败示例：
  ```json
  {
    "flow_id": "...",
    "status": "failed",
    "code": "ZH_COOKIE_EXPIRED",
    "message": "知乎 Cookie 已失效，请重新扫码登录",
    "metadata": {
      "session_token": "SESSIONID=7zKo...; JOID=VV0W...; ...",
      "cookie_sessionid": "7zKo...",
      "cookie_joid": "VV0W...",
      "cookie_osd": "VF4T...",
      "cookie_q_c1": "f5693a25|1765503157000|...",
      "cookie_d_c0": "kpYUP2UfhB...|1765503068",
      "cookie_unlock_ticket": "AEBCfDFt3AoX...",
      "last_failed_api": "GET https://www.zhihu.com/api/v4/articles/xxx"
    }
  }
  ```
- 外部系统收到 `code` 后决定是否重新触发 Flow。

### 3.3 API 版本管理

- Zhihu Web API 的具体实现按版本拆分在 `pkg/client/zhihu/web/sessionTokenClient/v4/*` 目录中，后续若知乎升级 `api/v5`，仅需新增 `v5` 目录并在入口注册。
- 通过 `zhihu_config.sessionToken.service.api_version`（或环境变量 `SESSIONTOKEN_ZHIHU_API_VERSION`）即可在不改代码的情况下切换版本，默认值为 `v4`。
- `/debug` 页面在创建 Flow 时会自动在 metadata 中注入 `api_version` 字段，便于 SDK/插件侧回溯；`sessiontoken_api` 日志也会输出 `version=v4`，方便排查。

### 3.4 会话复用策略

- Flow metadata 新增 `reuse_session` 字段（布尔字符串 `true/false`）。当其为 `true` 时，Orchestrator 在触发 Harvester 之前会查询缓存（`sessionToken:reuse:<hash>`）是否存在同租户/账号/Provider App 最近一次成功的 `CredentialPayload`。
- `CompleteFlowSuccess` 会自动将最新的 SessionToken/Cookies 写入上述缓存，并延长有效期（沿用 Flow TTL + `DefaultAuditTTL`）。后续同租户/账号请求勾选“复用”即可直接 `CompleteFlowSuccess`，无需再次打开浏览器登录。
- 若缓存不存在或 `reuse_session=false`，则继续走原有“提示用户登录→Harvester 抓取 Cookie”流程，不影响兼容性。

## 4. API 封装需求

### 4.1 公共约束

- **请求头**：必须传 `X-SessionToken`（即之前回调给外部的 `session_token`），服务端将其还原为 Cookie。
- **鉴权**：如果 SessionToken 服务检测到 token 缺失或过期，直接返回 `401` 并回调 `ZH_COOKIE_EXPIRED`。
- **Proxy/UA**：沿用 `config.yaml` 中 `zhihu_config.sessionToken.network`、`authenticator.default_user_agent`。

### 4.2 需要封装的接口

| 接口 | 目的 | 备注 |
| --- | --- | --- |
| `GET /zhihu/v1/me/followings` | 关注的专栏/话题/圆桌（分页） | 调用知乎 `https://www.zhihu.com/api/v4/people/{uid}/following-columns` 等 |
| `GET /zhihu/v1/channels/{channel_id}/articles` | 频道最新文章 | 支持 `limit/offset` |
| `GET /zhihu/v1/articles/{id}` | 登录后可见文章详情 | 直接代理 `api/v4/articles/{id}` |
| `POST /zhihu/v1/articles` | 发布文章 | body 包含标题、内容、可见性 |
| `POST /zhihu/v1/sanity/check` | 检查当前 Cookie 是否仍可访问 `/api/v4/me` | 作为心跳 |

**接口返回结构示例**

```json
{
  "request_id": "req-123",
  "data": {...},          // 与知乎原始数据保持一致或轻度裁剪
  "meta": {
    "source": "zhihu",
    "upstream": "https://www.zhihu.com/api/v4/...",
    "captured_at": "2025-12-12T01:38:00Z"
  }
}
```

**错误码补充**

| HTTP | code | message | 说明 |
| --- | --- | --- | --- |
| 400 | `ZH_BAD_REQUEST` | 参数不合法 | 转发知乎 4xx 时的默认 |
| 401 | `ZH_COOKIE_EXPIRED` | Cookie 已失效 | 同时触发回调 |
| 403 | `ZH_RISK_BLOCK` | 触发风控 | 记录 `risk_detail` |
| 5xx | `ZH_UPSTREAM_ERROR` | 知乎返回 5xx | 重试策略遵循 `retry_backoff` |

## 5. 回调协议

沿用 `session-token` Callback，但新增字段：

```json
{
  "flow_id": "stf_xxx",
  "status": "succeeded/failed",
  "provider_code": "zhihu",
  "provider_app_code": "zhihu_article",
  "tenant_uuid": "...",
  "state": "...",
  "metadata": {
    "session_token": "SESSIONID=7zKo...; JOID=VV0W...; osd=VF4T...; q_c1=f5693a25|1765503157000|...; d_c0=kpYUP2...|1765503068; unlock_ticket=AEBC...",
    "cookie_sessionid": "7zKo...",
    "cookie_joid": "VV0W...",
    "cookie_osd": "VF4T...",
    "cookie_q_c1": "f5693a25|1765503157000|...",
    "cookie_d_c0": "kpYUP2...|1765503068",
    "cookie_unlock_ticket": "AEBC...",
    "credentials_note": "Chrome 123 / Mac",
    "credentials_expires_hint": "2025-12-20T00:00:00Z"
  },
  "credentials": {
    "session_token": "SESSIONID=7zKo...; JOID=VV0W...; ...",            // 与 metadata.session_token 保持一致
    "cookies": [
      {"name": "SESSIONID", "value": "7zKo..."},
      {"name": "JOID", "value": "VV0W..."},
      {"name": "osd", "value": "VF4T..."},
      {"name": "q_c1", "value": "f5693a25|1765503157000|..."},
      {"name": "d_c0", "value": "kpYUP2...|1765503068"},
      {"name": "unlock_ticket", "value": "AEBC..."}
    ],
    "headers": null
  },
  "code": null,
  "message": null
}
```

当 Flow 失败或后续调用检测到过期时，`credentials` 为空，`code/message` 按上一节约定返回。

## 6. 本地调试流程

1. **准备环境**：按 `docs/develop/session-token/zhihu/debug.md` 启动 Redis + `make sessiontoken`。若需使用自动化 CLI，请在仓库根目录执行 `pnpm install`（或 `npm install`）与 `npx playwright install chromium` 以安装依赖/浏览器。
2. **创建 Flow**：在 `/debug` 选择 Zhihu 模板，点“创建 Flow”→ 复制 `authorize_url`，或直接运行 `pnpm sessiontoken:browser --flow <flow_id>` 让 CLI 自动拉起浏览器；在弹出的窗口完成知乎登录。
3. **采集 Cookie**：若未使用 CLI，可在浏览器 DevTools 控制台运行以下脚本回写：
   ```js
   (async () => {
     const cookie = document.cookie;
     await fetch('http://127.0.0.1:7070/debug/flows/<flow_id>/metadata', {
       method: 'POST',
       headers: {'Content-Type':'application/json'},
       body: JSON.stringify({
         session_token: cookie,
         cookie_sessionid: /SESSIONID=([^;]+)/.exec(cookie)?.[1] || '',
         cookie_joid: /JOID=([^;]+)/.exec(cookie)?.[1] || '',
         cookie_osd: /osd=([^;]+)/.exec(cookie)?.[1] || '',
         cookie_q_c1: /q_c1=([^;]+)/.exec(cookie)?.[1] || '',
         cookie_d_c0: /d_c0=([^;]+)/.exec(cookie)?.[1] || '',
         cookie_unlock_ticket: /unlock_ticket=([^;]+)/.exec(cookie)?.[1] || '',
         credentials_note: navigator.userAgent
       })
     });
   })();
   ```
   Playwright CLI 会在登录成功后自动将 Cookie 写入 `/debug/flows/<flow_id>/metadata`，终端打印“metadata updated”以及最新的 `/debug/callback` 记录。
4. **验证 Flow**：调试页“查询 Flow”应看到 `status=succeeded`；`最近回调` 出现从 `/debug/callback` 发送的 payload。
5. **封装接口调试**：启动 `go run ./cmd/sessiontoken` 后，通过 `curl` 调用 `GET /zhihu/v1/articles/{id}`，同时在 Header 带 `X-SessionToken`。可使用 `make sessiontoken-proxy`（待实现）让请求直连知乎并打印上游响应。
6. **失效模拟**：手动修改 metadata，填入错误 Cookie，调用 `GET /zhihu/v1/sanity/check` 验证服务是否返回 `ZH_COOKIE_EXPIRED` 并触发失败回调。

## 7. 交付与验收

1. **代码交付**：
   - Harvester 配置支持新的 `watch_cookies`；默认模板保存回调字段。
   - 新增 API 模块（推荐 `pkg/client/zhihu/web/sessionTokenClient/...`）。
   - 回调扩展 `code/message` 字段。
2. **文档**：
   - 本 PRD（当前文档）。
   - 更新 `docs/develop/session-token/zhihu/debug.md`，补充新的调试脚本。
   - 提供外部调用示例与错误码对照。
3. **验收标准**：
   - 本地可成功运行一次 Flow，回调含 Cookie。
   - 能用同一 `session_token` 调用封装接口获取关注频道/文章。
   - 将 Cookie 改为无效值后，接口返回 401 且收到 `ZH_COOKIE_EXPIRED` 回调。

## 8. API 版本扩展 Playbook

1. **代码结构**：在 `pkg/client/zhihu/web/sessionTokenClient/` 下创建新的版本目录（如 `v5/`），复制 `v4` 的 handler 骨架，仅改动上游路径/响应整形逻辑；公共工具函数仍放在 `v4` 目录或抽到 `internal`，避免重复。
2. **入口注册**：在 `pkg/client/zhihu/web/sessionTokenClient/client.go` 的 `buildRouter` 中新增 `case "v5": return v5.NewClient(...)`；如果需要临时灰度，可扩展 `resolveAPIVersion`，支持以 Flow metadata/Provider App 决定版本。
3. **配置/环境变量**：`zhihu_config.sessionToken.service.api_version` 用于设置默认版本；`SESSIONTOKEN_ZHIHU_API_VERSION` 可在运行期覆盖（含 `/debug` 页面）。更新 `config.example.yaml`、`.env.example` 与 README/quickstart 的映射，提示新版本可选值。
4. **调试工具**：`/debug` 页面模板需同步新增版本选项，并在创建 Flow 时把 `metadata.api_version` 写入，以确保插件/日志可以回溯；脚本 `sessiontoken-debug.sh` 也可支持 `--version`，方便 curl 时附带。
5. **测试/验收**：针对新增版本补充 httptest 覆盖基础成功/401/403/5xx 情形，并跑一遍“创建 Flow → 登录 → `/debug/callback` → API 调用 → 失效回调”流程确认兼容；必要时在日志中输出 `version=...` 字段方便监控。

## 9. 风险与TODO

- **Cookie 结构变动**：知乎可能更换字段，需要监控调试日志并快速更新 `watch_cookies`。
- **高频调用封禁**：对文章/频道接口需加节流与 IP 池支持，可复用 `zhihu_config.sessionToken.network`。
- **发布接口风控**：文章发布需 CSRF Token、`x-xsrf-token` 等，后续迭代再评估实现。
- **调试页安全**：`/debug/flows/<id>/metadata` 写入接口仅用于本地，可通过 `SESSIONTOKEN_DISABLE_DEBUG_PAGE` 阻止生产暴露。

--- 

> 负责人：SessionToken 团队  
> 截止目标：2025-12-31 完成联调上线
