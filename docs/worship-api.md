# 祭扫功能 API 文档

## 概述

祭扫功能提供了完整的线上祭扫体验，包括传统祭扫仪式（献花、点烛、上香、供奉）和现代化的祈福留言功能。

## API 接口

### 1. 献花功能

**接口地址：** `POST /api/v1/worship/memorials/{memorial_id}/flowers`

**请求参数：**
```json
{
  "flower_type": "chrysanthemum",  // 花卉类型：chrysanthemum|carnation|lily|rose
  "quantity": 3,                   // 数量
  "message": "愿您在天堂安好",      // 献花留言
  "is_scheduled": false,           // 是否定时送花
  "schedule_time": ""              // 定时时间（格式：2006-01-02 15:04:05）
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "献花成功",
  "data": null
}
```

### 2. 点烛功能

**接口地址：** `POST /api/v1/worship/memorials/{memorial_id}/candles`

**请求参数：**
```json
{
  "candle_type": "red",    // 蜡烛类型：red|white|yellow
  "duration": 60,          // 燃烧时长（分钟）
  "message": "为您点亮心灯"  // 点烛留言
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "点烛成功",
  "data": null
}
```

### 3. 续烛功能

**接口地址：** `PUT /api/v1/worship/memorials/{memorial_id}/candles/renew`

**请求参数：**
```json
{
  "additional_minutes": 30  // 续烛时长（分钟）
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "续烛成功",
  "data": null
}
```

### 4. 获取蜡烛状态

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/candles/status`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "active_candles_count": 2,
    "active_candles": [
      {
        "user_id": "user-123",
        "user_name": "张三",
        "candle_type": "red",
        "expire_time": "2024-01-01 15:30:00",
        "message": "为您点亮心灯",
        "lit_at": "2024-01-01 14:30:00"
      }
    ]
  }
}
```

### 5. 上香功能

**接口地址：** `POST /api/v1/worship/memorials/{memorial_id}/incense`

**请求参数：**
```json
{
  "incense_count": 3,        // 香柱数量：3或9
  "incense_type": "sandalwood", // 香的类型：sandalwood|agarwood|traditional
  "message": "愿您安息"       // 上香留言
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "上香成功",
  "data": null
}
```

### 6. 供奉供品

**接口地址：** `POST /api/v1/worship/memorials/{memorial_id}/tributes`

**请求参数：**
```json
{
  "tribute_type": "fruit",           // 供品类型：fruit|pastry|wine|tea
  "items": ["苹果", "香蕉", "橘子"], // 具体供品项目
  "message": "供奉您喜爱的水果"       // 供奉留言
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "供奉成功",
  "data": null
}
```

### 7. 创建祈福

**接口地址：** `POST /api/v1/worship/memorials/{memorial_id}/prayers`

**请求参数：**
```json
{
  "content": "愿您在天堂安好，保佑家人平安健康", // 祈福内容
  "is_public": true                              // 是否公开显示
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "祈福成功",
  "data": {
    "id": "prayer-123",
    "memorial_id": "memorial-456",
    "user_id": "user-789",
    "content": "愿您在天堂安好，保佑家人平安健康",
    "is_public": true,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 8. 创建留言

**接口地址：** `POST /api/v1/worship/memorials/{memorial_id}/messages`

**请求参数：**
```json
{
  "message_type": "text",        // 留言类型：text|audio|video
  "content": "想念您的音容笑貌", // 文字内容
  "media_url": "",               // 音频/视频URL
  "duration": 0                  // 音频/视频时长
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "留言成功",
  "data": {
    "id": "message-123",
    "memorial_id": "memorial-456",
    "user_id": "user-789",
    "message_type": "text",
    "content": "想念您的音容笑貌",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 9. 获取祭扫记录

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/records`

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
        "id": "record-123",
        "memorial_id": "memorial-456",
        "user_id": "user-789",
        "worship_type": "flower",
        "content": "{\"flower_type\":\"chrysanthemum\",\"quantity\":3,\"message\":\"愿您安好\"}",
        "created_at": "2024-01-01T10:00:00Z",
        "user": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        }
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

### 10. 获取祈福墙

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/prayer-wall`

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
        "id": "prayer-123",
        "memorial_id": "memorial-456",
        "user_id": "user-789",
        "content": "愿您在天堂安好",
        "is_public": true,
        "created_at": "2024-01-01T10:00:00Z",
        "user": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        }
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10
  }
}
```

### 11. 获取时光信箱

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/time-messages`

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
        "id": "message-123",
        "memorial_id": "memorial-456",
        "user_id": "user-789",
        "message_type": "audio",
        "content": "",
        "media_url": "https://example.com/audio.mp3",
        "duration": 60,
        "created_at": "2024-01-01T10:00:00Z",
        "user": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        }
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 10
  }
}
```

### 12. 获取祭扫统计

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/statistics`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "flower_count": 25,      // 献花次数
    "candle_count": 18,      // 点烛次数
    "incense_count": 12,     // 上香次数
    "tribute_count": 8,      // 供奉次数
    "prayer_count": 30,      // 祈福次数
    "message_count": 15,     // 留言次数
    "total_visits": 108,     // 总访问次数
    "unique_visitors": 45,   // 独立访客数
    "recent_visits": 20      // 最近7天访问次数
  }
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 请求参数错误 |
| 1002 | 用户未登录 |
| 1005 | 服务器内部错误 |
| 3001 | 纪念馆不存在 |
| 3002 | 无权访问此纪念馆 |

## 使用说明

### 花卉类型说明
- `chrysanthemum`: 菊花（传统祭扫花卉）
- `carnation`: 康乃馨（表达思念）
- `lily`: 百合花（纯洁高雅）
- `rose`: 玫瑰花（爱与怀念）

### 蜡烛类型说明
- `red`: 红蜡烛（传统祭祀）
- `white`: 白蜡烛（纯洁哀思）
- `yellow`: 黄蜡烛（温暖光明）

### 香的类型说明
- `sandalwood`: 檀香（清香淡雅）
- `agarwood`: 沉香（珍贵香料）
- `traditional`: 传统香（普通祭祀香）

### 供品类型说明
- `fruit`: 水果（苹果、香蕉、橘子等）
- `pastry`: 糕点（月饼、点心等）
- `wine`: 酒类（白酒、红酒等）
- `tea`: 茶类（绿茶、红茶等）

### 定时送花功能
用户可以设置定时送花，系统会在指定时间自动执行献花操作。适用于重要纪念日、生日、忌日等特殊时刻。

### 续烛功能
当用户点燃的蜡烛即将熄灭时，可以使用续烛功能延长燃烧时间，保持对逝者的持续缅怀。

### 祈福墙
公开的祈福内容会显示在祈福墙上，让更多人看到对逝者的美好祝愿。用户可以选择是否公开自己的祈福内容。

### 时光信箱
支持文字、语音、视频三种形式的留言，为用户提供多样化的情感表达方式。所有留言都会永久保存，成为珍贵的回忆。
## 
高级统计分析接口

### 13. 获取详细祭扫统计

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/detailed-statistics`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_records": 150,
    "unique_visitors": 25,
    "type_statistics": {
      "flower": 45,
      "candle": 30,
      "incense": 25,
      "tribute": 20,
      "prayer": 20,
      "message": 10
    },
    "monthly_trend": [
      {
        "month": "2024-01",
        "count": 20
      },
      {
        "month": "2024-02",
        "count": 35
      }
    ],
    "hourly_pattern": [
      {
        "hour": 0,
        "count": 2
      },
      {
        "hour": 9,
        "count": 15
      }
    ],
    "top_visitors": [
      {
        "user_id": "user-123",
        "user_name": "张三",
        "avatar_url": "https://example.com/avatar.jpg",
        "count": 25,
        "last_visit": "2024-01-01 10:00:00"
      }
    ],
    "recent_activity": [
      {
        "date": "2024-01-01",
        "worship_count": 5,
        "visitor_count": 3
      }
    ]
  }
}
```

### 14. 用户祭扫行为分析

**接口地址：** `GET /api/v1/worship/user/behavior-analysis`

**响应示例：**
```json
{
  "code": 0,
  "message": "分析完成",
  "data": {
    "user_id": "user-123",
    "total_worships": 50,
    "favorite_type": "flower",
    "active_hours": [9, 14, 20],
    "worship_frequency": {
      "flower": 20,
      "candle": 15,
      "incense": 10,
      "tribute": 3,
      "prayer": 2,
      "message": 0
    },
    "memorial_count": 5,
    "first_worship": "2023-06-01 10:00:00",
    "last_worship": "2024-01-01 15:30:00"
  }
}
```

### 15. 生成祭扫报告

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/report`

**查询参数：**
- `period`: 统计周期（week|month|quarter|year，默认month）

**响应示例：**
```json
{
  "code": 0,
  "message": "报告生成成功",
  "data": {
    "memorial_id": "memorial-456",
    "memorial_name": "张老爷子",
    "report_period": "month",
    "summary": {
      "total_records": 45,
      "unique_visitors": 12,
      "type_statistics": {
        "flower": 15,
        "candle": 10,
        "incense": 8,
        "tribute": 5,
        "prayer": 4,
        "message": 3
      },
      "period_start": "2023-12-01",
      "period_end": "2024-01-01"
    },
    "highlights": [
      "本月共收到45次祭扫",
      "有12位访客表达了思念",
      "最受欢迎的祭扫方式是献花"
    ],
    "recommendations": [
      "感谢大家的参与，让我们继续传承这份美好的纪念"
    ],
    "generated_at": "2024-01-01 16:00:00"
  }
}
```

## 智能功能接口

### 16. 祈福卡模板

**接口地址：** `GET /api/v1/worship/prayer-card-templates`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "template-1",
      "name": "传统祈福卡",
      "description": "古典雅致的传统祈福卡样式",
      "image_url": "/static/templates/traditional-prayer-card.jpg",
      "category": "traditional"
    }
  ]
}
```

### 17. 生成祈福卡

**接口地址：** `POST /api/v1/worship/generate-prayer-card`

**请求参数：**
```json
{
  "template_id": "template-1",
  "content": "愿您在天堂安好",
  "user_name": "张三"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "生成成功",
  "data": {
    "card_url": "/generated/prayer-cards/template-1-1704096000.jpg"
  }
}
```

### 18. 情感分析

**接口地址：** `POST /api/v1/worship/analyze-emotion`

**请求参数：**
```json
{
  "content": "想念您的音容笑貌，愿您在天堂安好"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "分析完成",
  "data": {
    "emotion": "nostalgic",
    "confidence": 0.85,
    "keywords": ["想念", "天堂", "安好"],
    "suggestion": "往昔的美好时光值得永远珍藏，这些回忆是您与逝者之间最珍贵的纽带。"
  }
}
```

### 19. 获取回复建议

**接口地址：** `GET /api/v1/worship/reply-suggestions`

**查询参数：**
- `message_type`: 留言类型（text|audio|video）
- `content`: 留言内容（可选）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "suggestions": [
      "您的话语充满了爱与思念",
      "相信逝者能感受到您的真挚情感",
      "这份深情让人动容"
    ]
  }
}
```

### 20. 获取留言创建提示

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/message-tips`

**查询参数：**
- `message_type`: 留言类型（text|audio|video）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "tips": [
      {
        "type": "suggestion",
        "content": "可以分享一些与逝者的美好回忆，或者表达您的思念之情"
      },
      {
        "type": "reminder",
        "content": "文字会永久保存，成为珍贵的纪念"
      },
      {
        "type": "encouragement",
        "content": "每一份真挚的表达都是对逝者最好的纪念"
      }
    ]
  }
}
```

### 21. 获取热门祈福内容

**接口地址：** `GET /api/v1/worship/popular-prayer-contents`

**查询参数：**
- `limit`: 返回数量限制（默认10，最大50）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "contents": [
      "愿您在天堂安好，保佑家人平安健康",
      "思念如潮水，愿您安息",
      "您的音容笑貌永远在我心中"
    ]
  }
}
```

### 22. 创建定时祈福

**接口地址：** `POST /api/v1/worship/scheduled-prayers`

**请求参数：**
```json
{
  "memorial_id": "memorial-456",
  "content": "愿您安息",
  "schedule_time": "2024-01-01T10:00:00Z",
  "is_recurring": true,
  "recurring_type": "yearly"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "定时祈福创建成功",
  "data": null
}
```

### 23. 纪念馆留言分析

**接口地址：** `GET /api/v1/worship/memorials/{memorial_id}/message-analytics`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "message_types": {
      "text": 25,
      "audio": 8,
      "video": 3
    },
    "daily_stats": {
      "2024-01-01": 5,
      "2024-01-02": 3
    },
    "active_users": 12,
    "total_messages": 36
  }
}
```

### 24. 内容审核

**接口地址：** `POST /api/v1/worship/messages/{message_id}/moderate`

**响应示例：**
```json
{
  "code": 0,
  "message": "审核完成",
  "data": {
    "is_approved": true,
    "reason": ""
  }
}
```

## 功能特色

### 智能情感分析
系统能够分析用户留言的情感倾向，识别悲伤、怀念、感恩等情感，并提供相应的回复建议，帮助用户更好地表达情感。

### 个性化祈福卡
提供多种精美的祈福卡模板，用户可以选择合适的样式生成个性化的祈福卡片，增强纪念仪式感。

### 行为分析洞察
通过分析用户的祭扫行为模式，了解用户偏好和活跃时段，为产品优化提供数据支持。

### 智能内容推荐
基于用户行为和情感分析，推荐合适的祈福内容和留言建议，降低用户表达门槛。

### 定时祈福提醒
支持设置定时祈福，在重要纪念日自动提醒用户参与祭扫活动，保持持续的纪念传承。

### 数据可视化报告
生成详细的祭扫统计报告，包含趋势分析、访客统计等，帮助用户了解纪念馆的活跃情况。

## 使用建议

### 最佳实践
1. **合理使用定时功能**：为重要纪念日设置定时祈福，保持纪念的连续性
2. **多样化表达方式**：结合文字、语音、视频等多种形式，丰富纪念内容
3. **关注情感健康**：利用情感分析功能，关注用户的心理状态
4. **数据驱动优化**：通过统计分析了解用户需求，持续改进服务

### 注意事项
1. **内容审核**：所有用户生成内容都会经过审核，确保内容健康正面
2. **隐私保护**：个人祭扫记录和行为分析数据严格保密
3. **技术限制**：部分AI功能可能存在识别误差，仅供参考
4. **服务稳定性**：定时功能依赖系统稳定运行，建议关注服务状态