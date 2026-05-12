.PHONY: build build-all clean test

# Build for current platform
build:
	go build -o duckduckgo-mcp .

# Build for all platforms
build-all:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/duckduckgo-mcp-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/duckduckgo-mcp-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/duckduckgo-mcp-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/duckduckgo-mcp-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/duckduckgo-mcp-windows-amd64.exe .

clean:
	rm -rf dist/
	rm -f duckduckgo-mcp

test:
	go test -v ./...
