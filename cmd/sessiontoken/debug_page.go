package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type debugPageData struct {
	DefaultToken    string
	DefaultCallback string
}

type callbackRecord struct {
	Timestamp    time.Time         `json:"timestamp"`
	Method       string            `json:"method"`
	FlowID       string            `json:"flow_id,omitempty"`
	ProviderCode string            `json:"provider_code,omitempty"`
	Query        string            `json:"query,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body"`
}

type callbackLogStore struct {
	mu      sync.Mutex
	limit   int
	records []callbackRecord
}

func newCallbackLogStore(limit int) *callbackLogStore {
	if limit <= 0 {
		limit = 20
	}
	return &callbackLogStore{limit: limit}
}

func (s *callbackLogStore) append(record callbackRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append([]callbackRecord{record}, s.records...)
	if len(s.records) > s.limit {
		s.records = s.records[:s.limit]
	}
}

func (s *callbackLogStore) list() []callbackRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	cloned := make([]callbackRecord, len(s.records))
	copy(cloned, s.records)
	return cloned
}

func (s *callbackLogStore) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = nil
}

const (
	debugFlowMetadataPrefix = "/debug/flows/"
)

var debugPageTemplate = template.Must(template.New("debug_page").Parse(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8" />
<title>SessionToken Debugger</title>
<style>
  body { font-family: Arial, sans-serif; margin: 24px; background:#f7f7f7; }
  h1 { margin-bottom: 8px; }
  section { background:#fff; padding:16px; margin-bottom:16px; border-radius:8px; box-shadow:0 1px 3px rgba(0,0,0,0.08); }
  label { display:block; margin-top:8px; font-weight:600; }
  input, textarea, select { width:100%; padding:8px; margin-top:4px; box-sizing:border-box; font-family:inherit; }
  textarea { min-height:80px; }
  button { margin-top:12px; padding:8px 16px; cursor:pointer; }
  pre { background:#111; color:#0f0; padding:12px; border-radius:8px; min-height:160px; overflow:auto; }
  .row { display:flex; gap:16px; }
  .row > div { flex:1; }
  .inline { display:flex; gap:12px; align-items:center; }
</style>
</head>
<body data-default-token="{{.DefaultToken}}" data-default-callback="{{.DefaultCallback}}">
<h1>SessionToken 调试台</h1>
<p>该页面仅用于本地调试，所有请求会直接调用当前服务的 <code>/session-token/flows</code> API，并内置 <code>/debug/callback</code> 用于 Mock 回调观察。</p>
<section>
  <h2>基础配置</h2>
  <label>Provider
    <select id="provider">
      <option value="zhihu" selected>知乎</option>
      <option value="redbook">小红书</option>
      <option value="google">Google / YouTube</option>
    </select>
  </label>
  <label>Provider App
    <select id="providerApp"></select>
  </label>
  <label id="apiVersionWrapper" style="display:none">
    API 版本
    <select id="apiVersion"></select>
  </label>
  <div id="customAppWrapper" style="display:none">
    <label>自定义 Provider App Code
      <input id="customAppInput" type="text" placeholder="如 zhihu_article" />
    </label>
  </div>
  <label>API Token
    <input id="token" type="password" placeholder="Bearer Token" />
  </label>
  <label>Callback URL
    <input id="callback" type="text" placeholder="http://127.0.0.1:7070/debug/callback" />
  </label>
  <div class="inline">
    <div style="flex:1">
      <label>State（可自定义，默认随机）
        <input id="state" type="text" />
      </label>
    </div>
    <div style="flex:1">
      <label>Tenant UUID
        <input id="tenant" type="text" value="tenant_debug" />
      </label>
    </div>
    <div style="flex:1">
      <label>Account ID
        <input id="account" type="text" value="acct_debug" />
      </label>
    </div>
  </div>
  <label>Metadata（JSON）
    <textarea id="metadata" data-autofill="1">{"login_mode":"pc"}</textarea>
  </label>
  <label class="inline">
    <input type="checkbox" id="reuseSession" checked />
    浏览器已登录则直接复用 Cookie（metadata.reuse_session）
  </label>
  <div class="row">
    <button type="button" id="applyPresetBtn">应用模板</button>
    <button type="button" onclick="createFlow()">创建 Flow</button>
  </div>
</section>
<section>
  <h2>Flow 详情</h2>
  <label>Flow ID
    <input id="flowId" type="text" placeholder="stf_xxx" />
  </label>
  <div class="row">
    <button type="button" onclick="pollFlow()">查询 Flow</button>
    <button type="button" onclick="copyAuthorizeURL()">复制 authorize_url</button>
  </div>
  <label>authorize_url</label>
  <input id="authorizeUrl" type="text" readonly />
</section>
<section>
  <h2>Playwright 自动调试器</h2>
  <p>若希望自动拉起 Chromium 并抓取知乎 Cookie，可在仓库根目录运行 <code>pnpm install</code> 与 <code>npx playwright install chromium</code> 后执行下方命令：</p>
  <label>CLI 命令</label>
  <div class="row">
    <input id="playwrightCommand" type="text" readonly />
    <button type="button" onclick="copyPlaywrightCommand()">复制命令</button>
  </div>
  <p style="font-size:12px;color:#666;">命令默认使用 <code>SESSIONTOKEN_API_TOKEN</code> 环境变量（默认 <code>dev-session-token</code>）。可在终端中替换 <code>&lt;flow_id&gt;</code> 或追加 <code>--api-token</code>、<code>--base</code> 参数。</p>
</section>
<section>
  <h2>API 调试（Beta）</h2>
  <p>选择 Provider 后可直接调用对应 API，默认会附带 <code>Authorization: Bearer ...</code>。如需 <code>X-SessionToken</code> 可以点击“从 Flow 填充”。</p>
  <label>API Endpoint
    <select id="apiEndpoint"></select>
  </label>
  <div id="apiInfo" style="font-size:12px;color:#666;margin-top:4px;">当前 Provider 暂无 API</div>
  <label>Query/Path 参数
    <input id="apiQuery" type="text" placeholder="limit=10&offset=0 或 article_id=123" />
  </label>
  <label>请求 Body（JSON，可选）
    <textarea id="apiBody" placeholder="{}"></textarea>
  </label>
  <label>Session Token（X-SessionToken）
    <input id="apiSessionToken" type="text" placeholder="点击下方按钮可从 Flow metadata 填充" />
  </label>
  <div class="row">
    <button type="button" onclick="loadSessionTokenFromFlow()">从 Flow 填充 SessionToken</button>
    <button type="button" onclick="sendApiRequest()">发送请求</button>
  </div>
  <pre id="apiResponse">等待 API 请求...</pre>
</section>
<section>
  <h2>输出</h2>
  <pre id="output">等待操作...</pre>
</section>
<section>
  <h2>Mock 回调 &amp; 日志</h2>
  <p>默认会将 Callback URL 设为 <code>/debug/callback</code>，便于在本地直接接收并查看 SessionToken 回调 payload。</p>
  <div class="row">
    <button type="button" onclick="refreshCallbackLogs()">刷新回调日志</button>
    <button type="button" onclick="clearCallbackLogs()">清空日志</button>
  </div>
  <label>最近回调（最多 20 条，新到旧）</label>
  <pre id="callbackLog">暂无回调...</pre>
</section>
<script>
const CUSTOM_APP_VALUE = '__custom__';
const apiBase = window.location.origin;
const defaults = {
  token: document.body.dataset.defaultToken || 'dev-session-token',
  callback: document.body.dataset.defaultCallback || (apiBase + '/debug/callback')
};
const providerSelect = document.getElementById('provider');
const providerAppSelect = document.getElementById('providerApp');
const apiVersionWrapper = document.getElementById('apiVersionWrapper');
const apiVersionSelect = document.getElementById('apiVersion');
const customAppWrapper = document.getElementById('customAppWrapper');
const customAppInput = document.getElementById('customAppInput');
const reuseInput = document.getElementById('reuseSession');
const tokenInput = document.getElementById('token');
const callbackInput = document.getElementById('callback');
const stateInput = document.getElementById('state');
const metadataInput = document.getElementById('metadata');
const tenantInput = document.getElementById('tenant');
const accountInput = document.getElementById('account');
const flowIdInput = document.getElementById('flowId');
const authorizeInput = document.getElementById('authorizeUrl');
const output = document.getElementById('output');
const callbackLogOutput = document.getElementById('callbackLog');
const applyPresetBtn = document.getElementById('applyPresetBtn');
const playwrightCommandInput = document.getElementById('playwrightCommand');
const apiEndpointSelect = document.getElementById('apiEndpoint');
const apiInfoLabel = document.getElementById('apiInfo');
const apiQueryInput = document.getElementById('apiQuery');
const apiBodyInput = document.getElementById('apiBody');
const apiSessionTokenInput = document.getElementById('apiSessionToken');
const apiResponseOutput = document.getElementById('apiResponse');

const providerCatalog = {
  zhihu: {
    label: '知乎',
    versions: [
      { value: 'v4', label: 'v4（默认，当前支持）' }
    ],
    defaults: { metadata: '{"login_mode":"pc","api_version":"v4"}', token: defaults.token, callback: defaults.callback },
    apps: [
      { code: 'zhihu_article', label: '知乎 Web（抓取文章）', metadata: '{"login_mode":"pc","api_version":"v4"}', token: defaults.token, callback: defaults.callback },
      { code: 'zhihu_mobile', label: '知乎移动端', metadata: '{"login_mode":"mobile","api_version":"v4"}', token: defaults.token, callback: defaults.callback }
    ]
  },
  redbook: {
    label: '小红书',
    versions: [],
    defaults: { metadata: '{"login_mode":"qr"}', token: defaults.token, callback: defaults.callback },
    apps: [
      { code: 'redbook_web', label: '小红书 Web', metadata: '{"login_mode":"qr"}', token: defaults.token, callback: defaults.callback },
      { code: 'redbook_creator', label: '小红书创作号', metadata: '{"login_mode":"pc"}', token: defaults.token, callback: defaults.callback }
    ]
  },
  google: {
    label: 'Google / YouTube',
    versions: [],
    defaults: { metadata: '{"login_mode":"oauth"}', token: defaults.token, callback: defaults.callback },
    apps: [
      { code: 'youtube_web', label: 'YouTube Web', metadata: '{"login_mode":"oauth"}', token: defaults.token, callback: defaults.callback }
    ]
  }
};

const apiCatalog = {
  zhihu: [
    {
      id: 'followings',
      label: 'GET /zhihu/v1/me/followings',
      method: 'GET',
      path: '/zhihu/v1/me/followings',
      sampleQuery: 'limit=10&offset=0',
      requiresSessionToken: true,
      description: '获取当前账号关注对象'
    },
    {
      id: 'channels',
      label: 'GET /zhihu/v1/channels/{channel_id}/articles',
      method: 'GET',
      path: '/zhihu/v1/channels/{channel_id}/articles',
      sampleQuery: 'channel_id=12345&limit=10&offset=0',
      requiresSessionToken: true,
      description: '查询某频道下文章列表'
    },
    {
      id: 'article_get',
      label: 'GET /zhihu/v1/articles/{article_id}',
      method: 'GET',
      path: '/zhihu/v1/articles/{article_id}',
      sampleQuery: 'article_id=1234567890',
      requiresSessionToken: true,
      description: '获取单篇文章详情'
    },
    {
      id: 'article_post',
      label: 'POST /zhihu/v1/articles',
      method: 'POST',
      path: '/zhihu/v1/articles',
      sampleBody: '{\n  "title": "示例标题",\n  "content": "<p>示例内容</p>"\n}',
      requiresSessionToken: true,
      description: '发布文章，content 支持 HTML'
    },
    {
      id: 'sanity_check',
      label: 'POST /zhihu/v1/sanity/check',
      method: 'POST',
      path: '/zhihu/v1/sanity/check',
      sampleBody: '{ }',
      requiresSessionToken: true,
      description: '调用知乎自检接口，验证 SessionToken 是否有效'
    }
  ],
  redbook: [],
  google: []
};

function apiEndpointStorageKey(provider) {
  return 'st_debug_api_endpoint_' + provider;
}

function rebuildApiEndpoints(preserveSelection) {
  if (!apiEndpointSelect) {
    return;
  }
  const provider = providerSelect.value;
  const endpoints = apiCatalog[provider] || [];
  apiEndpointSelect.innerHTML = '';
  if (!endpoints.length) {
    apiEndpointSelect.disabled = true;
    setApiInfoMessage('当前 Provider 暂无 API');
    resetApiResponse();
    return;
  }
  apiEndpointSelect.disabled = false;
  const savedValue = preserveSelection ? (localStorage.getItem(apiEndpointStorageKey(provider)) || '') : '';
  const currentValue = savedValue || apiEndpointSelect.value;
  endpoints.forEach((endpoint) => {
    const option = document.createElement('option');
    option.value = endpoint.id;
    option.textContent = endpoint.label;
    apiEndpointSelect.appendChild(option);
  });
  if (currentValue && endpoints.some(ep => ep.id === currentValue)) {
    apiEndpointSelect.value = currentValue;
  } else {
    apiEndpointSelect.value = endpoints[0].id;
  }
  updateApiInfo();
}

function getSelectedApiEndpoint() {
  if (!apiEndpointSelect) {
    return null;
  }
  const provider = providerSelect.value;
  const endpoints = apiCatalog[provider] || [];
  const selected = apiEndpointSelect.value;
  if (!selected) {
    return endpoints[0] || null;
  }
  return endpoints.find(ep => ep.id === selected) || endpoints[0] || null;
}

function setApiInfoMessage(message) {
  if (apiInfoLabel) {
    apiInfoLabel.textContent = message || '当前 Provider 暂无 API';
  }
}

function updateApiInfo() {
  const endpoint = getSelectedApiEndpoint();
  if (!endpoint) {
    setApiInfoMessage('当前 Provider 暂无 API');
    return;
  }
  localStorage.setItem(apiEndpointStorageKey(providerSelect.value), endpoint.id);
  const parts = [];
  if (endpoint.method && endpoint.path) {
    parts.push(endpoint.method + ' ' + endpoint.path);
  }
  if (endpoint.description) {
    parts.push('— ' + endpoint.description);
  }
  if (endpoint.requiresSessionToken) {
    parts.push('(需 SessionToken)');
  }
  setApiInfoMessage(parts.join(' '));
  if (apiQueryInput) {
    apiQueryInput.placeholder = endpoint.sampleQuery || 'limit=10&offset=0';
  }
  if (apiBodyInput) {
    apiBodyInput.placeholder = endpoint.sampleBody || '{}';
  }
}

function parseQueryInput(raw) {
  const params = {};
  if (!raw) {
    return params;
  }
  raw.split('&').forEach((segment) => {
    const piece = segment.trim();
    if (!piece) {
      return;
    }
    const idx = piece.indexOf('=');
    if (idx === -1) {
      params[decodeURIComponent(piece)] = '';
      return;
    }
    const key = decodeURIComponent(piece.slice(0, idx).trim());
    const value = decodeURIComponent(piece.slice(idx + 1).trim());
    params[key] = value;
  });
  return params;
}

function applyPathParams(path, params) {
  const usedKeys = new Set();
  if (!path) {
    return { path: '', usedKeys };
  }
  const replaced = path.replace(/\{([^}]+)\}/g, (match, key) => {
    const trimmed = key.trim();
    if (!trimmed) {
      return '';
    }
    if (!(trimmed in params) || params[trimmed] === '') {
      throw new Error('缺少路径参数: ' + trimmed);
    }
    usedKeys.add(trimmed);
    return encodeURIComponent(params[trimmed]);
  });
  return { path: replaced, usedKeys };
}

function buildQueryString(params, usedKeys) {
  const search = new URLSearchParams();
  Object.keys(params || {}).forEach((key) => {
    if (usedKeys.has(key)) {
      return;
    }
    const value = params[key];
    if (value !== undefined && value !== null && value !== '') {
      search.append(key, value);
    }
  });
  return search.toString();
}

function resetApiResponse() {
  if (apiResponseOutput) {
    apiResponseOutput.textContent = '等待 API 请求...';
  }
}

function setApiResponse(message) {
  if (apiResponseOutput) {
    apiResponseOutput.textContent = message;
  }
}

function providerDefaultField(provider, field) {
  const defaults = providerCatalog[provider]?.defaults || {};
  if (defaults && defaults[field]) {
    return defaults[field];
  }
  return '';
}

function resolveFieldDefault(field) {
  const option = providerAppSelect.options[providerAppSelect.selectedIndex];
  if (option && option.dataset && option.dataset[field]) {
    return option.dataset[field];
  }
  return providerDefaultField(providerSelect.value, field) || defaults[field] || '';
}

(function restore() {
  const lastProvider = localStorage.getItem('st_debug_provider');
  if (lastProvider && providerCatalog[lastProvider]) {
    providerSelect.value = lastProvider;
  }
  rebuildAppOptions(true);
  rebuildApiEndpoints(true);
  const storedToken = localStorage.getItem('st_debug_token');
  tokenInput.value = storedToken && storedToken.trim() ? storedToken : resolveFieldDefault('token');
  const storedCallback = localStorage.getItem('st_debug_callback');
  callbackInput.value = storedCallback && storedCallback.trim() ? storedCallback : resolveFieldDefault('callback');
  applyPreset(false, false);
  customAppInput.value = localStorage.getItem('st_debug_custom_app') || '';
  const storedReuse = localStorage.getItem('st_debug_reuse');
  reuseInput.checked = storedReuse === null ? true : storedReuse === '1';
  stateInput.value = randomState();
  metadataInput.addEventListener('input', () => metadataInput.dataset.autofill = '0');
  resetApiResponse();
})();

function rebuildAppOptions(preserveSelection) {
  const provider = providerSelect.value;
  const catalog = providerCatalog[provider] || { apps: [] };
  const savedApp = preserveSelection ? localStorage.getItem('st_debug_provider_app') : '';
  const currentValue = preserveSelection ? savedApp : providerAppSelect.value;
  providerAppSelect.innerHTML = '';
  const apps = catalog.apps || [];
  apps.forEach((app) => {
    const option = document.createElement('option');
    option.value = app.code;
    option.textContent = app.label + ' (' + app.code + ')';
    option.dataset.metadata = app.metadata || '';
    if (app.token) {
      option.dataset.token = app.token;
    }
    if (app.callback) {
      option.dataset.callback = app.callback;
    }
    providerAppSelect.appendChild(option);
  });
  const customOption = document.createElement('option');
  customOption.value = CUSTOM_APP_VALUE;
  customOption.textContent = '自定义...';
  providerAppSelect.appendChild(customOption);

  if (currentValue && (apps.some(a => a.code === currentValue) || currentValue === CUSTOM_APP_VALUE)) {
    providerAppSelect.value = currentValue;
  } else if (apps.length > 0) {
    providerAppSelect.value = apps[0].code;
  } else {
    providerAppSelect.value = CUSTOM_APP_VALUE;
  }
  toggleCustomAppInput();
  rebuildVersionOptions(preserveSelection);
}

function toggleCustomAppInput() {
  const isCustom = providerAppSelect.value === CUSTOM_APP_VALUE;
  customAppWrapper.style.display = isCustom ? 'block' : 'none';
}

function rebuildVersionOptions(preserveSelection) {
  const catalog = providerCatalog[providerSelect.value] || {};
  const versions = catalog.versions || [];
  if (!versions.length) {
    apiVersionWrapper.style.display = 'none';
    apiVersionSelect.innerHTML = '';
    return;
  }
  apiVersionWrapper.style.display = 'block';
  const savedVersion = preserveSelection ? localStorage.getItem('st_debug_api_version') : '';
  const currentValue = savedVersion || apiVersionSelect.value;
  apiVersionSelect.innerHTML = '';
  versions.forEach((version) => {
    const option = document.createElement('option');
    option.value = version.value;
    option.textContent = version.label || version.value;
    apiVersionSelect.appendChild(option);
  });
  if (currentValue && versions.some(v => v.value === currentValue)) {
    apiVersionSelect.value = currentValue;
  } else if (versions.length > 0) {
    apiVersionSelect.value = versions[0].value;
  }
}

function selectedVersion() {
  if (apiVersionWrapper.style.display === 'none' || !apiVersionSelect.value) {
    return '';
  }
  return apiVersionSelect.value;
}

function getProviderAppCode() {
  if (providerAppSelect.value === CUSTOM_APP_VALUE) {
    return customAppInput.value.trim();
  }
  return providerAppSelect.value;
}

function currentAppMetadata() {
  if (providerAppSelect.value === CUSTOM_APP_VALUE) {
    return providerCatalog[providerSelect.value]?.defaults?.metadata || '';
  }
  const option = providerAppSelect.options[providerAppSelect.selectedIndex];
  if (option && option.dataset.metadata) {
    return option.dataset.metadata;
  }
  const provider = providerSelect.value;
  return providerCatalog[provider]?.defaults?.metadata || '';
}

function log(msg) {
  const time = new Date().toISOString();
  output.textContent = '[' + time + ']\n' + msg;
}

function buildPlaywrightCommand() {
  const flowId = flowIdInput.value.trim() || '<flow_id>';
  let command = 'pnpm sessiontoken:browser --flow ' + flowId;
  if (apiBase && apiBase !== 'http://127.0.0.1:7070') {
    command += ' --base ' + apiBase;
  }
  return command;
}

function updatePlaywrightCommand() {
  if (playwrightCommandInput) {
    playwrightCommandInput.value = buildPlaywrightCommand();
  }
}

function copyPlaywrightCommand() {
  if (!playwrightCommandInput) {
    return;
  }
  navigator.clipboard.writeText(playwrightCommandInput.value).then(() => {
    log('Playwright CLI 命令已复制，可在终端运行。');
  }).catch((err) => {
    log('复制 CLI 命令失败: ' + err);
  });
}

function saveBasics() {
  localStorage.setItem('st_debug_token', tokenInput.value.trim());
  localStorage.setItem('st_debug_callback', callbackInput.value.trim());
  localStorage.setItem('st_debug_provider', providerSelect.value);
  localStorage.setItem('st_debug_provider_app', providerAppSelect.value);
  localStorage.setItem('st_debug_custom_app', customAppInput.value.trim());
  localStorage.setItem('st_debug_api_version', apiVersionSelect.value);
  localStorage.setItem('st_debug_reuse', reuseInput.checked ? '1' : '0');
}

function randomState() {
  return 'ui-debug-' + Math.random().toString(36).slice(2, 8);
}

function applyPreset(forceFields, showLog) {
  const defaultToken = resolveFieldDefault('token');
  if (forceFields || !tokenInput.value.trim()) {
    tokenInput.value = defaultToken;
  }
  const defaultCallback = resolveFieldDefault('callback');
  if (forceFields || !callbackInput.value.trim()) {
    callbackInput.value = defaultCallback;
  }
  const metadataTemplate = currentAppMetadata();
  if (metadataTemplate && (forceFields || metadataInput.dataset.autofill === '1' || !metadataInput.value.trim())) {
    metadataInput.value = metadataTemplate;
    metadataInput.dataset.autofill = '1';
  }
  if (showLog) {
    log('已应用 ' + providerSelect.value + ' / ' + (getProviderAppCode() || '默认') + ' 模板');
  }
}

function buildSessionTokenFromPieces(metadata) {
  if (!metadata) {
    return '';
  }
  const mapping = [
    ['cookie_sessionid', 'SESSIONID'],
    ['cookie_joid', 'JOID'],
    ['cookie_osd', 'osd'],
    ['cookie_q_c1', 'q_c1'],
    ['cookie_d_c0', 'd_c0'],
    ['cookie_unlock_ticket', 'unlock_ticket'],
    ['cookie_z_c0', 'z_c0'],
  ];
  const parts = [];
  mapping.forEach(([key, cookieName]) => {
    const value = metadata[key];
    if (value) {
      parts.push(cookieName + '=' + value);
    }
  });
  return parts.join('; ');
}

async function loadSessionTokenFromFlow() {
  try {
    saveBasics();
    const token = tokenInput.value.trim();
    const flowId = flowIdInput.value.trim();
    if (!token || !flowId) {
      alert('请先填写 Token 与 Flow ID');
      return;
    }
    const resp = await fetch(apiBase + '/session-token/flows/' + flowId, {
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    });
    const data = await resp.json();
    if (!resp.ok) {
      throw new Error(JSON.stringify(data));
    }
    const metadata = (data.flow && data.flow.metadata) || {};
    let sessionToken = metadata.session_token || '';
    if (!sessionToken) {
      sessionToken = buildSessionTokenFromPieces(metadata);
    }
    if (!sessionToken) {
      log('Flow metadata 中尚未写入 session_token，可稍后再试。');
      return;
    }
    apiSessionTokenInput.value = sessionToken;
    log('已从 Flow ' + flowId + ' 自动填充 SessionToken。');
  } catch (err) {
    log('加载 SessionToken 失败: ' + err);
  }
}

async function sendApiRequest() {
  try {
    saveBasics();
    const endpoint = getSelectedApiEndpoint();
    if (!endpoint) {
      alert('当前 Provider 暂无可调试的 API');
      return;
    }
    const token = tokenInput.value.trim();
    if (!token) {
      alert('请填写 API Token');
      return;
    }
    const params = parseQueryInput(apiQueryInput.value.trim());
    let finalPath = endpoint.path || '';
    let usedKeys = new Set();
    try {
      const applied = applyPathParams(finalPath, params);
      finalPath = applied.path;
      usedKeys = applied.usedKeys;
    } catch (err) {
      alert(err.message);
      return;
    }
    const queryString = buildQueryString(params, usedKeys);
    let url = apiBase + finalPath;
    if (queryString) {
      url += (url.includes('?') ? '&' : '?') + queryString;
    }
    let bodyPayload;
    const method = (endpoint.method || 'GET').toUpperCase();
    if (method !== 'GET' && method !== 'HEAD') {
      const rawBody = apiBodyInput.value.trim();
      if (rawBody) {
        try {
          bodyPayload = JSON.stringify(JSON.parse(rawBody));
        } catch (err) {
          alert('请求 Body 需要合法 JSON');
          return;
        }
      }
    }
    const headers = {
      'Authorization': 'Bearer ' + token,
    };
    if (bodyPayload) {
      headers['Content-Type'] = 'application/json';
    }
    const sessionToken = apiSessionTokenInput.value.trim();
    if (sessionToken) {
      headers['X-SessionToken'] = sessionToken;
    } else if (endpoint.requiresSessionToken) {
      log('提示：该接口需要 SessionToken，建议先点击“从 Flow 填充 SessionToken”。');
    }
    setApiResponse('请求中...');
    const resp = await fetch(url, {
      method,
      headers,
      body: bodyPayload,
    });
    const text = await resp.text();
    let formatted = text;
    try {
      formatted = JSON.stringify(JSON.parse(text), null, 2);
    } catch (_) {
      // ignore
    }
    setApiResponse('[HTTP ' + resp.status + '] ' + (resp.statusText || '') + '\n' + formatted);
  } catch (err) {
    setApiResponse('请求失败: ' + err);
  }
}

async function createFlow() {
  try {
    saveBasics();
    const token = tokenInput.value.trim();
    if (!token) {
      alert('请填写 API Token');
      return;
    }
    const providerAppCode = getProviderAppCode();
    if (!providerAppCode) {
      alert('请填写 Provider App Code');
      return;
    }
    const state = stateInput.value.trim() || randomState();
    stateInput.value = state;
    const metadataRaw = metadataInput.value.trim();
    let metadata = {};
    if (metadataRaw) {
      metadata = JSON.parse(metadataRaw);
    }
    const apiVersion = selectedVersion();
    if (apiVersion) {
      metadata.api_version = apiVersion;
    }
    metadata.reuse_session = reuseInput.checked ? "true" : "false";
    const payload = {
      provider_code: providerSelect.value,
      provider_app_code: providerAppCode,
      account_id: accountInput.value || 'acct_debug',
      tenant_uuid: tenantInput.value || 'tenant_debug',
      state: state,
      callback_url: callbackInput.value || defaults.callback,
      metadata: metadata,
    };
    log('创建 Flow 中...');
    const resp = await fetch(apiBase + '/session-token/flows', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token,
      },
      body: JSON.stringify(payload),
    });
    const data = await resp.json();
    if (!resp.ok) {
      throw new Error(JSON.stringify(data));
    }
    flowIdInput.value = data.flow.flow_id;
    authorizeInput.value = data.flow.authorize_url;
    updatePlaywrightCommand();
    log(JSON.stringify(data, null, 2));
  } catch (err) {
    log('创建失败: ' + err);
  }
}

window.createFlow = createFlow;

async function pollFlow() {
  try {
    saveBasics();
    const token = tokenInput.value.trim();
    const flowId = flowIdInput.value.trim();
    if (!token || !flowId) {
      alert('请填写 Token 与 Flow ID');
      return;
    }
    const resp = await fetch(apiBase + '/session-token/flows/' + flowId, {
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    });
    const data = await resp.json();
    if (!resp.ok) {
      throw new Error(JSON.stringify(data));
    }
    authorizeInput.value = data.flow.authorize_url || authorizeInput.value;
    log(JSON.stringify(data, null, 2));
  } catch (err) {
    log('查询失败: ' + err);
  }
}

window.pollFlow = pollFlow;

function copyAuthorizeURL() {
  const url = authorizeInput.value;
  if (!url) {
    alert('尚未生成 authorize_url');
    return;
  }
  navigator.clipboard.writeText(url).then(function() {
    log('authorize_url 已复制，可粘贴到浏览器调试');
  }).catch(function(err) {
    log('复制失败: ' + err);
  });
}

window.copyAuthorizeURL = copyAuthorizeURL;
window.loadSessionTokenFromFlow = loadSessionTokenFromFlow;
window.sendApiRequest = sendApiRequest;

async function refreshCallbackLogs(silent) {
  try {
    const resp = await fetch(apiBase + '/debug/callback');
    if (!resp.ok) {
      throw new Error('HTTP ' + resp.status);
    }
    const logs = await resp.json();
    callbackLogOutput.textContent = logs.length ? JSON.stringify(logs, null, 2) : '暂无回调...';
  } catch (err) {
    if (!silent) {
      log('获取 Mock 回调日志失败: ' + err);
    }
  }
}

async function clearCallbackLogs() {
  try {
    await fetch(apiBase + '/debug/callback', { method: 'DELETE' });
    callbackLogOutput.textContent = '已清空';
  } catch (err) {
    log('清空 Mock 回调日志失败: ' + err);
  }
}

function manualApplyPreset() {
  applyPreset(true, true);
}

flowIdInput.addEventListener('input', updatePlaywrightCommand);
providerSelect.addEventListener('change', () => {
  rebuildAppOptions(false);
  rebuildApiEndpoints(true);
  applyPreset(true, false);
  stateInput.value = randomState();
  resetApiResponse();
});
providerAppSelect.addEventListener('change', () => {
  toggleCustomAppInput();
  applyPreset(false, false);
});
if (apiEndpointSelect) {
  apiEndpointSelect.addEventListener('change', () => {
    updateApiInfo();
    resetApiResponse();
  });
}
apiVersionSelect.addEventListener('change', () => {
  localStorage.setItem('st_debug_api_version', apiVersionSelect.value);
});
reuseInput.addEventListener('change', () => {
  localStorage.setItem('st_debug_reuse', reuseInput.checked ? '1' : '0');
});
applyPresetBtn.addEventListener('click', manualApplyPreset);
refreshCallbackLogs(true);
setInterval(() => refreshCallbackLogs(true), 5000);
updatePlaywrightCommand();
</script>
</body>
</html>`))

func registerDebugPage(mux *http.ServeMux, log *logger.Logger, flowStore sessiontoken.FlowStore) {
	if os.Getenv("SESSIONTOKEN_DISABLE_DEBUG_PAGE") == "1" {
		return
	}
	if mux == nil || flowStore == nil {
		return
	}
	callbackStore := newCallbackLogStore(20)
	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		defaultCallback := firstNonEmptyEnv("POWERX_SESSION_TOKEN_CALLBACK_URL", "SESSIONTOKEN_DEBUG_CALLBACK_URL")
		if defaultCallback == "" {
			defaultCallback = "http://127.0.0.1:7070/debug/callback"
		}
		data := &debugPageData{
			DefaultToken:    resolveAPIToken(""),
			DefaultCallback: defaultCallback,
		}
		_ = debugPageTemplate.Execute(w, data)
	})
	mux.HandleFunc("/debug/callback", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			handleMockCallback(w, r, callbackStore, log)
		case http.MethodGet:
			respondJSON(w, http.StatusOK, callbackStore.list())
		case http.MethodDelete:
			callbackStore.clear()
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc(debugFlowMetadataPrefix, func(w http.ResponseWriter, r *http.Request) {
		handleDebugFlowRoute(w, r, flowStore, log)
	})
}

func handleMockCallback(w http.ResponseWriter, r *http.Request, store *callbackLogStore, log *logger.Logger) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "read body failed"})
		return
	}
	_ = r.Body.Close()
	formattedBody, flowID, provider := formatCallbackBody(body)
	record := callbackRecord{
		Timestamp:    time.Now(),
		Method:       r.Method,
		FlowID:       flowID,
		ProviderCode: provider,
		Query:        r.URL.RawQuery,
		Headers:      copyHeaders(r.Header),
		Body:         formattedBody,
	}
	store.append(record)
	if log != nil {
		log.InfoF("sessiontoken_debug: callback flow_id=%s provider=%s body=%s", flowID, provider, formattedBody)
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

func formatCallbackBody(body []byte) (string, string, string) {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return "", "", ""
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err == nil {
		flowID, _ := payload["flow_id"].(string)
		provider, _ := payload["provider_code"].(string)
		if pretty, err := json.MarshalIndent(payload, "", "  "); err == nil {
			return string(pretty), flowID, provider
		}
		return trimmed, flowID, provider
	}
	return trimmed, "", ""
}

func handleDebugFlowRoute(w http.ResponseWriter, r *http.Request, flowStore sessiontoken.FlowStore, log *logger.Logger) {
	flowID, segment, ok := parseDebugFlowPath(r.URL.Path)
	if !ok {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if segment != "metadata" {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		handleDebugFlowMetadataGet(w, r, flowStore, flowID)
	case http.MethodPost, http.MethodPut:
		handleDebugFlowMetadataUpdate(w, r, flowStore, flowID, log)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func parseDebugFlowPath(path string) (string, string, bool) {
	if !strings.HasPrefix(path, debugFlowMetadataPrefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(path, debugFlowMetadataPrefix)
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	flowID := strings.TrimSpace(parts[0])
	segment := strings.Trim(strings.TrimSpace(parts[1]), "/")
	if flowID == "" || segment == "" {
		return "", "", false
	}
	return flowID, segment, true
}

func handleDebugFlowMetadataGet(w http.ResponseWriter, r *http.Request, store sessiontoken.FlowStore, flowID string) {
	flow, err := store.Get(r.Context(), flowID)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "failed to load flow"
		if errors.Is(err, sessiontoken.ErrFlowNotFound) {
			status = http.StatusNotFound
			msg = "flow not found"
		}
		respondJSON(w, status, map[string]string{"error": msg})
		return
	}
	payload := map[string]any{
		"flow_id":        flow.FlowID,
		"provider_code":  flow.ProviderCode,
		"provider_app":   flow.ProviderAppCode,
		"tenant_uuid":    flow.TenantUUID,
		"account_id":     flow.AccountID,
		"state":          flow.State,
		"status":         flow.Status,
		"metadata":       cloneMetadata(flow.Metadata),
		"updated_at":     flow.UpdatedAt,
		"expires_at":     flow.ExpiresAt,
		"callback_url":   flow.CallbackURL,
		"retry_attempts": flow.RetryAttempts,
	}
	respondJSON(w, http.StatusOK, payload)
}

func handleDebugFlowMetadataUpdate(w http.ResponseWriter, r *http.Request, store sessiontoken.FlowStore, flowID string, log *logger.Logger) {
	if r.Body == nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "request body is empty"})
		return
	}
	defer r.Body.Close()
	var raw map[string]any
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid json: %v", err)})
		return
	}
	updates, err := normalizeMetadataPayload(raw)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	flow, err := store.Get(r.Context(), flowID)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "failed to load flow"
		if errors.Is(err, sessiontoken.ErrFlowNotFound) {
			status = http.StatusNotFound
			msg = "flow not found"
		}
		respondJSON(w, status, map[string]string{"error": msg})
		return
	}
	if len(updates) == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "metadata payload is empty"})
		return
	}
	changed := applyMetadataUpdates(flow, updates)
	if !changed {
		respondJSON(w, http.StatusOK, map[string]any{
			"flow_id":  flow.FlowID,
			"metadata": cloneMetadata(flow.Metadata),
		})
		return
	}
	if err := persistDebugFlow(r.Context(), store, flow); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist metadata"})
		return
	}
	if log != nil {
		log.InfoF("sessiontoken_debug: metadata_update flow_id=%s provider=%s keys=%s", flow.FlowID, flow.ProviderCode, strings.Join(sortedKeys(updates), ","))
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"flow_id":    flow.FlowID,
		"status":     flow.Status,
		"metadata":   cloneMetadata(flow.Metadata),
		"updated_at": flow.UpdatedAt,
	})
}

func normalizeMetadataPayload(raw map[string]any) (map[string]string, error) {
	if len(raw) == 0 {
		return nil, errors.New("metadata payload is empty")
	}
	normalized := make(map[string]string, len(raw))
	for k, v := range raw {
		key := strings.TrimSpace(k)
		if key == "" {
			continue
		}
		switch val := v.(type) {
		case nil:
			normalized[key] = ""
		case string:
			normalized[key] = val
		default:
			normalized[key] = fmt.Sprint(val)
		}
	}
	if len(normalized) == 0 {
		return nil, errors.New("metadata payload is empty")
	}
	return normalized, nil
}

func applyMetadataUpdates(flow *sessiontoken.Flow, updates map[string]string) bool {
	if flow == nil || len(updates) == 0 {
		return false
	}
	changed := false
	for key, value := range updates {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			continue
		}
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			if flow.Metadata != nil {
				if _, exists := flow.Metadata[normalizedKey]; exists {
					delete(flow.Metadata, normalizedKey)
					changed = true
				}
			}
			continue
		}
		if flow.Metadata == nil {
			flow.Metadata = make(map[string]string)
		}
		if existing, ok := flow.Metadata[normalizedKey]; !ok || existing != trimmed {
			flow.Metadata[normalizedKey] = trimmed
			changed = true
		}
	}
	if changed {
		flow.UpdatedAt = time.Now().UTC()
	}
	return changed
}

func persistDebugFlow(ctx context.Context, store sessiontoken.FlowStore, flow *sessiontoken.Flow) error {
	if store == nil || flow == nil {
		return errors.New("sessiontoken: flow store is nil")
	}
	now := time.Now().UTC()
	ttl := flow.ExpiresAt.Sub(now)
	if ttl <= 0 {
		ttl = time.Second
	}
	return store.Update(ctx, flow, ttl)
}

func sortedKeys(m map[string]string) []string {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func cloneMetadata(meta map[string]string) map[string]string {
	if len(meta) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(meta))
	for k, v := range meta {
		cloned[k] = v
	}
	return cloned
}

func copyHeaders(src http.Header) map[string]string {
	if len(src) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(src))
	for k, vals := range src {
		cloned[k] = strings.Join(vals, ", ")
	}
	return cloned
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
