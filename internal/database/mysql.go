package database

import (
	"fmt"
	"yun-nian-memorial/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// 根据环境设置日志级别
	logLevel := logger.Silent // 默认静默模式
	if cfg.LogLevel == "info" {
		logLevel = logger.Info
	} else if cfg.LogLevel == "warn" {
		logLevel = logger.Warn
	} else if cfg.LogLevel == "error" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// 注意：不在这里执行AutoMigrate
	// 数据库迁移应该通过以下方式执行：
	// 1. 使用 cmd/migrate/main.go 进行迁移
	// 2. 或者在启动服务时使用 -migrate 参数
	// 这样可以：
	// - 加快服务启动速度
	// - 更好地控制迁移时机
	// - 避免生产环境意外的表结构变更

	return db, nil
}
