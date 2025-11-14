package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupWorshipTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
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
		t.Skipf("Database not available for testing: %v", err)
	}

	db.AutoMigrate(&models.User{}, &models.Memorial{}, &models.WorshipRecord{},
		&models.Prayer{}, &models.Message{})

	router := gin.New()
	return router, db
}

func cleanupWorshipTestRouter(db *gorm.DB) {
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM prayers")
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

func TestOfferFlowersAPI(t *testing.T) {
	router, db := setupWorshipTestRouter(t)
	defer cleanupWorshipTestRouter(db)

	worshipService := services.NewWorshipService(db)
	controller := NewWorshipController(worshipService)

	// Create test user and memorial
	user := &models.User{
		ID:           "test-user-worship-1",
		WechatOpenID: "test-openid-worship-1",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-worship-1",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Setup route
	router.POST("/api/v1/memorials/:id/worship/flowers", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.OfferFlowers(c)
	})

	reqBody := services.OfferFlowersRequest{
		FlowerType: "chrysanthemum",
		Quantity:   3,
		Message:    "献上鲜花",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/memorials/"+memorial.ID+"/worship/flowers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}

func TestLightCandleAPI(t *testing.T) {
	router, db := setupWorshipTestRouter(t)
	defer cleanupWorshipTestRouter(db)

	worshipService := services.NewWorshipService(db)
	controller := NewWorshipController(worshipService)

	user := &models.User{
		ID:           "test-user-worship-2",
		WechatOpenID: "test-openid-worship-2",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-worship-2",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	router.POST("/api/v1/memorials/:id/worship/candle", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.LightCandle(c)
	})

	reqBody := services.LightCandleRequest{
		CandleType: "red",
		Duration:   60,
		Message:    "点燃蜡烛",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/memorials/"+memorial.ID+"/worship/candle", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}

func TestCreatePrayerAPI(t *testing.T) {
	router, db := setupWorshipTestRouter(t)
	defer cleanupWorshipTestRouter(db)

	worshipService := services.NewWorshipService(db)
	controller := NewWorshipController(worshipService)

	user := &models.User{
		ID:           "test-user-worship-3",
		WechatOpenID: "test-openid-worship-3",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-worship-3",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	router.POST("/api/v1/memorials/:id/prayers", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.CreatePrayer(c)
	})

	reqBody := services.CreatePrayerRequest{
		Content:  "愿您在天堂安好",
		IsPublic: true,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/memorials/"+memorial.ID+"/prayers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}

func TestGetWorshipRecordsAPI(t *testing.T) {
	router, db := setupWorshipTestRouter(t)
	defer cleanupWorshipTestRouter(db)

	worshipService := services.NewWorshipService(db)
	controller := NewWorshipController(worshipService)

	user := &models.User{
		ID:           "test-user-worship-4",
		WechatOpenID: "test-openid-worship-4",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-worship-4",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create test records
	record := &models.WorshipRecord{
		ID:          "worship-record-1",
		MemorialID:  memorial.ID,
		UserID:      user.ID,
		WorshipType: "flower",
		Content:     "{}",
	}
	db.Create(record)

	router.GET("/api/v1/memorials/:id/worship/records", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.GetWorshipRecords(c)
	})

	req, _ := http.NewRequest("GET", "/api/v1/memorials/"+memorial.ID+"/worship/records?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}

func TestGetPrayerWallAPI(t *testing.T) {
	router, db := setupWorshipTestRouter(t)
	defer cleanupWorshipTestRouter(db)

	worshipService := services.NewWorshipService(db)
	controller := NewWorshipController(worshipService)

	user := &models.User{
		ID:           "test-user-worship-5",
		WechatOpenID: "test-openid-worship-5",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-worship-5",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Create test prayer
	prayer := &models.Prayer{
		ID:         "prayer-test-1",
		MemorialID: memorial.ID,
		UserID:     user.ID,
		Content:    "祈福内容",
		IsPublic:   true,
	}
	db.Create(prayer)

	router.GET("/api/v1/memorials/:id/prayers", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.GetPrayerWall(c)
	})

	req, _ := http.NewRequest("GET", "/api/v1/memorials/"+memorial.ID+"/prayers?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}
