# 线上追思会 API 文档

## 概述

线上追思会功能为用户提供了完整的在线纪念活动体验，包括追思会创建、邀请管理、实时互动和录制回放等功能，让分散各地的亲友能够共同参与纪念活动。

## 追思会管理 API

### 1. 创建追思会

**接口地址：** `POST /api/v1/memorial-services/memorials/{memorial_id}`

**请求参数：**
```json
{
  "title": "张老爷子追思会",
  "description": "缅怀张老爷子的一生，分享美好回忆",
  "start_time": "2024-01-15T14:00:00Z",
  "end_time": "2024-01-15T16:00:00Z",
  "max_participants": 50,
  "is_public": false
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "service-123",
    "memorial_id": "memorial-456",
    "title": "张老爷子追思会",
    "description": "缅怀张老爷子的一生，分享美好回忆",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T16:00:00Z",
    "status": "scheduled",
    "max_participants": 50,
    "is_public": false,
    "invite_code": "ABC12345",
    "host_id": "user-789",
    "host": {
      "id": "user-789",
      "nickname": "张三",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "memorial": {
      "id": "memorial-456",
      "deceased_name": "张老爷子"
    },
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 2. 获取追思会列表

**接口地址：** `GET /api/v1/memorial-services/memorials/{memorial_id}`

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10，最大100）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "service-123",
        "memorial_id": "memorial-456",
        "title": "张老爷子追思会",
        "description": "缅怀张老爷子的一生，分享美好回忆",
        "start_time": "2024-01-15T14:00:00Z",
        "end_time": "2024-01-15T16:00:00Z",
        "status": "scheduled",
        "max_participants": 50,
        "is_public": false,
        "invite_code": "ABC12345",
        "host": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        },
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 3,
    "page": 1,
    "page_size": 10
  }
}
```

### 3. 获取追思会详情

**接口地址：** `GET /api/v1/memorial-services/{service_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "service-123",
    "memorial_id": "memorial-456",
    "title": "张老爷子追思会",
    "description": "缅怀张老爷子的一生，分享美好回忆",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T16:00:00Z",
    "status": "scheduled",
    "max_participants": 50,
    "is_public": false,
    "invite_code": "ABC12345",
    "recording_url": "",
    "host": {
      "id": "user-789",
      "nickname": "张三",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "memorial": {
      "id": "memorial-456",
      "deceased_name": "张老爷子",
      "avatar_url": "https://example.com/memorial-avatar.jpg"
    },
    "participants": [
      {
        "id": "participant-456",
        "service_id": "service-123",
        "user_id": "user-789",
        "role": "host",
        "status": "joined",
        "joined_at": "2024-01-01T10:00:00Z",
        "user": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        }
      }
    ],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 4. 更新追思会

**接口地址：** `PUT /api/v1/memorial-services/{service_id}`

**请求参数：**
```json
{
  "title": "张老爷子纪念追思会",
  "description": "更新后的描述",
  "start_time": "2024-01-15T15:00:00Z",
  "end_time": "2024-01-15T17:00:00Z",
  "max_participants": 100,
  "is_public": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "更新成功"
}
```

### 5. 删除追思会

**接口地址：** `DELETE /api/v1/memorial-services/{service_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "删除成功"
}
```

## 追思会控制 API

### 6. 开始追思会

**接口地址：** `POST /api/v1/memorial-services/{service_id}/start`

**响应示例：**
```json
{
  "code": 0,
  "message": "追思会已开始"
}
```

### 7. 结束追思会

**接口地址：** `POST /api/v1/memorial-services/{service_id}/end`

**响应示例：**
```json
{
  "code": 0,
  "message": "追思会已结束，正在生成录制视频"
}
```

### 8. 加入追思会

**接口地址：** `POST /api/v1/memorial-services/{service_id}/join`

**响应示例：**
```json
{
  "code": 0,
  "message": "加入成功"
}
```

### 9. 离开追思会

**接口地址：** `POST /api/v1/memorial-services/{service_id}/leave`

**响应示例：**
```json
{
  "code": 0,
  "message": "离开成功"
}
```

## 参与者管理 API

### 10. 邀请参与者

**接口地址：** `POST /api/v1/memorial-services/{service_id}/invite`

**请求参数：**
```json
{
  "user_ids": ["user-123", "user-456", "user-789"],
  "message": "诚邀您参加张老爷子的追思会，共同缅怀他的一生"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "邀请发送成功"
}
```

### 11. 响应邀请

**接口地址：** `POST /api/v1/memorial-services/invitations/{invitation_id}/respond`

**请求参数：**
```json
{
  "accept": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "邀请已接受"
}
```

## 聊天功能 API

### 12. 发送聊天消息

**接口地址：** `POST /api/v1/memorial-services/{service_id}/chat`

**请求参数：**
```json
{
  "message_type": "text",
  "content": "感谢大家参加追思会，让我们一起缅怀张老爷子"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "发送成功",
  "data": {
    "id": "chat-123",
    "service_id": "service-456",
    "user_id": "user-789",
    "message_type": "text",
    "content": "感谢大家参加追思会，让我们一起缅怀张老爷子",
    "timestamp": "2024-01-15T14:30:00Z",
    "user": {
      "id": "user-789",
      "nickname": "张三",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "created_at": "2024-01-15T14:30:00Z"
  }
}
```

### 13. 获取聊天消息

**接口地址：** `GET /api/v1/memorial-services/{service_id}/chat`

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认50，最大100）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": "chat-123",
        "service_id": "service-456",
        "user_id": "user-789",
        "message_type": "text",
        "content": "感谢大家参加追思会，让我们一起缅怀张老爷子",
        "timestamp": "2024-01-15T14:30:00Z",
        "user": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        }
      },
      {
        "id": "chat-124",
        "service_id": "service-456",
        "user_id": "user-456",
        "message_type": "text",
        "content": "张老爷子是一位慈祥的长者，我们永远怀念他",
        "timestamp": "2024-01-15T14:32:00Z",
        "user": {
          "id": "user-456",
          "nickname": "李四",
          "avatar_url": "https://example.com/avatar2.jpg"
        }
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 50
  }
}
```

## 数据字段说明

### 追思会状态 (status)
- `scheduled`: 已安排（未开始）
- `ongoing`: 进行中
- `completed`: 已完成
- `cancelled`: 已取消

### 参与者角色 (role)
- `host`: 主持人
- `co_host`: 协助主持人
- `participant`: 参与者

### 参与者状态 (participant status)
- `invited`: 已邀请
- `joined`: 已加入
- `left`: 已离开
- `removed`: 已移除

### 邀请状态 (invitation status)
- `pending`: 待处理
- `accepted`: 已接受
- `declined`: 已拒绝
- `expired`: 已过期

### 消息类型 (message_type)
- `text`: 文字消息
- `image`: 图片消息
- `emoji`: 表情消息

### 活动类型 (activity_type)
- `join`: 加入追思会
- `leave`: 离开追思会
- `worship`: 祭扫活动
- `speak`: 发言
- `share_screen`: 屏幕分享

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 请求参数错误 |
| 1002 | 用户未登录 |
| 1003 | 无权限操作 |
| 1004 | 资源不存在 |
| 1005 | 服务器内部错误 |
| 3001 | 纪念馆不存在 |
| 3002 | 无权访问此纪念馆 |

## 功能特色

### 1. 完整的追思会生命周期管理
- **创建与配置**：支持设置追思会时间、参与人数限制、公开性等
- **状态管理**：从安排到进行中再到完成的完整状态流转
- **权限控制**：主持人、协助主持人、参与者的分级权限管理

### 2. 灵活的邀请系统
- **批量邀请**：支持一次邀请多个用户
- **邀请码分享**：生成唯一邀请码便于分享
- **邀请响应**：被邀请者可以接受或拒绝邀请
- **过期机制**：邀请在追思会开始前1小时自动过期

### 3. 实时互动功能
- **聊天系统**：支持文字、图片、表情等多种消息类型
- **活动记录**：记录所有参与者的活动轨迹
- **同步祭扫**：在追思会中进行的祭扫活动会实时同步

### 4. 录制与回放
- **自动录制**：追思会结束后自动生成录制视频
- **异步处理**：录制视频生成采用异步处理，不影响用户体验
- **永久保存**：录制视频永久保存，供后续回看

### 5. 智能化管理
- **参与者状态跟踪**：实时跟踪参与者的加入、离开状态
- **时间管理**：自动处理追思会的开始和结束时间
- **容量控制**：支持设置最大参与人数限制

## 使用场景

### 1. 家族追思会
- 邀请分散各地的家族成员参与
- 共同分享对逝者的回忆和思念
- 进行集体祭扫活动

### 2. 朋友纪念会
- 同学、同事、朋友的纪念聚会
- 分享工作、学习中的美好回忆
- 表达对逝者的怀念之情

### 3. 周年纪念
- 逝世周年纪念活动
- 生日纪念活动
- 特殊节日纪念

## 最佳实践

### 1. 追思会规划
- **提前安排**：建议提前1-2周创建追思会并发送邀请
- **时间选择**：选择大多数参与者都方便的时间
- **人数控制**：根据纪念馆的重要性和影响力设置合适的参与人数

### 2. 邀请管理
- **个性化邀请**：为不同的参与者群体编写不同的邀请消息
- **及时跟进**：关注邀请响应情况，必要时进行电话确认
- **备用方案**：为重要参与者准备邀请码等备用加入方式

### 3. 活动主持
- **开场引导**：主持人应在开始时介绍追思会流程和注意事项
- **氛围营造**：通过分享回忆、播放音乐等方式营造庄重温馨的氛围
- **互动引导**：鼓励参与者分享回忆、表达思念

### 4. 技术准备
- **网络测试**：提前测试网络连接和设备功能
- **备用设备**：准备备用设备以防技术故障
- **录制确认**：确认录制功能正常，为无法参与的亲友保留回放

## 注意事项

### 1. 隐私保护
- 非公开追思会只有受邀者可以参与
- 聊天记录和录制视频仅参与者可见
- 个人信息严格保护，不会泄露给第三方

### 2. 内容审核
- 所有聊天消息都会经过内容审核
- 不当言论会被自动过滤或人工审核
- 维护庄重、和谐的纪念氛围

### 3. 技术限制
- 追思会最大支持100人同时在线
- 录制视频大小和时长有一定限制
- 网络不稳定可能影响参与体验

### 4. 服务稳定性
- 系统会定期维护以确保服务稳定
- 重要追思会建议提前测试功能
- 遇到技术问题可联系客服支持