# 云念纪念馆项目 Makefile

.PHONY: help build run test clean docker-up docker-down docker-logs docker-restart docker-build migrate migrate-seed seed reset-db fmt lint deps

# 默认目标
help:
	@echo "云念纪念馆项目 - 可用命令:"
	@echo "  build        - 构建应用"
	@echo "  run          - 运行应用"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  docker-up    - 启动Docker服务"
	@echo "  docker-down  - 停止Docker服务"
	@echo "  docker-logs  - 查看Docker日志"
	@echo "  migrate      - 执行数据库迁移"
	@echo "  migrate-seed - 数据库迁移并插入种子数据"
	@echo "  seed         - 插入种子数据"
	@echo "  reset-db     - 重置数据库（危险操作）"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 代码检查"
	@echo "  deps         - 安装依赖"

# 构建应用
build:
	@echo "构建云念纪念馆应用..."
	go mod tidy
	go build -o bin/yun-nian-memorial cmd/server/main.go

# 运行应用
run:
	@echo "启动云念纪念馆服务..."
	go run cmd/server/main.go

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	go clean

# 启动Docker服务
docker-up:
	@echo "启动Docker服务..."
	docker-compose up -d

# 停止Docker服务
docker-down:
	@echo "停止Docker服务..."
	docker-compose down

# 查看Docker日志
docker-logs:
	@echo "查看Docker服务日志..."
	docker-compose logs -f

# 重启Docker服务
docker-restart: docker-down docker-up

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t yun-nian-memorial:latest .

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	go run cmd/server/main.go -migrate

# 数据库迁移并插入种子数据
migrate-seed:
	@echo "执行数据库迁移并插入种子数据..."
	go run cmd/server/main.go -migrate -seed

# 插入种子数据
seed:
	@echo "插入种子数据..."
	go run cmd/server/main.go -seed

# 重置数据库（危险操作）
reset-db:
	@echo "重置数据库（危险操作）..."
	@read -p "确定要重置数据库吗？这将删除所有数据 (y/N): " confirm && [ "$$confirm" = "y" ]
	go run cmd/server/main.go -reset

# 格式化代码
fmt:
	@echo "格式化Go代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "执行代码检查..."
	golangci-lint run

# 安装依赖
deps:
	@echo "安装项目依赖..."
	go mod download
	go mod tidy