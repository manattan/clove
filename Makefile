.PHONY: help install build test clean fmt lint

# バイナリ名
BINARY_NAME := clove

# GOBIN の自動検出 (go env GOBIN が設定されていればそれを使い、なければ ~/go/bin)
GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(HOME)/go/bin
endif

# デフォルトターゲット
.DEFAULT_GOAL := help

## help: このヘルプを表示
help:
	@echo "使用可能なターゲット:"
	@echo "  make install   - バイナリをビルドして $(GOBIN) にインストール"
	@echo "  make build     - ローカルにバイナリをビルド (./clove)"
	@echo "  make test      - すべてのテストを実行"
	@echo "  make fmt       - コードをフォーマット"
	@echo "  make lint      - 静的解析を実行 (golangci-lint が必要)"
	@echo "  make clean     - ビルド成果物を削除"

## install: バイナリをビルドして GOBIN にインストール
install:
	@echo "==> Installing $(BINARY_NAME) to $(GOBIN)..."
	go build -o $(GOBIN)/$(BINARY_NAME) .
	@echo "==> Installed: $(GOBIN)/$(BINARY_NAME)"
	@echo "==> Run '$(BINARY_NAME) help' to get started"

## build: ローカルにバイナリをビルド
build:
	@echo "==> Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .
	@echo "==> Built: ./$(BINARY_NAME)"

## test: すべてのテストを実行
test:
	@echo "==> Running tests..."
	go test -v ./...

## test-short: 統合テストをスキップして実行
test-short:
	@echo "==> Running unit tests only..."
	go test -short -v ./...

## test-coverage: カバレッジ付きでテストを実行
test-coverage:
	@echo "==> Running tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "==> Coverage report: coverage.html"

## fmt: コードをフォーマット
fmt:
	@echo "==> Formatting code..."
	go fmt ./...

## lint: 静的解析を実行 (golangci-lint が必要)
lint:
	@echo "==> Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

## clean: ビルド成果物を削除
clean:
	@echo "==> Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -f main.go.bak
	@echo "==> Cleaned"

## uninstall: GOBIN からバイナリを削除
uninstall:
	@echo "==> Uninstalling $(BINARY_NAME) from $(GOBIN)..."
	rm -f $(GOBIN)/$(BINARY_NAME)
	@echo "==> Uninstalled"
