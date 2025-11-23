# Variables
BINARY_NAME=kube-waste
DIST_DIR=dist

# .PHONY ensures these targets are treated as commands, not files
.PHONY: all clean build-linux build-windows build-macos run

# Default target: Build everything
all: clean build-linux build-windows build-macos
	@echo "All builds completed! Check the '$(DIST_DIR)' folder."

# Create the distribution directory
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# Build for Linux (Standard Servers)
build-linux: | $(DIST_DIR)
	@echo "Building for Linux (AMD64)..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 main.go

# Build for Windows
build-windows: | $(DIST_DIR)
	@echo "Building for Windows (AMD64)..."
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

# Build for macOS (Intel & Apple Silicon)
build-macos: | $(DIST_DIR)
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	@echo "Building for macOS (Apple Silicon/M1)..."
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(DIST_DIR)

# Development shortcut: Run locally using local .kubeconfig if present
run:
	go run main.go
