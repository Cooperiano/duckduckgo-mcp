# @cooperiano/darkdark-search

Free web search MCP server using DuckDuckGo.

## Installation

```bash
npm install -g @cooperiano/darkdark-search
```

## Usage

### As MCP Server

Add to your Claude Code configuration:

```json
{
  "mcpServers": {
    "darkdark-search": {
      "command": "npx",
      "args": ["-y", "@cooperiano/darkdark-search"]
    }
  }
}
```

### MCP Tools

- `search_web` - Search the web using DuckDuckGo
  - `query` (required): Search query
  - `limit` (optional): Maximum number of results (default: 10)

## Examples

```bash
# In Claude Code
> Use search_web tool to search for "Go programming language"
```

## License

MIT
