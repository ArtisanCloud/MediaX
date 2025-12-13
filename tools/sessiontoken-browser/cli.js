#!/usr/bin/env node
const readline = require('readline');
const fetch = require('node-fetch');
const { chromium } = require('playwright');

const WATCHED_COOKIES = [
  'SESSIONID',
  'JOID',
  'osd',
  'q_c1',
  'd_c0',
  'unlock_ticket',
  'z_c0'
];

async function main() {
  const args = parseArgs(process.argv.slice(2));
  if (!args.flow) {
    throw new Error('缺少参数 --flow <flow_id>');
  }
  const base = args.base || process.env.POWERX_SESSION_TOKEN_BASE_URL || process.env.SESSIONTOKEN_BASE_URL || 'http://127.0.0.1:7070';
  const apiToken = args.token || process.env.SESSIONTOKEN_API_TOKEN || process.env.POWERX_SESSION_TOKEN_API_TOKEN || 'dev-session-token';
  const flow = await fetchFlow(base, args.flow, apiToken);
  const authorizeUrl = flow.authorize_url || flow.authorizeUrl;
  if (!authorizeUrl) {
    throw new Error('Flow 中缺少 authorize_url，可尝试重新创建 Flow');
  }
  console.log('================ Playwright 调试器 ================');
  console.log(`目标 Flow: ${flow.flow_id}`);
  console.log(`授权地址: ${authorizeUrl}`);
  console.log('\n请按以下步骤执行:');
  console.log('1. 将会打开一个 Chromium 窗口，请在其中完成知乎登录（与插件行为一致）。');
  console.log('2. 登录成功后不要关闭窗口，返回终端并按 Enter，我们会自动抓取 Cookie 并写回 Flow。');
  console.log('===================================================\n');

  let browser;
  try {
    browser = await chromium.launch({ headless: args.headless || false });
    const context = await browser.newContext();
    const page = await context.newPage();
    await page.goto(authorizeUrl, { waitUntil: 'load' });
    await waitForEnter('完成登录并确保知乎页面已显示登录状态后，按 Enter 继续...');
    const cookies = await context.cookies();
    const userAgent = await page.evaluate(() => navigator.userAgent);
    const metadata = buildMetadata(cookies, userAgent);
    await browser.close();
    browser = null;
    console.log('\n抓取到以下核心字段:');
    Object.entries(metadata)
      .filter(([key]) => key === 'session_token' || key.startsWith('cookie_'))
      .forEach(([key, value]) => {
        console.log(`  ${key}: ${value ? value.slice(0, 60) + (value.length > 60 ? '...' : '') : ''}`);
      });
    await postMetadata(base, args.flow, metadata);
    console.log('\n已写入 Flow metadata，正在查询最新状态...');
    const updated = await fetchFlow(base, args.flow, apiToken);
    console.log(`Flow 状态: ${updated.status}`);
    await printLatestCallback(base, args.flow);
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}

function parseArgs(argv) {
  const result = {};
  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (!arg.startsWith('--')) {
      continue;
    }
    const key = arg.slice(2);
    const value = argv[i + 1] && !argv[i + 1].startsWith('--') ? argv[i + 1] : undefined;
    switch (key) {
      case 'flow':
        result.flow = value;
        if (value) i += 1;
        break;
      case 'base':
        result.base = value;
        if (value) i += 1;
        break;
      case 'api-token':
        result.token = value;
        if (value) i += 1;
        break;
      case 'headless':
        result.headless = true;
        i -= 1;
        break;
      default:
        console.warn(`未知参数: --${key}`);
        if (value) i += 1;
    }
  }
  return result;
}

async function fetchFlow(base, flowId, apiToken) {
  const res = await fetch(`${base}/session-token/flows/${flowId}`, {
    headers: {
      'Authorization': `Bearer ${apiToken}`,
      'Accept': 'application/json'
    }
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`获取 Flow 失败 (${res.status}): ${text}`);
  }
  const data = await res.json();
  if (!data.flow) {
    throw new Error('响应中缺少 flow 字段');
  }
  return data.flow;
}

function buildMetadata(cookies, userAgent) {
  const zhihuCookies = cookies.filter(cookie => /zhihu\.com$/i.test(cookie.domain || ''));
  if (!zhihuCookies.length) {
    throw new Error('未捕获到 zhihu.com 域下的 Cookie，请确认已完成登录');
  }
  const cookieMap = new Map();
  zhihuCookies.forEach(cookie => {
    cookieMap.set(cookie.name, cookie.value);
  });
  const sessionToken = zhihuCookies.map(cookie => `${cookie.name}=${cookie.value}`).join('; ');
  const metadata = {
    session_token: sessionToken,
    credentials_note: userAgent || 'Playwright',
    credentials_expires_hint: '',
    reuse_session: 'true'
  };
  WATCHED_COOKIES.forEach(name => {
    const key = `cookie_${name.toLowerCase()}`;
    metadata[key] = cookieMap.get(name) || '';
  });
  return metadata;
}

async function postMetadata(base, flowId, metadata) {
  const res = await fetch(`${base}/debug/flows/${flowId}/metadata`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(metadata)
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`写入 metadata 失败 (${res.status}): ${text}`);
  }
  await res.json().catch(() => null);
}

async function printLatestCallback(base, flowId) {
  try {
    const res = await fetch(`${base}/debug/callback`);
    if (!res.ok) {
      console.warn('获取回调日志失败');
      return;
    }
    const logs = await res.json();
    if (!Array.isArray(logs) || logs.length === 0) {
      console.log('暂未看到 /debug/callback 记录，可稍后刷新调试页。');
      return;
    }
    const record = logs.find(item => item.flow_id === flowId) || logs[0];
    console.log('\n最近的 /debug/callback 记录:');
    console.log(JSON.stringify(record, null, 2));
  } catch (err) {
    console.warn('读取回调失败: ', err.message);
  }
}

function waitForEnter(promptMessage) {
  return new Promise(resolve => {
    const rl = readline.createInterface({ input: process.stdin, output: process.stdout });
    rl.question(`${promptMessage}\n`, () => {
      rl.close();
      resolve();
    });
  });
}

main().catch(err => {
  console.error(`执行失败: ${err.message}`);
  process.exitCode = 1;
});
