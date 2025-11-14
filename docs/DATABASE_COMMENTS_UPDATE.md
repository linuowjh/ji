# 数据库字段中文注释更新

## 📋 更新进度

### ✅ 已完成的模型 (4个)

1. **user.go** - 用户模型 ✅
   - 所有字段已添加中文注释
   - 包括：用户ID、微信OpenID、昵称、头像、手机号、状态等

2. **memorial.go** - 纪念馆模型 ✅
   - Memorial - 纪念馆主表
   - MediaFile - 媒体文件表
   - MemorialFamily - 纪念馆家族关联表
   - 所有字段已添加中文注释

3. **worship.go** - 祭扫模型 ✅
   - WorshipRecord - 祭扫记录表
   - Prayer - 祈福表
   - Message - 留言表
   - 所有字段已添加中文注释

### ⏳ 待更新的模型 (7个)

4. **family.go** - 家族模型
   - Family - 家族表
   - FamilyMember - 家族成员表
   - MemorialReminder - 纪念日提醒表
   - FamilyGenealogy - 家族谱系表
   - FamilyStory - 家族故事表
   - FamilyTradition - 家族传统表

5. **visitor.go** - 访客模型
   - VisitorRecord - 访客记录表

6. **album.go** - 相册模型
   - Album - 相册表
   - AlbumPhoto - 相册照片表
   - LifeStory - 生平故事表
   - LifeStoryMedia - 生平故事媒体表
   - Timeline - 时间轴表

7. **memorial_service.go** - 追思会模型
   - MemorialService - 追思会表
   - MemorialServiceParticipant - 参与者表
   - ServiceActivity - 活动表
   - ServiceInvitation - 邀请表
   - ServiceRecording - 录制表
   - ServiceChat - 聊天表

8. **privacy.go** - 隐私模型
   - VisitorPermissionSetting - 访客权限设置表
   - VisitorBlacklist - 访客黑名单表
   - AccessRequest - 访问请求表

9. **system_config.go** - 系统配置模型
   - SystemConfig - 系统配置表
   - FestivalConfig - 节日配置表
   - TemplateConfig - 模板配置表
   - DataBackup - 数据备份表
   - SystemLog - 系统日志表
   - SystemMonitor - 系统监控表

10. **premium_service.go** - 增值服务模型
    - PremiumPackage - 高级套餐表
    - UserSubscription - 用户订阅表
    - MemorialUpgrade - 纪念馆升级表
    - CustomTemplate - 定制模板表
    - StorageUsage - 存储使用表
    - PaymentOrder - 支付订单表
    - ServiceUsageLog - 服务使用日志表

11. **exclusive_service.go** - 专属服务模型
    - ExclusiveService - 专属服务表
    - ServiceBooking - 服务预订表
    - DataExportRequest - 数据导出请求表
    - PhotoRestoreRequest - 照片修复请求表
    - CustomDesignRequest - 定制设计请求表
    - ServiceReview - 服务评价表
    - ServiceStaff - 服务人员表

## 📝 注释格式规范

### GORM标签中添加comment

```go
// 正确格式
ID string `json:"id" gorm:"primaryKey;type:varchar(36);comment:用户ID"`

// 状态字段注释格式
Status int `json:"status" gorm:"default:1;comment:状态:1正常 0禁用"`

// 枚举类型注释格式
FileType string `json:"file_type" gorm:"type:varchar(20);not null;comment:文件类型:image图片 video视频 audio音频"`

// 时间字段注释
CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
```

### 注释内容规范

1. **简洁明了**: 使用简短的中文描述字段用途
2. **枚举说明**: 对于有固定值的字段，列出所有可能的值及其含义
3. **单位说明**: 对于数值字段，说明单位（如：字节、秒、分钟等）
4. **关系说明**: 对于外键字段，说明关联的表
5. **状态说明**: 对于状态字段，列出所有状态值及其含义

## 🔄 更新数据库

### 方法1: 重新运行迁移（推荐）

```bash
# 删除所有表并重新创建（危险！会丢失数据）
go run cmd/migrate/main.go -action=reset

# 或者只运行迁移（会更新表结构）
go run cmd/migrate/main.go -action=migrate
```

### 方法2: 手动执行ALTER TABLE

如果不想删除数据，可以手动执行ALTER TABLE语句：

```sql
-- 示例：为users表添加注释
ALTER TABLE users MODIFY COLUMN id VARCHAR(36) COMMENT '用户ID';
ALTER TABLE users MODIFY COLUMN wechat_openid VARCHAR(100) NOT NULL COMMENT '微信OpenID';
ALTER TABLE users MODIFY COLUMN nickname VARCHAR(50) COMMENT '用户昵称';
-- ... 其他字段
```

## 📊 更新统计

- **总模型文件**: 11个
- **已完成**: 4个 (36%)
- **待完成**: 7个 (64%)
- **预计总字段数**: 约200+个
- **已添加注释**: 约40个字段

## 🎯 下一步行动

### 立即行动
1. ✅ 已完成核心模型的注释（用户、纪念馆、祭扫）
2. ⏳ 继续完成剩余7个模型文件的注释
3. ⏳ 重新运行数据库迁移以应用注释

### 优先级
1. **高优先级**: family.go, visitor.go（核心功能）
2. **中优先级**: album.go, memorial_service.go（常用功能）
3. **低优先级**: privacy.go, system_config.go, premium_service.go, exclusive_service.go（辅助功能）

## 💡 提示

- 所有注释都使用中文，便于中文团队理解
- 注释会在数据库表结构中显示，方便DBA和开发人员查看
- 使用 `SHOW FULL COLUMNS FROM table_name` 可以查看字段注释
- 注释不影响程序运行，只是增强可读性

## ✅ 验证方法

### 查看表结构和注释

```sql
-- 查看users表的字段注释
SHOW FULL COLUMNS FROM users;

-- 查看memorials表的字段注释
SHOW FULL COLUMNS FROM memorials;

-- 查看worship_records表的字段注释
SHOW FULL COLUMNS FROM worship_records;
```

### 预期输出示例

```
+---------------+--------------+------+-----+---------+-------+---------------------------------+
| Field         | Type         | Null | Key | Default | Extra | Comment                         |
+---------------+--------------+------+-----+---------+-------+---------------------------------+
| id            | varchar(36)  | NO   | PRI | NULL    |       | 用户ID                          |
| wechat_openid | varchar(100) | NO   | UNI | NULL    |       | 微信OpenID                      |
| nickname      | varchar(50)  | YES  |     | NULL    |       | 用户昵称                        |
+---------------+--------------+------+-----+---------+-------+---------------------------------+
```

## 📚 参考资料

- [GORM 文档 - 字段标签](https://gorm.io/zh_CN/docs/models.html#%E5%AD%97%E6%AE%B5%E6%A0%87%E7%AD%BE)
- [MySQL 字段注释](https://dev.mysql.com/doc/refman/8.0/en/create-table.html)

---

**更新日期**: 2024年11月14日
**状态**: 进行中 (36% 完成)
