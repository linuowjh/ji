package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
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
		&models.VisitorRecord{}, &models.FamilyMember{}, &models.MemorialFamily{})

	router := gin.New()
	return router, db
}

func cleanupTestRouter(db *gorm.DB) {
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM visitor_records")
	db.Exec("DELETE FROM memorial_families")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

func TestCreateMemorialAPI(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestRouter(db)

	memorialService := services.NewMemorialService(db)
	controller := NewMemorialController(memorialService)

	// Create test user
	user := &models.User{
		ID:           "test-user-api-1",
		WechatOpenID: "test-openid-api-1",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	// Setup route
	router.POST("/api/v1/memorials", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.CreateMemorial(c)
	})

	birthDate := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	reqBody := services.CreateMemorialRequest{
		DeceasedName:   "测试逝者",
		BirthDate:      &birthDate,
		DeathDate:      &deathDate,
		Biography:      "生平简介",
		ThemeStyle:     "traditional",
		TombstoneStyle: "marble",
		PrivacyLevel:   1,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/memorials", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "创建成功", response.Message)
}

func TestGetMemorialAPI(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestRouter(db)

	memorialService := services.NewMemorialService(db)
	controller := NewMemorialController(memorialService)

	// Create test user and memorial
	user := &models.User{
		ID:           "test-user-api-2",
		WechatOpenID: "test-openid-api-2",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-api-1",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		ThemeStyle:   "traditional",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Setup route
	router.GET("/api/v1/memorials/:id", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.GetMemorial(c)
	})

	req, _ := http.NewRequest("GET", "/api/v1/memorials/"+memorial.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "获取成功", response.Message)
}

func TestUpdateMemorialAPI(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestRouter(db)

	memorialService := services.NewMemorialService(db)
	controller := NewMemorialController(memorialService)

	// Create test user and memorial
	user := &models.User{
		ID:           "test-user-api-3",
		WechatOpenID: "test-openid-api-3",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial := &models.Memorial{
		ID:           "test-memorial-api-2",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者",
		ThemeStyle:   "traditional",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial)

	// Setup route
	router.PUT("/api/v1/memorials/:id", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.UpdateMemorial(c)
	})

	reqBody := services.UpdateMemorialRequest{
		DeceasedName: "更新的逝者姓名",
		Biography:    "更新的生平简介",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/memorials/"+memorial.ID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}

func TestGetMemorialListAPI(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestRouter(db)

	memorialService := services.NewMemorialService(db)
	controller := NewMemorialController(memorialService)

	// Create test user and memorials
	user := &models.User{
		ID:           "test-user-api-4",
		WechatOpenID: "test-openid-api-4",
		Nickname:     "Test User",
		Status:       1,
	}
	db.Create(user)

	memorial1 := &models.Memorial{
		ID:           "test-memorial-api-3",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者1",
		PrivacyLevel: 1,
		Status:       1,
	}
	memorial2 := &models.Memorial{
		ID:           "test-memorial-api-4",
		CreatorID:    user.ID,
		DeceasedName: "测试逝者2",
		PrivacyLevel: 1,
		Status:       1,
	}
	db.Create(memorial1)
	db.Create(memorial2)

	// Setup route
	router.GET("/api/v1/memorials", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		controller.GetMemorialList(c)
	})

	req, _ := http.NewRequest("GET", "/api/v1/memorials?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 0, response.Code)
}

func TestUnauthorizedAccess(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestRouter(db)

	memorialService := services.NewMemorialService(db)
	controller := NewMemorialController(memorialService)

	// Setup route without setting user_id
	router.POST("/api/v1/memorials", controller.CreateMemorial)

	reqBody := services.CreateMemorialRequest{
		DeceasedName: "测试逝者",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/memorials", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1002, response.Code)
	assert.Contains(t, response.Message, "未登录")
}
