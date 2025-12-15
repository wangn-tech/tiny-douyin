.PHONY: help wire build run test clean

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

wire: ## 生成 Wire 依赖注入代码
	@echo "生成 Wire 代码..."
	@cd internal/wire && wire
	@echo "✅ Wire 代码生成完成"

build: wire ## 编译项目
	@echo "编译项目..."
	@go build -o bin/tiny-douyin main.go
	@echo "✅ 编译完成: bin/tiny-douyin"

run: wire ## 运行项目
	@echo "启动服务..."
	@go run main.go

test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告: coverage.html"

clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -f bin/tiny-douyin
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

deps: ## 安装依赖
	@echo "安装依赖..."
	@go mod download
	@go install github.com/google/wire/cmd/wire@latest
	@echo "✅ 依赖安装完成"

fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...
	@echo "✅ 代码格式化完成"

lint: ## 代码检查
	@echo "代码检查..."
	@golangci-lint run
	@echo "✅ 代码检查完成"

.DEFAULT_GOAL := help
