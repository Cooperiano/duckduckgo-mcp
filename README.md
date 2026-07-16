# DuckDuckGo MCP Server · DuckDuckGo MCP 服务器

[![npm version](https://img.shields.io/npm/v/mcp-duckduckgo)](https://www.npmjs.com/package/mcp-duckduckgo)
[![license](https://img.shields.io/npm/l/mcp-duckduckgo)](LICENSE)
[![platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-blue)]()

A free, no-API-key web search MCP server with parallel crawling, AI-powered research ranking, and auto proxy detection. Works out of the box — no signup, no rate limits, no CAPTCHAs.

免费的、无需 API Key 的网页搜索 MCP 服务器，支持并行爬取、AI 驱动的研究排序和代理自动检测。开箱即用——无需注册、无速率限制、无验证码。

---

## Features · 功能特性

|  |  |
|---------|---------|
| **4 tools · 4 个工具** | `search`, `search_and_crawl`, `research`, `fetch` |
| **Instant Answer · 即时答案** | DDG API integration for definitions, Wikipedia abstracts, and encyclopedia lookups · 集成 DDG API，返回定义、维基百科摘要和百科查询 |
| **Proxy auto-detect · 代理自动检测** | Automatically discovers local HTTP proxies (Clash, V2Ray, etc.) · 自动发现本地 HTTP 代理（Clash、V2Ray 等） |
| **Region filtering · 区域过滤** | Localized results via locale codes · 通过区域代码获取本地化结果（`cn-zh`、`us-en`、`jp-jp` 等） |
| **Time filtering · 时间过滤** | Restrict results to past day/week/month/year · 将结果限制在过去一天/一周/一月/一年内 |
| **Safe search · 安全搜索** | `off`, `moderate`, `strict` |
| **Parallel crawling · 并行爬取** | Fetches multiple pages concurrently · 同时抓取多个页面 |
| **AI-powered ranking · AI 驱动排序** | `research` tool scores by keyword, quality, authority, and relevance · `research` 工具按关键词、质量、权威性、相关性评分排序 |
| **Zero config · 零配置** | No API keys, no signup, no rate limits · 无需 API Key、无需注册、无速率限制 |
| **Cross-platform · 跨平台** | Pre-built binaries for Linux, macOS, Windows (amd64 + arm64) · 预编译二进制文件，支持 Linux、macOS、Windows |

---

## Quick Start · 快速开始

```bash
npm install -g mcp-duckduckgo
```

Add to your Claude Code `.mcp.json` · 添加到 Claude Code 的 `.mcp.json`：

```json
{
  "mcpServers": {
    "duckduckgo-mcp": {
      "command": "npx",
      "args": ["-y", "mcp-duckduckgo"]
    }
  }
}
```

Or use a pre-built binary · 或使用预编译二进制文件：

```json
{
  "mcpServers": {
    "duckduckgo-mcp": {
      "command": "/path/to/duckduckgo-mcp-linux-amd64",
      "args": []
    }
  }
}
```

> **For users in China · 中国用户注意：** The server auto-detects local proxies on ports 7892, 7890, 7891, 1080, 1087. If your proxy runs on a different port, set `HTTPS_PROXY` or `HTTP_PROXY`. · 服务器会自动检测本地代理端口（7892、7890、7891、1080、1087）。如果代理运行在其他端口，请设置 `HTTPS_PROXY` 或 `HTTP_PROXY` 环境变量。

---

## Tools · 工具

### 1. `search`

Basic web search. Returns titles, URLs, and snippets. Includes Instant Answer from DDG API (definitions, Wikipedia abstracts).
基础网页搜索。返回标题、URL 和摘要。包含 DDG API 的即时答案（定义、维基百科摘要）。

```json
{
  "name": "search",
  "arguments": {
    "query": "latest AI news 2026",
    "limit": 10
  }
}
```

| Parameter · 参数 | Type · 类型 | Default · 默认值 | Max · 最大 | Description · 说明 |
|-----------|------|---------|-----|-------------|
| `query` | string | required · 必填 | — | Search query · 搜索关键词 |
| `limit` | number | 10 | 20 | Number of results · 结果数量 |
| `type` | string | `"text"` | — | `"text"`, `"news"`, or `"image"` |
| `region` | string | `"cn-zh"` | — | Locale code. Use `"us-en"` for US, `""` for global · 区域代码。美国用 `"us-en"`，全局用 `""` |
| `time` | string | `""` | — | `"d"` (day · 天), `"w"` (week · 周), `"m"` (month · 月), `"y"` (year · 年) |
| `safe` | string | `"moderate"` | — | `"off"`, `"moderate"`, or `"strict"` |
| `instant_answer` | boolean | `true` | — | Fetch Instant Answer from DDG API · 从 DDG API 获取即时答案 |

---

### 2. `search_and_crawl`

Search + crawl all result pages in parallel. Returns full extracted content from each source (navigation, ads, scripts stripped).
搜索 + 并行爬取所有结果页面。返回每个来源的完整提取内容（去除导航、广告、脚本）。

```json
{
  "name": "search_and_crawl",
  "arguments": {
    "query": "best JavaScript frameworks 2026",
    "count": 5,
    "maxContentLength": 3000
  }
}
```

| Parameter · 参数 | Type · 类型 | Default · 默认值 | Max · 最大 | Description · 说明 |
|-----------|------|---------|-----|-------------|
| `query` | string | required · 必填 | — | Search query · 搜索关键词 |
| `count` | number | 5 | 10 | Number of results to crawl · 要爬取的结果数量 |
| `maxContentLength` | number | 3000 | 10000 | Max characters per page · 每页最大字符数 |
| `type` | string | `"text"` | — | `"text"`, `"news"`, or `"image"` |
| `region` | string | `"cn-zh"` | — | Locale code · 区域代码 |
| `time` | string | `""` | — | `"d"`, `"w"`, `"m"`, `"y"` |
| `safe` | string | `"moderate"` | — | `"off"`, `"moderate"`, `"strict"` |

---

### 3. `research`

Best for answering questions. Searches, crawls in parallel, then **ranks results by relevance** using a 4-factor scoring model:
最适合回答问题。搜索、并行爬取，然后使用 4 因子评分模型**按相关性排序结果**：

| Scoring Factor · 评分因子 | Weight · 权重 | Description · 说明 |
|----------------|--------|-------------|
| **Keywords · 关键词** | 30% | Keyword match count in content · 内容中关键词匹配数量 |
| **Content Quality · 内容质量** | 25% | Word count, structure, penalty for paywalls · 字数、结构、付费墙扣分 |
| **Source Authority · 来源权威** | 20% | Domain reputation score · 域名声誉评分（见下表） |
| **Relevance · 相关性** | 25% | Question-type matching (how-to, what-is, why, comparison) · 问题类型匹配（操作指南、定义、原因、对比） |

```json
{
  "name": "research",
  "arguments": {
    "question": "How does Starlink help Ukraine in the war?",
    "count": 5,
    "maxContentLength": 3000
  }
}
```

| Parameter · 参数 | Type · 类型 | Default · 默认值 | Max · 最大 | Description · 说明 |
|-----------|------|---------|-----|-------------|
| `question` | string | required · 必填 | — | Research question · 研究问题 |
| `count` | number | 5 | 10 | Number of sources to analyze · 分析的来源数量 |
| `maxContentLength` | number | 3000 | 10000 | Max characters per page · 每页最大字符数 |
| `region` | string | `"cn-zh"` | — | Locale code · 区域代码 |
| `time` | string | `""` | — | `"d"`, `"w"`, `"m"`, `"y"` |
| `safe` | string | `"moderate"` | — | `"off"`, `"moderate"`, `"strict"` |
| `instant_answer` | boolean | `true` | — | Fetch Instant Answer from DDG API · 从 DDG API 获取即时答案 |

---

### 4. `fetch`

Fetch and extract clean text from a single URL. Strips navigation, ads, scripts, and boilerplate — useful for reading a specific page in depth.
抓取并提取单个 URL 的纯文本内容。去除导航、广告、脚本和模板内容——适合深度阅读特定页面。

```json
{
  "name": "fetch",
  "arguments": {
    "url": "https://example.com/article",
    "maxContentLength": 5000
  }
}
```

| Parameter · 参数 | Type · 类型 | Default · 默认值 | Max · 最大 | Description · 说明 |
|-----------|------|---------|-----|-------------|
| `url` | string | required · 必填 | — | URL to fetch · 要抓取的 URL |
| `maxContentLength` | number | 3000 | 10000 | Max characters to return · 返回的最大字符数 |

---

## Installation · 安装

### NPM (Recommended · 推荐)

```bash
npm install -g mcp-duckduckgo
```

The package automatically downloads the correct pre-built binary for your platform. Requires Node.js ≥ 14.
安装时自动下载适配当前平台的预编译二进制文件。需要 Node.js ≥ 14。

### Pre-built Binaries · 预编译二进制文件

Download from · 下载地址：[GitHub Releases](https://github.com/Cooperiano/duckduckgo-mcp/releases)

| Platform · 平台 | Architecture · 架构 | Binary · 文件名 |
|----------|-------------|--------|
| Linux | amd64 | `duckduckgo-mcp-linux-amd64` |
| Linux | arm64 | `duckduckgo-mcp-linux-arm64` |
| macOS | amd64 | `duckduckgo-mcp-darwin-amd64` |
| macOS | arm64 | `duckduckgo-mcp-darwin-arm64` |
| Windows | amd64 | `duckduckgo-mcp-windows-amd64.exe` |

### From Source · 从源码构建

Requires Go ≥ 1.21 · 需要 Go ≥ 1.21。

```bash
git clone https://github.com/Cooperiano/duckduckgo-mcp.git
cd duckduckgo-mcp
make build        # current platform · 当前平台
make build-all    # all platforms → dist/ · 所有平台
```

---

## Proxy Auto-Detection · 代理自动检测

The server automatically discovers HTTP proxies in this order · 服务器按以下顺序自动发现 HTTP 代理：

1. **Environment variables · 环境变量** — `HTTPS_PROXY`, `https_proxy`, `HTTP_PROXY`, `http_proxy`, `ALL_PROXY`, `all_proxy`
2. **Common local ports · 常用本地端口** — tests actual proxy connectivity on · 测试实际代理连通性：7892, 7890, 7891, 1080, 1087

This means tools like Clash Verge, V2Ray, and Shadowsocks are picked up automatically. No configuration needed.
这意味着 Clash Verge、V2Ray、Shadowsocks 等工具会被自动识别，无需任何配置。

To force a specific proxy · 强制指定代理：

```bash
HTTPS_PROXY=http://127.0.0.1:7892 npx -y mcp-duckduckgo
```

---

## Domain Authority Scores · 域名权威评分

High-authority domains get better rankings in the `research` tool · 高权威域名在 `research` 工具中获得更高排名：

| Domain · 域名 | Score · 评分 |
|--------|-------|
| Wikipedia | 10 |
| Reuters, AP News, BBC | 9 |
| NYT, WSJ, Bloomberg | 8 |
| .gov sites · 政府网站 | 8 |
| MDN, W3C | 8 |
| .edu sites · 教育网站 | 7 |
| Stack Overflow | 7 |
| GitHub | 6 |
| TechCrunch, Wired, Ars Technica | 6 |
| Dev.to | 5 |
| Medium | 4 |
| Reddit | 3 |
| Unknown domains · 未知域名 | 2 |

---

## Comparison · 对比

| Feature · 功能 | DuckDuckGo MCP | Tavily | Google Search API | Brave Search API |
|---------|---------------|--------|-------------------|------------------|
| **Free · 免费** | ✅ Yes · 是 | ❌ No · 否 | ❌ No · 否 | ❌ No · 否 (limited free tier · 有限免费额度) |
| **API Key · API 密钥** | Not needed · 不需要 | Required · 需要 | Required · 需要 | Required · 需要 |
| **Rate Limits · 速率限制** | None · 无 | Yes · 有 | Yes · 有 | Yes · 有 |
| **CAPTCHA · 验证码** | No · 无 | No · 无 | Yes · 有 | No · 无 |
| **Parallel Crawl · 并行爬取** | ✅ Yes · 是 | ✅ Yes · 是 | ❌ No · 否 | ❌ No · 否 |
| **AI Result Ranking · AI 排序** | ✅ Yes · 是 | ✅ Yes · 是 | ❌ No · 否 | ❌ No · 否 |
| **Proxy Auto-Detect · 代理自动检测** | ✅ Yes · 是 | ❌ No · 否 | ❌ No · 否 | ❌ No · 否 |
| **Instant Answer · 即时答案** | ✅ Yes · 是 | ❌ No · 否 | ❌ No · 否 | ❌ No · 否 |
| **Single-URL Fetch · 单页抓取** | ✅ Yes · 是 | ❌ No · 否 | ❌ No · 否 | ❌ No · 否 |

---

## License · 许可证

MIT © [Cooperiano](https://github.com/Cooperiano)
