# 云念 - 网上扫墓祭奠小程序

「云念」是一个基于微信小程序的数字化纪念服务平台，旨在为用户提供突破时空限制的线上祭奠服务。

## 项目特色

- 🕯️ **庄重缅怀** - 还原传统祭扫仪式感，提供献花、点烛、上香等功能
- 🏠 **家族联动** - 支持家族圈创建，让分散各地的亲友共同参与纪念
- 📱 **便捷使用** - 基于微信小程序，随时随地表达思念
- 💾 **永久保存** - 数字化存储纪念内容，传承家族情感

## 技术架构

### 后端技术栈
- **语言**: Golang 1.19+
- **框架**: Gin Web框架
- **数据库**: MySQL 8.0 + Redis 6.0
- **ORM**: GORM
- **认证**: JWT
- **容器化**: Docker + Docker Compose

### 前端技术栈
- **平台**: 微信小程序
- **开发**: 微信小程序原生开发
- **UI组件**: WeUI

## 快速开始

### 环境要求
- Go 1.19+
- Docker & Docker Compose
- MySQL 8.0
- Redis 6.0

### 本地开发

1. 克隆项目
```bash
git clone <repository-url>
cd yun-nian-memorial
```

2. 启动服务
```bash
# 使用Docker Compose启动所有服务
docker-compose up -d

# 或者本地开发模式
go mod tidy
go run cmd/server/main.go
```

3. 访问服务
- API服务: http://localhost:8080
- 健康检查: http://localhost:8080/health

### 环境变量配置

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| MYSQL_HOST | MySQL主机 | localhost |
| MYSQL_PORT | MySQL端口 | 3306 |
| MYSQL_USERNAME | MySQL用户名 | root |
| MYSQL_PASSWORD | MySQL密码 | - |
| MYSQL_DATABASE | 数据库名 | yun_nian_memorial |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| JWT_SECRET | JWT密钥 | yun-nian-memorial-secret |
| WECHAT_APP_ID | 微信小程序AppID | - |
| WECHAT_APP_SECRET | 微信小程序AppSecret | - |

## API文档

### 认证相关
- `POST /api/v1/auth/wechat-login` - 微信登录

### 纪念馆相关
- `GET /api/v1/memorials` - 获取纪念馆列表
- `POST /api/v1/memorials` - 创建纪念馆
- `GET /api/v1/memorials/:id` - 获取纪念馆详情
- `PUT /api/v1/memorials/:id` - 更新纪念馆

### 祭扫相关
- `POST /api/v1/worship/flower` - 献花
- `POST /api/v1/worship/candle` - 点烛
- `POST /api/v1/worship/incense` - 上香
- `POST /api/v1/worship/prayer` - 祈福

### 家族相关
- `GET /api/v1/families` - 获取家族列表
- `POST /api/v1/families` - 创建家族圈

## 项目结构

```
.
├── cmd/
│   └── server/          # 应用入口
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   └── router/          # 路由配置
├── scripts/
│   └── init.sql         # 数据库初始化脚本
├── docker-compose.yml   # Docker编排文件
├── Dockerfile          # Docker镜像构建文件
└── README.md           # 项目说明
```

## 开发规范

### 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 添加必要的注释和文档

### 提交规范
- feat: 新功能
- fix: 修复bug
- docs: 文档更新
- style: 代码格式调整
- refactor: 代码重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

## 许可证

本项目采用 MIT 许可证。