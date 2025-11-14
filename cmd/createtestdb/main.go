package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件")
	}

	// 获取数据库配置
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	if host == "" || port == "" || username == "" {
		log.Fatal("错误: 缺少必要的数据库配置")
	}

	// 连接到MySQL服务器（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接MySQL服务器失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("无法连接到MySQL服务器: %v", err)
	}

	log.Printf("成功连接到MySQL服务器: %s:%s", host, port)

	// 创建测试数据库
	testDatabase := "yun_nian_memorial_test"
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", testDatabase)
	_, err = db.Exec(createDBSQL)
	if err != nil {
		log.Fatalf("创建测试数据库失败: %v", err)
	}

	log.Printf("✅ 测试数据库 '%s' 创建成功（或已存在）", testDatabase)

	// 验证数据库是否存在
	var dbName string
	err = db.QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", testDatabase).Scan(&dbName)
	if err != nil {
		log.Fatalf("验证测试数据库失败: %v", err)
	}

	log.Printf("✅ 验证成功: 测试数据库 '%s' 已存在", dbName)
}
