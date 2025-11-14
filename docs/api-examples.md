# API 接口测试示例

## 认证相关接口

### 1. 微信小程序登录

**接口地址：** `POST /api/v1/auth/wechat-login`

**请求示例：**
```bash
curl -X POST http://localhost:8080/api/v1/auth/wechat-login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test_code_123",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg"
  }'
```

**响应示例：**
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user-uuid-123",
      "wechat_openid": "test_openid_123",
      "nickname": "测试用户",
      "avatar_url": "https://example.com/avatar.jpg",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "is_new_user": true
  }
}
```

## 用户相关接口

### 2. 获取用户信息

**接口地址：** `GET /api/v1/users/profile`

**请求示例：**
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "user-uuid-123",
    "wechat_openid": "test_openid_123",
    "nickname": "测试用户",
    "avatar_url": "https://example.com/avatar.jpg",
    "phone": "13800138000",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. 更新用户信息

**接口地址：** `PUT /api/v1/users/profile`

**请求示例：**
```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新昵称",
    "phone": "13800138000"
  }'
```

**响应示例：**
```json
{
  "code": 0,
  "message": "更新成功"
}
```

### 4. 获取用户纪念馆列表

**接口地址：** `GET /api/v1/users/memorials`

**请求示例：**
```bash
curl -X GET "http://localhost:8080/api/v1/users/memorials?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "memorial-uuid-123",
        "creator_id": "user-uuid-123",
        "deceased_name": "张老爷子",
        "birth_date": "1950-01-01",
        "death_date": "2020-12-31",
        "biography": "张老爷子是一位慈祥的长者...",
        "theme_style": "traditional",
        "privacy_level": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 5. 获取用户祭扫记录

**接口地址：** `GET /api/v1/users/worship-records`

**请求示例：**
```bash
curl -X GET "http://localhost:8080/api/v1/users/worship-records?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "worship-uuid-123",
        "memorial_id": "memorial-uuid-123",
        "user_id": "user-uuid-123",
        "worship_type": "flower",
        "content": "{\"flower_type\": \"菊花\", \"count\": 3, \"message\": \"爷爷，我们想您了\"}",
        "created_at": "2024-01-01T00:00:00Z",
        "memorial": {
          "id": "memorial-uuid-123",
          "deceased_name": "张老爷子"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

## 错误响应示例

### 认证失败
```json
{
  "code": 1002,
  "message": "认证令牌无效"
}
```

### 权限不足
```json
{
  "code": 1003,
  "message": "权限不足"
}
```

### 用户不存在
```json
{
  "code": 2001,
  "message": "用户不存在"
}
```

### 参数错误
```json
{
  "code": 1001,
  "message": "请求参数错误: code字段为必填项"
}
```

## 测试流程

1. **获取微信登录code**（在实际微信小程序中获取）
2. **调用登录接口**获取JWT Token
3. **使用Token访问需要认证的接口**
4. **测试权限控制**（访问其他用户的资源应该被拒绝）

## 注意事项

1. 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer TOKEN`
2. Token有效期为7天，过期后需要重新登录
3. 微信登录需要配置正确的AppID和AppSecret
4. 在开发环境中，可以使用测试用户数据进行接口测试