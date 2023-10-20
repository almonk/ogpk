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

sha:
	@sha256sum $(BUILD_DIR)/* | sed 's/build\///g' | sed 's/  /  /g'

formula:
	@CPU_ARCHS="amd64 arm64"; \
	echo "class Ogpk < Formula"; \
	echo "  desc \"CLI tool to fetch OpenGraph data from a URL\""; \
	echo "  homepage \"https://github.com/almonk/$(BINARY_NAME)\""; \
	echo "  version \"$(VERSION)\""; \
	for ARCH in $$CPU_ARCHS; do \
		FILENAME=$(BINARY_NAME)-$(VERSION)-darwin-$$ARCH; \
		LOCAL_PATH=$(BUILD_DIR)/$$FILENAME; \
		if [ ! -f $$LOCAL_PATH ]; then \
			echo "Error: Binary not found at $$LOCAL_PATH. Ensure you have built it."; \
			continue; \
		fi; \
		SHA256=$$(shasum -a 256 $$LOCAL_PATH | awk '{print $$1}'); \
		echo "  if Hardware::CPU.arm? && \"$$ARCH\" == \"arm64\""; \
		echo "    url \"https://github.com/almonk/$(BINARY_NAME)/releases/download/$(VERSION)/$$FILENAME\""; \
		echo "    sha256 \"$$SHA256\""; \
		echo "  elsif \"$$ARCH\" == \"amd64\""; \
		echo "    url \"https://github.com/almonk/$(BINARY_NAME)/releases/download/$(VERSION)/$$FILENAME\""; \
		echo "    sha256 \"$$SHA256\""; \
		echo "  end"; \
	done; \
	echo "  def install"; \
	echo "    if Hardware::CPU.arm?"; \
	echo "      bin.install \"$(BINARY_NAME)-$(VERSION)-darwin-arm64\" => \"$(BINARY_NAME)\""; \
	echo "    else"; \
	echo "      bin.install \"$(BINARY_NAME)-$(VERSION)-darwin-amd64\" => \"$(BINARY_NAME)\""; \
	echo "    end"; \
	echo "  end"; \
	echo "end";



.PHONY: all darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 install clean
