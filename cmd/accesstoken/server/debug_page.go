package main

import (
	"html/template"
	"net/http"
)

type debugPageData struct {
	APIToken        string
	ProvidersJSON   template.JS
	DefaultProvider string
	DefaultApp      string
	DefaultMode     string
	DefaultCallback string
	StorageBackend  string
	ListenAddr      string
	FlowTTLSeconds  int
	PublicListen    bool
}

var debugPageTemplate = template.Must(template.New("accesstoken_debug").Parse(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <title>AccessToken 调试台</title>
  <style>
    body { font-family: -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif; margin: 0; padding: 24px; background: #f4f5f7; color:#1a1a1a; }
    h1, h2 { margin-bottom: 8px; }
    section { background:#fff; margin-bottom:20px; padding:20px; border-radius:12px; box-shadow:0 1px 3px rgba(0,0,0,0.08); }
    label { display:block; font-size:13px; font-weight:600; margin-top:12px; }
    input, select, textarea { width:100%; padding:8px; margin-top:4px; box-sizing:border-box; font-family: inherit; border:1px solid #c8d1dc; border-radius:6px; }
    textarea { min-height:140px; font-family: 'SFMono-Regular',Consolas,'Liberation Mono',Menlo,monospace; }
    button { padding:8px 16px; border:none; border-radius:6px; background:#0052cc; color:#fff; cursor:pointer; }
    button.secondary { background:#e0e0e0; color:#1a1a1a; }
    pre { background:#0b2545; color:#d6f0ff; padding:12px; border-radius:8px; overflow:auto; min-height:120px; }
    .row { display:flex; gap:16px; flex-wrap:wrap; }
    .row > div { flex:1; min-width:200px; }
    .actions { margin-top:12px; display:flex; gap:12px; flex-wrap:wrap; }
    .note { font-size:12px; color:#555; margin-top:4px; }
    .note strong { color:#1a1a1a; }
    table { width:100%; border-collapse:collapse; font-size:13px; }
    th, td { border-bottom:1px solid #f0f0f0; text-align:left; padding:6px 4px; vertical-align:top; }
    code { background:#f0f0f0; padding:2px 4px; border-radius:4px; }
    .flash { margin-bottom:16px; padding:10px 14px; border-radius:8px; background:#e6f4ff; color:#084298; display:none; }
    .flash.show { display:block; }
    .flash.error { background:#fdecea; color:#a61b1b; }
    .inline-input { display:flex; gap:8px; align-items:center; }
    .inline-input input { flex:1; }
    .inline-input button { flex:0 0 auto; white-space:nowrap; }
    .checkbox-inline { display:flex; align-items:center; gap:6px; margin-top:8px; }
    .top-bar { display:flex; justify-content:space-between; align-items:center; background:#0b2545; color:#fff; padding:12px 20px; border-radius:10px; margin-bottom:20px; }
    .top-bar .meta { font-size:13px; color:#c0d3ea; margin-left:8px; }
    .top-actions { display:flex; align-items:center; gap:10px; }
    .token-chip { background:#09213a; border:1px solid #1f4c7a; padding:6px 10px; border-radius:999px; font-size:13px; }
    .banner { padding:10px 14px; border-radius:8px; margin-bottom:18px; font-size:13px; display:none; }
    .banner.warning { background:#fff4e5; color:#8a5100; border:1px solid #f5c97a; }
    .banner.danger { background:#fdecea; color:#a61b1b; border:1px solid #f5a3a3; }
    .provider-card { margin-top:16px; border:1px solid #dbe3f0; border-radius:10px; padding:12px 16px; background:#f8fbff; }
    .provider-card-header { display:flex; justify-content:space-between; align-items:flex-start; gap:12px; }
    .provider-card-title { font-size:16px; font-weight:600; }
    .provider-card-subtitle { font-size:12px; color:#4a6b9c; margin-top:2px; }
    .provider-card-tags { display:flex; gap:8px; flex-wrap:wrap; }
    .provider-card-tags span { background:#e3e9f4; color:#2c3e68; padding:4px 8px; border-radius:999px; font-size:12px; }
    .provider-card-body { display:flex; gap:12px; flex-wrap:wrap; margin:12px 0 8px; }
    .provider-card-body div { flex:1; min-width:160px; font-size:12px; color:#4a5568; }
    .provider-card-body code { display:block; margin-top:4px; background:#fff; border:1px solid #dbe3f0; padding:4px 6px; border-radius:6px; }
    .provider-card pre { background:#0b2545; color:#d6f0ff; min-height:80px; }
  </style>
</head>
<body data-api-token="{{.APIToken}}" data-default-provider="{{.DefaultProvider}}" data-default-app="{{.DefaultApp}}" data-default-mode="{{.DefaultMode}}" data-default-callback="{{.DefaultCallback}}" data-storage-backend="{{.StorageBackend}}" data-flow-ttl="{{.FlowTTLSeconds}}" data-listen-addr="{{.ListenAddr}}" data-public-listen="{{if .PublicListen}}true{{else}}false{{end}}">
  <div class="top-bar">
    <div>
      <strong>AccessToken 调试台</strong>
      <span class="meta">监听 {{.ListenAddr}} · Flow TTL {{.FlowTTLSeconds}} 秒 · 存储 {{.StorageBackend}}</span>
    </div>
    <div class="top-actions">
      <span id="tokenChip" class="token-chip">Token: -</span>
      <button type="button" class="secondary" onclick="openTokenDialog()">更新 API Token</button>
    </div>
  </div>
  <div id="storageBanner" class="banner warning"></div>
  <div id="publicBanner" class="banner danger"></div>
  <h1>AccessToken 调试台</h1>
  <p>该页面复用了 SessionToken 调试模式，可在浏览器内选择 Provider/App，直接调用当前服务的 AccessToken API，并观察 <code>/debug/callback</code>。建议首次访问时通过 URL 附带 <code>?api_token=&lt;值&gt;</code> 或点击右上角按钮写入浏览器存储，后续请求会自动携带。</p>
  <div id="flashMessage" class="flash"></div>

  <section>
    <h2>1. 基础配置</h2>
    <div class="row">
      <div>
        <label>Provider</label>
        <select id="providerSelect"></select>
      </div>
      <div>
        <label>Provider App</label>
        <select id="providerAppSelect"></select>
      </div>
      <div>
        <label>授权模式</label>
        <select id="providerModeSelect"></select>
      </div>
      <div>
        <label>API 版本</label>
        <input id="apiVersionInput" type="text" value="v1" />
      </div>
    </div>
    <div class="row">
      <div>
        <label>Config Path</label>
        <input id="configPathInput" type="text" />
      </div>
      <div>
        <label>OAuth Key</label>
        <input id="oauthKeyInput" type="text" readonly />
      </div>
    </div>
    <div class="provider-card" id="providerCard">
      <div class="provider-card-header">
        <div>
          <div class="provider-card-title" id="providerCardTitle">-</div>
          <div class="provider-card-subtitle" id="providerCardVersion">API v-</div>
        </div>
        <div class="provider-card-tags">
          <span id="providerCardModeTag">模式：-</span>
          <span id="providerCardOAuthKeyTag">OAuth Key：-</span>
        </div>
      </div>
      <div class="provider-card-body">
        <div>
          Provider Code
          <code id="providerCardProviderCode">-</code>
        </div>
        <div>
          App Code
          <code id="providerCardAppCode">-</code>
        </div>
        <div>
          Config Path
          <code id="providerCardConfigPath">config.yaml</code>
        </div>
      </div>
      <label>默认模板</label>
      <pre id="providerCardTemplate">{}</pre>
    </div>
    <label>Callback URL</label>
    <input id="callbackInput" type="text" />
    <div class="actions">
      <button type="button" onclick="applyTemplate()">同步模板</button>
      <button type="button" class="secondary" onclick="startOAuth()">发起授权</button>
      <button type="button" class="secondary" onclick="loadTokens(true)">刷新授权记录</button>
    </div>
    <p class="note">默认回调为 <code>{{.DefaultCallback}}</code>，如实际监听地址不同，可在此修改并复制给第三方。</p>
  </section>

  <section>
    <h2>2. AccessToken 解析</h2>
    <p class="note">用于确认 Token 来源（payload/env/config）、TTL 与 http_debug 等运行时信息，排查 401/配置差异时先点击解析再调用 API。</p>
    <textarea id="tokenPayload"></textarea>
    <div class="actions">
      <button type="button" onclick="invoke('/accesstoken/token', 'tokenPayload', 'tokenOutput')">解析 AccessToken</button>
    </div>
    <pre id="tokenOutput">// 输出 AccessToken 来源、脱敏 token 及 TTL</pre>
  </section>

  <section>
    <h2>3. 授权记录</h2>
    <p class="note"><strong>使用指南：</strong>点击“发起授权”完成 OAuth 后，回调提示会包含 <code>flow_id</code>，表格也会展示最近记录。即使服务重启，也可以在下方输入 Flow ID（或从 Redis key <code>accesstoken:oauth:flow:&lt;id&gt;</code> / <code>accesstoken:oauth:&lt;provider&gt;:&lt;app&gt;:&lt;mode&gt;</code> 中获取）后点击“加载 Flow ID”重新载入。</p>
    <p class="note">示例命令：<code>redis-cli GET accesstoken:oauth:flow:oauth-abc123</code>（返回 value 即主 key，再 <code>GET</code> 该 key 可看到完整授权记录）。</p>
    <div class="row">
      <div>
        <label>通过 Flow ID 回填</label>
        <div class="inline-input">
          <input id="flowIdInput" type="text" placeholder="oauth-xxxxxx" />
          <button type="button" class="secondary" onclick="loadTokenByFlow()">加载 Flow ID</button>
          <button type="button" class="secondary" onclick="loadFlowIndexList()">列出 Redis Flow</button>
        </div>
        <p class="note">Flow ID 可在授权成功提示、回调日志或 Redis key 中获取。</p>
      </div>
    </div>
    <table id="tokenTable">
      <thead>
        <tr><th>Flow ID</th><th>状态</th><th>有效期</th><th>操作</th></tr>
      </thead>
      <tbody></tbody>
    </table>
  </section>

  <section>
    <h2>4. API 调用</h2>
    <p class="note"><strong>像 SessionToken 调试页一样</strong>，先在下方表单选择 API/参数，页面会自动同步 JSON；若需要自定义参数，可直接编辑文本框。</p>
    <div class="row">
      <div>
        <label>API Action</label>
        <select id="apiActionSelect">
          <option value="videos.list">videos.list（按 ID/榜单）</option>
          <option value="search.list">search.list（关键词/频道）</option>
          <option value="playlists.list">playlists.list（频道/我的列表）</option>
        </select>
      </div>
      <div>
        <label>Part</label>
        <input id="callPartInput" type="text" placeholder="snippet,statistics" />
      </div>
      <div>
        <label>IDs / Chart</label>
        <input id="callIdsInput" type="text" placeholder="dQw4w9WgXcQ" />
      </div>
    </div>
    <div class="row">
      <div>
        <label>Query</label>
        <input id="callQueryInput" type="text" placeholder="MediaX demo" />
      </div>
      <div>
        <label>Channel ID</label>
        <input id="callChannelInput" type="text" placeholder="UC_x5XG1..." />
      </div>
      <div>
        <label>Max Results (1-50)</label>
        <input id="callMaxResultsInput" type="number" min="1" max="50" value="5" />
      </div>
    </div>
    <div class="row">
      <div>
        <div class="checkbox-inline">
          <input id="callMineCheckbox" type="checkbox" />
          <label for="callMineCheckbox">Playlists Mine</label>
        </div>
        <div class="checkbox-inline">
          <input id="callSearchMineCheckbox" type="checkbox" />
          <label for="callSearchMineCheckbox">Search Mine</label>
        </div>
      </div>
    </div>
    <textarea id="callPayload"></textarea>
    <div class="actions">
      <button type="button" onclick="executeAPICall()">执行 API</button>
      <button type="button" class="secondary" onclick="loadCallbacks()">刷新回调记录</button>
      <button type="button" class="secondary" onclick="clearCallbacks()">清空回调</button>
    </div>
    <pre id="callOutput">// API 响应与错误输出</pre>
  </section>

  <section>
    <h2>5. OAuth 回调日志</h2>
    <p class="note">将 Google/OAuth 回调临时指向 <code>/debug/callback</code>，即可在此查看最近 50 条记录。</p>
    <table id="callbackTable">
      <thead>
        <tr><th>时间 (UTC)</th><th>方法</th><th>Query</th><th>Body</th></tr>
      </thead>
      <tbody></tbody>
    </table>
  </section>

  <script>
    const TOKEN_STORAGE_KEY = 'accesstoken_api_token';
    const defaultProvider = document.body.dataset.defaultProvider || '';
    const defaultApp = document.body.dataset.defaultApp || '';
    const defaultMode = document.body.dataset.defaultMode || '';
    const storageBackend = document.body.dataset.storageBackend || 'memory';
    const isPublicListen = (document.body.dataset.publicListen || '').toLowerCase() === 'true';
    const flowTTLSeconds = parseInt(document.body.dataset.flowTtl || '0', 10) || 0;
    const listenAddr = document.body.dataset.listenAddr || '';
    const providers = {{.ProvidersJSON}};
    const state = {
      provider: defaultProvider || (providers[0]?.code || ''),
      app: defaultApp || (providers[0]?.apps?.[0]?.code || ''),
      mode: defaultMode || (providers[0]?.apps?.[0]?.modes?.[0]?.key || ''),
      flowCursor: ''
    };
    let apiToken = '';

    const providerSelect = document.getElementById('providerSelect');
    const appSelect = document.getElementById('providerAppSelect');
    const modeSelect = document.getElementById('providerModeSelect');
    const configPathInput = document.getElementById('configPathInput');
    const oauthKeyInput = document.getElementById('oauthKeyInput');
    const callbackInput = document.getElementById('callbackInput');
    const apiVersionInput = document.getElementById('apiVersionInput');
    const flashBox = document.getElementById('flashMessage');
    const flowIdInput = document.getElementById('flowIdInput');
    const apiActionSelect = document.getElementById('apiActionSelect');
    const callPartInput = document.getElementById('callPartInput');
    const callIdsInput = document.getElementById('callIdsInput');
    const callQueryInput = document.getElementById('callQueryInput');
    const callChannelInput = document.getElementById('callChannelInput');
    const callMaxResultsInput = document.getElementById('callMaxResultsInput');
    const callMineCheckbox = document.getElementById('callMineCheckbox');
    const callSearchMineCheckbox = document.getElementById('callSearchMineCheckbox');
    const storageBanner = document.getElementById('storageBanner');
    const publicBanner = document.getElementById('publicBanner');
    const tokenChip = document.getElementById('tokenChip');
    const providerCard = document.getElementById('providerCard');
    const providerCardTitle = document.getElementById('providerCardTitle');
    const providerCardVersion = document.getElementById('providerCardVersion');
    const providerCardModeTag = document.getElementById('providerCardModeTag');
    const providerCardOAuthKeyTag = document.getElementById('providerCardOAuthKeyTag');
    const providerCardProviderCode = document.getElementById('providerCardProviderCode');
    const providerCardAppCode = document.getElementById('providerCardAppCode');
    const providerCardConfigPath = document.getElementById('providerCardConfigPath');
    const providerCardTemplate = document.getElementById('providerCardTemplate');
    let tokenCache = [];
    let flashTimer = null;

    function maskToken(token) {
      token = (token || '').trim();
      if (!token) return '-';
      if (token.length <= 4) return '***';
      return token.slice(0, 2) + '***' + token.slice(-2);
    }

    function isRedbookProvider() {
      const currentApp = getCurrentApp();
      const providerCode = (currentApp?.provider_code || state.provider || '').toLowerCase();
      return providerCode === 'redbook_juguang';
    }

    function executeAPICall() {
      invoke('/accesstoken/call', 'callPayload', 'callOutput');
    }

    function updateTokenChip() {
      if (!tokenChip) return;
      tokenChip.textContent = 'Token: ' + (apiToken ? maskToken(apiToken) : '未设置');
    }

    function initAPIToken() {
      const stored = (window.localStorage.getItem(TOKEN_STORAGE_KEY) || '').trim();
      const params = new URLSearchParams(window.location.search);
      const queryToken = (params.get('api_token') || '').trim();
      if (queryToken) {
        apiToken = queryToken;
        window.localStorage.setItem(TOKEN_STORAGE_KEY, apiToken);
        params.delete('api_token');
        const next = params.toString();
        const nextURL = window.location.pathname + (next ? '?' + next : '') + window.location.hash;
        window.history.replaceState({}, '', nextURL);
      } else if (stored) {
        apiToken = stored;
      } else {
        apiToken = (document.body.dataset.apiToken || '').trim();
      }
      updateTokenChip();
    }

    function openTokenDialog() {
      const next = prompt('请输入与 ACCESSTOKEN_API_TOKEN 一致的 Token', apiToken);
      if (next === null) {
        return;
      }
      apiToken = next.trim();
      if (apiToken) {
        window.localStorage.setItem(TOKEN_STORAGE_KEY, apiToken);
        showFlash('API Token 已更新，后续请求会自动携带。', false);
      } else {
        window.localStorage.removeItem(TOKEN_STORAGE_KEY);
        showFlash('已清空 API Token，请重新设置后再发起请求。', true);
      }
      updateTokenChip();
    }

    function ensureAPIToken() {
      if (!apiToken) {
        throw new Error('缺少 API Token：请点击右上角“更新 API Token”输入，或在 URL 附带 ?api_token=...');
      }
      return apiToken;
    }

    function authHeaders() {
      return { 'Authorization': 'Bearer ' + ensureAPIToken() };
    }

    function authorizedFetch(url, init = {}) {
      try {
        const headers = Object.assign({}, init.headers || {}, authHeaders());
        return fetch(url, Object.assign({}, init, { headers }));
      } catch (err) {
        return Promise.reject(err);
      }
    }

    function setupEnvBanners() {
      if (storageBackend !== 'redis' && storageBanner) {
        storageBanner.innerHTML = 'Flow 数据当前存储在内存中（storage_backend=memory），服务重启将清空。请配置 <code>ACCESSTOKEN_REDIS_*</code> 以启用持久化。' +
          (flowTTLSeconds ? '<br/>Flow TTL 默认 ' + flowTTLSeconds + ' 秒，可通过 <code>ACCESSTOKEN_FLOW_TTL_SECONDS</code> 覆盖。' : '');
        storageBanner.style.display = 'block';
      }
      if (isPublicListen && publicBanner) {
        const addrLabel = listenAddr || '(unknown)';
        publicBanner.innerHTML = '监听地址 <code>' + addrLabel + '</code> 对公网开放，请确保仅在受控网络使用或设置 <code>ACCESSTOKEN_LISTEN_ADDR=127.0.0.1:7071</code>。';
        publicBanner.style.display = 'block';
      }
    }

    const apiActionPresets = {
      'videos.list': {
        part: 'snippet,statistics',
        ids: '',
        chart: 'mostPopular',
        query: '',
        channel_id: '',
        max_results: 5,
        mine: false,
        search_mine: false
      },
      'search.list': {
        part: 'snippet',
        ids: '',
        chart: '',
        query: 'MediaX demo',
        channel_id: '',
        max_results: 5,
        mine: false,
        search_mine: false
      },
      'playlists.list': {
        part: 'snippet',
        ids: '',
        chart: '',
        query: '',
        channel_id: 'UC_x5XG1OV2P6uZZ5FSM9Ttw',
        max_results: 5,
        mine: false,
        search_mine: false
      }
    };

    if (configPathInput && providerCard) {
      configPathInput.addEventListener('input', () => renderProviderCard());
    }

    function getCurrentProvider() {
      return providers.find(p => p.code === state.provider);
    }

    function getCurrentApp() {
      const provider = getCurrentProvider();
      if (!provider) return null;
      return provider.apps.find(a => a.code === state.app) || provider.apps[0] || null;
    }

    function getCurrentMode() {
      const app = getCurrentApp();
      if (!app || !app.modes) return null;
      return app.modes.find(m => m.key === state.mode) || app.modes[0] || null;
    }

    function populateProviderSelect() {
      providerSelect.innerHTML = '';
      providers.forEach((p) => {
        const option = document.createElement('option');
        option.value = p.code;
        option.textContent = p.name;
        providerSelect.appendChild(option);
      });
      providerSelect.value = state.provider;
      providerSelect.onchange = () => {
        state.provider = providerSelect.value;
        state.app = '';
        populateAppSelect();
      };
      populateAppSelect();
    }

    function populateAppSelect() {
      const provider = getCurrentProvider();
      const apps = provider?.apps || [];
      appSelect.innerHTML = '';
      apps.forEach((app) => {
        const option = document.createElement('option');
        option.value = app.code;
        option.textContent = app.name || app.code;
        appSelect.appendChild(option);
      });
      if (!apps.find((a) => a.code === state.app)) {
        state.app = apps[0]?.code || '';
      }
      appSelect.value = state.app;
      appSelect.disabled = apps.length <= 1;
      appSelect.onchange = () => {
        state.app = appSelect.value;
        populateModeSelect();
      };
      populateModeSelect();
    }

    function populateModeSelect() {
      const app = getCurrentApp();
      const modes = app?.modes || [];
      modeSelect.innerHTML = '';
      modes.forEach((mode) => {
        const option = document.createElement('option');
        option.value = mode.key;
        option.textContent = mode.label || mode.key;
        modeSelect.appendChild(option);
      });
      if (!modes.find((m) => m.key === state.mode)) {
        state.mode = modes[0]?.key || '';
      }
      modeSelect.value = state.mode;
      modeSelect.disabled = modes.length <= 1;
      modeSelect.onchange = () => {
        state.mode = modeSelect.value;
        updateFormFromApp();
      };
      updateFormFromApp();
    }

    function updateFormFromApp() {
      const app = getCurrentApp();
      const mode = getCurrentMode();
      configPathInput.value = app?.config_path || '';
      oauthKeyInput.value = mode?.oauth_key || app?.oauth_key || '';
      apiVersionInput.value = app?.api_version || apiVersionInput.value || 'v1';
      applyTemplate();
      state.flowCursor = '';
      loadTokens(false);
    }

    function selectionSnapshot() {
      const app = getCurrentApp();
      const mode = getCurrentMode();
      const configPath = configPathInput.value.trim();
      const providerCode = app?.provider_code || state.provider || '';
      const appCode = app?.app_key || app?.code || state.app || '';
      const modeKey = state.mode || mode?.key || '';
      return { app, mode, configPath, providerCode, appCode, modeKey };
    }

    function renderProviderCard(preview) {
      if (!providerCard) return;
      const snapshot = selectionSnapshot();
      const provider = getCurrentProvider();
      const { app, mode, providerCode, appCode, configPath, modeKey } = snapshot;
      const providerName = provider?.name || provider?.code || '-';
      const appName = app?.name || app?.code || '-';
      providerCardTitle.textContent = providerName && appName ? providerName + ' / ' + appName : providerName || appName || '-';
      providerCardVersion.textContent = 'API ' + (app?.api_version || 'v1');
      providerCardModeTag.textContent = '模式：' + (mode?.label || mode?.key || '-');
      const oauthKey = (mode?.oauth_key || app?.oauth_key || '').trim();
      providerCardOAuthKeyTag.textContent = 'OAuth Key：' + (oauthKey || '未配置');
      providerCardProviderCode.textContent = providerCode || '-';
      providerCardAppCode.textContent = app?.code || appCode || '-';
      providerCardConfigPath.textContent = configPath || app?.config_path || '';
      const templateData = preview || {
        provider_code: providerCode,
        provider_app: appCode,
        provider_auth_mode: modeKey,
        config_path: configPath,
        api_version: app?.api_version || '',
        oauth_key: oauthKey
      };
      if (providerCardTemplate) {
        providerCardTemplate.textContent = JSON.stringify(templateData, null, 2);
      }
    }

    function applyTemplate() {
      const snapshot = selectionSnapshot();
      const { providerCode, appCode, modeKey, configPath, app, mode } = snapshot;
      const token = {
        provider_code: providerCode,
        provider_app: appCode,
        provider_auth_mode: modeKey,
        config_path: configPath,
        access_token: "",
        access_token_ttl: 3600
      };
      const call = {
        provider_code: providerCode,
        provider_app: appCode,
        provider_auth_mode: modeKey,
        config_path: configPath,
        action: "videos.list",
        part: "snippet,statistics",
        ids: "",
        max_results: 5
      };
      document.getElementById('tokenPayload').value = JSON.stringify(token, null, 2);
      setCallFormValues(call);
      syncCallPayloadFromForm();
      renderProviderCard({
        provider_code: providerCode,
        provider_app: appCode,
        provider_auth_mode: modeKey,
        config_path: configPath,
        api_version: app?.api_version || '',
        oauth_key: (mode?.oauth_key || app?.oauth_key || '').trim()
      });
      showFlash('已根据当前 Provider/App 同步 JSON 模板。', false);
    }

    async function startOAuth() {
      const configPath = configPathInput.value.trim();
      const currentApp = getCurrentApp();
      const providerCode = currentApp?.provider_code || state.provider || '';
      const appCode = currentApp?.code || state.app || '';
      const modeKey = state.mode || getCurrentMode()?.key || '';
      if (!providerCode || !appCode) {
        alert('请先选择 Provider 与 App');
        return;
      }
      try {
        const payload = {
          provider_code: providerCode,
          provider_app: appCode,
          provider_auth_mode: modeKey,
          config_path: configPath
        };
        const res = await authorizedFetch('/accesstoken/oauth/start', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        const data = await res.json();
        if (data.auth_url) {
          window.open(data.auth_url, '_blank', 'noopener,noreferrer');
        } else {
          showFlash('接口未返回授权地址', true);
        }
      } catch (err) {
        showFlash('发起授权失败: ' + err.message, true);
      }
    }

    async function loadTokens(force) {
      const currentApp = getCurrentApp();
      const providerCode = currentApp?.provider_code || state.provider || '';
      const appCode = currentApp?.code || state.app || '';
      const modeKey = state.mode || getCurrentMode()?.key || '';
      if (!providerCode) {
        tokenCache = [];
        renderTokenTable();
        return;
      }
      try {
        const params = new URLSearchParams({
          provider_code: providerCode,
          provider_app: appCode,
          provider_auth_mode: modeKey
        });
        const res = await authorizedFetch('/accesstoken/flows?' + params.toString());
        if (!res.ok) {
          if (force) {
            const text = await res.text();
            showFlash('加载授权记录失败: ' + text, true);
          }
          return;
        }
        const data = await res.json();
        tokenCache = data.flows || [];
        state.flowCursor = data.next_cursor || '';
        renderTokenTable();
        if (force && tokenCache.length === 0) {
          showFlash('未找到授权记录，请先完成 OAuth 或使用 Flow ID 回填。', true);
        }
      } catch (err) {
        if (force) {
          showFlash('加载授权记录失败: ' + err.message, true);
        }
      }
    }

    function renderTokenTable() {
      const tbody = document.querySelector('#tokenTable tbody');
      tbody.innerHTML = '';
      if (!tokenCache || tokenCache.length === 0) {
        const tr = document.createElement('tr');
        const td = document.createElement('td');
        td.colSpan = 4;
        td.textContent = '暂无授权记录，请点击“发起授权”或使用 Flow ID 回填。';
        tr.appendChild(td);
        tbody.appendChild(tr);
        return;
      }
      tokenCache.forEach((token) => {
        const tr = document.createElement('tr');
        const expireText = token.expire_at ? new Date(token.expire_at).toLocaleString() : '-';
        const statusText = (token.status || token.source || '-');
        const maskedAccount = token.masked_account ? ' · ' + token.masked_account : '';
        tr.innerHTML = '<td>' + (token.flow_id || '-') + '</td>' +
                       '<td>' + statusText + maskedAccount + '</td>' +
                       '<td>' + expireText + '</td>';
        const actionTd = document.createElement('td');
        const btn = document.createElement('button');
        btn.textContent = '填充';
        btn.type = 'button';
        btn.onclick = () => fillFlowToken(token.flow_id);
        actionTd.appendChild(btn);
        tr.appendChild(actionTd);
        tbody.appendChild(tr);
      });
    }

    function applyTokenRecord(flowPayload, silent) {
      if (!flowPayload || !flowPayload.access_token) {
        showFlash('记录中缺少 AccessToken', true);
        return;
      }
      try {
        const tokenArea = document.getElementById('tokenPayload');
        const payload = JSON.parse(tokenArea.value || '{}');
        payload.access_token = flowPayload.access_token;
        if (flowPayload.expires_in) {
          payload.access_token_ttl = flowPayload.expires_in;
        }
        tokenArea.value = JSON.stringify(payload, null, 2);

        const callArea = document.getElementById('callPayload');
        const callPayload = JSON.parse(callArea.value || '{}');
        callPayload.access_token = flowPayload.access_token;
        callArea.value = JSON.stringify(callPayload, null, 2);
        if (silent) {
          showFlash('已自动填充最新 AccessToken，可直接解析/调用 API。', false);
        } else {
          showFlash('已填充最新 AccessToken，可直接调用接口。', false);
        }
      } catch (err) {
        showFlash('填充 AccessToken 失败: ' + err.message, true);
      }
    }

    function setCallFormValues(payload) {
      try {
        apiActionSelect.value = payload.action || 'videos.list';
      } catch (e) {
        apiActionSelect.value = 'videos.list';
      }
      const action = apiActionSelect.value;
      const defaults = apiActionPresets[action] || apiActionPresets['videos.list'];
      callPartInput.value = payload.part || defaults.part || '';
      callIdsInput.value = payload.ids || defaults.ids || '';
      callQueryInput.value = payload.query || defaults.query || '';
      callChannelInput.value = payload.channel_id || defaults.channel_id || '';
      callMaxResultsInput.value = payload.max_results || defaults.max_results || 5;
      callMineCheckbox.checked = !!payload.mine;
      callSearchMineCheckbox.checked = !!payload.search_mine;
    }

    function buildCallPayloadFromForm() {
      const configPath = configPathInput.value.trim();
      const currentApp = getCurrentApp();
      const providerCode = currentApp?.provider_code || state.provider || '';
      const appCode = currentApp?.code || state.app || '';
      const modeKey = state.mode || getCurrentMode()?.key || '';
      const idsValue = callIdsInput.value.trim();
      const chartValue = idsValue ? '' : 'mostPopular';
      const action = apiActionSelect.value;
      return {
        provider_code: providerCode,
        provider_app: appCode,
        provider_auth_mode: modeKey,
        config_path: configPath,
        action: action,
        part: callPartInput.value.trim(),
        ids: idsValue,
        chart: chartValue,
        query: callQueryInput.value.trim(),
        channel_id: callChannelInput.value.trim(),
        max_results: parseInt(callMaxResultsInput.value || '5', 10) || 5,
        mine: callMineCheckbox.checked,
        search_mine: callSearchMineCheckbox.checked
      };
    }

    function syncCallPayloadFromForm() {
      const payload = buildCallPayloadFromForm();
      document.getElementById('callPayload').value = JSON.stringify(payload, null, 2);
    }

    function applyActionPreset(action) {
      const preset = apiActionPresets[action] || apiActionPresets['videos.list'];
      callPartInput.value = preset.part || '';
      callIdsInput.value = preset.ids || '';
      callQueryInput.value = preset.query || '';
      callChannelInput.value = preset.channel_id || '';
      callMaxResultsInput.value = preset.max_results || 5;
      callMineCheckbox.checked = !!preset.mine;
      callSearchMineCheckbox.checked = !!preset.search_mine;
      syncCallPayloadFromForm();
    }

    async function replayFlow(flowId) {
      const res = await authorizedFetch('/accesstoken/flow/replay', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ flow_id: flowId })
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'flow replay failed');
      }
      return res.json();
    }

    async function fillFlowToken(flowId, silent) {
      if (!flowId) {
        showFlash('Flow ID 缺失', true);
        return;
      }
      try {
        const replay = await replayFlow(flowId);
        applyTokenRecord(replay.payload || {}, silent);
      } catch (err) {
        showFlash('加载 Flow 失败: ' + err.message, true);
      }
    }

    async function loadTokenByFlow() {
      const flowId = (flowIdInput?.value || '').trim();
      if (!flowId) {
        showFlash('请输入 Flow ID', true);
        return;
      }
      try {
        const replay = await replayFlow(flowId);
        const payload = replay.payload || {};
        tokenCache = [{
          flow_id: replay.flow_id,
          status: replay.status,
          expire_at: payload.flow_expire_at || payload.expire_at,
          masked_account: payload.provider_code ? payload.provider_code + '/' + (payload.provider_app || '-') : '',
          storage_backend: replay.storage_backend
        }];
        renderTokenTable();
        applyTokenRecord(payload, true);
        showFlash('已加载 Flow ID 对应的授权记录，可点击“填充”复用。', false);
      } catch (err) {
        showFlash('加载 Flow ID 失败: ' + err.message, true);
      }
    }

    async function loadFlowIndexList() {
      try {
        const params = new URLSearchParams({ limit: '100' });
        if (state.flowCursor) {
          params.set('cursor', state.flowCursor);
        }
        const currentApp = getCurrentApp();
        const providerCode = currentApp?.provider_code || state.provider || '';
        const appCode = currentApp?.code || state.app || '';
        const modeKey = state.mode || getCurrentMode()?.key || '';
        if (providerCode) params.set('provider_code', providerCode);
        if (appCode) params.set('provider_app', appCode);
        if (modeKey) params.set('provider_auth_mode', modeKey);
        const res = await authorizedFetch('/accesstoken/flows?' + params.toString());
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        const data = await res.json();
        tokenCache = data.flows || [];
        state.flowCursor = data.next_cursor || '';
        renderTokenTable();
        if (!tokenCache.length) {
          showFlash('Redis 中暂未找到 Flow 记录，可先完成 OAuth。', true);
        } else {
          showFlash('已从 Flow 列表加载最新记录，可直接点击“填充”复用 Token。', false);
        }
      } catch (err) {
        showFlash('加载 Redis Flow 列表失败: ' + err.message, true);
      }
    }

    async function invoke(endpoint, payloadId, outputId) {
      const textarea = document.getElementById(payloadId);
      const output = document.getElementById(outputId);
      try {
        const payload = JSON.parse(textarea.value || '{}');
        const res = await authorizedFetch(endpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        const text = await res.text();
        output.textContent = text;
      } catch (err) {
        output.textContent = '提交失败: ' + err.message;
      }
    }

    async function loadCallbacks() {
      try {
        const res = await authorizedFetch('/api/callbacks');
        const data = await res.json();
        const tbody = document.querySelector('#callbackTable tbody');
        tbody.innerHTML = '';
        (data.records || []).forEach((item) => {
          const tr = document.createElement('tr');
          tr.innerHTML = '<td>' + new Date(item.timestamp).toISOString() + '</td>' +
                         '<td>' + (item.method || '-') + '</td>' +
                         '<td>' + (item.query || '-') + '</td>' +
                         '<td><pre>' + (item.body || '-') + '</pre></td>';
          tbody.appendChild(tr);
        });
      } catch (err) {
        showFlash('加载回调失败: ' + err.message, true);
      }
    }

    async function clearCallbacks() {
      try {
        await authorizedFetch('/api/callbacks/clear', { method: 'POST' });
        loadCallbacks();
      } catch (err) {
        showFlash('清空回调失败: ' + err.message, true);
      }
    }

    function showFlash(message, isError) {
      if (!flashBox) {
        return;
      }
      flashBox.textContent = message;
      flashBox.classList.add('show');
      if (isError) {
        flashBox.classList.add('error');
      } else {
        flashBox.classList.remove('error');
      }
      if (flashTimer) {
        clearTimeout(flashTimer);
      }
      flashTimer = setTimeout(() => {
        flashBox.classList.remove('show');
      }, 3500);
    }

    function attachCallFormListeners() {
    if (!apiActionSelect) {
      return;
    }
    apiActionSelect.addEventListener('change', () => {
      applyActionPreset(apiActionSelect.value);
    });
    if (redbookAdvertiserInput) {
      redbookAdvertiserInput.addEventListener('input', () => {
        if (isRedbookAction(apiActionSelect.value)) {
          syncCallPayloadFromForm();
        }
      });
    }
      [callPartInput, callIdsInput, callQueryInput, callChannelInput, callMaxResultsInput].forEach((el) => {
        if (!el) return;
        el.addEventListener('input', () => syncCallPayloadFromForm());
      });
      [callMineCheckbox, callSearchMineCheckbox].forEach((el) => {
        if (!el) return;
        el.addEventListener('change', () => syncCallPayloadFromForm());
      });
    }

    function init() {
      initAPIToken();
      setupEnvBanners();
      populateProviderSelect();
      callbackInput.value = document.body.dataset.defaultCallback || '';
      loadCallbacks();
      attachCallFormListeners();
    applyActionPreset(apiActionSelect.value);
    }

    init();
  </script>
</body>
</html>`))

func (s *accessTokenServer) handleDebugPage(w http.ResponseWriter, r *http.Request) {
	token := s.extractAPIToken(r)
	if token == "" {
		token = s.apiToken
	}
	data := debugPageData{
		APIToken:        token,
		ProvidersJSON:   s.providersJSON,
		DefaultProvider: s.defaultProviderUI,
		DefaultApp:      s.defaultApp,
		DefaultMode:     s.defaultMode,
		DefaultCallback: s.defaultCallbackURL,
		StorageBackend:  s.storageBackend,
		ListenAddr:      s.listenAddr,
		FlowTTLSeconds:  s.flowTTLSeconds,
		PublicListen:    s.listenAddrPublic,
	}
	if err := debugPageTemplate.Execute(w, data); err != nil {
		s.logger.ErrorF("accesstoken-server: render debug page failed: %v", err)
	}
}
