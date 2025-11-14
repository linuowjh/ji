package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/controllers"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupE2ETest(t *testing.T) (*gin.Engine, *gorm.DB, *config.Config) {
	gin.SetMode(gin.TestMode)

	host := getEnvOrDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvOrDefault("MYSQL_PORT", "3306")
	user := getEnvOrDefault("MYSQL_USERNAME", "root")
	password := getEnvOrDefault("MYSQL_PASSWORD", "root")
	database := getEnvOrDefault("MYSQL_DATABASE", "yun_nian_memorial") + "_test"
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Database not available for E2E testing: %v", err)
	}

	// Auto migrate all tables
	db.AutoMigrate(
		&models.User{},
		&models.Memorial{},
		&models.WorshipRecord{},
		&models.Prayer{},
		&models.Message{},
		&models.VisitorRecord{},
		&models.Family{},
		&models.FamilyMember{},
		&models.MemorialFamily{},
	)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			ExpireTime: 3600,
		},
	}

	router := gin.New()
	return router, db, cfg
}

func cleanupE2ETest(db *gorm.DB) {
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM prayers")
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM visitor_records")
	db.Exec("DELETE FROM family_members")
	db.Exec("DELETE FROM memorial_families")
	db.Exec("DELETE FROM families")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

// TestCompleteMemorialCreationFlow tests the complete flow of creating and managing a memorial
func TestCompleteMemorialCreationFlow(t *testing.T) {
	router, db, cfg := setupE2ETest(t)
	defer cleanupE2ETest(db)

	// Initialize services
	userService := services.NewUserService(db, cfg)
	memorialService := services.NewMemorialService(db)

	// Initialize controllers
	_ = controllers.NewUserController(userService) // userController not used in this test
	memorialController := controllers.NewMemorialController(memorialService)

	// Step 1: Create a user
	user := &models.User{
		ID:           "e2e-user-1",
		WechatOpenID: "e2e-openid-1",
		Nickname:     "E2E Test User",
		Status:       1,
	}
	db.Create(user)

	// Step 2: Create a memorial
	router.POST("/api/v1/memorials", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		memorialController.CreateMemorial(c)
	})

	birthDate := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	createReq := services.CreateMemorialRequest{
		DeceasedName:   "E2E测试逝者",
		BirthDate:      &birthDate,
		DeathDate:      &deathDate,
		Biography:      "这是一个端到端测试的生平简介",
		ThemeStyle:     "traditional",
		TombstoneStyle: "marble",
		Epitaph:        "永远怀念",
		PrivacyLevel:   1,
	}

	jsonData, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/api/v1/memorials", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createResp controllers.APIResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.Equal(t, 0, createResp.Code)

	// Extract memorial ID from response
	memorialData := createResp.Data.(map[string]interface{})
	memorialID := memorialData["id"].(string)

	// Step 3: Get memorial details
	router.GET("/api/v1/memorials/:id", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		memorialController.GetMemorial(c)
	})

	req, _ = http.NewRequest("GET", "/api/v1/memorials/"+memorialID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResp controllers.APIResponse
	json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.Equal(t, 0, getResp.Code)

	// Step 4: Update memorial
	router.PUT("/api/v1/memorials/:id", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		memorialController.UpdateMemorial(c)
	})

	updateReq := services.UpdateMemorialRequest{
		Biography: "更新后的生平简介",
	}

	jsonData, _ = json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", "/api/v1/memorials/"+memorialID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify update
	var memorial models.Memorial
	db.First(&memorial, "id = ?", memorialID)
	assert.Equal(t, "更新后的生平简介", memorial.Biography)
}

// TestCompleteWorshipFlow tests the complete worship flow
func TestCompleteWorshipFlow(t *testing.T) {
	router, db, _ := setupE2ETest(t)
	defer cleanupE2ETest(db)

	// Initialize services
	worshipService := services.NewWorshipService(db)

	// Initialize controllers
	worshipController := controllers.NewWorshipController(worshipService)

	// Create test user and memorial
	user := &models.User{
		ID:           "e2e-user-2",
		WechatOpenID: "e2e-openid-2",
		Nickname:     "E2E Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "e2e-memorial-1",
		CreatorID:    user.ID,
		DeceasedName: "E2E测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Step 1: Offer flowers
	router.POST("/api/v1/memorials/:id/worship/flowers", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		worshipController.OfferFlowers(c)
	})

	flowerReq := services.OfferFlowersRequest{
		FlowerType: "chrysanthemum",
		Quantity:   3,
		Message:    "献上鲜花，表达思念",
	}

	jsonData, _ := json.Marshal(flowerReq)
	req, _ := http.NewRequest("POST", "/api/v1/memorials/"+memorial.ID+"/worship/flowers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 2: Light candle
	router.POST("/api/v1/memorials/:id/worship/candle", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		worshipController.LightCandle(c)
	})

	candleReq := services.LightCandleRequest{
		CandleType: "red",
		Duration:   60,
		Message:    "点燃蜡烛，照亮前路",
	}

	jsonData, _ = json.Marshal(candleReq)
	req, _ = http.NewRequest("POST", "/api/v1/memorials/"+memorial.ID+"/worship/candle", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 3: Create prayer
	router.POST("/api/v1/memorials/:id/prayers", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		worshipController.CreatePrayer(c)
	})

	prayerReq := services.CreatePrayerRequest{
		Content:  "愿您在天堂安好",
		IsPublic: true,
	}

	jsonData, _ = json.Marshal(prayerReq)
	req, _ = http.NewRequest("POST", "/api/v1/memorials/"+memorial.ID+"/prayers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 4: Get worship records
	router.GET("/api/v1/memorials/:id/worship/records", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		worshipController.GetWorshipRecords(c)
	})

	req, _ = http.NewRequest("GET", "/api/v1/memorials/"+memorial.ID+"/worship/records?page=1&page_size=10", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var recordsResp controllers.APIResponse
	json.Unmarshal(w.Body.Bytes(), &recordsResp)
	assert.Equal(t, 0, recordsResp.Code)

	// Verify we have worship records
	var count int64
	db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorial.ID).Count(&count)
	assert.GreaterOrEqual(t, count, int64(2)) // At least flowers and candle
}

// TestUserJourneyFlow tests a complete user journey
func TestUserJourneyFlow(t *testing.T) {
	_, db, cfg := setupE2ETest(t)
	defer cleanupE2ETest(db)

	// Initialize services
	userService := services.NewUserService(db, cfg)
	memorialService := services.NewMemorialService(db)
	worshipService := services.NewWorshipService(db)

	// Step 1: Create user
	user := &models.User{
		ID:           "e2e-user-3",
		WechatOpenID: "e2e-openid-3",
		Nickname:     "完整流程用户",
		Status:       1,
	}
	db.Create(user)

	// Step 2: User creates a memorial
	birthDate := time.Date(1945, 5, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2019, 10, 15, 0, 0, 0, 0, time.UTC)

	createReq := &services.CreateMemorialRequest{
		DeceasedName:   "亲爱的祖父",
		BirthDate:      &birthDate,
		DeathDate:      &deathDate,
		Biography:      "祖父是一位伟大的教育家",
		ThemeStyle:     "traditional",
		TombstoneStyle: "marble",
		Epitaph:        "桃李满天下",
		PrivacyLevel:   1,
	}

	memorial, err := memorialService.CreateMemorial(user.ID, createReq)
	assert.NoError(t, err)
	assert.NotNil(t, memorial)

	// Step 3: User performs worship activities
	flowerReq := &services.OfferFlowersRequest{
		FlowerType: "chrysanthemum",
		Quantity:   9,
		Message:    "祖父，我很想念您",
	}
	err = worshipService.OfferFlowers(user.ID, memorial.ID, flowerReq)
	assert.NoError(t, err)

	// Step 4: User creates a prayer
	prayerReq := &services.CreatePrayerRequest{
		Content:  "愿祖父在天堂安好，保佑家人平安健康",
		IsPublic: true,
	}
	prayer, err := worshipService.CreatePrayer(user.ID, memorial.ID, prayerReq)
	assert.NoError(t, err)
	assert.NotNil(t, prayer)

	// Step 5: User creates a message
	messageReq := &services.CreateMessageRequest{
		MessageType: "text",
		Content:     "祖父，今天是您的生日，我们全家都很想念您",
	}
	message, err := worshipService.CreateMessage(user.ID, memorial.ID, messageReq)
	assert.NoError(t, err)
	assert.NotNil(t, message)

	// Step 6: Get user statistics
	stats, err := userService.GetUserStatistics(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), stats["memorial_count"])
	assert.GreaterOrEqual(t, stats["worship_count"].(int64), int64(1))

	// Step 7: Get memorial details
	retrievedMemorial, err := memorialService.GetMemorial(user.ID, memorial.ID)
	assert.NoError(t, err)
	assert.Equal(t, "亲爱的祖父", retrievedMemorial.DeceasedName)

	// Step 8: Get worship statistics
	worshipStats, err := worshipService.GetWorshipStatistics(user.ID, memorial.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, worshipStats["flower_count"].(int64), int64(1))
	assert.GreaterOrEqual(t, worshipStats["prayer_count"].(int64), int64(1))
}

// TestMultiUserInteraction tests interaction between multiple users
func TestMultiUserInteraction(t *testing.T) {
	_, db, _ := setupE2ETest(t)
	defer cleanupE2ETest(db)

	memorialService := services.NewMemorialService(db)
	worshipService := services.NewWorshipService(db)

	// Create two users
	user1 := &models.User{
		ID:           "e2e-user-4",
		WechatOpenID: "e2e-openid-4",
		Nickname:     "用户1",
		Status:       1,
	}
	user2 := &models.User{
		ID:           "e2e-user-5",
		WechatOpenID: "e2e-openid-5",
		Nickname:     "用户2",
		Status:       1,
	}
	db.Create(user1)
	db.Create(user2)

	// User1 creates a memorial
	createReq := &services.CreateMemorialRequest{
		DeceasedName: "共同的亲人",
		Biography:    "一位受人尊敬的长辈",
		PrivacyLevel: 1, // Public to family
	}

	memorial, err := memorialService.CreateMemorial(user1.ID, createReq)
	assert.NoError(t, err)

	// User2 visits and worships
	_, err = memorialService.GetMemorial(user2.ID, memorial.ID)
	assert.NoError(t, err)

	// User2 offers flowers
	flowerReq := &services.OfferFlowersRequest{
		FlowerType: "lily",
		Quantity:   3,
		Message:    "来自用户2的思念",
	}
	err = worshipService.OfferFlowers(user2.ID, memorial.ID, flowerReq)
	assert.NoError(t, err)

	// User2 creates a prayer
	prayerReq := &services.CreatePrayerRequest{
		Content:  "愿您安息",
		IsPublic: true,
	}
	_, err = worshipService.CreatePrayer(user2.ID, memorial.ID, prayerReq)
	assert.NoError(t, err)

	// Verify worship records from both users
	records, total, err := worshipService.GetWorshipRecords(user1.ID, memorial.ID, 1, 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.NotEmpty(t, records)

	// Verify prayer wall
	prayers, total, err := worshipService.GetPrayerWall(user1.ID, memorial.ID, 1, 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.NotEmpty(t, prayers)
}
