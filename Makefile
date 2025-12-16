.PHONY: help wire build run dev stop restart clean test docker-up docker-down

# 默认目标
.DEFAULT_GOAL := help

# 应用名称和路径
APP_NAME := tiny-douyin
BIN_PATH := ./$(APP_NAME)
PID_FILE := app.pid

# 颜色定义
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RESET  := \033[0m

help: ## 显示帮助信息
	@echo "$(GREEN)可用命令:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-12s$(RESET) %s\n", $$1, $$2}'

wire: ## 生成 Wire 依赖注入代码
	@echo "$(GREEN)生成 Wire 代码...$(RESET)"
	@wire gen ./internal/wire/
	@echo "$(GREEN)✅ Wire 代码生成完成$(RESET)"

build: wire ## 编译项目
	@echo "$(GREEN)编译项目...$(RESET)"
	@go build -o $(BIN_PATH) main.go
	@echo "$(GREEN)✅ 编译完成: $(BIN_PATH)$(RESET)"

run: build ## 编译并运行项目（前台）
	@echo "$(GREEN)启动服务...$(RESET)"
	@$(BIN_PATH) --env=dev

dev: build ## 编译并后台运行项目
	@echo "$(GREEN)后台启动服务...$(RESET)"
	@nohup $(BIN_PATH) --env=dev > app.log 2>&1 & echo $$! > $(PID_FILE)
	@sleep 1
	@if [ -f $(PID_FILE) ]; then \
		echo "$(GREEN)✅ 服务已启动，PID: $$(cat $(PID_FILE))$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  服务启动失败$(RESET)"; \
	fi

stop: ## 停止后台运行的服务
	@if [ -f $(PID_FILE) ]; then \
		echo "$(GREEN)停止服务 PID: $$(cat $(PID_FILE))...$(RESET)"; \
		kill $$(cat $(PID_FILE)) 2>/dev/null || echo "$(YELLOW)进程已不存在$(RESET)"; \
		rm -f $(PID_FILE); \
		echo "$(GREEN)✅ 服务已停止$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  未找到 PID 文件，尝试强制停止...$(RESET)"; \
		pkill -f $(APP_NAME) || echo "$(YELLOW)没有运行中的服务$(RESET)"; \
	fi

restart: stop dev ## 重启服务

status: ## 查看服务状态
	@if [ -f $(PID_FILE) ]; then \
		PID=$$(cat $(PID_FILE)); \
		if ps -p $$PID > /dev/null 2>&1; then \
			echo "$(GREEN)✅ 服务运行中，PID: $$PID$(RESET)"; \
			ps -p $$PID -o pid,ppid,%cpu,%mem,etime,cmd; \
		else \
			echo "$(YELLOW)⚠️  PID 文件存在但进程不存在$(RESET)"; \
			rm -f $(PID_FILE); \
		fi \
	else \
		echo "$(YELLOW)⚠️  服务未运行$(RESET)"; \
	fi

logs: ## 查看应用日志（实时）
	@tail -f ./tmp/logs/tiny-douyin.log

logs-app: ## 查看启动日志
	@tail -f app.log

test: ## 运行测试
	@echo "$(GREEN)运行测试...$(RESET)"
	@go test -v ./...

test-cover: ## 运行测试并生成覆盖率报告
	@echo "$(GREEN)运行测试并生成覆盖率...$(RESET)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ 覆盖率报告: coverage.html$(RESET)"

clean: ## 清理构建文件和临时文件
	@echo "$(GREEN)清理文件...$(RESET)"
	@rm -f $(BIN_PATH) $(PID_FILE) app.log
	@rm -f coverage.out coverage.html
	@rm -rf ./tmp/uploads/*
	@echo "$(GREEN)✅ 清理完成$(RESET)"

docker-up: ## 启动 Docker 服务（MySQL、Redis、MinIO、RabbitMQ）
	@echo "$(GREEN)启动 Docker 服务...$(RESET)"
	@docker compose up -d
	@echo "$(GREEN)✅ Docker 服务已启动$(RESET)"
	@docker compose ps

docker-down: ## 停止 Docker 服务
	@echo "$(GREEN)停止 Docker 服务...$(RESET)"
	@docker compose down
	@echo "$(GREEN)✅ Docker 服务已停止$(RESET)"

docker-restart: docker-down docker-up ## 重启 Docker 服务

docker-logs: ## 查看 Docker 服务日志
	@docker compose logs -f

docker-clean: ## 清理 Docker 数据卷
	@echo "$(YELLOW)⚠️  这将删除所有数据！$(RESET)"
	@read -p "确认删除？(y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@docker compose down -v
	@echo "$(GREEN)✅ Docker 数据已清理$(RESET)"

deps: ## 安装 Go 依赖和工具
	@echo "$(GREEN)安装依赖...$(RESET)"
	@go mod download
	@go install github.com/google/wire/cmd/wire@latest
	@echo "$(GREEN)✅ 依赖安装完成$(RESET)"

fmt: ## 格式化代码
	@echo "$(GREEN)格式化代码...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ 代码格式化完成$(RESET)"

tidy: ## 整理 go.mod
	@echo "$(GREEN)整理依赖...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✅ 依赖整理完成$(RESET)"

init: deps docker-up ## 初始化项目（安装依赖 + 启动 Docker）
	@echo "$(GREEN)✅ 项目初始化完成$(RESET)"
	@echo "$(YELLOW)提示：运行 'make dev' 启动应用$(RESET)"

all: clean build ## 清理并重新编译
