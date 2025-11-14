package tests

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupPerformanceTest(t *testing.T) (*gorm.DB, *config.Config) {
	host := getEnvOrDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvOrDefault("MYSQL_PORT", "3306")
	user := getEnvOrDefault("MYSQL_USERNAME", "root")
	password := getEnvOrDefault("MYSQL_PASSWORD", "root")
	database := getEnvOrDefault("MYSQL_DATABASE", "yun_nian_memorial") + "_test"
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Database not available for performance testing: %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Memorial{},
		&models.WorshipRecord{},
		&models.Prayer{},
		&models.Message{},
	)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			ExpireTime: 3600,
		},
	}

	return db, cfg
}

func cleanupPerformanceTest(db *gorm.DB) {
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM prayers")
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

// TestMemorialCreationPerformance tests memorial creation performance
func TestMemorialCreationPerformance(t *testing.T) {
	db, _ := setupPerformanceTest(t)
	defer cleanupPerformanceTest(db)

	memorialService := services.NewMemorialService(db)

	// Create test user
	user := &models.User{
		ID:           "perf-user-1",
		WechatOpenID: "perf-openid-1",
		Nickname:     "Performance Test User",
		Status:       1,
	}
	db.Create(user)

	// Measure time to create 100 memorials
	startTime := time.Now()
	successCount := 0

	for i := 0; i < 100; i++ {
		req := &services.CreateMemorialRequest{
			DeceasedName: fmt.Sprintf("测试逝者%d", i),
			Biography:    "性能测试生平简介",
			PrivacyLevel: 1,
		}

		_, err := memorialService.CreateMemorial(user.ID, req)
		if err == nil {
			successCount++
		}
	}

	duration := time.Since(startTime)

	t.Logf("Created %d memorials in %v", successCount, duration)
	t.Logf("Average time per memorial: %v", duration/100)

	assert.Equal(t, 100, successCount)
	assert.Less(t, duration, 10*time.Second, "Should create 100 memorials in less than 10 seconds")
}

// TestConcurrentWorshipOperations tests concurrent worship operations
func TestConcurrentWorshipOperations(t *testing.T) {
	db, _ := setupPerformanceTest(t)
	defer cleanupPerformanceTest(db)

	worshipService := services.NewWorshipService(db)

	// Create test user and memorial
	user := &models.User{
		ID:           "perf-user-2",
		WechatOpenID: "perf-openid-2",
		Nickname:     "Performance Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "perf-memorial-1",
		CreatorID:    user.ID,
		DeceasedName: "性能测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Test concurrent worship operations
	concurrentUsers := 50
	operationsPerUser := 10

	var wg sync.WaitGroup
	startTime := time.Now()
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()

			for j := 0; j < operationsPerUser; j++ {
				req := &services.OfferFlowersRequest{
					FlowerType: "chrysanthemum",
					Quantity:   1,
					Message:    fmt.Sprintf("并发测试 用户%d 操作%d", userIndex, j),
				}

				err := worshipService.OfferFlowers(user.ID, memorial.ID, req)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	totalOperations := concurrentUsers * operationsPerUser
	t.Logf("Completed %d/%d concurrent worship operations in %v", successCount, totalOperations, duration)
	t.Logf("Throughput: %.2f operations/second", float64(successCount)/duration.Seconds())

	assert.GreaterOrEqual(t, successCount, totalOperations*9/10, "At least 90% of operations should succeed")
}

// TestQueryPerformance tests query performance with large dataset
func TestQueryPerformance(t *testing.T) {
	db, _ := setupPerformanceTest(t)
	defer cleanupPerformanceTest(db)

	_ = services.NewMemorialService(db) // memorialService not used in this test
	worshipService := services.NewWorshipService(db)

	// Create test user
	user := &models.User{
		ID:           "perf-user-3",
		WechatOpenID: "perf-openid-3",
		Nickname:     "Performance Test User",
		Status:       1,
	}
	db.Create(user)

	// Create memorial
	memorial := &models.Memorial{
		ID:           "perf-memorial-2",
		CreatorID:    user.ID,
		DeceasedName: "查询性能测试",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create 1000 worship records
	t.Log("Creating 1000 worship records...")
	for i := 0; i < 1000; i++ {
		record := &models.WorshipRecord{
			ID:          fmt.Sprintf("perf-record-%d", i),
			MemorialID:  memorial.ID,
			UserID:      user.ID,
			WorshipType: "flower",
			Content:     "{}",
			CreatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
		}
		db.Create(record)
	}

	// Test query performance
	startTime := time.Now()
	records, total, err := worshipService.GetWorshipRecords(user.ID, memorial.ID, 1, 20)
	duration := time.Since(startTime)

	assert.NoError(t, err)
	assert.Equal(t, int64(1000), total)
	assert.Len(t, records, 20)

	t.Logf("Query 20 records from 1000 total in %v", duration)
	assert.Less(t, duration, 100*time.Millisecond, "Query should complete in less than 100ms")

	// Test statistics query performance
	startTime = time.Now()
	stats, err := worshipService.GetWorshipStatistics(user.ID, memorial.ID)
	duration = time.Since(startTime)

	assert.NoError(t, err)
	assert.NotNil(t, stats)

	t.Logf("Statistics query completed in %v", duration)
	assert.Less(t, duration, 200*time.Millisecond, "Statistics query should complete in less than 200ms")
}

// TestMemoryUsage tests memory usage under load
func TestMemoryUsage(t *testing.T) {
	db, _ := setupPerformanceTest(t)
	defer cleanupPerformanceTest(db)

	memorialService := services.NewMemorialService(db)

	// Create test user
	user := &models.User{
		ID:           "perf-user-4",
		WechatOpenID: "perf-openid-4",
		Nickname:     "Memory Test User",
		Status:       1,
	}
	db.Create(user)

	// Create multiple memorials and query them
	t.Log("Creating 500 memorials...")
	for i := 0; i < 500; i++ {
		req := &services.CreateMemorialRequest{
			DeceasedName: fmt.Sprintf("内存测试%d", i),
			Biography:    "内存测试生平简介",
			PrivacyLevel: 1,
		}
		memorialService.CreateMemorial(user.ID, req)
	}

	// Query all memorials multiple times
	t.Log("Querying memorials 10 times...")
	for i := 0; i < 10; i++ {
		_, _, err := memorialService.GetMemorialList(user.ID, 1, 50)
		assert.NoError(t, err)
	}

	t.Log("Memory test completed successfully")
}

// BenchmarkMemorialCreation benchmarks memorial creation
func BenchmarkMemorialCreation(b *testing.B) {
	host := getEnvOrDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvOrDefault("MYSQL_PORT", "3306")
	user := getEnvOrDefault("MYSQL_USERNAME", "root")
	password := getEnvOrDefault("MYSQL_PASSWORD", "root")
	database := getEnvOrDefault("MYSQL_DATABASE", "yun_nian_memorial") + "_test"
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("Database not available: %v", err)
	}

	memorialService := services.NewMemorialService(db)

	user := &models.User{
		ID:           "bench-user-1",
		WechatOpenID: "bench-openid-1",
		Nickname:     "Benchmark User",
		Status:       1,
	}
	db.Create(user)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &services.CreateMemorialRequest{
			DeceasedName: fmt.Sprintf("基准测试%d", i),
			Biography:    "基准测试生平简介",
			PrivacyLevel: 1,
		}
		memorialService.CreateMemorial(user.ID, req)
	}

	b.StopTimer()
	db.Exec("DELETE FROM memorials WHERE creator_id = ?", user.ID)
	db.Exec("DELETE FROM users WHERE id = ?", user.ID)
}

// BenchmarkWorshipOperation benchmarks worship operations
func BenchmarkWorshipOperation(b *testing.B) {
	host := getEnvOrDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvOrDefault("MYSQL_PORT", "3306")
	user := getEnvOrDefault("MYSQL_USERNAME", "root")
	password := getEnvOrDefault("MYSQL_PASSWORD", "root")
	database := getEnvOrDefault("MYSQL_DATABASE", "yun_nian_memorial") + "_test"
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("Database not available: %v", err)
	}

	worshipService := services.NewWorshipService(db)

	user := &models.User{
		ID:           "bench-user-2",
		WechatOpenID: "bench-openid-2",
		Nickname:     "Benchmark User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "bench-memorial-1",
		CreatorID:    user.ID,
		DeceasedName: "基准测试",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &services.OfferFlowersRequest{
			FlowerType: "chrysanthemum",
			Quantity:   1,
			Message:    "基准测试",
		}
		worshipService.OfferFlowers(user.ID, memorial.ID, req)
	}

	b.StopTimer()
	db.Exec("DELETE FROM worship_records WHERE memorial_id = ?", memorial.ID)
	db.Exec("DELETE FROM memorials WHERE id = ?", memorial.ID)
	db.Exec("DELETE FROM users WHERE id = ?", user.ID)
}
