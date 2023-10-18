BINARY_NAME=ogpk

VERSION=0.1.3

BUILD_DIR=build

GOFLAGS := -ldflags="-X main.version=$(VERSION)"

all: darwin-amd64 darwin-arm64 linux-amd64 linux-arm64

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-arm64

linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64

linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64

install:
	cp $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-`uname -s | tr A-Z a-z`-`uname -m` /usr/local/bin/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)

.PHONY: all darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 install clean
