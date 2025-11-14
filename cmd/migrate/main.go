package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// 定义命令行参数
	action := flag.String("action", "migrate", "操作类型: migrate(迁移), seed(种子数据), reset(重置), drop(删除)")
	flag.Parse()

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件，使用默认配置")
	}

	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db, err := database.InitMySQL(cfg.Database.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}
	defer sqlDB.Close()

	// 创建迁移管理器
	migration := database.NewMigration(db)

	// 执行操作
	switch *action {
	case "migrate":
		log.Println("执行数据库迁移...")
		if err := migration.AutoMigrate(); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		if err := migration.CreateIndexes(); err != nil {
			log.Fatalf("创建索引失败: %v", err)
		}
		log.Println("✅ 数据库迁移成功完成")

	case "seed":
		log.Println("插入种子数据...")
		if err := migration.SeedData(); err != nil {
			log.Fatalf("插入种子数据失败: %v", err)
		}
		log.Println("✅ 种子数据插入成功")

	case "reset":
		fmt.Print("⚠️  警告: 此操作将删除所有数据！确认继续？(yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			log.Println("操作已取消")
			os.Exit(0)
		}

		log.Println("重置数据库...")
		if err := migration.Reset(); err != nil {
			log.Fatalf("重置数据库失败: %v", err)
		}
		log.Println("✅ 数据库重置成功")

	case "drop":
		fmt.Print("⚠️  警告: 此操作将删除所有表！确认继续？(yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			log.Println("操作已取消")
			os.Exit(0)
		}

		log.Println("删除所有表...")
		if err := migration.DropAllTables(); err != nil {
			log.Fatalf("删除表失败: %v", err)
		}
		log.Println("✅ 所有表删除成功")

	default:
		log.Fatalf("未知操作: %s", *action)
	}
}
