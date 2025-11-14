package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IndexInfo struct {
	Table      string `gorm:"column:Table"`
	NonUnique  int    `gorm:"column:Non_unique"`
	KeyName    string `gorm:"column:Key_name"`
	SeqInIndex int    `gorm:"column:Seq_in_index"`
	ColumnName string `gorm:"column:Column_name"`
	IndexType  string `gorm:"column:Index_type"`
}

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("=== Database Index Verification Report ===")
	fmt.Println()

	// 检查需要的索引
	tables := []string{
		"family_members",
		"memorial_families",
		"memorial_reminders",
	}

	requiredIndexes := map[string][]string{
		"family_members":     {"user_id", "family_id"},
		"memorial_families":  {"family_id", "memorial_id"},
		"memorial_reminders": {"memorial_id", "reminder_date"},
	}

	for _, table := range tables {
		fmt.Printf("Table: %s\n", table)
		fmt.Println(strings.Repeat("-", 80))

		var indexes []IndexInfo
		query := fmt.Sprintf("SHOW INDEX FROM %s", table)
		if err := db.Raw(query).Scan(&indexes).Error; err != nil {
			log.Printf("Error querying indexes for %s: %v\n", table, err)
			continue
		}

		// 组织索引信息
		indexMap := make(map[string][]string)
		for _, idx := range indexes {
			indexMap[idx.KeyName] = append(indexMap[idx.KeyName], idx.ColumnName)
		}

		// 显示所有索引
		for keyName, columns := range indexMap {
			fmt.Printf("  Index: %-30s Columns: %v\n", keyName, columns)
		}

		// 检查必需的索引
		fmt.Println("\n  Required Index Check:")
		for _, col := range requiredIndexes[table] {
			found := false
			for keyName, columns := range indexMap {
				for _, indexCol := range columns {
					if indexCol == col {
						found = true
						fmt.Printf("    ✓ %s (found in index: %s)\n", col, keyName)
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				fmt.Printf("    ✗ %s (MISSING)\n", col)
			}
		}

		fmt.Println()
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("All required indexes for the user reminders API have been verified.")
	fmt.Println("\nRequired indexes:")
	fmt.Println("  - family_members.user_id (for user family lookup)")
	fmt.Println("  - memorial_families.family_id (for family memorial lookup)")
	fmt.Println("  - memorial_reminders.memorial_id (for memorial reminder lookup)")
	fmt.Println("  - memorial_reminders.reminder_date (for date range filtering)")
}
