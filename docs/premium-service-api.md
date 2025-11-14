# 增值服务 API 文档

## 概述

本文档描述了增值服务相关的API接口，包括高级纪念馆服务、套餐订阅、定制模板、存储管理等功能。

## 高级套餐管理 API

### 1. 获取高级套餐列表

**接口地址：** `GET /api/v1/premium/packages`

**请求参数：**
- `package_type` (query, optional): 套餐类型 (`memorial`, `service`, `storage`)
- `active_only` (query, optional): 是否只获取激活的套餐，默认为 `true`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "pkg-premium",
      "package_name": "高级版",
      "package_type": "memorial",
      "description": "提供更多定制化功能和存储空间",
      "features": "[\"500MB存储空间\", \"高级主题模板\", \"定制墓碑样式\"]",
      "price": 99.00,
      "duration": 365,
      "storage_size": 524288000,
      "is_active": true,
      "sort_order": 2,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 2. 获取套餐详情

**接口地址：** `GET /api/v1/premium/packages/:id`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "pkg-premium",
    "package_name": "高级版",
    "package_type": "memorial",
    "description": "提供更多定制化功能和存储空间",
    "features": "[\"500MB存储空间\", \"高级主题模板\", \"定制墓碑样式\", \"老照片修复\", \"优先客服支持\"]",
    "price": 99.00,
    "duration": 365,
    "storage_size": 524288000,
    "is_active": true
  }
}
```

### 3. 创建高级套餐（管理员）

**接口地址：** `POST /api/v1/admin/premium/packages`

**权限要求：** 管理员

**请求体：**
```json
{
  "package_name": "企业版",
  "package_type": "memorial",
  "description": "适合企业使用的纪念馆服务",
  "features": "[\"10GB存储空间\", \"无限纪念馆\", \"专属客服\"]",
  "price": 999.00,
  "duration": 365,
  "storage_size": 10737418240,
  "is_active": true,
  "sort_order": 5
}
```

### 4. 更新套餐信息（管理员）

**接口地址：** `PUT /api/v1/admin/premium/packages/:id`

**权限要求：** 管理员

**请求体：**
```json
{
  "price": 89.00,
  "description": "更新后的描述"
}
```

### 5. 初始化默认套餐（管理员）

**接口地址：** `POST /api/v1/admin/premium/packages/init`

**权限要求：** 管理员

## 用户订阅管理 API

### 1. 订阅套餐

**接口地址：** `POST /api/v1/premium/subscribe`

**权限要求：** 登录用户

**请求体：**
```json
{
  "package_id": "pkg-premium",
  "memorial_id": "memorial-xxx"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "订阅成功",
  "data": {
    "id": "sub-xxx",
    "user_id": "user-xxx",
    "package_id": "pkg-premium",
    "memorial_id": "memorial-xxx",
    "status": "active",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2025-01-01T00:00:00Z",
    "auto_renew": false,
    "payment_amount": 99.00,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 获取用户订阅列表

**接口地址：** `GET /api/v1/premium/subscriptions`

**权限要求：** 登录用户

**请求参数：**
- `status` (query, optional): 订阅状态 (`active`, `expired`, `cancelled`)

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "sub-xxx",
      "user_id": "user-xxx",
      "package_id": "pkg-premium",
      "status": "active",
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2025-01-01T00:00:00Z",
      "auto_renew": false,
      "payment_amount": 99.00,
      "package": {
        "id": "pkg-premium",
        "package_name": "高级版",
        "price": 99.00
      }
    }
  ]
}
```

### 3. 获取订阅详情

**接口地址：** `GET /api/v1/premium/subscriptions/:id`

**权限要求：** 登录用户（仅能查看自己的订阅）

### 4. 取消订阅

**接口地址：** `POST /api/v1/premium/subscriptions/:id/cancel`

**权限要求：** 登录用户

**响应示例：**
```json
{
  "code": 0,
  "message": "取消成功"
}
```

### 5. 续订

**接口地址：** `POST /api/v1/premium/subscriptions/:id/renew`

**权限要求：** 登录用户

**响应示例：**
```json
{
  "code": 0,
  "message": "续订成功"
}
```

## 纪念馆升级 API

### 1. 升级纪念馆

**接口地址：** `POST /api/v1/premium/memorial/upgrade`

**权限要求：** 登录用户

**请求体：**
```json
{
  "memorial_id": "memorial-xxx",
  "subscription_id": "sub-xxx",
  "upgrade_type": "theme",
  "upgrade_data": {
    "theme_id": "theme-luxury",
    "custom_colors": {
      "primary": "#8B4513",
      "secondary": "#D2691E"
    }
  }
}
```

**升级类型：**
- `theme`: 主题升级
- `tombstone`: 墓碑样式升级
- `storage`: 存储空间升级
- `feature`: 功能升级

**响应示例：**
```json
{
  "code": 0,
  "message": "升级成功"
}
```

### 2. 获取纪念馆升级记录

**接口地址：** `GET /api/v1/premium/memorial/:memorial_id/upgrades`

**权限要求：** 登录用户

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "upgrade-xxx",
      "memorial_id": "memorial-xxx",
      "subscription_id": "sub-xxx",
      "upgrade_type": "theme",
      "upgrade_data": "{\"theme_id\": \"theme-luxury\"}",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "subscription": {
        "id": "sub-xxx",
        "package": {
          "package_name": "高级版"
        }
      }
    }
  ]
}
```

## 定制模板 API

### 1. 创建定制模板

**接口地址：** `POST /api/v1/premium/templates`

**权限要求：** 登录用户（需要高级订阅）

**请求体：**
```json
{
  "memorial_id": "memorial-xxx",
  "template_type": "theme",
  "template_name": "我的专属主题",
  "template_data": "{\"background\": \"custom-bg.jpg\", \"colors\": {\"primary\": \"#8B4513\"}}",
  "preview_url": "https://example.com/preview.jpg"
}
```

**模板类型：**
- `theme`: 主题模板
- `tombstone`: 墓碑模板
- `layout`: 布局模板

**响应示例：**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "template-xxx",
    "user_id": "user-xxx",
    "memorial_id": "memorial-xxx",
    "template_type": "theme",
    "template_name": "我的专属主题",
    "status": "draft",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 获取用户定制模板列表

**接口地址：** `GET /api/v1/premium/templates`

**权限要求：** 登录用户

**请求参数：**
- `template_type` (query, optional): 模板类型

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "template-xxx",
      "template_type": "theme",
      "template_name": "我的专属主题",
      "preview_url": "https://example.com/preview.jpg",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 3. 更新定制模板

**接口地址：** `PUT /api/v1/premium/templates/:id`

**权限要求：** 登录用户

**请求体：**
```json
{
  "template_name": "更新后的名称",
  "status": "active"
}
```

### 4. 删除定制模板

**接口地址：** `DELETE /api/v1/premium/templates/:id`

**权限要求：** 登录用户

## 存储管理 API

### 1. 获取用户存储使用情况

**接口地址：** `GET /api/v1/premium/storage`

**权限要求：** 登录用户

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "storage": {
      "id": "storage-xxx",
      "user_id": "user-xxx",
      "used_space": 52428800,
      "total_space": 524288000,
      "file_count": 150,
      "last_updated": "2024-01-01T12:00:00Z"
    },
    "usage_percent": 10.0,
    "available_space": 471859200
  }
}
```

**说明：**
- `used_space`: 已使用空间（字节）
- `total_space`: 总空间（字节）
- `usage_percent`: 使用百分比
- `available_space`: 可用空间（字节）

## 服务使用统计 API

### 1. 获取服务使用统计

**接口地址：** `GET /api/v1/premium/usage/stats`

**权限要求：** 登录用户

**请求参数：**
- `days_ago` (query, optional): 统计最近N天，默认为 `30`
- `start_time` (query, optional): 开始时间 (RFC3339格式)
- `end_time` (query, optional): 结束时间 (RFC3339格式)

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_usage": 45,
    "by_service": {
      "subscription": 2,
      "memorial_upgrade": 5,
      "custom_template": 10,
      "photo_restore": 28
    },
    "start_time": "2023-12-01T00:00:00Z",
    "end_time": "2024-01-01T00:00:00Z"
  }
}
```

## 套餐功能对比

| 功能 | 基础版 | 高级版 | 尊享版 |
|------|--------|--------|--------|
| 存储空间 | 100MB | 500MB | 2GB |
| 纪念馆数量 | 3个 | 10个 | 无限 |
| 主题模板 | 基础 | 高级 | 全部 |
| 定制墓碑 | ❌ | ✅ | ✅ |
| 老照片修复 | ❌ | 5次/年 | 无限 |
| 专属追思会 | ❌ | ❌ | ✅ |
| 数据备份 | ❌ | ❌ | ✅ |
| 客服支持 | 标准 | 优先 | 专属 |
| 价格 | 免费 | ¥99/年 | ¥299/年 |

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 请求参数错误 |
| 1002 | 用户未登录 |
| 1003 | 权限不足 |
| 1004 | 资源不存在 |
| 1005 | 内部服务器错误 |
| 3001 | 套餐不存在 |
| 3002 | 套餐已下架 |
| 3003 | 已订阅该套餐 |
| 3004 | 订阅已失效 |
| 3005 | 订阅已过期 |
| 3006 | 存储空间不足 |

## 注意事项

1. **订阅管理**
   - 订阅成功后立即生效
   - 取消订阅不会立即失效，会在到期后自动失效
   - 续订会从当前到期日期延长，如果已过期则从当前时间开始计算

2. **存储管理**
   - 基础版用户默认100MB存储空间
   - 升级套餐会立即增加存储空间
   - 取消订阅后，存储空间会在到期后恢复到基础版
   - 如果使用空间超过基础版限制，需要清理文件或续订

3. **定制模板**
   - 定制模板功能仅对高级版及以上用户开放
   - 每个用户最多可创建10个定制模板
   - 模板可以在多个纪念馆之间共享使用

4. **支付流程**
   - 本API文档仅描述订阅管理接口
   - 实际支付需要集成微信支付或支付宝
   - 支付成功后需要调用订阅接口创建订阅记录

5. **自动续费**
   - 用户可以开启自动续费功能
   - 系统会在到期前3天尝试自动扣费
   - 扣费失败会通知用户手动续费

## 使用示例

### 完整订阅流程

1. **查看套餐列表**
```bash
GET /api/v1/premium/packages
```

2. **选择套餐并订阅**
```bash
POST /api/v1/premium/subscribe
{
  "package_id": "pkg-premium",
  "memorial_id": "memorial-xxx"
}
```

3. **升级纪念馆**
```bash
POST /api/v1/premium/memorial/upgrade
{
  "memorial_id": "memorial-xxx",
  "subscription_id": "sub-xxx",
  "upgrade_type": "theme",
  "upgrade_data": {
    "theme_id": "theme-luxury"
  }
}
```

4. **查看存储使用情况**
```bash
GET /api/v1/premium/storage
```

5. **创建定制模板**
```bash
POST /api/v1/premium/templates
{
  "template_type": "theme",
  "template_name": "我的专属主题",
  "template_data": "{...}"
}
```
