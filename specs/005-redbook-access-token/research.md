# Research Summary

## Decision 1: Reuse现有 cmd/accesstoken Server 作为聚光调试入口
- **Rationale**: Google/Bilibili 已证明该服务可统一配置解析、OAuth 回调、Redis Flow 缓存与 `/debug` UI，无需另起服务。
- **Alternatives Considered**:
  - *独立的 redbook 调试服务器*：增加维护成本，与既有调试体验不一致 → 否决。
  - *Playground-only 脚本*：缺少 Flow 管理与 UI，不满足文档要求 → 否决。

## Decision 2: 选取账户余额 API 作为首个聚光调试用例
- **Rationale**: `JuGuangAccountClient.GetAccountBalance` 依赖最少、不涉及复杂 payload，可验证 AccessToken 生命周期与日志脱敏。
- **Alternatives Considered**:
  - *笔记发布 API*：依赖素材上传，流程复杂，不适合作为首个验证接口。
  - *数据报表 API*：需要分页与时间窗口，易受数据延迟影响。

## Decision 3: 配置模板沿用 redbook_juguang_config 并补充 scope/oauth 字段
- **Rationale**: 当前 `config.example.yaml` 已有聚光节点，只需补充 OAuth 字段即可让 `/debug` 读取；保持与 Google/Bilibili 的统一层级。
- **Alternatives Considered**:
  - *新增独立配置文件*：违背 Config-Layered Security，不利于部署；
  - *通过环境变量单独注入 JSON*：难以在文档中示范和版本管理。
