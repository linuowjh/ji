package middleware

import (
	"net/http"
	"strings"
	"yun-nian-memorial/internal/utils"

	"github.com/gin-gonic/gin"
)

// SQL注入防护中间件
func SQLInjectionProtection() gin.HandlerFunc {
	sanitizer := utils.NewSQLSanitizer()
	
	return func(c *gin.Context) {
		// 检查URL参数
		for key, values := range c.Request.URL.Query() {
			for _, value := range values {
				if sanitizer.ContainsDangerousContent(value) {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    1007,
						"message": "请求参数包含不安全内容",
						"field":   key,
					})
					c.Abort()
					return
				}
			}
		}

		// 检查表单数据
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if err := c.Request.ParseForm(); err == nil {
				for key, values := range c.Request.PostForm {
					for _, value := range values {
						if sanitizer.ContainsDangerousContent(value) {
							c.JSON(http.StatusBadRequest, gin.H{
								"code":    1007,
								"message": "表单数据包含不安全内容",
								"field":   key,
							})
							c.Abort()
							return
						}
					}
				}
			}
		}

		c.Next()
	}
}

// XSS防护中间件
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置XSS防护头
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'none'; frame-ancestors 'none'; form-action 'self'; base-uri 'self';")
		
		c.Next()
	}
}

// HTTPS重定向中间件
func HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("X-Forwarded-Proto") == "http" {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		
		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")
		
		// XSS防护
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// 强制HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}

// IP白名单中间件（用于管理员接口）
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		// 如果没有配置白名单，则允许所有IP
		if len(allowedIPs) == 0 {
			c.Next()
			return
		}
		
		// 检查IP是否在白名单中
		allowed := false
		for _, ip := range allowedIPs {
			if ip == clientIP || ip == "*" {
				allowed = true
				break
			}
			
			// 支持CIDR格式的IP段（简单实现）
			if strings.Contains(ip, "/") {
				// 这里可以实现更复杂的CIDR匹配逻辑
				// 暂时简化处理
				if strings.HasPrefix(clientIP, strings.Split(ip, "/")[0][:strings.LastIndex(strings.Split(ip, "/")[0], ".")]) {
					allowed = true
					break
				}
			}
		}
		
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    1003,
				"message": "访问被拒绝：IP地址不在允许列表中",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// 请求大小限制中间件
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    1008,
				"message": "请求体过大",
			})
			c.Abort()
			return
		}
		
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// 用户代理检查中间件（防止恶意爬虫）
func UserAgentFilter() gin.HandlerFunc {
	// 恶意用户代理黑名单
	blacklistedUserAgents := []string{
		"sqlmap",
		"nikto",
		"nessus",
		"openvas",
		"nmap",
		"masscan",
		"zap",
		"w3af",
		"burp",
		"acunetix",
		"appscan",
	}
	
	return func(c *gin.Context) {
		userAgent := strings.ToLower(c.Request.UserAgent())
		
		// 检查是否为空用户代理
		if userAgent == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    1009,
				"message": "缺少用户代理信息",
			})
			c.Abort()
			return
		}
		
		// 检查黑名单
		for _, blacklisted := range blacklistedUserAgents {
			if strings.Contains(userAgent, blacklisted) {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    1010,
					"message": "访问被拒绝：不允许的用户代理",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// API密钥验证中间件（用于第三方接口）
func APIKeyAuth(validAPIKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}
		
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1011,
				"message": "缺少API密钥",
			})
			c.Abort()
			return
		}
		
		// 验证API密钥
		valid := false
		for _, validKey := range validAPIKeys {
			if apiKey == validKey {
				valid = true
				break
			}
		}
		
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1012,
				"message": "无效的API密钥",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// 请求日志中间件（安全日志）
func SecurityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求信息（脱敏处理）
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		
		// 检查是否为可疑请求
		suspicious := false
		reasons := []string{}
		
		// 检查是否包含可疑路径
		suspiciousPaths := []string{
			"admin", "wp-admin", "phpmyadmin", "mysql", "sql",
			"config", "backup", "test", "debug", "api/v1/admin",
		}
		
		for _, suspiciousPath := range suspiciousPaths {
			if strings.Contains(strings.ToLower(path), suspiciousPath) {
				suspicious = true
				reasons = append(reasons, "suspicious_path")
				break
			}
		}
		
		// 检查是否为内网IP访问敏感接口
		if strings.Contains(path, "/admin/") && !utils.IsPrivateIP(clientIP) {
			suspicious = true
			reasons = append(reasons, "external_admin_access")
		}
		
		// 如果是可疑请求，记录详细日志
		if suspicious {
			logMsg := utils.SecureLog(
				"Suspicious request: IP=%s, Method=%s, Path=%s, UserAgent=%s, Reasons=%v",
				clientIP, method, path, userAgent, reasons,
			)
			// 这里可以集成到日志系统
			println("SECURITY WARNING:", logMsg)
		}
		
		c.Next()
	}
}