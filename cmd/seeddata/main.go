package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/database"

	"github.com/joho/godotenv"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
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

	// 执行SQL
	log.Println("开始插入测试数据...")

	// 禁用外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 清空表
	tables := []string{"family_activities", "memorial_reminders", "messages", "prayers", "worship_records", "memorial_families", "family_members", "families", "memorials", "users"}
	for _, table := range tables {
		log.Printf("清空表: %s", table)
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
	}

	// 启用外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 读取并执行SQL文件
	sqlContent, err := os.ReadFile("scripts/insert_test_data.sql")
	if err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	// 分割SQL语句并执行
	sqlStatements := strings.Split(string(sqlContent), ";")
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") || strings.HasPrefix(stmt, "SET FOREIGN_KEY_CHECKS") || strings.HasPrefix(stmt, "TRUNCATE") {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("警告: 执行SQL失败: %v\nSQL: %s", err, stmt[:min(len(stmt), 100)])
		}
	}

	log.Println("✅ 测试数据插入成功")
	fmt.Println("\n测试账号信息:")
	fmt.Println("用户1: 张三 (test-user-1)")
	fmt.Println("用户2: 李四 (test-user-2)")
	fmt.Println("用户3: 王五 (test-user-3)")
	fmt.Println("用户4: 赵六 (test-user-4)")
	fmt.Println("\n家族圈:")
	fmt.Println("张氏家族 (test-family-1) - 邀请码: ZHANG001")
	fmt.Println("  成员: 张三(管理员), 王五, 赵六")
	fmt.Println("  纪念馆: 张老爷子, 王爷爷")
	fmt.Println("\n李氏家族 (test-family-2) - 邀请码: LI002")
	fmt.Println("  成员: 李四(管理员), 王五")
	fmt.Println("  纪念馆: 李奶奶")
}
