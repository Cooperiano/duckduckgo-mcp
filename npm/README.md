# @cooperiano/duckduckgo-mcp

Free web search MCP server using DuckDuckGo with AI-powered research capabilities.

## Features

- **3 Powerful Tools** - Search, crawl, and research with AI-powered ranking
- **Parallel Crawling** - Fetch multiple pages simultaneously
- **Smart Ranking** - Research tool ranks results by relevance to your question
- **No API Keys** - Works out of the box, no rate limits

## Installation

```bash
npm install -g @cooperiano/duckduckgo-mcp
```

## Usage

### As MCP Server

Add to your Claude Code configuration:

```json
{
  "mcpServers": {
    "duckduckgo-mcp": {
      "command": "npx",
      "args": ["-y", "@cooperiano/duckduckgo-mcp"]
    }
  }
}
```

### MCP Tools

#### `search`

Quick web search. Returns titles, URLs, and snippets.

```json
{
  "name": "search",
  "arguments": {
    "query": "latest AI news 2026",
    "limit": 10,
    "news": "true"
  }
}
```

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `query` | string | required | - | Search query |
| `limit` | number | 10 | 20 | Number of results |
| `news` | string | - | - | Set to "true" for news only |

---

#### `search_and_crawl`

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

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `query` | string | required | - | Search query |
| `count` | number | 5 | 10 | Number of results to crawl |
| `maxContentLength` | number | 3000 | 10000 | Max characters per page |

---

#### `research`

Best for answering questions. Searches, crawls in parallel, then **ranks results by relevance** using AI-powered scoring (keywords, content quality, source authority, relevance).

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

## Examples

```bash
# In Claude Code
> Use search tool to find "latest AI news"
> Use search_and_crawl to research "React vs Vue"
> Use research tool to answer "How do quantum computers work?"
```

## License

MIT
