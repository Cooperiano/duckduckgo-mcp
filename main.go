package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// DuckDuckGo search result
type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// Search web using DuckDuckGo HTML
func searchWeb(query string, limit int) ([]SearchResult, error) {
	url := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", strings.ReplaceAll(query, " ", "+"))

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := make([]byte, 0, 1024*1024)
	tmp := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(tmp)
		buf = append(buf, tmp[:n]...)
		if err != nil {
			break
		}
	}

	html := string(buf)
	return parseResults(html, limit), nil
}

// Parse DuckDuckGo HTML results
func parseResults(html string, limit int) []SearchResult {
	results := []SearchResult{}

	// Find result blocks
	re := regexp.MustCompile(`(?s)<a rel="nofollow" class="result__a" href="([^"]*)">(.*?)</a>.*?<a class="result__snippet"[^>]*>(.*?)</a>`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(results) >= limit {
			break
		}

		title := cleanHTML(match[2])
		url := strings.ReplaceAll(match[1], "//duckduckgo.com/l/?uddg=", "")
		url = strings.Split(url, "&")[0]
		snippet := cleanHTML(match[3])

		if title != "" && url != "" {
			results = append(results, SearchResult{
				Title:   title,
				URL:     url,
				Snippet: snippet,
			})
		}
	}

	return results
}

// Clean HTML tags
func cleanHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func main() {
	s := server.NewMCPServer("darkdark-search", "1.0.0")

	// Add search tool
	searchTool := mcp.NewTool("search_web",
		mcp.WithDescription("Search the web using DuckDuckGo"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results (default: 10)"),
		),
	)

	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.Params.Arguments["query"].(string)
		limit := 10
		if l, ok := request.Params.Arguments["limit"].(float64); ok {
			limit = int(l)
		}

		results, err := searchWeb(query, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		// Format results
		var output strings.Builder
		output.WriteString(fmt.Sprintf("Found %d results:\n\n", len(results)))
		for i, r := range results {
			output.WriteString(fmt.Sprintf("%d. %s\n", i+1, r.Title))
			output.WriteString(fmt.Sprintf("   URL: %s\n", r.URL))
			output.WriteString(fmt.Sprintf("   %s\n\n", r.Snippet))
		}

		return mcp.NewToolResultText(output.String()), nil
	})

	// Start stdio server
	stdioServer := server.NewStdioServer(s)
	log.Println("DarkDark Search MCP server starting on stdio")
	if err := stdioServer.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
