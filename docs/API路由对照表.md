# API路由对照表

## 常见混淆的路由

### ❌ 不存在的路由

```
GET /api/v1/memorials/recent  ❌ 此路由不存在
```

### ✅ 正确的路由

如果你想获取：

#### 1. 用户最近活动
```
GET /api/v1/users/activities
```
**说明**：获取当前用户的最近活动记录（祭扫、家族活动等）

**响应示例**：
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "type": "worship",
      "action": "flower",
      "memorial": {
        "id": "memorial_id",
        "deceased_name": "张三"
      },
      "created_at": "2025-11-14T10:00:00Z"
    }
  ]
}
```

#### 2. 用户的纪念馆列表
```
GET /api/v1/users/memorials
```
**说明**：获取当前用户创建的所有纪念馆

**响应示例**：
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "memorial_id",
        "deceased_name": "张三",
        "created_at": "2025-11-14T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

#### 3. 所有纪念馆列表（公开的）
```
GET /api/v1/memorials
```
**说明**：获取所有公开的纪念馆列表

**查询参数**：
- `page` - 页码（默认1）
- `page_size` - 每页数量（默认10）
- `keyword` - 搜索关键词（可选）

#### 4. 用户仪表板（包含最近活动）
```
GET /api/v1/users/dashboard
```
**说明**：获取用户仪表板数据，包含统计信息、最近活动、纪念馆列表等

**响应示例**：
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "statistics": {
      "memorial_count": 5,
      "worship_count": 20
    },
    "recent_activities": [...],
    "memorials": [...],
    "families": [...]
  }
}
```

## 完整的API路由列表

### 认证相关（无需Token）

```
POST /api/v1/auth/wechat-login  # 微信登录
```

### 用户相关（需要Token）

```
GET    /api/v1/users/profile              # 获取用户信息
PUT    /api/v1/users/profile              # 更新用户信息
GET    /api/v1/users/memorials            # 获取用户的纪念馆列表
GET    /api/v1/users/worship-records      # 获取用户的祭扫记录
GET    /api/v1/users/dashboard            # 获取用户仪表板
GET    /api/v1/users/statistics           # 获取用户统计信息
GET    /api/v1/users/activities           # 获取用户最近活动 ⭐
GET    /api/v1/users/memorial-details     # 获取用户纪念馆详情
GET    /api/v1/users/families             # 获取用户的家族圈
GET    /api/v1/users/memorials/:memorial_id/visitors  # 获取纪念馆访客
```

### 纪念馆相关（需要Token）

```
GET    /api/v1/memorials                  # 获取纪念馆列表
POST   /api/v1/memorials                  # 创建纪念馆
GET    /api/v1/memorials/:id              # 获取纪念馆详情
PUT    /api/v1/memorials/:id              # 更新纪念馆
DELETE /api/v1/memorials/:id              # 删除纪念馆
GET    /api/v1/memorials/:id/visitors     # 获取纪念馆访客
PUT    /api/v1/memorials/:id/tombstone-style  # 更新墓碑样式
PUT    /api/v1/memorials/:id/epitaph      # 更新墓志铭
```

### 祭扫相关（需要Token）

```
POST   /api/v1/worship/memorials/:memorial_id/flowers   # 献花
POST   /api/v1/worship/memorials/:memorial_id/candles   # 点烛
PUT    /api/v1/worship/memorials/:memorial_id/candles/renew  # 续烛
GET    /api/v1/worship/memorials/:memorial_id/candles/status # 蜡烛状态
POST   /api/v1/worship/memorials/:memorial_id/incense   # 上香
POST   /api/v1/worship/memorials/:memorial_id/tributes  # 供品
POST   /api/v1/worship/memorials/:memorial_id/prayers   # 祈福
POST   /api/v1/worship/memorials/:memorial_id/messages  # 留言
```

### 家族圈相关（需要Token）

```
GET    /api/v1/families                   # 获取家族列表
POST   /api/v1/families                   # 创建家族
GET    /api/v1/families/:id               # 获取家族详情
PUT    /api/v1/families/:id               # 更新家族
DELETE /api/v1/families/:id               # 删除家族
GET    /api/v1/families/:id/members       # 获取家族成员
POST   /api/v1/families/:id/invite        # 邀请成员
```

### 相册相关（需要Token）

```
POST   /api/v1/albums/memorials/:memorial_id  # 创建相册
GET    /api/v1/albums/memorials/:memorial_id  # 获取相册列表
GET    /api/v1/albums/:id                     # 获取相册详情
PUT    /api/v1/albums/:id                     # 更新相册
DELETE /api/v1/albums/:id                     # 删除相册
POST   /api/v1/albums/:id/photos              # 添加照片
```

### 生平故事相关（需要Token）

```
POST   /api/v1/stories/memorials/:memorial_id  # 创建生平故事
GET    /api/v1/stories/memorials/:memorial_id  # 获取生平故事列表
GET    /api/v1/stories/:id                     # 获取故事详情
PUT    /api/v1/stories/:id                     # 更新故事
DELETE /api/v1/stories/:id                     # 删除故事
```

### 追思会相关（需要Token）

```
POST   /api/v1/memorial-services/memorials/:memorial_id  # 创建追思会
GET    /api/v1/memorial-services/memorials/:memorial_id  # 获取追思会列表
GET    /api/v1/memorial-services/:id                     # 获取追思会详情
POST   /api/v1/memorial-services/:id/start               # 开始追思会
POST   /api/v1/memorial-services/:id/end                 # 结束追思会
POST   /api/v1/memorial-services/:id/join                # 加入追思会
```

## 路由命名规则

### 资源路由

```
GET    /api/v1/{resource}           # 获取资源列表
POST   /api/v1/{resource}           # 创建资源
GET    /api/v1/{resource}/:id       # 获取单个资源
PUT    /api/v1/{resource}/:id       # 更新资源
DELETE /api/v1/{resource}/:id       # 删除资源
```

### 嵌套资源路由

```
GET    /api/v1/{parent}/:parent_id/{child}  # 获取父资源下的子资源列表
POST   /api/v1/{parent}/:parent_id/{child}  # 在父资源下创建子资源
```

### 动作路由

```
POST   /api/v1/{resource}/:id/{action}  # 对资源执行特定动作
```

## 常见错误

### 1. 路由不存在

```bash
GET /api/v1/memorials/recent
→ 404 Not Found 或 认证错误
```

**解决**：使用正确的路由 `/api/v1/users/activities`

### 2. 缺少认证Token

```bash
GET /api/v1/users/activities
→ {"code":1002,"message":"请提供认证令牌"}
```

**解决**：在请求头中添加Token
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/users/activities
```

### 3. 路径参数错误

```bash
GET /api/v1/memorials/  # 缺少ID
→ 返回列表而不是详情
```

**解决**：
- 获取列表：`GET /api/v1/memorials`
- 获取详情：`GET /api/v1/memorials/:id`

## 如何查找正确的路由

### 方法1：查看路由配置

```bash
# 查看路由定义
cat internal/router/router.go
```

### 方法2：查看API文档

访问：http://localhost:8080/api/v1

### 方法3：查看控制器

```bash
# 查找控制器方法
grep -r "func.*Controller" internal/controllers/
```

### 方法4：查看日志

启动服务时会输出所有路由：
```
[GIN-debug] GET    /api/v1/users/activities
[GIN-debug] GET    /api/v1/memorials
...
```

## 测试路由

### 使用curl

```bash
# 健康检查（无需认证）
curl http://localhost:8080/health

# 登录获取Token
curl -X POST http://localhost:8080/api/v1/auth/wechat-login \
  -H "Content-Type: application/json" \
  -d '{"code":"test_code"}'

# 使用Token访问API
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/users/activities
```

### 使用Postman

1. 创建新请求
2. 设置URL：`http://localhost:8080/api/v1/users/activities`
3. 添加Header：`Authorization: Bearer YOUR_TOKEN`
4. 发送请求

## 总结

### `/memorials/recent` 不存在 ❌

这个路由在当前系统中不存在。

### 正确的替代路由 ✅

根据你的需求，应该使用：

1. **用户最近活动**：`GET /api/v1/users/activities`
2. **用户纪念馆列表**：`GET /api/v1/users/memorials`
3. **用户仪表板**：`GET /api/v1/users/dashboard`

### 如何使用

```bash
# 1. 登录获取Token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/wechat-login \
  -H "Content-Type: application/json" \
  -d '{"code":"test"}' | jq -r '.data.token')

# 2. 获取用户最近活动
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/users/activities
```

**建议使用 `/api/v1/users/activities` 来获取用户的最近活动！** ✅
