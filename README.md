# DarkDark Search MCP Server

A free web search MCP server using DuckDuckGo.

## Features

- Search the web using DuckDuckGo (no API key required)
- MCP protocol support
- HTTP server with SSE transport

## Usage

### Build

```bash
cd /home/julian/projects/darkdark
go build -o darkdark-server .
```

### Run

```bash
./darkdark-server
# Server starts on port 8080
```

### MCP Configuration

Add to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "darkdark-search": {
      "command": "/home/julian/projects/darkdark/darkdark-server",
      "args": []
    }
  }
}
```

### MCP Tools

- `search_web` - Search the web
  - `query` (required): Search query
  - `limit` (optional): Maximum results (default: 10)
