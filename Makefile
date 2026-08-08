# Variables
BINARY_NAME=kube-waste
DIST_DIR=dist
VERSION ?= $(shell git describe --tags --abbrev=0)

# .PHONY ensures these targets are treated as commands, not files
.PHONY: all clean build-linux build-windows build-macos checksums release run

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

# Write a SHA256SUMS checksum file for the built binaries
checksums: | $(DIST_DIR)
	@echo "Writing checksums..."
	cd $(DIST_DIR) && sha256sum * > SHA256SUMS

# Build everything, checksum it, and publish a GitHub release
# Requires: gh CLI authenticated (see https://cli.github.com)
release: all checksums
	@echo "Publishing release $(VERSION)..."
	gh release create $(VERSION) $(DIST_DIR)/* \
		--title "kube-waste $(VERSION)" \
		--notes "$$(cat .github/release-notes.md)"

# Development shortcut: Run locally using local .kubeconfig if present
run:
	go run main.go
