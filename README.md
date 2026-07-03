# DuckDuckGo MCP Server

A powerful web search MCP server using DuckDuckGo. No API key required.

## Features

- **3 Powerful Tools** - Search, crawl, and research with AI-powered ranking
- **Parallel Crawling** - Fetch multiple pages simultaneously
- **Smart Ranking** - Research tool ranks results by relevance to your question
- **News Search** - Filter results to news articles only
- **No API Keys** - Works out of the box, no rate limits
- **Cross-platform** - Linux, macOS, Windows

## Tools

### 1. `search`

Quick web search. Returns titles, URLs, and snippets.

```json
{
  "name": "search",
  "arguments": {
    "query": "latest AI news 2026",
    "limit": 10
  }
}
```

**Parameters:**
| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `query` | string | required | - | Search query |
| `limit` | number | 10 | 20 | Number of results |
| `news` | string | - | - | Set to "true" for news only |

---

### 2. `search_and_crawl`

Search + crawl all result pages in parallel. Get full content from each source.

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

**Parameters:**
| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `query` | string | required | - | Search query |
| `count` | number | 5 | 10 | Number of results to crawl |
| `maxContentLength` | number | 3000 | 10000 | Max characters per page |

---

### 3. `research`

Best for answering questions. Searches, crawls in parallel, then **ranks results by relevance** using:

| Scoring Factor | Weight | Description |
|----------------|--------|-------------|
| **Keywords** | 30% | How many question keywords appear in content |
| **Content Quality** | 25% | Length, structure, no paywalls |
| **Source Authority** | 20% | Domain reputation (Wikipedia=10, Reuters=9, etc.) |
| **Relevance** | 25% | Question type matching (how-to, what, why, etc.) |

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

## Installation

### NPM (Recommended)

```bash
npm install -g mcp-duckduckgo
```

Then add to your `~/.mcp.json`:

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

### From Source

## Usage

### Build

```bash
go build -o duckduckgo-mcp .
```

### Run

```bash
./duckduckgo-mcp
```

### MCP Configuration

Add to your Claude Code `.mcp.json`:

```json
{
  "mcpServers": {
    "duckduckgo-mcp": {
      "command": "/path/to/duckduckgo-mcp",
      "args": []
    }
  }
}
```

## Domain Authority Scores

High-authority domains get better rankings in research:

| Domain | Score |
|--------|-------|
| Wikipedia | 10 |
| Reuters, AP News, BBC | 9 |
| NYT, WSJ, Bloomberg | 8 |
| .gov sites | 8 |
| .edu sites | 7 |
| Stack Overflow | 7 |
| MDN, W3C | 8 |
| TechCrunch, Wired, Ars | 6 |
| Medium | 4 |
| Reddit | 3 |

## Comparison with Alternatives

| Feature | DuckDuckGo MCP | Tavily | Google Search API |
|---------|---------------|--------|-------------------|
| **Free** | Yes | No | No |
| **API Key** | Not needed | Required | Required |
| **Rate Limits** | None | Yes | Yes |
| **CAPTCHA Issues** | No | No | Yes |
| **Parallel Crawl** | Yes | Yes | No |
| **Result Ranking** | Yes | Yes | No |

## License

MIT
