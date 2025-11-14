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
	database := os.Getenv("MYSQL_DATABASE")

	if host == "" || port == "" || username == "" || database == "" {
		log.Fatal("错误: 缺少必要的数据库配置")
	}

	// 连接到数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	log.Printf("成功连接到数据库: %s", database)

	// 删除现有的service_reviews表
	log.Println("删除现有的service_reviews表...")
	_, err = db.Exec("DROP TABLE IF EXISTS service_reviews")
	if err != nil {
		log.Printf("删除表失败: %v", err)
	} else {
		log.Println("✅ service_reviews表已删除")
	}

	log.Println("✅ 修复完成，现在可以重新运行迁移")
}
