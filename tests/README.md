# 云念纪念馆测试文档

## 测试概述

本项目包含完整的测试套件，涵盖单元测试、集成测试、端到端测试和性能测试。

## 测试结构

```
tests/
├── e2e_test.go           # 端到端测试
├── performance_test.go   # 性能测试
└── README.md            # 测试文档

internal/
├── services/
│   ├── memorial_service_test.go    # 纪念馆服务单元测试
│   ├── user_service_test.go        # 用户服务单元测试
│   └── worship_service_test.go     # 祭扫服务单元测试
└── controllers/
    ├── memorial_controller_test.go  # 纪念馆控制器集成测试
    └── worship_controller_test.go   # 祭扫控制器集成测试
```

## 测试环境配置

### 数据库配置

测试需要MySQL数据库，配置如下：

```
Host: 127.0.0.1
Port: 3306
Database: yun_nian_test
Username: root
Password: root
```

### 创建测试数据库

```sql
CREATE DATABASE IF NOT EXISTS yun_nian_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 运行测试

### 运行所有测试

```bash
go test ./... -v
```

### 运行单元测试

```bash
# 运行所有服务单元测试
go test ./internal/services/... -v

# 运行特定服务的测试
go test ./internal/services -run TestMemorialService -v
go test ./internal/services -run TestUserService -v
go test ./internal/services -run TestWorshipService -v
```

### 运行集成测试

```bash
# 运行所有控制器集成测试
go test ./internal/controllers/... -v

# 运行特定控制器的测试
go test ./internal/controllers -run TestMemorialController -v
go test ./internal/controllers -run TestWorshipController -v
```

### 运行端到端测试

```bash
# 运行所有E2E测试
go test ./tests -run TestComplete -v
go test ./tests -run TestUserJourney -v
go test ./tests -run TestMultiUser -v
```

### 运行性能测试

```bash
# 运行性能测试
go test ./tests -run TestPerformance -v

# 运行基准测试
go test ./tests -bench=. -benchmem
```

## 测试覆盖率

### 生成覆盖率报告

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out

# 查看覆盖率
go tool cover -func=coverage.out

# 生成HTML覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

## 测试说明

### 单元测试

单元测试专注于测试单个函数和方法的功能：

- **memorial_service_test.go**: 测试纪念馆的创建、更新、查询等核心功能
- **user_service_test.go**: 测试用户管理、JWT认证、统计等功能
- **worship_service_test.go**: 测试祭扫操作、祈福、留言等功能

### 集成测试

集成测试验证API接口的完整性：

- **memorial_controller_test.go**: 测试纪念馆相关API接口
- **worship_controller_test.go**: 测试祭扫相关API接口

### 端到端测试

端到端测试模拟真实用户场景：

1. **TestCompleteMemorialCreationFlow**: 测试完整的纪念馆创建和管理流程
2. **TestCompleteWorshipFlow**: 测试完整的祭扫流程
3. **TestUserJourneyFlow**: 测试用户从注册到使用的完整旅程
4. **TestMultiUserInteraction**: 测试多用户交互场景

### 性能测试

性能测试评估系统在负载下的表现：

1. **TestMemorialCreationPerformance**: 测试纪念馆创建性能
2. **TestConcurrentWorshipOperations**: 测试并发祭扫操作
3. **TestQueryPerformance**: 测试查询性能
4. **TestMemoryUsage**: 测试内存使用情况

## 测试最佳实践

### 1. 测试隔离

每个测试都应该独立运行，不依赖其他测试的状态：

```go
func TestExample(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    // 测试代码
}
```

### 2. 使用断言

使用testify库进行断言，提高测试可读性：

```go
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.NotNil(t, result)
```

### 3. 测试数据清理

确保每个测试后清理测试数据：

```go
defer cleanupTestDB(db)
```

### 4. 跳过不可用的测试

当测试环境不可用时，优雅地跳过测试：

```go
if err != nil {
    t.Skip("Database not available for testing")
}
```

## 持续集成

### GitHub Actions配置示例

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: yun_nian_test
        ports:
          - 3306:3306
        options: --health-cmd="mysqladmin ping" --health-interval=10s --health-timeout=5s --health-retries=3
    
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.19
      
      - name: Run tests
        run: |
          go test ./... -v -coverprofile=coverage.out
          go tool cover -func=coverage.out
```

## 故障排查

### 数据库连接失败

如果测试因数据库连接失败而跳过：

1. 确认MySQL服务正在运行
2. 检查数据库配置是否正确
3. 确认测试数据库已创建

### 测试超时

如果性能测试超时：

1. 检查数据库性能
2. 调整测试参数（减少测试数据量）
3. 增加测试超时时间

### 并发测试失败

如果并发测试失败：

1. 检查数据库连接池配置
2. 确认数据库支持并发操作
3. 检查是否有死锁或竞态条件

## 测试指标

### 目标指标

- **单元测试覆盖率**: > 80%
- **集成测试覆盖率**: > 70%
- **E2E测试覆盖率**: 核心业务流程100%
- **性能测试**:
  - 纪念馆创建: < 100ms/个
  - 祭扫操作: < 50ms/次
  - 查询操作: < 100ms
  - 并发支持: > 100 QPS

## 贡献指南

添加新测试时，请遵循以下规范：

1. 测试函数命名: `Test<功能名称>`
2. 使用描述性的测试名称
3. 添加必要的注释说明测试目的
4. 确保测试可以独立运行
5. 清理测试数据
6. 更新本文档

## 联系方式

如有测试相关问题，请联系开发团队。
