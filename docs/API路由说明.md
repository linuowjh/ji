# API路由说明

## ✅ 问题已解决

### 原始问题
- **错误路径**: `/api/v1/auth/wechat/login` ❌
- **正确路径**: `/api/v1/auth/wechat-login` ✅

### 修复内容
已修改 `miniprogram/app.js` 中的登录接口路径。

## 🔐 认证相关API

### 微信登录
```
POST /api/v1/auth/wechat-login
```

**请求体**:
```json
{
  "code": "微信登录code"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "jwt_token",
    "user": {
      "id": "user_id",
      "nickname": "用户昵称",
      "avatar": "头像URL"
    }
  }
}
```

**错误响应**（AppID未配置）:
```json
{
  "code": 1005,
  "message": "微信登录失败: invalid appid"
}
```

## 📝 配置微信AppID

要使登录功能正常工作，需要配置微信小程序的AppID和AppSecret。

### 1. 获取微信小程序凭证

1. 登录[微信公众平台](https://mp.weixin.qq.com/)
2. 进入"开发" -> "开发管理" -> "开发设置"
3. 复制AppID和AppSecret

### 2. 配置环境变量

编辑 `.env` 文件：

```bash
# 微信小程序配置
WECHAT_APP_ID=your_wechat_app_id        # 替换为你的AppID
WECHAT_APP_SECRET=your_wechat_app_secret # 替换为你的AppSecret
```

### 3. 重启服务

```bash
# 停止当前服务（Ctrl+C）
# 重新启动
go run cmd/server/main.go
```

## 🧪 测试登录接口

### 使用curl测试

```bash
# 测试接口（会返回AppID错误，这是正常的）
curl -X POST http://localhost:8080/api/v1/auth/wechat-login \
  -H "Content-Type: application/json" \
  -d '{"code":"test_code"}'
```

### 在小程序中测试

```javascript
// 在小程序中调用
wx.login({
  success: res => {
    wx.request({
      url: 'http://localhost:8080/api/v1/auth/wechat-login',
      method: 'POST',
      data: { code: res.code },
      success: response => {
        console.log('登录响应:', response.data)
      }
    })
  }
})
```

## 📋 完整API列表

### 认证相关
- `POST /api/v1/auth/wechat-login` - 微信登录

### 用户相关（需要认证）
- `GET /api/v1/users/profile` - 获取用户信息
- `PUT /api/v1/users/profile` - 更新用户信息
- `GET /api/v1/users/memorials` - 获取用户的纪念馆
- `GET /api/v1/users/worship-records` - 获取用户的祭扫记录

### 纪念馆相关（需要认证）
- `GET /api/v1/memorials` - 获取纪念馆列表
- `POST /api/v1/memorials` - 创建纪念馆
- `GET /api/v1/memorials/:id` - 获取纪念馆详情
- `PUT /api/v1/memorials/:id` - 更新纪念馆
- `DELETE /api/v1/memorials/:id` - 删除纪念馆
- `GET /api/v1/memorials/:id/visitors` - 获取访客记录
- `GET /api/v1/memorials/:id/statistics` - 获取统计信息

### 祭扫相关（需要认证）
- `POST /api/v1/worship` - 创建祭扫记录
- `GET /api/v1/worship/memorials/:memorial_id` - 获取纪念馆祭扫记录
- `GET /api/v1/worship/memorials/:memorial_id/statistics` - 获取祭扫统计
- `GET /api/v1/worship/user/history` - 获取用户祭扫历史

### 祈福相关（需要认证）
- `POST /api/v1/prayers` - 创建祈福
- `GET /api/v1/prayers/memorials/:memorial_id` - 获取纪念馆祈福列表
- `PUT /api/v1/prayers/:id` - 更新祈福
- `DELETE /api/v1/prayers/:id` - 删除祈福

### 留言相关（需要认证）
- `POST /api/v1/messages` - 创建留言
- `GET /api/v1/messages/memorials/:memorial_id` - 获取纪念馆留言
- `PUT /api/v1/messages/:id` - 更新留言
- `DELETE /api/v1/messages/:id` - 删除留言

### 家族圈相关（需要认证）
- `GET /api/v1/families` - 获取家族列表
- `POST /api/v1/families` - 创建家族
- `GET /api/v1/families/:id` - 获取家族详情
- `PUT /api/v1/families/:id` - 更新家族
- `DELETE /api/v1/families/:id` - 删除家族
- `GET /api/v1/families/:id/members` - 获取家族成员
- `POST /api/v1/families/:id/invite` - 邀请成员

### 相册相关（需要认证）
- `POST /api/v1/albums/memorials/:memorial_id` - 创建相册
- `GET /api/v1/albums/memorials/:memorial_id` - 获取相册列表
- `GET /api/v1/albums/:id` - 获取相册详情
- `PUT /api/v1/albums/:id` - 更新相册
- `DELETE /api/v1/albums/:id` - 删除相册
- `POST /api/v1/albums/:id/photos` - 添加照片

### 生平故事相关（需要认证）
- `POST /api/v1/stories/memorials/:memorial_id` - 创建生平故事
- `GET /api/v1/stories/memorials/:memorial_id` - 获取生平故事列表
- `GET /api/v1/stories/:id` - 获取故事详情
- `PUT /api/v1/stories/:id` - 更新故事
- `DELETE /api/v1/stories/:id` - 删除故事

### 追思会相关（需要认证）
- `POST /api/v1/memorial-services/memorials/:memorial_id` - 创建追思会
- `GET /api/v1/memorial-services/memorials/:memorial_id` - 获取追思会列表
- `GET /api/v1/memorial-services/:id` - 获取追思会详情
- `POST /api/v1/memorial-services/:id/start` - 开始追思会
- `POST /api/v1/memorial-services/:id/end` - 结束追思会
- `POST /api/v1/memorial-services/:id/join` - 加入追思会
- `POST /api/v1/memorial-services/:id/leave` - 离开追思会

### 隐私设置相关（需要认证）
- `POST /api/v1/privacy/memorials/settings` - 设置纪念馆隐私
- `GET /api/v1/privacy/memorials/:memorial_id/settings` - 获取隐私设置
- `GET /api/v1/privacy/memorials/:memorial_id/access` - 检查访问权限
- `POST /api/v1/privacy/memorials/:memorial_id/request-access` - 请求访问权限

### 管理员相关（需要管理员权限）
- `GET /api/v1/admin/users` - 获取用户列表
- `GET /api/v1/admin/users/:id` - 获取用户详情
- `POST /api/v1/admin/users/manage` - 管理用户
- `GET /api/v1/admin/content/pending` - 获取待审核内容
- `POST /api/v1/admin/content/moderate` - 审核内容
- `GET /api/v1/admin/stats` - 获取系统统计

## 🔑 认证说明

### JWT Token

大部分API需要在请求头中携带JWT Token：

```
Authorization: Bearer <your_jwt_token>
```

### 获取Token

通过微信登录接口获取Token：

```bash
POST /api/v1/auth/wechat-login
```

### Token使用示例

```javascript
wx.request({
  url: 'http://localhost:8080/api/v1/memorials',
  method: 'GET',
  header: {
    'Authorization': `Bearer ${token}`
  },
  success: res => {
    console.log(res.data)
  }
})
```

## 📊 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

### 错误响应

```json
{
  "code": 1001,
  "message": "错误信息"
}
```

### 错误码说明

- `0` - 成功
- `1001` - 参数错误
- `1002` - 未授权
- `1003` - 禁止访问
- `1004` - 资源不存在
- `1005` - 业务逻辑错误
- `1006` - 服务器内部错误

## 🔄 小程序API封装

小程序中已封装了统一的请求方法，位于 `miniprogram/utils/api.js`：

```javascript
import { request } from '../../utils/api'

// GET请求
request('/memorials', 'GET')
  .then(data => console.log(data))
  .catch(err => console.error(err))

// POST请求
request('/memorials', 'POST', {
  name: '纪念馆名称',
  description: '描述'
})
  .then(data => console.log(data))
  .catch(err => console.error(err))
```

## ✅ 验证清单

- [x] 登录接口路径已修正
- [x] 小程序代码已更新
- [x] 接口可以正常访问
- [ ] 配置微信AppID（需要真实的小程序凭证）
- [ ] 测试完整的登录流程

## 🎯 下一步

1. **配置微信小程序凭证**
   - 获取AppID和AppSecret
   - 更新 `.env` 文件
   - 重启服务

2. **测试登录功能**
   - 在小程序中测试登录
   - 验证Token获取
   - 测试需要认证的API

3. **开发其他功能**
   - 使用已有的API
   - 开发新的业务逻辑

**API路由问题已解决！** ✅
