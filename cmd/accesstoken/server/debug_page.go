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
  </style>
</head>
<body data-api-token="{{.APIToken}}" data-default-provider="{{.DefaultProvider}}" data-default-app="{{.DefaultApp}}" data-default-mode="{{.DefaultMode}}" data-default-callback="{{.DefaultCallback}}">
  <h1>AccessToken 调试台</h1>
  <p>该页面复用了 SessionToken 调试模式，可在浏览器内选择 Provider/App，直接调用当前服务的 AccessToken API，并观察 <code>/debug/callback</code>。</p>
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
        <tr><th>Flow ID</th><th>Token Source</th><th>有效期</th><th>操作</th></tr>
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
      <button type="button" onclick="invoke('/accesstoken/call', 'callPayload', 'callOutput')">执行 API</button>
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
    const apiToken = document.body.dataset.apiToken;
    const defaultProvider = document.body.dataset.defaultProvider || '';
    const defaultApp = document.body.dataset.defaultApp || '';
    const defaultMode = document.body.dataset.defaultMode || '';
    const providers = {{.ProvidersJSON}};
    const state = {
      provider: defaultProvider || (providers[0]?.code || ''),
      app: defaultApp || (providers[0]?.apps?.[0]?.code || ''),
      mode: defaultMode || (providers[0]?.apps?.[0]?.modes?.[0]?.key || '')
    };

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
    let tokenCache = [];
    let flashTimer = null;

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
      loadTokens(false);
    }

    function applyTemplate() {
      const configPath = configPathInput.value.trim();
      const currentApp = getCurrentApp();
      const providerCode = currentApp?.provider_code || state.provider || '';
      const appCode = currentApp?.app_key || currentApp?.code || state.app || '';
      const modeKey = state.mode || getCurrentMode()?.key || '';
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
        const res = await fetch('/api/oauth/start', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + apiToken
          },
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
          alert('接口未返回授权地址');
        }
      } catch (err) {
        alert('发起授权失败: ' + err);
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
        const res = await fetch('/api/oauth/tokens?' + params.toString(), {
          headers: {
            'Authorization': 'Bearer ' + apiToken
          }
        });
        if (!res.ok) {
          if (force) {
            const text = await res.text();
            alert('加载授权记录失败: ' + text);
          }
          return;
        }
        const data = await res.json();
        tokenCache = data.tokens || [];
        renderTokenTable();
        if (force && tokenCache.length === 0) {
          showFlash('未找到授权记录，请先完成 OAuth 或使用 Flow ID 回填。', true);
        }
      } catch (err) {
        if (force) {
          alert('加载授权记录失败: ' + err);
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
        tr.innerHTML = '<td>' + (token.flow_id || '-') + '</td>' +
                       '<td>' + (token.source || token.token_type || '-') + '</td>' +
                       '<td>' + expireText + '</td>';
        const actionTd = document.createElement('td');
        const btn = document.createElement('button');
        btn.textContent = '填充';
        btn.type = 'button';
        btn.onclick = () => applyTokenRecord(token);
        actionTd.appendChild(btn);
        tr.appendChild(actionTd);
        tbody.appendChild(tr);
      });
    }

    function applyTokenRecord(token, silent) {
      if (!token || !token.access_token) {
        alert('记录中缺少 AccessToken');
        return;
      }
      try {
        const tokenArea = document.getElementById('tokenPayload');
        const payload = JSON.parse(tokenArea.value || '{}');
        payload.access_token = token.access_token;
        if (token.expires_in) {
          payload.access_token_ttl = token.expires_in;
        }
        tokenArea.value = JSON.stringify(payload, null, 2);

        const callArea = document.getElementById('callPayload');
        const callPayload = JSON.parse(callArea.value || '{}');
        callPayload.access_token = token.access_token;
        callArea.value = JSON.stringify(callPayload, null, 2);
        if (silent) {
          showFlash('已自动填充最新 AccessToken，可直接解析/调用 API。', false);
        } else {
          alert('已填充最新 AccessToken，可直接调用接口。');
        }
      } catch (err) {
        alert('填充 AccessToken 失败: ' + err);
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
      return {
        provider_code: providerCode,
        provider_app: appCode,
        provider_auth_mode: modeKey,
        config_path: configPath,
        action: apiActionSelect.value,
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

    async function loadTokenByFlow() {
      const flowId = (flowIdInput?.value || '').trim();
      if (!flowId) {
        alert('请输入 Flow ID');
        return;
      }
      try {
        const params = new URLSearchParams({ flow_id: flowId });
        const res = await fetch('/api/oauth/tokens?' + params.toString(), {
          headers: { 'Authorization': 'Bearer ' + apiToken }
        });
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        const data = await res.json();
        tokenCache = data.tokens || [];
        renderTokenTable();
        if (!tokenCache.length) {
          showFlash('未找到 Flow ID 对应的授权记录，请确认是否仍在有效期。', true);
        } else {
          showFlash('已加载 Flow ID 对应的授权记录，可点击“填充”复用。', false);
          applyTokenRecord(tokenCache[0], true);
        }
      } catch (err) {
        showFlash('加载 Flow ID 失败: ' + err, true);
      }
    }

    async function loadFlowIndexList() {
      try {
        const params = new URLSearchParams({ limit: '100' });
        const res = await fetch('/api/oauth/flow-indexes?' + params.toString(), {
          headers: { 'Authorization': 'Bearer ' + apiToken }
        });
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        const data = await res.json();
        tokenCache = data.tokens || [];
        renderTokenTable();
        if (!tokenCache.length) {
          showFlash('Redis 中暂未找到 Flow 记录，可先完成 OAuth。', true);
        } else {
          showFlash('已从 Redis 加载 Flow 列表，可直接点击“填充”复用 Token。', false);
        }
      } catch (err) {
        showFlash('加载 Redis Flow 列表失败: ' + err, true);
      }
    }

    async function invoke(endpoint, payloadId, outputId) {
      const textarea = document.getElementById(payloadId);
      const output = document.getElementById(outputId);
      try {
        const payload = JSON.parse(textarea.value || '{}');
        const res = await fetch(endpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + apiToken
          },
          body: JSON.stringify(payload)
        });
        const text = await res.text();
        output.textContent = text;
      } catch (err) {
        output.textContent = '提交失败: ' + err;
      }
    }

    async function loadCallbacks() {
      const res = await fetch('/api/callbacks', {
        headers: { 'Authorization': 'Bearer ' + apiToken }
      });
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
    }

    async function clearCallbacks() {
      await fetch('/api/callbacks/clear', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + apiToken }
      });
      loadCallbacks();
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
	data := debugPageData{
		APIToken:        s.apiToken,
		ProvidersJSON:   s.providersJSON,
		DefaultProvider: s.defaultProviderUI,
		DefaultApp:      s.defaultApp,
		DefaultMode:     s.defaultMode,
		DefaultCallback: s.defaultCallbackURL,
	}
	if err := debugPageTemplate.Execute(w, data); err != nil {
		s.logger.ErrorF("accesstoken-server: render debug page failed: %v", err)
	}
}
