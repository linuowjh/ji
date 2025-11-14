package services

import (
	"fmt"
	"os"
	"testing"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用环境变量或默认配置
	host := getEnvOrDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvOrDefault("MYSQL_PORT", "3306")
	user := getEnvOrDefault("MYSQL_USERNAME", "root")
	password := getEnvOrDefault("MYSQL_PASSWORD", "root")
	database := getEnvOrDefault("MYSQL_DATABASE", "yun_nian_memorial") + "_test"
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Database not available for testing: %v", err)
	}

	// Auto migrate test tables
	db.AutoMigrate(&models.User{}, &models.Memorial{}, &models.WorshipRecord{}, 
		&models.VisitorRecord{}, &models.FamilyMember{}, &models.MemorialFamily{})

	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func cleanupTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM visitor_records")
	db.Exec("DELETE FROM memorial_families")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

func TestCreateMemorial(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewMemorialService(db)

	// Create test user
	user := &models.User{
		ID:           "test-user-1",
		WechatOpenID: "test-openid-1",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	birthDate := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	req := &CreateMemorialRequest{
		DeceasedName:   "张三",
		BirthDate:      &birthDate,
		DeathDate:      &deathDate,
		Biography:      "生平简介",
		ThemeStyle:     "traditional",
		TombstoneStyle: "marble",
		Epitaph:        "墓志铭",
		PrivacyLevel:   1,
	}

	memorial, err := service.CreateMemorial(user.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, memorial)
	assert.Equal(t, "张三", memorial.DeceasedName)
	assert.Equal(t, "traditional", memorial.ThemeStyle)
	assert.Equal(t, 1, memorial.PrivacyLevel)
}

func TestCreateMemorialValidation(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewMemorialService(db)

	user := &models.User{
		ID:           "test-user-2",
		WechatOpenID: "test-openid-2",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	// Test empty name
	req := &CreateMemorialRequest{
		DeceasedName: "",
	}
	_, err := service.CreateMemorial(user.ID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "姓名不能为空")

	// Test invalid date range
	birthDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	req = &CreateMemorialRequest{
		DeceasedName: "测试",
		BirthDate:    &birthDate,
		DeathDate:    &deathDate,
	}
	_, err = service.CreateMemorial(user.ID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "出生日期不能晚于逝世日期")
}

func TestGetMemorial(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewMemorialService(db)

	// Create test user and memorial
	user := &models.User{
		ID:           "test-user-3",
		WechatOpenID: "test-openid-3",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-1",
		CreatorID:    user.ID,
		DeceasedName: "李四",
		ThemeStyle:   "traditional",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Get memorial
	result, err := service.GetMemorial(user.ID, memorial.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "李四", result.DeceasedName)
}

func TestUpdateMemorial(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewMemorialService(db)

	user := &models.User{
		ID:           "test-user-4",
		WechatOpenID: "test-openid-4",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-2",
		CreatorID:    user.ID,
		DeceasedName: "王五",
		ThemeStyle:   "traditional",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Update memorial
	req := &UpdateMemorialRequest{
		DeceasedName: "王五（更新）",
		Biography:    "更新的生平简介",
	}

	err := service.UpdateMemorial(user.ID, memorial.ID, req)
	assert.NoError(t, err)

	// Verify update
	var updated models.Memorial
	db.First(&updated, "id = ?", memorial.ID)
	assert.Equal(t, "王五（更新）", updated.DeceasedName)
	assert.Equal(t, "更新的生平简介", updated.Biography)
}

func TestGetTombstoneStyles(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewMemorialService(db)

	styles := service.GetTombstoneStyles()

	assert.NotEmpty(t, styles)
	assert.GreaterOrEqual(t, len(styles), 4)
	
	// Check for expected styles
	hasMarble := false
	for _, style := range styles {
		if style.ID == "marble" {
			hasMarble = true
			assert.Equal(t, "汉白玉", style.Name)
		}
	}
	assert.True(t, hasMarble)
}

func TestGetThemeStyles(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewMemorialService(db)

	styles := service.GetThemeStyles()

	assert.NotEmpty(t, styles)
	assert.GreaterOrEqual(t, len(styles), 3)
	
	// Check for expected styles
	hasTraditional := false
	for _, style := range styles {
		if style.ID == "traditional" {
			hasTraditional = true
			assert.Equal(t, "中式传统", style.Name)
		}
	}
	assert.True(t, hasTraditional)
}
