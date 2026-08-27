BINARY := commitkit
BUILD_DIR := dist

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.buildDate=$(BUILD_DATE)

.PHONY: setup build install test lint clean release

setup:
	go mod tidy
	$(MAKE) build
	$(MAKE) install
	commitkit install-hook

build:
	go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BINARY) \
		./

install:
	go install \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		./

test:
	go test ./...

lint:
	go vet ./...

release:
	@mkdir -p $(BUILD_DIR)

	GOOS=darwin GOARCH=arm64 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY)-darwin-arm64 \
		./

	GOOS=darwin GOARCH=amd64 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY)-darwin-amd64 \
		./

	GOOS=linux GOARCH=amd64 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY)-linux-amd64 \
		./

	GOOS=linux GOARCH=arm64 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY)-linux-arm64 \
		./

	GOOS=windows GOARCH=amd64 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe \
		./

clean:
	rm -rf $(BINARY) $(BUILD_DIR)