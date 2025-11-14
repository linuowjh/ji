# 云念纪念馆数据库设计文档

## 概述

本文档描述了云念纪念馆系统的数据库设计，包括表结构、索引、关系和数据迁移策略。

## 数据库表结构

### 核心业务表

#### 1. 用户表 (users)
存储用户基本信息和微信登录凭证。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| wechat_openid | varchar(100) | 微信OpenID，唯一 |
| wechat_unionid | varchar(100) | 微信UnionID |
| nickname | varchar(50) | 用户昵称 |
| avatar_url | varchar(255) | 头像URL |
| phone | varchar(20) | 手机号 |
| status | tinyint | 状态：1正常，0禁用 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

#### 2. 纪念馆表 (memorials)
存储逝者纪念馆的基本信息。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| creator_id | varchar(36) | 创建者ID，外键 |
| deceased_name | varchar(50) | 逝者姓名 |
| birth_date | date | 出生日期 |
| death_date | date | 逝世日期 |
| biography | text | 生平简介 |
| avatar_url | varchar(255) | 逝者照片URL |
| theme_style | varchar(50) | 主题风格 |
| tombstone_style | varchar(50) | 墓碑样式 |
| epitaph | text | 墓志铭 |
| privacy_level | tinyint | 隐私级别：1家族可见，2私密 |
| status | tinyint | 状态 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

#### 3. 祭扫记录表 (worship_records)
记录用户的祭扫行为。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| user_id | varchar(36) | 用户ID，外键 |
| worship_type | varchar(20) | 祭扫类型：flower/candle/incense/tribute/prayer |
| content | json | 祭扫内容详情 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

### 家族功能表

#### 4. 家族表 (families)
存储家族圈信息。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| name | varchar(100) | 家族名称 |
| creator_id | varchar(36) | 创建者ID，外键 |
| description | text | 家族描述 |
| invite_code | varchar(20) | 邀请码，唯一 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

#### 5. 家族成员表 (family_members)
存储家族成员关系。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| family_id | varchar(36) | 家族ID，外键 |
| user_id | varchar(36) | 用户ID，外键 |
| role | varchar(20) | 角色：admin/member |
| joined_at | timestamp | 加入时间 |

### 内容管理表

#### 6. 媒体文件表 (media_files)
存储纪念馆相关的媒体文件。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| file_type | varchar(20) | 文件类型：image/video/audio |
| file_url | varchar(255) | 文件URL |
| file_name | varchar(255) | 文件名 |
| file_size | bigint | 文件大小 |
| description | text | 文件描述 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

#### 7. 祈福表 (prayers)
存储用户的祈福内容。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| user_id | varchar(36) | 用户ID，外键 |
| content | text | 祈福内容 |
| is_public | tinyint(1) | 是否公开 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

#### 8. 留言表 (messages)
存储用户的语音、视频、文字留言。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| user_id | varchar(36) | 用户ID，外键 |
| message_type | varchar(20) | 留言类型：text/audio/video |
| content | text | 文字内容 |
| media_url | varchar(255) | 媒体文件URL |
| duration | int | 音频/视频时长(秒) |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

### 功能辅助表

#### 9. 纪念日提醒表 (memorial_reminders)
存储纪念日提醒设置。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| reminder_type | varchar(20) | 提醒类型：birthday/death_anniversary/festival |
| reminder_date | date | 提醒日期 |
| title | varchar(100) | 提醒标题 |
| content | text | 提醒内容 |
| is_active | tinyint(1) | 是否激活 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

#### 10. 访客记录表 (visitor_records)
记录纪念馆访客信息。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| visitor_id | varchar(36) | 访客ID，外键 |
| visit_time | timestamp | 访问时间 |
| ip_address | varchar(45) | IP地址 |

#### 11. 纪念馆家族关联表 (memorial_families)
存储纪念馆与家族的关联关系。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键，UUID |
| memorial_id | varchar(36) | 纪念馆ID，外键 |
| family_id | varchar(36) | 家族ID，外键 |
| created_at | timestamp | 创建时间 |

## 索引设计

### 主要索引

1. **用户表索引**
   - `idx_users_wechat_openid` (唯一索引)
   - `idx_users_status`
   - `idx_users_created_at`

2. **纪念馆表索引**
   - `idx_memorials_creator_status` (复合索引)
   - `idx_memorials_privacy_status` (复合索引)
   - `idx_memorials_created_at`

3. **祭扫记录表索引**
   - `idx_worship_memorial_time` (复合索引)
   - `idx_worship_user_time` (复合索引)
   - `idx_worship_type`

4. **家族相关索引**
   - `idx_families_invite_code` (唯一索引)
   - `idx_family_members_family_role` (复合索引)
   - `unique_family_user` (唯一复合索引)

## 外键约束

系统使用外键约束确保数据完整性：

- 所有关联表都设置了适当的外键约束
- 使用 `ON DELETE CASCADE` 确保数据一致性
- 主要关联关系：
  - 纪念馆 → 用户 (创建者)
  - 祭扫记录 → 纪念馆、用户
  - 家族成员 → 家族、用户
  - 媒体文件 → 纪念馆

## 数据迁移系统

### 迁移管理器功能

1. **自动迁移** (`AutoMigrate`)
   - 自动创建和更新表结构
   - 支持增量迁移
   - 保持数据完整性

2. **索引管理** (`CreateIndexes`)
   - 创建性能优化索引
   - 支持复合索引
   - 错误容忍机制

3. **种子数据** (`SeedData`)
   - 插入测试数据
   - 支持开发环境初始化
   - 防重复插入

4. **数据库重置** (`Reset`)
   - 完全重建数据库
   - 仅用于开发环境
   - 包含安全确认

### 使用方法

```bash
# 执行数据库迁移
make migrate

# 迁移并插入种子数据
make migrate-seed

# 仅插入种子数据
make seed

# 重置数据库（危险操作）
make reset-db
```

## 性能优化

### 查询优化

1. **分页查询优化**
   - 使用游标分页避免深度分页
   - 合理使用 LIMIT 和 OFFSET

2. **索引策略**
   - 为常用查询字段创建索引
   - 使用复合索引优化多条件查询
   - 避免过多索引影响写入性能

3. **数据类型优化**
   - 使用适当的数据类型
   - VARCHAR 长度合理设置
   - 使用 TINYINT 替代 BOOLEAN

### 存储优化

1. **软删除策略**
   - 使用 `deleted_at` 字段实现软删除
   - 保留数据完整性
   - 支持数据恢复

2. **JSON 字段使用**
   - 祭扫记录使用 JSON 存储灵活内容
   - 减少表结构复杂度
   - 支持动态字段扩展

## 安全考虑

1. **数据加密**
   - 敏感信息加密存储
   - 使用 HTTPS 传输
   - 数据库连接加密

2. **访问控制**
   - 基于角色的权限控制
   - 纪念馆隐私级别控制
   - 家族成员权限管理

3. **数据备份**
   - 定期数据备份
   - 支持数据导出
   - 灾难恢复计划

## 扩展性设计

1. **水平扩展**
   - 支持读写分离
   - 分库分表准备
   - 缓存层设计

2. **垂直扩展**
   - 模块化表设计
   - 松耦合关系
   - 易于功能扩展

这个数据库设计为云念纪念馆系统提供了稳定、高效、可扩展的数据存储基础。