package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)


// Pre-compiled regexes for content extraction
var (
	reScript = regexp.MustCompile(`(?s)<script[^>]*>.*?</script>`)
	reStyle  = regexp.MustCompile(`(?s)<style[^>]*>.*?</style>`)
	reNav    = regexp.MustCompile(`(?s)<nav[^>]*>.*?</nav>`)
	reFooter = regexp.MustCompile(`(?s)<footer[^>]*>.*?</footer>`)
	reHeader = regexp.MustCompile(`(?s)<header[^>]*>.*?</header>`)
	reAside  = regexp.MustCompile(`(?s)<aside[^>]*>.*?</aside>`)
)

// Domain authority scores for research ranking
var domainAuthority = map[string]float64{
	"wikipedia.org":         10.0,
	"reuters.com":           9.0,
	"apnews.com":            9.0,
	"bbc.com":               9.0,
	"nytimes.com":           8.0,
	"wsj.com":               8.0,
	"bloomberg.com":         8.0,
	"gov":                   8.0,
	"edu":                   7.0,
	"techcrunch.com":        6.0,
	"wired.com":             6.0,
	"arstechnica.com":       6.0,
	"medium.com":            4.0,
	"reddit.com":            3.0,
	"stackoverflow.com":     7.0,
	"github.com":            6.0,
	"dev.to":                5.0,
	"mdn.mozilla.org":       8.0,
	"developer.mozilla.org": 8.0,
	"w3.org":                8.0,
}

// testProxyConnectivity verifies an HTTP proxy actually routes traffic
func testProxyConnectivity(proxyURL *url.URL) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

// Get proxy URL from environment variables, with fallback to common proxy ports
func getProxyURL() *url.URL {
	// Check environment variables first
	for _, envVar := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy", "ALL_PROXY", "all_proxy"} {
		if proxyStr := os.Getenv(envVar); proxyStr != "" {
			if u, err := url.Parse(proxyStr); err == nil {
				return u
			}
		}
	}

	// Fallback: try common proxy ports (test actual HTTP proxy connectivity)
	commonPorts := []string{"7892", "7890", "7891", "1080", "1087"}
	for _, port := range commonPorts {
		addr := fmt.Sprintf("127.0.0.1:%s", port)
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s", addr))
		if err != nil {
			continue
		}
		if testProxyConnectivity(proxyURL) {
			log.Printf("Auto-detected proxy: %s", proxyURL.String())
			return proxyURL
		}
	}

	return nil
}

// DuckDuckGo search result
type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// Crawled result with full content
type CrawledResult struct {
	SearchResult
	Content string
	WordCount int
}

// Research result with scoring
type ResearchResult struct {
	CrawledResult
	Score       float64
	KeywordScore float64
	QualityScore float64
	AuthorityScore float64
	RelevanceScore float64
}

// Search options
type SearchOptions struct {
	Query    string
	Limit    int
	NewsOnly bool
	Region   string
	Time     string
}

// Search web using DuckDuckGo HTML
func searchWeb(opts SearchOptions) ([]SearchResult, error) {
	baseURL := "https://html.duckduckgo.com/html/"
	query := strings.ReplaceAll(opts.Query, " ", "+")

	searchURL := fmt.Sprintf("%s?q=%s", baseURL, query)

	// Add news filter if requested
	if opts.NewsOnly {
		searchURL += "&iar=news"
	}

	// Add region filter
	if opts.Region != "" {
		searchURL += "&kl=" + opts.Region
	}

	// Add time filter
	if opts.Time != "" {
		searchURL += "&df=" + opts.Time
	}

	// Create HTTP client with proxy support
	client := &http.Client{Timeout: 30 * time.Second}

	// Check for proxy environment variables
	proxyURL := getProxyURL()
	if proxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		log.Printf("Using proxy: %s", proxyURL.String())
	}

	req, err := http.NewRequest("GET", searchURL, nil)
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
	return parseResults(html, opts.Limit), nil
}

// Crawl a single page and extract main content
func crawlPage(pageURL string, maxLength int) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	proxyURL := getProxyURL()
	if proxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	buf := make([]byte, 0, maxLength*2)
	tmp := make([]byte, 4096)
	for len(buf) < maxLength*2 {
		n, err := resp.Body.Read(tmp)
		buf = append(buf, tmp[:n]...)
		if err != nil || len(buf) >= maxLength*2 {
			break
		}
	}

	content := extractMainContent(string(buf))
	if len(content) > maxLength {
		content = content[:maxLength]
	}

	return content, nil
}

// Extract main content from HTML (remove nav, ads, etc)
func extractMainContent(html string) string {
	// Remove non-content elements
	html = reScript.ReplaceAllString(html, "")
	html = reStyle.ReplaceAllString(html, "")
	html = reNav.ReplaceAllString(html, "")
	html = reFooter.ReplaceAllString(html, "")
	html = reHeader.ReplaceAllString(html, "")
	html = reAside.ReplaceAllString(html, "")

	// Remove common ad/class patterns
	adRe := regexp.MustCompile(`(?s)<[^>]*(class|id)="[^"]*(ad|advertisement|sidebar|navigation|menu|footer|header)[^"]*"[^>]*>.*?</[^>]+>`)
	html = adRe.ReplaceAllString(html, "")

	// Remove all HTML tags
	tagRe := regexp.MustCompile(`<[^>]*>`)
	html = tagRe.ReplaceAllString(html, "")

	// Clean up whitespace
	html = strings.Join(strings.Fields(html), " ")
	return strings.TrimSpace(html)
}

// Search and crawl multiple pages in parallel
func searchAndCrawl(opts SearchOptions, maxContentLength int) ([]CrawledResult, error) {
	results, err := searchWeb(opts)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	crawled := make([]CrawledResult, 0, len(results))
	errChan := make(chan error, len(results))

	for _, r := range results {
		wg.Add(1)
		go func(sr SearchResult) {
			defer wg.Done()

			content, err := crawlPage(sr.URL, maxContentLength)
			if err != nil {
				log.Printf("Failed to crawl %s: %v", sr.URL, err)
				return
			}

			mu.Lock()
			crawled = append(crawled, CrawledResult{
				SearchResult: sr,
				Content:      content,
				WordCount:    len(strings.Fields(content)),
			})
			mu.Unlock()
		}(r)
	}

	wg.Wait()
	close(errChan)

	return crawled, nil
}

// Calculate domain authority score
func getDomainAuthority(pageURL string) float64 {
	u, err := url.Parse(pageURL)
	if err != nil {
		return 0
	}

	domain := strings.ToLower(u.Host)

	// Check for subdomain matches
	for authDomain, score := range domainAuthority {
		if strings.HasSuffix(domain, authDomain) {
			return score
		}
	}

	// Base score for unknown domains
	return 2.0
}

// Calculate relevance score based on question type
func calculateRelevanceScore(question, content string) float64 {
	question = strings.ToLower(question)
	content = strings.ToLower(content)

	score := 0.0

	// Question type detection
	howTo := strings.HasPrefix(question, "how ") || strings.HasPrefix(question, "how to")
	whatIs := strings.HasPrefix(question, "what is") || strings.HasPrefix(question, "what are")
	whyIs := strings.HasPrefix(question, "why ") || strings.HasPrefix(question, "why does")
	compare := strings.Contains(question, " vs ") || strings.Contains(question, " versus ") ||
	           strings.Contains(question, " difference ") || strings.Contains(question, " compare ")

	// Look for answer patterns
	if howTo {
		if strings.Contains(content, "step") || strings.Contains(content, "first") ||
		   strings.Contains(content, "then") || strings.Contains(content, "follow") {
			score += 1.5
		}
		if strings.Contains(content, "tutorial") || strings.Contains(content, "guide") {
			score += 1.0
		}
	}

	if whatIs {
		if strings.Contains(content, "is a") || strings.Contains(content, "is an") ||
		   strings.Contains(content, "refers to") || strings.Contains(content, "defined") {
			score += 1.5
		}
		if strings.Contains(content, "definition") {
			score += 1.0
		}
	}

	if whyIs {
		if strings.Contains(content, "because") || strings.Contains(content, "due to") ||
		   strings.Contains(content, "reason") || strings.Contains(content, "cause") {
			score += 1.5
		}
	}

	if compare {
		if strings.Contains(content, "however") || strings.Contains(content, "whereas") ||
		   strings.Contains(content, "while") || strings.Contains(content, "on the other hand") {
			score += 1.5
		}
		if strings.Contains(content, "difference") || strings.Contains(content, "similar") {
			score += 1.0
		}
	}

	return score
}

// Research with AI-powered ranking
func research(question string, count int, maxContentLength int, region string, timeFilter string) ([]ResearchResult, error) {
	opts := SearchOptions{
		Query:  question,
		Limit:  count,
		Region: region,
		Time:   timeFilter,
	}

	crawled, err := searchAndCrawl(opts, maxContentLength)
	if err != nil {
		return nil, err
	}

	// Extract keywords from question
	questionWords := strings.Fields(strings.ToLower(question))
	keywords := make(map[string]bool)
	for _, word := range questionWords {
		if len(word) > 3 { // Skip short words
			keywords[word] = true
		}
	}

	results := make([]ResearchResult, 0, len(crawled))

	for _, c := range crawled {
		contentLower := strings.ToLower(c.Content)

		// Keyword score (30%)
		keywordCount := 0
		for kw := range keywords {
			if strings.Contains(contentLower, kw) {
				keywordCount++
			}
		}
		keywordScore := float64(keywordCount) / float64(len(keywords)) * 10

		// Quality score (25%)
		qualityScore := 0.0
		wc := c.WordCount

		switch {
		case wc < 100:
			qualityScore = 1.0
		case wc < 300:
			qualityScore = 3.0
		case wc < 800:
			qualityScore = 6.0
		case wc < 2000:
			qualityScore = 8.0
		default:
			qualityScore = 7.0
		}

		// Penalize low quality indicators
		if strings.Contains(contentLower, "subscribe") || strings.Contains(contentLower, "premium") {
			qualityScore -= 2
		}

		// Authority score (20%)
		authorityScore := getDomainAuthority(c.URL)

		// Relevance score (25%)
		relevanceScore := calculateRelevanceScore(question, c.Content)

		// Total score (normalized to 0-10)
		totalScore := (keywordScore*0.3 + qualityScore*0.25 + authorityScore*0.2 + relevanceScore*0.25)

		results = append(results, ResearchResult{
			CrawledResult:   c,
			Score:           totalScore,
			KeywordScore:    keywordScore,
			QualityScore:    qualityScore,
			AuthorityScore:  authorityScore,
			RelevanceScore:  relevanceScore,
		})
	}

	// Sort by score (descending)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
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
		rawURL := strings.ReplaceAll(match[1], "//duckduckgo.com/l/?uddg=", "")
		decodedURL, _ := url.QueryUnescape(rawURL)
		resultURL := strings.Split(decodedURL, "&")[0]
		snippet := cleanHTML(match[3])

		if title != "" && resultURL != "" {
			results = append(results, SearchResult{
				Title:   title,
				URL:     resultURL,
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
	s := server.NewMCPServer("duckduckgo-search", "3.0.0")

	// Tool 1: search - Basic web search
	searchTool := mcp.NewTool("search",
		mcp.WithDescription("Search the web using DuckDuckGo. Returns titles, URLs, and snippets."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results (default: 10, max: 20)"),
		),
		mcp.WithString("news",
			mcp.Description("Set to 'true' to search news only"),
		),
		mcp.WithString("region",
			mcp.Description("Region/locale code for localized results (default: 'cn-zh' for China). Use 'us-en' for US, 'jp-jp' for Japan, '' (empty) for global."),
		),
		mcp.WithString("time",
			mcp.Description("Time filter: 'd' (day), 'w' (week), 'm' (month), 'y' (year). Empty for any time."),
		),
	)

	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.Params.Arguments["query"].(string)
		limit := 10
		if l, ok := request.Params.Arguments["limit"].(float64); ok {
			limit = int(l)
		}
		if limit > 20 {
			limit = 20
		}

		newsOnly := false
		if newsStr, ok := request.Params.Arguments["news"].(string); ok && newsStr == "true" {
			newsOnly = true
		}

		region := "cn-zh" // default to China
		if r, ok := request.Params.Arguments["region"].(string); ok {
			region = r
		}
		timeFilter := ""
		if t, ok := request.Params.Arguments["time"].(string); ok {
			timeFilter = t
		}

		opts := SearchOptions{
			Query:    query,
			Limit:    limit,
			NewsOnly: newsOnly,
			Region:   region,
			Time:     timeFilter,
		}

		results, err := searchWeb(opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("## Found %d results\n\n", len(results)))
		for i, r := range results {
			output.WriteString(fmt.Sprintf("**%d. [%s](%s)**\n", i+1, r.Title, r.URL))
			output.WriteString(fmt.Sprintf("> %s\n\n", r.Snippet))
		}

		return mcp.NewToolResultText(output.String()), nil
	})

	// Tool 2: search_and_crawl - Search and crawl full content
	crawlTool := mcp.NewTool("search_and_crawl",
		mcp.WithDescription("Search the web and crawl full content from each result page in parallel."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of results to crawl (default: 5, max: 10)"),
		),
		mcp.WithNumber("maxContentLength",
			mcp.Description("Maximum characters per page (default: 3000, max: 10000)"),
		),
		mcp.WithString("region",
			mcp.Description("Region/locale code for localized results (default: 'cn-zh' for China). Use 'us-en' for US, 'jp-jp' for Japan, '' (empty) for global."),
		),
		mcp.WithString("time",
			mcp.Description("Time filter: 'd' (day), 'w' (week), 'm' (month), 'y' (year). Empty for any time."),
		),
	)

	s.AddTool(crawlTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.Params.Arguments["query"].(string)
		count := 5
		if c, ok := request.Params.Arguments["count"].(float64); ok {
			count = int(c)
		}
		if count > 10 {
			count = 10
		}

		maxLength := 3000
		if ml, ok := request.Params.Arguments["maxContentLength"].(float64); ok {
			maxLength = int(ml)
		}
		if maxLength > 10000 {
			maxLength = 10000
		}

		region := "cn-zh" // default to China
		if r, ok := request.Params.Arguments["region"].(string); ok {
			region = r
		}
		timeFilter := ""
		if t, ok := request.Params.Arguments["time"].(string); ok {
			timeFilter = t
		}

		opts := SearchOptions{
			Query:  query,
			Limit:  count,
			Region: region,
			Time:   timeFilter,
		}

		results, err := searchAndCrawl(opts, maxLength)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search and crawl failed: %v", err)), nil
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("## Crawled %d pages\n\n", len(results)))
		for i, r := range results {
			output.WriteString(fmt.Sprintf("### %d. %s\n", i+1, r.Title))
			output.WriteString(fmt.Sprintf("**URL:** %s\n", r.URL))
			output.WriteString(fmt.Sprintf("**Word Count:** %d\n\n", r.WordCount))
			output.WriteString(fmt.Sprintf("%s\n\n", r.Content))
			output.WriteString("---\n\n")
		}

		return mcp.NewToolResultText(output.String()), nil
	})

	// Tool 3: research - Search, crawl, and rank by relevance
	researchTool := mcp.NewTool("research",
		mcp.WithDescription("Research a question by searching, crawling, and ranking results by relevance using AI-powered scoring."),
		mcp.WithString("question",
			mcp.Required(),
			mcp.Description("Research question"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of sources to analyze (default: 5, max: 10)"),
		),
		mcp.WithNumber("maxContentLength",
			mcp.Description("Maximum characters per page (default: 3000, max: 10000)"),
		),
		mcp.WithString("region",
			mcp.Description("Region/locale code for localized results (default: 'cn-zh' for China). Use 'us-en' for US, 'jp-jp' for Japan, '' (empty) for global."),
		),
		mcp.WithString("time",
			mcp.Description("Time filter: 'd' (day), 'w' (week), 'm' (month), 'y' (year). Empty for any time."),
		),
	)

	s.AddTool(researchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		question := request.Params.Arguments["question"].(string)
		count := 5
		if c, ok := request.Params.Arguments["count"].(float64); ok {
			count = int(c)
		}
		if count > 10 {
			count = 10
		}

		maxLength := 3000
		if ml, ok := request.Params.Arguments["maxContentLength"].(float64); ok {
			maxLength = int(ml)
		}
		if maxLength > 10000 {
			maxLength = 10000
		}

		region := "cn-zh"
		if r, ok := request.Params.Arguments["region"].(string); ok {
			region = r
		}
		timeFilter := ""
		if t, ok := request.Params.Arguments["time"].(string); ok {
			timeFilter = t
		}

		results, err := research(question, count, maxLength, region, timeFilter)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Research failed: %v", err)), nil
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("# Research Results for: %s\n\n", question))
		output.WriteString(fmt.Sprintf("Analyzed %d sources, ranked by relevance:\n\n", len(results)))

		for i, r := range results {
			output.WriteString(fmt.Sprintf("## %d. %s - Score: %.1f/10\n", i+1, r.Title, r.Score))
			output.WriteString(fmt.Sprintf("**URL:** %s\n", r.URL))
			output.WriteString(fmt.Sprintf("**Word Count:** %d\n", r.WordCount))
			output.WriteString(fmt.Sprintf("**Scores:** Keywords: %.1f | Quality: %.1f | Authority: %.1f | Relevance: %.1f\n\n",
				r.KeywordScore, r.QualityScore, r.AuthorityScore, r.RelevanceScore))
			output.WriteString(fmt.Sprintf("**Content:**\n%s\n\n", r.Content))
			output.WriteString("---\n\n")
		}

		return mcp.NewToolResultText(output.String()), nil
	})

	// Tool 4: fetch - Fetch and extract clean content from a single URL
	fetchTool := mcp.NewTool("fetch",
		mcp.WithDescription("Fetch and extract clean text content from a single URL. Strips navigation, ads, scripts, and other non-content elements. Useful for reading a specific page in depth."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to fetch and extract content from"),
		),
		mcp.WithNumber("maxContentLength",
			mcp.Description("Maximum characters to return (default: 3000, max: 10000)"),
		),
	)

	s.AddTool(fetchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := request.Params.Arguments["url"].(string)
		maxLength := 3000
		if ml, ok := request.Params.Arguments["maxContentLength"].(float64); ok {
			maxLength = int(ml)
		}
		if maxLength > 10000 {
			maxLength = 10000
		}

		content, err := crawlPage(url, maxLength)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Fetch failed: %v", err)), nil
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("# Fetched: %s\n\n", url))
		output.WriteString(fmt.Sprintf("**Word Count:** %d\n\n", len(strings.Fields(content))))
		output.WriteString(content)

		return mcp.NewToolResultText(output.String()), nil
	})

	// Start stdio server
	stdioServer := server.NewStdioServer(s)
	log.Println("DuckDuckGo Search MCP server v3.0.0 starting on stdio")
	log.Println("Tools available: search, search_and_crawl, research, fetch")
	if err := stdioServer.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
