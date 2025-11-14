package middleware

import (
	"net/http"
	"strings"
	"yun-nian-memorial/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTAuth JWT认证中间件
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "请提供认证令牌",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 解析JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "认证令牌无效",
			})
			c.Abort()
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("user_id", userID)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    1002,
					"message": "认证令牌解析失败",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "认证令牌解析失败",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 权限验证中间件
func RequirePermission(db *gorm.DB, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "用户未登录",
			})
			c.Abort()
			return
		}

		// 检查用户状态
		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    2001,
				"message": "用户不存在",
			})
			c.Abort()
			return
		}

		if user.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    2002,
				"message": "用户已被禁用",
			})
			c.Abort()
			return
		}

		// 这里可以扩展更复杂的权限检查逻辑
		// 目前简单检查用户状态即可
		c.Set("user", user)
		c.Next()
	}
}

// RequireMemorialAccess 纪念馆访问权限中间件
func RequireMemorialAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "用户未登录",
			})
			c.Abort()
			return
		}

		memorialID := c.Param("id")
		if memorialID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    1001,
				"message": "纪念馆ID不能为空",
			})
			c.Abort()
			return
		}

		// 查询纪念馆信息
		var memorial models.Memorial
		if err := db.Where("id = ? AND status = ?", memorialID, 1).First(&memorial).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    3001,
					"message": "纪念馆不存在",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    1005,
					"message": "查询纪念馆失败",
				})
			}
			c.Abort()
			return
		}

		// 检查访问权限
		hasAccess := false

		// 1. 创建者有完全访问权限
		if memorial.CreatorID == userID.(string) {
			hasAccess = true
			c.Set("access_level", "owner")
		} else if memorial.PrivacyLevel == 1 {
			// 2. 家族可见的纪念馆，检查是否为家族成员
			var count int64
			db.Table("memorial_families mf").
				Joins("JOIN family_members fm ON mf.family_id = fm.family_id").
				Where("mf.memorial_id = ? AND fm.user_id = ?", memorialID, userID).
				Count(&count)
			
			if count > 0 {
				hasAccess = true
				c.Set("access_level", "family")
			}
		}
		// 私密纪念馆(privacy_level = 2)只有创建者可以访问

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    3002,
				"message": "无权访问此纪念馆",
			})
			c.Abort()
			return
		}

		c.Set("memorial", memorial)
		c.Next()
	}
}

// RequireFamilyAccess 家族访问权限中间件
func RequireFamilyAccess(db *gorm.DB, requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "用户未登录",
			})
			c.Abort()
			return
		}

		familyID := c.Param("id")
		if familyID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    1001,
				"message": "家族ID不能为空",
			})
			c.Abort()
			return
		}

		// 查询家族成员信息
		var member models.FamilyMember
		err := db.Where("family_id = ? AND user_id = ?", familyID, userID).First(&member).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    4003,
					"message": "您不是该家族成员",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    1005,
					"message": "查询家族成员失败",
				})
			}
			c.Abort()
			return
		}

		// 如果需要管理员权限
		if requireAdmin && member.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    1003,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Set("family_member", member)
		c.Set("is_family_admin", member.Role == "admin")
		c.Next()
	}
}