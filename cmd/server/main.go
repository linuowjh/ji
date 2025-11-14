package main

import (
	"flag"
	"log"
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/database"
	"yun-nian-memorial/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件，使用默认配置")
	}

	// 解析命令行参数
	var (
		migrate = flag.Bool("migrate", false, "执行数据库迁移")
		seed    = flag.Bool("seed", false, "插入种子数据")
		reset   = flag.Bool("reset", false, "重置数据库（危险操作）")
	)
	flag.Parse()

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitMySQL(cfg.Database.MySQL)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	// 创建迁移管理器
	migration := database.NewMigration(db)

	// 处理命令行操作
	if *reset {
		log.Println("执行数据库重置...")
		if err := migration.Reset(); err != nil {
			log.Fatal("数据库重置失败:", err)
		}
		log.Println("数据库重置完成")
		return
	}

	if *migrate {
		log.Println("执行数据库迁移...")
		if err := migration.AutoMigrate(); err != nil {
			log.Fatal("数据库迁移失败:", err)
		}
		if err := migration.CreateIndexes(); err != nil {
			log.Fatal("创建索引失败:", err)
		}
		log.Println("数据库迁移完成")
		
		if *seed {
			if err := migration.SeedData(); err != nil {
				log.Fatal("插入种子数据失败:", err)
			}
		}
		return
	}

	if *seed {
		log.Println("插入种子数据...")
		if err := migration.SeedData(); err != nil {
			log.Fatal("插入种子数据失败:", err)
		}
		log.Println("种子数据插入完成")
		return
	}

	// 初始化Redis（可选）
	rdb, err := database.InitRedis(cfg.Database.Redis)
	if err != nil {
		log.Println("警告: Redis连接失败，将在无缓存模式下运行:", err)
		rdb = nil
	}

	// 初始化路由
	r := router.Setup(db, rdb, cfg)

	// 启动服务器
	log.Printf("云念纪念馆服务启动中，端口: %s", cfg.Server.Port)
	log.Printf("健康检查: http://localhost:%s/health", cfg.Server.Port)
	log.Printf("API文档: http://localhost:%s/api/v1", cfg.Server.Port)
	
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}