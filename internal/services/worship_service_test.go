package services

import (
	"fmt"
	"testing"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupWorshipTestDB(t *testing.T) *gorm.DB {
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

	db.AutoMigrate(&models.User{}, &models.Memorial{}, &models.WorshipRecord{},
		&models.Prayer{}, &models.Message{})

	return db
}

func cleanupWorshipTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM prayers")
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

func TestOfferFlowers(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	// Create test user and memorial
	user := &models.User{
		ID:           "test-user-1",
		WechatOpenID: "test-openid-1",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-1",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Offer flowers
	req := &OfferFlowersRequest{
		FlowerType: "chrysanthemum",
		Quantity:   3,
		Message:    "献上鲜花，表达思念",
	}

	err := service.OfferFlowers(user.ID, memorial.ID, req)

	assert.NoError(t, err)

	// Verify record created
	var record models.WorshipRecord
	db.Where("memorial_id = ? AND user_id = ? AND worship_type = ?", 
		memorial.ID, user.ID, "flower").First(&record)
	assert.NotEmpty(t, record.ID)
}

func TestLightCandle(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-2",
		WechatOpenID: "test-openid-2",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-2",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Light candle
	req := &LightCandleRequest{
		CandleType: "red",
		Duration:   60,
		Message:    "点燃蜡烛，照亮前路",
	}

	err := service.LightCandle(user.ID, memorial.ID, req)

	assert.NoError(t, err)

	// Verify record created
	var record models.WorshipRecord
	db.Where("memorial_id = ? AND user_id = ? AND worship_type = ?", 
		memorial.ID, user.ID, "candle").First(&record)
	assert.NotEmpty(t, record.ID)
}

func TestOfferIncense(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-3",
		WechatOpenID: "test-openid-3",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-3",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Offer incense
	req := &OfferIncenseRequest{
		IncenseCount: 3,
		IncenseType:  "sandalwood",
		Message:      "上香祈福",
	}

	err := service.OfferIncense(user.ID, memorial.ID, req)

	assert.NoError(t, err)

	// Verify record created
	var record models.WorshipRecord
	db.Where("memorial_id = ? AND user_id = ? AND worship_type = ?", 
		memorial.ID, user.ID, "incense").First(&record)
	assert.NotEmpty(t, record.ID)
}

func TestCreatePrayer(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-4",
		WechatOpenID: "test-openid-4",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-4",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create prayer
	req := &CreatePrayerRequest{
		Content:  "愿您在天堂安好",
		IsPublic: true,
	}

	prayer, err := service.CreatePrayer(user.ID, memorial.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, prayer)
	assert.Equal(t, "愿您在天堂安好", prayer.Content)
	assert.True(t, prayer.IsPublic)
}

func TestCreateMessage(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-5",
		WechatOpenID: "test-openid-5",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-5",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create text message
	req := &CreateMessageRequest{
		MessageType: "text",
		Content:     "想念您",
	}

	message, err := service.CreateMessage(user.ID, memorial.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "text", message.MessageType)
	assert.Equal(t, "想念您", message.Content)
}

func TestGetWorshipRecords(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-6",
		WechatOpenID: "test-openid-6",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-6",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create test records
	record1 := &models.WorshipRecord{
		ID:          "record-1",
		MemorialID:  memorial.ID,
		UserID:      user.ID,
		WorshipType: "flower",
		Content:     "{}",
	}
	record2 := &models.WorshipRecord{
		ID:          "record-2",
		MemorialID:  memorial.ID,
		UserID:      user.ID,
		WorshipType: "candle",
		Content:     "{}",
	}
	db.Create(record1)
	db.Create(record2)

	// Get records
	records, total, err := service.GetWorshipRecords(user.ID, memorial.ID, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, records, 2)
}

func TestGetPrayerWall(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-7",
		WechatOpenID: "test-openid-7",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-7",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create test prayers
	prayer1 := &models.Prayer{
		ID:         "prayer-1",
		MemorialID: memorial.ID,
		UserID:     user.ID,
		Content:    "祈福1",
		IsPublic:   true,
	}
	prayer2 := &models.Prayer{
		ID:         "prayer-2",
		MemorialID: memorial.ID,
		UserID:     user.ID,
		Content:    "祈福2",
		IsPublic:   false, // Private, should not appear
	}
	db.Create(prayer1)
	db.Create(prayer2)

	// Get prayer wall
	prayers, total, err := service.GetPrayerWall(user.ID, memorial.ID, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total) // Only public prayers
	assert.Len(t, prayers, 1)
	assert.Equal(t, "祈福1", prayers[0].Content)
}

func TestRenewCandle(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-8",
		WechatOpenID: "test-openid-8",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-8",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Light a candle first
	req := &LightCandleRequest{
		CandleType: "red",
		Duration:   30,
		Message:    "点燃蜡烛",
	}
	service.LightCandle(user.ID, memorial.ID, req)

	// Renew candle
	err := service.RenewCandle(user.ID, memorial.ID, 30)

	assert.NoError(t, err)
}

func TestGetWorshipStatistics(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	user := &models.User{
		ID:           "test-user-9",
		WechatOpenID: "test-openid-9",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-9",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create test records
	record1 := &models.WorshipRecord{
		ID:          "record-3",
		MemorialID:  memorial.ID,
		UserID:      user.ID,
		WorshipType: "flower",
		Content:     "{}",
		CreatedAt:   time.Now(),
	}
	record2 := &models.WorshipRecord{
		ID:          "record-4",
		MemorialID:  memorial.ID,
		UserID:      user.ID,
		WorshipType: "candle",
		Content:     "{}",
		CreatedAt:   time.Now(),
	}
	db.Create(record1)
	db.Create(record2)

	// Get statistics
	stats, err := service.GetWorshipStatistics(user.ID, memorial.ID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(1), stats["flower_count"])
	assert.Equal(t, int64(1), stats["candle_count"])
	assert.Equal(t, int64(2), stats["total_visits"])
}

func TestAnalyzeMessageEmotion(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	// Test sad emotion
	result, err := service.AnalyzeMessageEmotion("我很想念您，心里很难过")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, []string{"sad", "nostalgic"}, result.Emotion)

	// Test grateful emotion
	result, err = service.AnalyzeMessageEmotion("感谢您的教诲和恩情")
	assert.NoError(t, err)
	assert.Equal(t, "grateful", result.Emotion)

	// Test peaceful emotion
	result, err = service.AnalyzeMessageEmotion("愿您安息，一切安好")
	assert.NoError(t, err)
	assert.Equal(t, "peaceful", result.Emotion)
}

func TestGetPrayerCardTemplates(t *testing.T) {
	db := setupWorshipTestDB(t)
	defer cleanupWorshipTestDB(db)

	service := NewWorshipService(db)

	templates := service.GetPrayerCardTemplates()

	assert.NotEmpty(t, templates)
	assert.GreaterOrEqual(t, len(templates), 3)
	
	// Check template structure
	for _, template := range templates {
		assert.NotEmpty(t, template.ID)
		assert.NotEmpty(t, template.Name)
		assert.NotEmpty(t, template.Category)
	}
}
