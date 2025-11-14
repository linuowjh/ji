package services

import (
	"fmt"
	"testing"
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
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
		&models.Family{}, &models.FamilyMember{})

	return db
}

func cleanupUserTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM family_members")
	db.Exec("DELETE FROM families")
	db.Exec("DELETE FROM worship_records")
	db.Exec("DELETE FROM memorials")
	db.Exec("DELETE FROM users")
}

func TestGetUserInfo(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	// Create test user
	user := &models.User{
		ID:           "test-user-1",
		WechatOpenID: "test-openid-1",
		Nickname:     "测试用户",
		AvatarURL:    "http://example.com/avatar.jpg",
		Status:       1,
	}
	db.Create(user)

	// Get user info
	result, err := service.GetUserInfo(user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "测试用户", result.Nickname)
	assert.Equal(t, "test-openid-1", result.WechatOpenID)
}

func TestGetUserInfoNotFound(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	// Try to get non-existent user
	_, err := service.GetUserInfo("non-existent-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户不存在")
}

func TestUpdateUserInfo(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	user := &models.User{
		ID:           "test-user-2",
		WechatOpenID: "test-openid-2",
		Nickname:     "原始昵称",
		Status:       1,
	}
	db.Create(user)

	// Update user info
	err := service.UpdateUserInfo(user.ID, "新昵称", "13800138000")

	assert.NoError(t, err)

	// Verify update
	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	assert.Equal(t, "新昵称", updated.Nickname)
	assert.Equal(t, "13800138000", updated.Phone)
}

func TestGenerateJWT(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	userID := "test-user-jwt"
	token, err := service.generateJWT(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	validatedUserID, err := service.ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestValidateJWTInvalid(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	// Try to validate invalid token
	_, err := service.ValidateJWT("invalid-token")

	assert.Error(t, err)
}

func TestGetUserMemorials(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	user := &models.User{
		ID:           "test-user-3",
		WechatOpenID: "test-openid-3",
		Nickname:     "测试用户",
		Status:       1,
	}
	db.Create(user)

	// Create test memorials
	memorial1 := &models.Memorial{
		ID:           "memorial-1",
		CreatorID:    user.ID,
		DeceasedName: "纪念馆1",
		Status:       1,
	}
	memorial2 := &models.Memorial{
		ID:           "memorial-2",
		CreatorID:    user.ID,
		DeceasedName: "纪念馆2",
		Status:       1,
	}
	db.Create(memorial1)
	db.Create(memorial2)

	// Get user memorials
	memorials, total, err := service.GetUserMemorials(user.ID, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, memorials, 2)
}

func TestGetUserStatistics(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	user := &models.User{
		ID:           "test-user-4",
		WechatOpenID: "test-openid-4",
		Nickname:     "测试用户",
		Status:       1,
	}
	db.Create(user)

	// Create test data
	memorial := &models.Memorial{
		ID:           "memorial-3",
		CreatorID:    user.ID,
		DeceasedName: "纪念馆",
		Status:       1,
	}
	db.Create(memorial)

	worship := &models.WorshipRecord{
		ID:          "worship-1",
		MemorialID:  memorial.ID,
		UserID:      user.ID,
		WorshipType: "flower",
	}
	db.Create(worship)

	// Get statistics
	stats, err := service.GetUserStatistics(user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(1), stats["memorial_count"])
	assert.Equal(t, int64(1), stats["worship_count"])
}

func TestGenerateInviteCode(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireTime: 3600,
		},
	}
	service := NewUserService(db, cfg)

	code := service.GenerateInviteCode()

	assert.NotEmpty(t, code)
	assert.Len(t, code, 8) // hex encoded 4 bytes = 8 characters
}
