# 纪念相册和生平故事 API 文档

## 概述

纪念相册和生平故事功能为用户提供了丰富的纪念内容管理能力，包括照片相册管理、生平故事记录和时间轴展示，让纪念更加生动和完整。

## 纪念相册 API

### 1. 创建相册

**接口地址：** `POST /api/v1/albums/memorials/{memorial_id}`

**请求参数：**
```json
{
  "title": "童年时光",
  "description": "记录美好的童年回忆",
  "cover_url": "https://example.com/cover.jpg",
  "is_public": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "album-123",
    "memorial_id": "memorial-456",
    "title": "童年时光",
    "description": "记录美好的童年回忆",
    "cover_url": "https://example.com/cover.jpg",
    "is_public": true,
    "sort_order": 0,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

### 2. 获取相册列表

**接口地址：** `GET /api/v1/albums/memorials/{memorial_id}`

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
        "id": "album-123",
        "memorial_id": "memorial-456",
        "title": "童年时光",
        "description": "记录美好的童年回忆",
        "cover_url": "https://example.com/cover.jpg",
        "is_public": true,
        "sort_order": 0,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10
  }
}
```

### 3. 获取相册详情

**接口地址：** `GET /api/v1/albums/{album_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "album-123",
    "memorial_id": "memorial-456",
    "title": "童年时光",
    "description": "记录美好的童年回忆",
    "cover_url": "https://example.com/cover.jpg",
    "is_public": true,
    "photos": [
      {
        "id": "photo-789",
        "album_id": "album-123",
        "photo_url": "https://example.com/photo1.jpg",
        "thumbnail_url": "https://example.com/thumb1.jpg",
        "caption": "5岁生日照片",
        "taken_date": "1990-06-15T00:00:00Z",
        "location": "北京",
        "sort_order": 0,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 4. 更新相册

**接口地址：** `PUT /api/v1/albums/{album_id}`

**请求参数：**
```json
{
  "title": "美好童年",
  "description": "更新后的描述",
  "cover_url": "https://example.com/new-cover.jpg",
  "is_public": false
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "更新成功"
}
```

### 5. 删除相册

**接口地址：** `DELETE /api/v1/albums/{album_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "删除成功"
}
```

### 6. 添加照片到相册

**接口地址：** `POST /api/v1/albums/{album_id}/photos`

**请求参数：**
```json
{
  "photo_url": "https://example.com/photo.jpg",
  "thumbnail_url": "https://example.com/thumb.jpg",
  "caption": "这是一张珍贵的照片",
  "taken_date": "1990-06-15T00:00:00Z",
  "location": "上海"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "添加成功",
  "data": {
    "id": "photo-789",
    "album_id": "album-123",
    "photo_url": "https://example.com/photo.jpg",
    "thumbnail_url": "https://example.com/thumb.jpg",
    "caption": "这是一张珍贵的照片",
    "taken_date": "1990-06-15T00:00:00Z",
    "location": "上海",
    "sort_order": 0,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 7. 更新照片信息

**接口地址：** `PUT /api/v1/albums/photos/{photo_id}`

**请求参数：**
```json
{
  "caption": "更新后的照片说明",
  "taken_date": "1990-06-15T00:00:00Z",
  "location": "广州",
  "thumbnail_url": "https://example.com/new-thumb.jpg"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "更新成功"
}
```

### 8. 删除照片

**接口地址：** `DELETE /api/v1/albums/photos/{photo_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "删除成功"
}
```

## 生平故事 API

### 9. 创建生平故事

**接口地址：** `POST /api/v1/stories/memorials/{memorial_id}`

**请求参数：**
```json
{
  "title": "求学时光",
  "content": "那是一个充满希望的年代，他怀着对知识的渴望踏进了校园...",
  "story_date": "1985-09-01T00:00:00Z",
  "age_at_time": 18,
  "location": "北京大学",
  "category": "education",
  "is_public": true,
  "media": [
    {
      "media_type": "image",
      "media_url": "https://example.com/school.jpg",
      "thumbnail_url": "https://example.com/school-thumb.jpg",
      "caption": "大学校园照片"
    }
  ]
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "story-123",
    "memorial_id": "memorial-456",
    "title": "求学时光",
    "content": "那是一个充满希望的年代，他怀着对知识的渴望踏进了校园...",
    "story_date": "1985-09-01T00:00:00Z",
    "age_at_time": 18,
    "location": "北京大学",
    "category": "education",
    "is_public": true,
    "author_id": "user-789",
    "author": {
      "id": "user-789",
      "nickname": "张三",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "media": [
      {
        "id": "media-456",
        "life_story_id": "story-123",
        "media_type": "image",
        "media_url": "https://example.com/school.jpg",
        "thumbnail_url": "https://example.com/school-thumb.jpg",
        "caption": "大学校园照片",
        "sort_order": 0
      }
    ],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 10. 获取生平故事列表

**接口地址：** `GET /api/v1/stories/memorials/{memorial_id}`

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
        "id": "story-123",
        "memorial_id": "memorial-456",
        "title": "求学时光",
        "content": "那是一个充满希望的年代...",
        "story_date": "1985-09-01T00:00:00Z",
        "age_at_time": 18,
        "location": "北京大学",
        "category": "education",
        "is_public": true,
        "author": {
          "id": "user-789",
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar.jpg"
        },
        "media": [
          {
            "id": "media-456",
            "media_type": "image",
            "media_url": "https://example.com/school.jpg",
            "thumbnail_url": "https://example.com/school-thumb.jpg",
            "caption": "大学校园照片"
          }
        ],
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 8,
    "page": 1,
    "page_size": 10
  }
}
```

### 11. 按分类获取生平故事

**接口地址：** `GET /api/v1/stories/memorials/{memorial_id}/by-category`

**查询参数：**
- `category`: 故事分类（childhood|youth|career|family|achievement|other）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "story-123",
      "title": "求学时光",
      "content": "那是一个充满希望的年代...",
      "category": "education",
      "story_date": "1985-09-01T00:00:00Z",
      "author": {
        "nickname": "张三"
      },
      "media": []
    }
  ]
}
```

### 12. 获取生平故事详情

**接口地址：** `GET /api/v1/stories/{story_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "story-123",
    "memorial_id": "memorial-456",
    "title": "求学时光",
    "content": "那是一个充满希望的年代，他怀着对知识的渴望踏进了校园...",
    "story_date": "1985-09-01T00:00:00Z",
    "age_at_time": 18,
    "location": "北京大学",
    "category": "education",
    "is_public": true,
    "author": {
      "id": "user-789",
      "nickname": "张三",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "media": [
      {
        "id": "media-456",
        "media_type": "image",
        "media_url": "https://example.com/school.jpg",
        "thumbnail_url": "https://example.com/school-thumb.jpg",
        "caption": "大学校园照片",
        "sort_order": 0
      }
    ],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 13. 更新生平故事

**接口地址：** `PUT /api/v1/stories/{story_id}`

**请求参数：**
```json
{
  "title": "大学求学时光",
  "content": "更新后的故事内容...",
  "story_date": "1985-09-01T00:00:00Z",
  "age_at_time": 18,
  "location": "北京大学",
  "category": "education",
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

### 14. 删除生平故事

**接口地址：** `DELETE /api/v1/stories/{story_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "删除成功"
}
```

## 时间轴 API

### 15. 创建时间轴事件

**接口地址：** `POST /api/v1/timelines/memorials/{memorial_id}`

**请求参数：**
```json
{
  "title": "出生",
  "description": "在一个春暖花开的日子里来到这个世界",
  "event_date": "1967-03-15T00:00:00Z",
  "event_type": "birth",
  "icon_url": "https://example.com/birth-icon.png",
  "is_public": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "timeline-123",
    "memorial_id": "memorial-456",
    "title": "出生",
    "description": "在一个春暖花开的日子里来到这个世界",
    "event_date": "1967-03-15T00:00:00Z",
    "event_type": "birth",
    "icon_url": "https://example.com/birth-icon.png",
    "is_public": true,
    "author_id": "user-789",
    "author": {
      "id": "user-789",
      "nickname": "张三",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 16. 获取时间轴

**接口地址：** `GET /api/v1/timelines/memorials/{memorial_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "timeline-123",
      "memorial_id": "memorial-456",
      "title": "出生",
      "description": "在一个春暖花开的日子里来到这个世界",
      "event_date": "1967-03-15T00:00:00Z",
      "event_type": "birth",
      "icon_url": "https://example.com/birth-icon.png",
      "is_public": true,
      "author": {
        "nickname": "张三"
      },
      "created_at": "2024-01-01T10:00:00Z"
    },
    {
      "id": "timeline-124",
      "title": "入学",
      "description": "开始了求学之路",
      "event_date": "1973-09-01T00:00:00Z",
      "event_type": "education",
      "author": {
        "nickname": "李四"
      }
    }
  ]
}
```

### 17. 删除时间轴事件

**接口地址：** `DELETE /api/v1/timelines/{timeline_id}`

**响应示例：**
```json
{
  "code": 0,
  "message": "删除成功"
}
```

## 数据字段说明

### 生平故事分类 (category)
- `childhood`: 童年时光
- `youth`: 青春岁月
- `career`: 事业发展
- `family`: 家庭生活
- `achievement`: 成就荣誉
- `other`: 其他

### 时间轴事件类型 (event_type)
- `birth`: 出生
- `education`: 教育求学
- `career`: 事业发展
- `marriage`: 婚姻家庭
- `achievement`: 重要成就
- `death`: 逝世
- `other`: 其他重要事件

### 媒体文件类型 (media_type)
- `image`: 图片
- `video`: 视频
- `audio`: 音频

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

### 1. 丰富的内容管理
- **多媒体支持**：支持图片、视频、音频等多种媒体格式
- **分类管理**：按照人生阶段和主题对内容进行分类
- **时间排序**：按照时间顺序展示人生轨迹

### 2. 协作编辑
- **多人贡献**：家族成员可以共同添加和编辑内容
- **权限控制**：支持公开和私密内容设置
- **作者标识**：记录每个内容的贡献者

### 3. 时间轴展示
- **可视化时间线**：直观展示人生重要节点
- **事件分类**：不同类型事件使用不同图标标识
- **详细描述**：每个事件都可以添加详细说明

### 4. 相册管理
- **相册分组**：按主题创建不同相册
- **照片标注**：支持为每张照片添加说明和拍摄信息
- **封面设置**：为相册设置精美封面

## 使用建议

### 最佳实践
1. **内容组织**：建议按照时间顺序和主题分类来组织内容
2. **媒体质量**：上传高质量的图片和视频，确保纪念效果
3. **详细描述**：为每个故事和照片添加详细的背景信息
4. **家族参与**：鼓励家族成员共同参与内容创建

### 注意事项
1. **隐私保护**：合理设置内容的公开性，保护隐私
2. **内容审核**：所有内容都会经过审核，确保健康正面
3. **存储限制**：注意媒体文件的大小和数量限制
4. **权限管理**：只有有权限的用户才能编辑和删除内容