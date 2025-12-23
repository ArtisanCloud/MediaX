package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
)

type debugPageData struct {
	Providers       template.JS `json:"providers"`
	DefaultAPIToken string      `json:"default_api_token"`
	DefaultConfig   string      `json:"default_config"`
}

type debugProvider struct {
	Code string             `json:"code"`
	Name string             `json:"name"`
	Apps []debugProviderApp `json:"apps"`
}

type debugProviderApp struct {
	Code  string              `json:"code"`
	Name  string              `json:"name"`
	Modes []debugProviderMode `json:"modes"`
}

type debugProviderMode struct {
	Key             string `json:"key"`
	Label           string `json:"label"`
	AppID           string `json:"app_id"`
	ConfigPath      string `json:"config_path"`
	ProviderCode    string `json:"provider_code"`
	AppCode         string `json:"app_code"`
	ProviderType    string `json:"provider_type"`
	IdentifierLabel string `json:"identifier_label"`
	SupportsAPI     bool   `json:"supports_api"`
	SupportsMessage bool   `json:"supports_message"`
	RiskRequired    bool   `json:"risk_required,omitempty"`
	RiskReady       bool   `json:"risk_ready,omitempty"`
	RiskMissing     string `json:"risk_missing,omitempty"`
}

func buildProviderMeta(localCfg *config.LocalConfig, defaultConfig string) template.JS {
	result := make([]debugProvider, 0)
	if localCfg == nil || localCfg.ClientTokenProviders == nil {
		data, _ := json.Marshal(result)
		return template.JS(data)
	}
	for _, provider := range localCfg.ClientTokenProviders.Providers {
		if provider == nil {
			continue
		}
		metaProvider := debugProvider{
			Code: provider.Code,
			Name: provider.Name,
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			appMeta := debugProviderApp{
				Code: app.Code,
				Name: app.Name,
			}
			for _, mode := range app.AuthModes {
				metaMode := buildDebugProviderMode(app, mode, defaultConfig)
				if metaMode == nil {
					continue
				}
				appMeta.Modes = append(appMeta.Modes, *metaMode)
			}
			if len(appMeta.Modes) > 0 {
				metaProvider.Apps = append(metaProvider.Apps, appMeta)
			}
		}
		if len(metaProvider.Apps) > 0 {
			result = append(result, metaProvider)
		}
	}
	data, _ := json.Marshal(result)
	return template.JS(data)
}

func buildDebugProviderMode(app *config.ClientTokenProviderApp, mode *config.ClientTokenAuthMode, defaultConfig string) *debugProviderMode {
	if mode == nil {
		return nil
	}
	label := strings.TrimSpace(mode.Label)
	if label == "" {
		label = mode.Key
	}
	meta := &debugProviderMode{
		Key:             mode.Key,
		Label:           label,
		ConfigPath:      defaultConfig,
		ProviderCode:    effectiveClientTokenProviderCode(app, mode),
		AppCode:         strings.TrimSpace(app.Code),
		SupportsAPI:     true,
		SupportsMessage: false,
	}
	switch {
	case mode.WechatOfficialAccountConfig != nil:
		cfg := mode.WechatOfficialAccountConfig
		meta.ProviderType = string(providerTypeWechat)
		meta.AppID = cfg.AppIDValue()
		meta.IdentifierLabel = "AppID"
		meta.SupportsMessage = true
	case mode.ByteDanceDouYinConfig != nil:
		douyinCfg := mode.ByteDanceDouYinConfig
		douyinCfg.NormalizeClientTokenCredentials()
		cred := douyinCfg.ClientTokenCredential()
		clientKey := ""
		if cred != nil {
			clientKey = cred.ClientKeyValue()
		}
		if clientKey == "" && douyinCfg.ClientConfig != nil && douyinCfg.ClientConfig.OAuthConfig != nil {
			clientKey = strings.TrimSpace(douyinCfg.ClientConfig.OAuthConfig.ClientID)
		}
		meta.ProviderType = string(providerTypeByteDance)
		meta.AppID = clientKey
		meta.IdentifierLabel = "ClientKey"
		meta.RiskRequired = true
		meta.RiskReady = douyinCfg.RiskControlReady()
		if !meta.RiskReady {
			meta.RiskMissing = strings.Join(douyinCfg.MissingRiskControlFields(), ", ")
		}
	default:
		return nil
	}
	return meta
}

func effectiveClientTokenProviderCode(app *config.ClientTokenProviderApp, mode *config.ClientTokenAuthMode) string {
	if mode != nil {
		if code := strings.TrimSpace(mode.ProviderCode); code != "" {
			return code
		}
	}
	if app != nil {
		return app.ProviderCodeValue()
	}
	return ""
}

var clientTokenDebugTemplate = template.Must(template.New("clienttoken_debug").Parse(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <title>ClientToken 调试台</title>
  <style>
    body { font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif; margin:0; padding:0; background:#f5f6f8; color:#222; }
    header { background:#0f172a; color:#fff; padding:24px; }
    main { padding:24px; max-width:1200px; margin:0 auto; }
    section { background:#fff; border-radius:12px; padding:20px; margin-bottom:20px; box-shadow:0 8px 18px rgba(15,23,42,0.08); }
    h1 { margin:0; font-size:28px; }
    h2 { margin-top:0; border-left:4px solid #2563eb; padding-left:12px; }
    label { display:block; font-weight:600; margin-top:12px; }
    select, input, textarea { width:100%; margin-top:6px; padding:10px 12px; border:1px solid #d0d7de; border-radius:6px; box-sizing:border-box; font-family:inherit; }
    textarea { min-height:90px; resize:vertical; }
    button { padding:10px 16px; border:none; border-radius:6px; cursor:pointer; background:#2563eb; color:#fff; font-weight:600; }
    button.secondary { background:#475569; }
    button:disabled { opacity:.6; cursor:not-allowed; }
    .grid { display:grid; gap:16px; }
    .grid.two { grid-template-columns: repeat(auto-fit, minmax(240px,1fr)); }
    .row { display:flex; gap:12px; flex-wrap:wrap; }
    .row > * { flex:1; min-width:200px; }
    pre { background:#0f172a; color:#86efac; padding:12px; border-radius:8px; min-height:120px; overflow:auto; font-size:13px; }
    .muted { color:#64748b; font-size:13px; margin-top:4px; }
    .card { background:#0f172a; color:#fff; padding:16px; border-radius:10px; display:flex; justify-content:space-between; align-items:center; }
    .card h3 { margin:0 0 4px; font-size:20px; }
    .callbacks { max-height:260px; overflow:auto; border:1px solid #e2e8f0; border-radius:8px; }
    .callbacks table { width:100%; border-collapse:collapse; font-size:13px; }
    .callbacks th, .callbacks td { padding:8px; border-bottom:1px solid #e2e8f0; text-align:left; }
    .tag { display:inline-block; padding:2px 8px; border-radius:999px; font-size:12px; background:#e2e8f0; color:#0f172a; }
  </style>
</head>
<body>
  <header>
    <h1>ClientToken 调试台</h1>
    <p>默认监听 <code>http://127.0.0.1:7072/debug</code>，用于微信公众号 / 抖音 ClientToken 缓存、API 调试与日志排查。</p>
  </header>
  <main data-providers='{{.Providers}}' data-token='{{.DefaultAPIToken}}' data-config='{{.DefaultConfig}}'>
    <section>
      <h2>Provider 选择</h2>
      <div class="row">
        <div>
          <label>Provider</label>
          <select id="providerSelect"></select>
        </div>
        <div>
          <label>Provider App</label>
          <select id="appSelect"></select>
        </div>
        <div>
          <label>授权 Mode</label>
          <select id="modeSelect"></select>
        </div>
      </div>
      <p class="muted" id="providerMeta"></p>
      <p class="muted" id="riskWarning" style="display:none;color:#b91c1c;font-weight:600;"></p>
      <div class="row">
        <div>
          <label>config.yaml 路径</label>
          <input id="configPath" type="text" />
          <p class="muted">可覆盖默认配置文件路径，配合 <code>CLIENTTOKEN_CONFIG</code>。</p>
        </div>
        <div>
          <label>API Token</label>
          <input id="apiToken" type="password" />
          <p class="muted">默认 <code>dev-clienttoken</code>，可使用 <code>CLIENTTOKEN_API_TOKEN</code> 覆盖。</p>
        </div>
      </div>
    </section>

    <section>
      <div class="card">
        <div>
          <h3>Token 缓存</h3>
          <p id="tokenStatus" class="muted" style="color:#cbd5f5;">尚未刷新</p>
        </div>
        <div class="row" style="justify-content:flex-end;">
          <button onclick="refreshToken()">刷新 Token</button>
          <button class="secondary" onclick="loadCache()">查看缓存</button>
        </div>
      </div>
      <pre id="tokenOutput">// token result</pre>
    </section>

    <section>
      <h2>API 调试</h2>
      <label>常用 API 模板（会根据 Provider 显示微信或抖音示例）</label>
      <select id="apiPreset">
        <option value="">自定义（手动填写下面的请求参数）</option>
      </select>
      <div class="grid two">
        <div>
          <label id="apiActionLabel">微信 API 路径</label>
          <input id="apiAction" type="text" value="cgi-bin/getcallbackip" placeholder="cgi-bin/getcallbackip" />
        </div>
        <div>
          <label>HTTP Method</label>
          <select id="apiMethod">
            <option value="GET">GET</option>
            <option value="POST">POST</option>
          </select>
        </div>
      </div>
      <label>Query 参数（示例：next_openid=123）</label>
      <input id="apiQuery" type="text" />
      <label>Body（JSON，POST 时可选）</label>
      <textarea id="apiBody" placeholder="{}"></textarea>
      <button onclick="invokeAPI()">调用 API</button>
      <pre id="apiOutput">// api result</pre>
    </section>

    <section id="messageSection">
      <h2>消息验证 / 回调</h2>
      <div class="grid two">
        <div>
          <label>Signature</label>
          <input id="msgSignature" type="text" />
        </div>
        <div>
          <label>Timestamp</label>
          <input id="msgTimestamp" type="text" />
        </div>
        <div>
          <label>Nonce</label>
          <input id="msgNonce" type="text" />
        </div>
        <div>
          <label>EchoStr</label>
          <input id="msgEcho" type="text" />
        </div>
      </div>
      <button onclick="validateMessage()">验签</button>
      <pre id="messageOutput">// message result</pre>
      <h3>回调日志</h3>
      <div class="row">
        <button class="secondary" onclick="fetchCallbacks()">刷新日志</button>
        <button class="secondary" onclick="clearCallbacks()">清空日志</button>
      </div>
      <div class="callbacks" id="callbackTable"></div>
    </section>
  </main>

<script>
const PROVIDER_TYPE_WECHAT = 'wechat';
const PROVIDER_TYPE_DOUYIN = 'byte_dance_douyin';

const API_PRESETS = {
  wechat: [
    { label: '系统：获取 callback IP', action: 'cgi-bin/getcallbackip', method: 'GET', query: '', body: '' },
    { label: '用户：粉丝列表', action: 'cgi-bin/user/get', method: 'GET', query: 'next_openid=', body: '' },
    { label: '草稿：创建', action: 'cgi-bin/draft/add', method: 'POST', query: '', body: '{"articles":[{"title":"","author":"","digest":"","content":"","content_source_url":""}]}' },
    { label: '草稿：列表', action: 'cgi-bin/draft/batchget', method: 'POST', query: '', body: '{"offset":0,"count":10,"no_content":0}' },
    { label: '素材：统计', action: 'cgi-bin/material/get_materialcount', method: 'GET', query: '', body: '' },
    { label: '菜单：获取自定义菜单', action: 'cgi-bin/get_current_selfmenu_info', method: 'GET', query: '', body: '' },
  ],
  douyin: [
    { label: '内容：视频列表', action: 'open_api/1/content/video/list/', method: 'POST', query: '', body: '{"page":1,"size":10}' },
    { label: '内容：任务列表', action: 'open_api/1/content/task/list/', method: 'POST', query: '', body: '{"page":1,"size":10}' },
    { label: '工具：票据校验', action: 'open_api/1/tools/ticket/verify/', method: 'POST', query: '', body: '{"ticket_type":"API_TICKET"}' },
  ],
};

const state = {
  providers: [],
  tokenInput: document.getElementById('apiToken'),
  configInput: document.getElementById('configPath'),
  providerSelect: document.getElementById('providerSelect'),
  appSelect: document.getElementById('appSelect'),
  modeSelect: document.getElementById('modeSelect'),
  providerMeta: document.getElementById('providerMeta'),
  riskWarning: document.getElementById('riskWarning'),
  apiActionInput: document.getElementById('apiAction'),
  tokenOutput: document.getElementById('tokenOutput'),
  tokenStatus: document.getElementById('tokenStatus'),
  apiOutput: document.getElementById('apiOutput'),
  messageOutput: document.getElementById('messageOutput'),
  callbackTable: document.getElementById('callbackTable'),
  apiPreset: document.getElementById('apiPreset'),
  messageSection: document.getElementById('messageSection'),
};

function initPage() {
  const root = document.querySelector('main');
  try {
    state.providers = JSON.parse(root.dataset.providers || '[]');
  } catch (err) {
    state.providers = [];
  }
  state.tokenInput.value = root.dataset.token || '';
  state.configInput.value = root.dataset.config || 'config.yaml';
  state.providerSelect.addEventListener('change', () => renderAppOptions());
  state.appSelect.addEventListener('change', () => renderModeOptions());
  state.modeSelect.addEventListener('change', () => updateProviderMeta());
  if (state.apiPreset) {
    state.apiPreset.addEventListener('change', () => applyApiPreset());
    renderApiPresets();
  }
  if (state.apiActionInput) {
    state.apiActionInput.addEventListener('input', () => {
      state.apiActionInput.dataset.userEdited = '1';
    });
  }
  renderProviderOptions();
  fetchCallbacks();
  loadCache();
}

function renderApiPresets() {
  if (!state.apiPreset) return;
  const selection = currentSelection();
  const blocked = isRiskBlocked(selection && selection.mode);
  if (blocked) {
    state.apiPreset.innerHTML = '<option value="">Douyin API 模板已禁用（配置 device_id + risk_info 后启用）</option>';
    state.apiPreset.disabled = true;
    return;
  }
  state.apiPreset.disabled = false;
  const presets = getCurrentPresets();
  state.apiPreset.innerHTML = '<option value="">自定义（手动填写下面的请求参数）</option>';
  presets.forEach((preset, idx) => {
    const opt = document.createElement('option');
    opt.value = String(idx);
    opt.textContent = preset.label;
    state.apiPreset.appendChild(opt);
  });
}

function applyApiPreset() {
  if (!state.apiPreset) return;
  const value = state.apiPreset.value;
  if (value === '') {
    return;
  }
  const presets = getCurrentPresets();
  const preset = presets[parseInt(value, 10)];
  if (!preset) return;
  const selection = currentSelection();
  const providerType = selection && selection.mode ? selection.mode.provider_type : '';
  if (state.apiActionInput) {
    state.apiActionInput.value = preset.action;
    state.apiActionInput.dataset.userEdited = '';
    state.apiActionInput.dataset.lastProviderType = providerType || '';
  } else {
    document.getElementById('apiAction').value = preset.action;
  }
  document.getElementById('apiMethod').value = preset.method;
  document.getElementById('apiQuery').value = preset.query;
  document.getElementById('apiBody').value = preset.body;
}

function renderProviderOptions() {
  state.providerSelect.innerHTML = '';
  state.providers.forEach((item, idx) => {
    const opt = document.createElement('option');
    opt.value = String(idx);
    opt.textContent = item.name || item.code || 'provider';
    state.providerSelect.appendChild(opt);
  });
  renderAppOptions();
}


function renderAppOptions() {
  state.appSelect.innerHTML = '';
  const provider = getProvider();
  if (!provider) {
    state.appSelect.disabled = true;
    state.modeSelect.disabled = true;
    updateProviderMeta();
    return;
  }
  state.appSelect.disabled = false;
  provider.apps.forEach((app, idx) => {
    const opt = document.createElement('option');
    opt.value = String(idx);
    opt.textContent = app.name || app.code || 'app';
    state.appSelect.appendChild(opt);
  });
  renderModeOptions();
}

function renderModeOptions() {
  state.modeSelect.innerHTML = '';
  const selection = getApp();
  if (!selection) {
    state.modeSelect.disabled = true;
    updateProviderMeta();
    return;
  }
  state.modeSelect.disabled = false;
  selection.modes.forEach((mode, idx) => {
    const opt = document.createElement('option');
    opt.value = String(idx);
    const label = mode.label || mode.key || 'default';
    opt.textContent = label;
    state.modeSelect.appendChild(opt);
  });
  updateProviderMeta();
}

function getProvider() {
  const idx = parseInt(state.providerSelect.value || '0', 10);
  return state.providers[idx] || null;
}

function getApp() {
  const provider = getProvider();
  if (!provider || !provider.apps || provider.apps.length === 0) {
    return null;
  }
  const idx = parseInt(state.appSelect.value || '0', 10);
  return provider.apps[idx] || provider.apps[0];
}

function currentSelection() {
  const provider = getProvider();
  const app = getApp();
  if (!provider || !app || !app.modes || app.modes.length === 0) {
    return null;
  }
  const idx = parseInt(state.modeSelect.value || '0', 10);
  const mode = app.modes[idx] || app.modes[0];
  return { provider, app, mode };
}

function updateProviderMeta() {
  const selection = currentSelection();
  if (!selection) {
    state.providerMeta.textContent = '未找到 Provider 配置';
    toggleMessageSection(true);
    updateRiskWarning(null);
    updateApiLabel(null);
    syncApiActionField(null);
    renderApiPresets();
    return;
  }
  const mode = selection.mode;
  const providerType = mode.provider_type || '';
  const configPath = state.configInput.value || mode.config_path || 'config.yaml';
  const providerLabel = selection.provider.name || selection.provider.code || 'Provider';
  const appLabel = selection.app.name || selection.app.code || 'App';
  const modeLabel = mode.label || mode.key || 'Mode';
  const identifier = mode.identifier_label || 'AppID';
  const appId = mode.app_id || '-';
  state.providerMeta.innerHTML =
    providerLabel + ' / ' + appLabel + ' / ' + modeLabel +
    ' · ' + identifier + ': <code>' + appId + '</code> • config_path: <code>' + configPath + '</code>';
  updateRiskWarning(selection);
  updateApiLabel(providerType);
  syncApiActionField(providerType);
  renderApiPresets();
  toggleMessageSection(providerType === PROVIDER_TYPE_WECHAT);
}

function toggleMessageSection(enable) {
  if (!state.messageSection) return;
  state.messageSection.style.display = enable ? 'block' : 'none';
}

function getCurrentPresets() {
  const selection = currentSelection();
  if (!selection) return API_PRESETS.wechat;
  const type = selection.mode.provider_type;
  if (type === PROVIDER_TYPE_DOUYIN) {
    return API_PRESETS.douyin;
  }
  return API_PRESETS.wechat;
}

function isRiskBlocked(mode) {
  if (!mode) return false;
  if (mode.provider_type !== PROVIDER_TYPE_DOUYIN) {
    return false;
  }
  if (!mode.risk_required) {
    return false;
  }
  return !mode.risk_ready;
}

function updateRiskWarning(selection) {
  if (!state.riskWarning) return;
  const blocked = selection && isRiskBlocked(selection.mode);
  if (blocked) {
    const missing = selection.mode.risk_missing || 'device_id, risk_info';
    state.riskWarning.style.display = 'block';
    state.riskWarning.textContent = '当前未配置 ' + missing + '，已禁用 Douyin API 模板。';
  } else {
    state.riskWarning.style.display = 'none';
    state.riskWarning.textContent = '';
  }
}

function updateApiLabel(providerType) {
  const label = document.getElementById('apiActionLabel');
  if (!label) return;
  if (providerType === PROVIDER_TYPE_DOUYIN) {
    label.textContent = 'Douyin API 路径';
  } else {
    label.textContent = '微信 API 路径';
  }
}

function syncApiActionField(providerType) {
  const input = state.apiActionInput;
  if (!input) {
    return;
  }
  const defaultWechat = 'cgi-bin/getcallbackip';
  const defaultDouyin = 'open_api/1/content/video/list/';
  const defaultValue = providerType === PROVIDER_TYPE_DOUYIN ? defaultDouyin : defaultWechat;
  input.placeholder = defaultValue;
  const userEdited = input.dataset.userEdited === '1';
  if (!userEdited || input.value === defaultWechat || input.value === defaultDouyin || input.value.trim() === '') {
    input.value = defaultValue;
    input.dataset.userEdited = '';
  }
}

function withAuthHeaders(extra) {
  const token = state.tokenInput.value.trim();
  return Object.assign({
    'Content-Type': 'application/json',
    'Authorization': token ? 'Bearer ' + token : ''
  }, extra || {});
}

async function handleResponse(resp) {
  const text = await resp.text();
  let data = {};
  if (text) {
    try {
      data = JSON.parse(text);
    } catch (err) {
      data = { raw: text };
    }
  }
  if (!resp.ok) {
    const error = new Error((data && data.error) || resp.statusText);
    error.payload = data;
    throw error;
  }
  return data;
}

function showJSON(element, data) {
  element.textContent = JSON.stringify(data, null, 2);
}

function showRequestError(target, err) {
  if (err && err.payload && Object.keys(err.payload).length > 0) {
    showJSON(target, err.payload);
  } else {
    showJSON(target, { error: err.message || '请求失败' });
  }
}

async function refreshToken() {
  const sel = currentSelection();
  if (!sel) return;
  const meta = sel.mode;
  state.tokenStatus.textContent = '刷新中…';
  const payload = {
    provider_code: meta.provider_code,
    provider_app: meta.app_code,
    provider_auth_mode: meta.key,
    config_path: state.configInput.value || meta.config_path,
  };
  try {
    const resp = await fetch('/client-token/token', {
      method: 'POST',
      headers: withAuthHeaders(),
      body: JSON.stringify(payload),
    }).then(handleResponse);
    showJSON(state.tokenOutput, resp);
    state.tokenStatus.textContent = '最近刷新：' + new Date().toLocaleString();
  } catch (err) {
    state.tokenStatus.textContent = '刷新失败';
    showRequestError(state.tokenOutput, err);
  }
}

async function loadCache() {
  const sel = currentSelection();
  if (!sel) return;
  const meta = sel.mode;
  const params = new URLSearchParams({
    provider_code: meta.provider_code || '',
    provider_app: meta.app_code || '',
    provider_auth_mode: meta.key || '',
    config_path: state.configInput.value || meta.config_path || '',
  });
  try {
    const resp = await fetch('/client-token/cache?' + params.toString(), {
      headers: withAuthHeaders(),
    }).then(handleResponse);
    showJSON(state.tokenOutput, resp);
    if (resp.cached) {
      state.tokenStatus.textContent = '缓存命中，剩余 TTL ' + (resp.ttl_seconds ?? '-') + ' 秒';
    } else {
      state.tokenStatus.textContent = '未命中缓存';
    }
  } catch (err) {
    showRequestError(state.tokenOutput, err);
    state.tokenStatus.textContent = '读取缓存失败';
  }
}

async function invokeAPI() {
  const sel = currentSelection();
  if (!sel) return;
  const meta = sel.mode;
  const payload = {
    provider_code: meta.provider_code,
    provider_app: meta.app_code,
    provider_auth_mode: meta.key,
    config_path: state.configInput.value || meta.config_path,
    action: document.getElementById('apiAction').value.trim(),
    method: document.getElementById('apiMethod').value.trim(),
    query: document.getElementById('apiQuery').value.trim(),
    body: document.getElementById('apiBody').value.trim() || '{}',
  };
  try {
    const resp = await fetch('/client-token/call', {
      method: 'POST',
      headers: withAuthHeaders(),
      body: JSON.stringify(payload),
    }).then(handleResponse);
    showJSON(state.apiOutput, resp);
  } catch (err) {
    showRequestError(state.apiOutput, err);
  }
}

async function validateMessage() {
  const sel = currentSelection();
  if (!sel) return;
  const meta = sel.mode;
  const payload = {
    provider_code: meta.provider_code,
    provider_app: meta.app_code,
    provider_auth_mode: meta.key,
    config_path: state.configInput.value || meta.config_path,
    signature: document.getElementById('msgSignature').value.trim(),
    timestamp: document.getElementById('msgTimestamp').value.trim(),
    nonce: document.getElementById('msgNonce').value.trim(),
    echostr: document.getElementById('msgEcho').value.trim(),
  };
  try {
    const resp = await fetch('/client-token/message/validate', {
      method: 'POST',
      headers: withAuthHeaders(),
      body: JSON.stringify(payload),
    }).then(handleResponse);
    showJSON(state.messageOutput, resp);
  } catch (err) {
    showRequestError(state.messageOutput, err);
  }
}

async function fetchCallbacks() {
  try {
    const resp = await fetch('/client-token/message/callbacks', {
      headers: withAuthHeaders(),
    }).then(handleResponse);
    renderCallbackTable(resp.records || []);
  } catch (err) {
    const details = (err.payload && JSON.stringify(err.payload)) || err.message;
    state.callbackTable.innerHTML = '<p class="muted">加载失败：' + details + '</p>';
  }
}

async function clearCallbacks() {
  try {
    await fetch('/client-token/message/callbacks', {
      method: 'DELETE',
      headers: withAuthHeaders(),
    }).then(handleResponse);
    fetchCallbacks();
  } catch (err) {
    const details = (err.payload && JSON.stringify(err.payload)) || err.message;
    state.callbackTable.innerHTML = '<p class="muted">清空失败：' + details + '</p>';
  }
}

function renderCallbackTable(records) {
  if (!records.length) {
    state.callbackTable.innerHTML = '<p class="muted" style="padding:12px;">暂无回调日志，可向 <code>/client-token/message/callback</code> POST 消息。</p>';
    return;
  }
  let html = '<table><thead><tr><th>时间</th><th>Method</th><th>Query</th><th>Body</th></tr></thead><tbody>';
  records.forEach(rec => {
    html += '<tr>' +
      '<td>' + new Date(rec.timestamp).toLocaleString() + '</td>' +
      '<td><span class="tag">' + (rec.method || '-') + '</span></td>' +
      '<td><code>' + (rec.query || '-') + '</code></td>' +
      '<td><pre>' + (rec.body || '').replace(/</g,'&lt;').substring(0, 300) + '</pre></td>' +
      '</tr>';
  });
  html += '</tbody></table>';
  state.callbackTable.innerHTML = html;
}

initPage();
</script>
</body>
</html>
`))

func (s *clientTokenServer) handleDebugPage(w http.ResponseWriter, _ *http.Request) {
	data := debugPageData{
		Providers:       buildProviderMeta(s.localConfig, s.configPath),
		DefaultAPIToken: s.apiToken,
		DefaultConfig:   s.configPath,
	}
	if err := clientTokenDebugTemplate.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("render debug page: %v", err), http.StatusInternalServerError)
	}
}
