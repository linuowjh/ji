# 索引创建错误修复文档

## 问题描述

在运行数据库迁移时，创建索引时出现SQL语法错误：

```
Error 1064 (42000): You have an error in your SQL syntax; 
check the manual that corresponds to your MySQL server version 
for the right syntax to use near 'IF NOT EXISTS idx_service_chats_timestamp 
ON service_chats(service_id, timestamp' at line 1
```

## 问题原因

**根本原因**：MySQL不支持`CREATE INDEX IF NOT EXISTS`语法。

- `CREATE INDEX IF NOT EXISTS`是PostgreSQL的语法
- MySQL 5.7及更早版本不支持这个语法
- MySQL 8.0虽然支持某些`IF NOT EXISTS`语法，但对于索引创建仍然不支持

## 解决方案

### 1. 移除`IF NOT EXISTS`子句

将所有索引创建语句从：
```sql
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column)
```

改为：
```sql
CREATE INDEX idx_name ON table_name(column)
```

### 2. 添加错误处理

在索引创建循环中添加智能错误处理：

```go
for _, indexSQL := range indexes {
    if err := m.db.Exec(indexSQL).Error; err != nil {
        // 如果索引已存在，跳过错误
        if strings.Contains(err.Error(), "Duplicate key name") || 
           strings.Contains(err.Error(), "already exists") {
            skipCount++
        } else {
            log.Printf("创建索引失败: %s, 错误: %v", indexSQL, err)
        }
    } else {
        successCount++
    }
}
```

### 3. 添加统计信息

在索引创建完成后输出统计信息：

```go
log.Printf("数据库索引创建完成 - 成功: %d, 跳过: %d", successCount, skipCount)
```

## 修改的文件

**文件**: `internal/database/migration.go`

### 修改内容

1. **添加strings包导入**：
```go
import (
    "fmt"
    "log"
    "strings"  // 新增
    "yun-nian-memorial/internal/models"
    "gorm.io/gorm"
)
```

2. **修改CreateIndexes函数**：
   - 移除所有`IF NOT EXISTS`子句
   - 添加错误分类处理
   - 添加成功/跳过计数器
   - 改进日志输出

## 测试结果

运行迁移后的输出：

```
2025/11/14 15:36:50 数据库索引创建完成 - 成功: 38, 跳过: 11
2025/11/14 15:36:50 ✅ 数据库迁移成功完成
```

**说明**：
- 成功创建38个新索引
- 跳过11个已存在的索引
- 没有任何错误

## MySQL索引语法对比

### PostgreSQL（支持）
```sql
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column);
```

### MySQL 5.7及以下（不支持）
```sql
-- 错误的语法
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column);

-- 正确的语法
CREATE INDEX idx_name ON table_name(column);
```

### MySQL 8.0（部分支持）
```sql
-- 表创建支持
CREATE TABLE IF NOT EXISTS table_name (...);

-- 索引创建仍不支持
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column);  -- 错误
```

## 最佳实践

### 1. 跨数据库兼容性

如果需要支持多种数据库，应该：
- 检测数据库类型
- 使用对应的SQL语法
- 或者使用ORM的抽象层

### 2. 索引创建策略

**方案A：忽略重复错误**（当前方案）
```go
if strings.Contains(err.Error(), "Duplicate key name") {
    // 跳过
}
```

**方案B：先检查后创建**
```go
var count int64
db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_name = ? AND index_name = ?", 
    tableName, indexName).Scan(&count)
if count == 0 {
    db.Exec(createIndexSQL)
}
```

**方案C：使用GORM的Migrator**
```go
if !db.Migrator().HasIndex(&Model{}, "idx_name") {
    db.Migrator().CreateIndex(&Model{}, "idx_name")
}
```

### 3. 推荐方案

对于当前项目，**方案A**（忽略重复错误）是最简单有效的：
- 代码简洁
- 性能好（不需要额外查询）
- 易于维护
- 适合初始化和更新场景

## 相关问题

### ServiceReview表问题

在修复索引创建问题的同时，也解决了ServiceReview表的创建问题：

**问题**：GORM将`booking_id`识别为`longtext`而不是`varchar(36)`

**解决**：使用`size:36`标签而不是`type:varchar(36)`

```go
// 修改前
BookingID  string    `gorm:"type:varchar(36);not null;uniqueIndex" json:"booking_id"`

// 修改后
BookingID  string    `gorm:"size:36;not null;uniqueIndex" json:"booking_id"`
```

## 总结

通过移除MySQL不支持的`IF NOT EXISTS`语法，并添加智能错误处理，成功解决了索引创建失败的问题。现在数据库迁移可以正常运行，所有表和索引都能正确创建。
