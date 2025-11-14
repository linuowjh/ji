# 系统配置和维护 API 文档

## 概述

本文档描述了系统配置和维护相关的API接口，包括祭扫节日配置、模板配置、数据备份、系统监控和日志管理功能。

## 祭扫节日配置 API

### 1. 获取祭扫节日配置列表

**接口地址：** `GET /api/v1/system/festivals`

**请求参数：**
- `active_only` (query, optional): 是否只获取激活的节日，默认为 `true`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "festival-qingming",
      "name": "清明节",
      "festival_date": "04-05",
      "description": "清明节是中国传统的祭祖节日",
      "reminder_days": 3,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 2. 创建祭扫节日配置

**接口地址：** `POST /api/v1/admin/festivals`

**权限要求：** 管理员

**请求体：**
```json
{
  "name": "春节",
  "festival_date": "01-01",
  "description": "春节是中国最重要的传统节日",
  "reminder_days": 3,
  "is_active": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "festival-xxx",
    "name": "春节",
    "festival_date": "01-01",
    "description": "春节是中国最重要的传统节日",
    "reminder_days": 3,
    "is_active": true
  }
}
```

### 3. 更新祭扫节日配置

**接口地址：** `PUT /api/v1/admin/festivals/:id`

**权限要求：** 管理员

**请求体：**
```json
{
  "name": "清明节",
  "description": "更新后的描述",
  "reminder_days": 5
}
```

### 4. 删除祭扫节日配置

**接口地址：** `DELETE /api/v1/admin/festivals/:id`

**权限要求：** 管理员

### 5. 获取即将到来的节日

**接口地址：** `GET /api/v1/system/festivals/upcoming`

**请求参数：**
- `days_ahead` (query, optional): 提前多少天，默认为 `7`

## 模板配置 API

### 1. 获取模板配置列表

**接口地址：** `GET /api/v1/system/templates`

**请求参数：**
- `template_type` (query, optional): 模板类型 (`theme`, `tombstone`, `prayer`)
- `active_only` (query, optional): 是否只获取激活的模板，默认为 `true`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "theme-traditional",
      "template_type": "theme",
      "template_name": "中式传统",
      "template_data": "{\"background\": \"traditional-bg.jpg\"}",
      "preview_url": "https://example.com/preview.jpg",
      "is_premium": false,
      "sort_order": 1,
      "is_active": true
    }
  ]
}
```

### 2. 创建模板配置

**接口地址：** `POST /api/v1/admin/templates`

**权限要求：** 管理员

**请求体：**
```json
{
  "template_type": "theme",
  "template_name": "现代简约",
  "template_data": "{\"background\": \"modern-bg.jpg\"}",
  "preview_url": "https://example.com/preview.jpg",
  "is_premium": false,
  "sort_order": 4
}
```

### 3. 更新模板配置

**接口地址：** `PUT /api/v1/admin/templates/:id`

**权限要求：** 管理员

### 4. 删除模板配置

**接口地址：** `DELETE /api/v1/admin/templates/:id`

**权限要求：** 管理员

## 系统配置 API

### 1. 获取系统配置列表

**接口地址：** `GET /api/v1/admin/system/configs`

**权限要求：** 管理员

**请求参数：**
- `config_type` (query, optional): 配置类型

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "sys-max-memorial",
      "config_key": "max_memorial_per_user",
      "config_value": "10",
      "config_type": "system",
      "description": "每个用户最多可创建的纪念馆数量",
      "is_active": true
    }
  ]
}
```

### 2. 设置系统配置

**接口地址：** `POST /api/v1/admin/system/configs`

**权限要求：** 管理员

**请求体：**
```json
{
  "config_key": "max_memorial_per_user",
  "config_value": "20",
  "config_type": "system",
  "description": "每个用户最多可创建的纪念馆数量"
}
```

### 3. 初始化默认配置

**接口地址：** `POST /api/v1/admin/system/configs/init`

**权限要求：** 管理员

### 4. 导出配置

**接口地址：** `GET /api/v1/admin/system/configs/export`

**权限要求：** 管理员

**请求参数：**
- `config_type` (query, optional): 配置类型 (`festival`, `template`, `system`, `all`)

**响应：** JSON 文件下载

## 数据备份 API

### 1. 创建数据备份

**接口地址：** `POST /api/v1/admin/backups`

**权限要求：** 管理员

**请求体：**
```json
{
  "backup_type": "full"
}
```

**备份类型：**
- `full`: 完整备份
- `incremental`: 增量备份（最近24小时）
- `user`: 用户数据备份

**响应示例：**
```json
{
  "code": 0,
  "message": "备份任务已创建，正在后台处理",
  "data": {
    "id": "backup-xxx",
    "backup_type": "full",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 获取备份列表

**接口地址：** `GET /api/v1/admin/backups`

**权限要求：** 管理员

**请求参数：**
- `page` (query, optional): 页码，默认为 `1`
- `page_size` (query, optional): 每页数量，默认为 `20`
- `backup_type` (query, optional): 备份类型

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "backup-xxx",
        "backup_type": "full",
        "backup_path": "/backups/full_backup_20240101_120000.zip",
        "file_size": 1048576,
        "status": "completed",
        "created_by": "admin-user-id",
        "created_at": "2024-01-01T12:00:00Z",
        "completed_at": "2024-01-01T12:05:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 3. 获取备份详情

**接口地址：** `GET /api/v1/admin/backups/:id`

**权限要求：** 管理员

### 4. 下载备份文件

**接口地址：** `GET /api/v1/admin/backups/:id/download`

**权限要求：** 管理员

**响应：** ZIP 文件下载

### 5. 恢复备份

**接口地址：** `POST /api/v1/admin/backups/:id/restore`

**权限要求：** 管理员

### 6. 删除备份

**接口地址：** `DELETE /api/v1/admin/backups/:id`

**权限要求：** 管理员

### 7. 清理旧备份

**接口地址：** `POST /api/v1/admin/backups/clean`

**权限要求：** 管理员

**请求参数：**
- `keep_count` (query, optional): 保留最近N个备份，默认为 `10`

### 8. 获取备份统计

**接口地址：** `GET /api/v1/admin/backups/stats`

**权限要求：** 管理员

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_backups": 15,
    "completed_backups": 12,
    "failed_backups": 3,
    "total_size": 104857600,
    "total_size_mb": 100.0,
    "last_backup_time": "2024-01-01T12:00:00Z"
  }
}
```

## 系统监控 API

### 1. 获取系统健康状态

**接口地址：** `GET /api/v1/admin/monitor/health`

**权限要求：** 管理员

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "database": {
      "status": "healthy",
      "open_connections": 5,
      "in_use": 2,
      "idle": 3
    },
    "memory": {
      "status": "healthy",
      "alloc_mb": 50.5,
      "total_alloc": 100.2,
      "sys_mb": 75.3,
      "num_gc": 10
    },
    "cpu": {
      "status": "healthy",
      "num_cpu": 8,
      "num_goroutine": 25
    },
    "api": {
      "status": "healthy",
      "avg_response_ms": 45.2,
      "recent_requests": 100
    },
    "overall_status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### 2. 获取仪表板统计

**接口地址：** `GET /api/v1/admin/monitor/dashboard`

**权限要求：** 管理员

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "memory": {
      "count": 60,
      "average": 52.3,
      "min": 45.0,
      "max": 65.5,
      "unit": "MB"
    },
    "cpu": {
      "count": 60,
      "average": 25.0,
      "min": 20,
      "max": 35,
      "unit": "goroutines"
    },
    "api": {
      "count": 1000,
      "average": 45.2,
      "min": 10.0,
      "max": 200.0,
      "unit": "ms"
    },
    "health": {
      "overall_status": "healthy"
    }
  }
}
```

### 3. 获取监控指标

**接口地址：** `GET /api/v1/admin/monitor/metrics`

**权限要求：** 管理员

**请求参数：**
- `metric_type` (query, optional): 指标类型 (`cpu`, `memory`, `disk`, `api`, `database`)
- `start_time` (query, optional): 开始时间 (RFC3339格式)
- `end_time` (query, optional): 结束时间 (RFC3339格式)
- `limit` (query, optional): 限制数量，默认为 `100`

### 4. 获取指标统计

**接口地址：** `GET /api/v1/admin/monitor/metrics/stats`

**权限要求：** 管理员

**请求参数：**
- `metric_type` (query, required): 指标类型
- `duration_hours` (query, optional): 时间范围（小时），默认为 `24`

### 5. 获取指标趋势

**接口地址：** `GET /api/v1/admin/monitor/metrics/trend`

**权限要求：** 管理员

**请求参数：**
- `metric_type` (query, required): 指标类型
- `duration_hours` (query, optional): 时间范围（小时），默认为 `24`
- `interval_minutes` (query, optional): 时间间隔（分钟），默认为 `60`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "timestamp": "2024-01-01T00:00:00Z",
      "value": 50.5,
      "count": 10
    },
    {
      "timestamp": "2024-01-01T01:00:00Z",
      "value": 52.3,
      "count": 12
    }
  ]
}
```

### 6. 清理旧监控数据

**接口地址：** `POST /api/v1/admin/monitor/metrics/clean`

**权限要求：** 管理员

**请求参数：**
- `days_to_keep` (query, optional): 保留最近N天的数据，默认为 `7`

## 日志管理 API

### 1. 获取日志列表

**接口地址：** `GET /api/v1/admin/logs`

**权限要求：** 管理员

**请求参数：**
- `log_level` (query, optional): 日志级别 (`info`, `warning`, `error`, `critical`)
- `log_type` (query, optional): 日志类型 (`admin`, `system`, `security`, `api`)
- `user_id` (query, optional): 用户ID
- `action` (query, optional): 操作关键词
- `ip_address` (query, optional): IP地址
- `start_time` (query, optional): 开始时间
- `end_time` (query, optional): 结束时间
- `page` (query, optional): 页码，默认为 `1`
- `page_size` (query, optional): 每页数量，默认为 `20`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "log-xxx",
        "log_level": "info",
        "log_type": "admin",
        "user_id": "admin-user-id",
        "action": "create_memorial",
        "details": "{\"memorial_id\": \"xxx\"}",
        "ip_address": "192.168.1.1",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 2. 获取日志统计

**接口地址：** `GET /api/v1/admin/logs/stats`

**权限要求：** 管理员

**请求参数：**
- `start_time` (query, optional): 开始时间
- `end_time` (query, optional): 结束时间

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_logs": 1000,
    "by_level": {
      "info": 800,
      "warning": 150,
      "error": 45,
      "critical": 5
    },
    "by_type": {
      "admin": 200,
      "system": 300,
      "security": 100,
      "api": 400
    },
    "recent_errors": [
      {
        "id": "log-xxx",
        "log_level": "error",
        "action": "database_connection_failed",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

### 3. 搜索日志

**接口地址：** `GET /api/v1/admin/logs/search`

**权限要求：** 管理员

**请求参数：**
- `keyword` (query, required): 搜索关键词
- `page` (query, optional): 页码，默认为 `1`
- `page_size` (query, optional): 每页数量，默认为 `20`

### 4. 获取安全警报

**接口地址：** `GET /api/v1/admin/logs/security-alerts`

**权限要求：** 管理员

**请求参数：**
- `page` (query, optional): 页码，默认为 `1`
- `page_size` (query, optional): 每页数量，默认为 `20`

### 5. 清理旧日志

**接口地址：** `POST /api/v1/admin/logs/clean`

**权限要求：** 管理员

**请求参数：**
- `days_to_keep` (query, optional): 保留最近N天的日志，默认为 `30`

### 6. 导出日志

**接口地址：** `GET /api/v1/admin/logs/export`

**权限要求：** 管理员

**请求参数：** 与获取日志列表相同

**响应：** JSON 文件下载

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 请求参数错误 |
| 1002 | 用户未登录 |
| 1003 | 权限不足 |
| 1004 | 资源不存在 |
| 1005 | 内部服务器错误 |

## 注意事项

1. 所有管理员接口都需要管理员权限
2. 备份和恢复操作可能需要较长时间，建议异步处理
3. 日志和监控数据会定期清理，建议定期导出重要数据
4. 系统配置修改后可能需要重启服务才能生效
5. 数据备份文件应妥善保管，避免泄露敏感信息
