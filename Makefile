.PHONY: build build-all clean

# Build for current platform
build:
	go build -o darkdark-server .

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/darkdark-server-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/darkdark-server-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/darkdark-server-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/darkdark-server-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/darkdark-server-windows-amd64.exe .

clean:
	rm -rf dist/
	rm -f darkdark-server
