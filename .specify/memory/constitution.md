<!--
Sync Impact Report
Version: 1.1.0 -> 1.1.1
Modified Principles: Provider Adapter Parity; Config-Layered Security; Token Lifecycle Discipline; Observability & Error Traceability; Testable Modularity & SessionToken Readiness
Added Sections: None
Removed Sections: None
Templates: ✅ Not required (translation only)
Follow-up TODOs: None
-->

# MediaX 宪章

## 核心原则

### I. 提供方适配器一致性（Provider Adapter Parity）
所有平台接入必须遵循统一的适配器目录：provider 目录使用英文小驼峰命名、每个产品都暴露 `core/client.go`，且只能通过 MediaX 工厂完成依赖注入。共享逻辑集中在 provider 级 `core/` 包中，产品差异放在 `accessTokenClient`/`clientTokenClient`，禁止出现随手定义的常量或结构体。

### II. 配置分层安全（Config-Layered Security）
所有凭证与开关都要通过 `pkg/client/config` 下的强类型结构体输入，继承 `ClientConfig`/`OAuthConfig` 并注明 yaml/json tag、默认值与 TTL/保留策略。配置文件只能引用环境变量或机密管理系统，不得写入明文 token、脚本 ID 等敏感信息；每次变更都需声明存储介质与逐出策略，以保障多租户隔离。

### III. Token 生命周期规范（Token Lifecycle Discipline）
AccessToken、ClientToken 与 SessionToken 操作必须复用 `kernel.BaseClient` 提供的 HTTP、重试与刷新检测能力。若需要重写（如抖音刷新逻辑），需在代码内注明原因，并在修改请求后 clone `http.Request` 以保留 Body。凭证获取统一委托给外部回调，默认 Redis 缓存可替换且必须可审计。

### IV. 可观测性与错误可追溯（Observability & Error Traceability）
所有 IO 都必须调用 `client.Logger.WithContext(ctx)` 并记录 `provider`、`api`、`tenant_uuid`、`account_id`、`flow_id` 与重试次数，敏感字段只保留前后缀。错误链使用 `fmt.Errorf("feature: %w", err)` 等形式保留上游信息；任何离开 SDK 的回调都必须附带含 timestamp/nonce 的 HMAC-SHA256 签名以阻断重放。

### V. 可测试模块化与 SessionToken 准备度（Testable Modularity & SessionToken Readiness）
导出的类型/方法需要 GoDoc 注释，每个纯逻辑模块在合入前必须具备 `*_test.go` 覆盖。SessionToken Flow 必须包含 `pending → authorizing → succeeded/failed` 的状态机测试，以及签名/脱敏辅助工具。规格与计划文档要明确列出将新增的测试包，方便评审基于可度量的覆盖率做 Gate。

## 架构与运行约束
- **仅允许 BaseClient HTTP**：所有外部请求必须经过 `kernel.BaseClient.HttpHelper`，严禁手动实例化 `http.Client`，以确保重试、签名与日志统一生效。
- **多租户上下文**：日志、缓存键与 Flow 存储必须携带 `tenant_uuid`、`account_id`；Flow 记录在完成后依然要保留审计 TTL。
- **SessionToken 五大支柱**：provider 配置需完整描述 `Service`、`Authenticator`、`Harvester`、`Callback`、`Network`，以便适配器调度代理池、验证码脚本与回调验签。
- **外发请求签名**：所有回调/外部钩子必须带有 HMAC-SHA256 签名与时间戳/随机数，缺失或验签失败需要立即拒绝。
- **可扩展的存储层**：默认使用 Redis，但接口必须保持实现无关；如需自定义，必须说明重试/退避策略并避免锁定具体存储行为。

## 配置与秘钥
- **集中式 schema**：所有 provider 字段放在 `pkg/client/config/<provider>.go`，并通过注释说明默认值及用途。
- **秘钥管理**：配置只能引用环境变量或机密管理平台，仓库文件禁止存放明文 token、cookie 或脚本内容。
- **TTL 约束**：Token/Session 相关配置必须声明存储介质、TTL 与逐出策略，避免缓存无限增长。
- **文档联动**：新增配置结构需要同步 README/规格中的示例，让运维能在 sandbox/prod 一致地注入参数。

## SessionToken 与 Flow 纪律
- **持久化 Flow**：每个 Flow 代表一次模拟登录，必须写入共享存储（推荐 Redis），并具备 TTL/索引以支撑审计与重试。
- **生命周期窗口**：Flow 完成后在合规窗口内可查询，长期 Session 由业务独立的表负责，过期后重新创建 Flow。
- **策略注入**：Authenticator/Harvester 策略通过配置提供入口 URL、脚本、UA、代理池等，外部脚本仅引用 ID/URL 而非嵌入代码。
- **遥测与脱敏**：Flow 日志与指标必须包含 `flow_id`、`provider_code`、tenant 标识，并在输出前完成脱敏。

## 开发流程与质量门禁
- Go 代码保持 `gofmt`，导出符号需有注释。包命名遵循 provider/产品模板（如 `pkg/client/google/youtube/accessTokenClient`）。
- 实施前需在计划中说明接口/配置变更，提交 PR 前必须本地运行 `go test ./...` 并确保日志脱敏。
- 重试/Token 刷新逻辑必须输出结构化日志，记录触发原因及前后尝试次数。
- 功能分支（`feat/<topic>`、`fix/<topic>`）至少一次评审，评审需对照本宪章及文档/规格的更新。
- 代码伴随文档：新增 provider 需要在 README/配置注释中描述默认值、存储假设及 BaseClient 覆写点。

## 治理
本宪章高于所有临时约定。任何修订都要通过 RFC 文档说明动机、迁移与验证计划。Pull Request 如需偏离或扩展既有行为，必须在描述中引用相关章节；违反适配器一致性、配置纪律、可观测性或测试要求的更改不得合并。运行期指引（日志字段、重试策略、Flow 存储）需与本宪章及 `docs/plan/session_token_client.md` 保持同步。

**Version**: 1.1.1 | **Ratified**: 2025-12-09 | **Last Amended**: 2025-12-09
