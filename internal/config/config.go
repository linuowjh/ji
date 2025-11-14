package config

import (
	"os"
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	JWT        JWTConfig        `json:"jwt"`
	Wechat     WechatConfig     `json:"wechat"`
	COS        COSConfig        `json:"cos"`
	Encryption EncryptionConfig `json:"encryption"`
	Security   SecurityConfig   `json:"security"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Mode string `json:"mode"`
}

type DatabaseConfig struct {
	MySQL MySQLConfig `json:"mysql"`
	Redis RedisConfig `json:"redis"`
}

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	LogLevel string `json:"log_level"` // silent, error, warn, info
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type JWTConfig struct {
	Secret     string `json:"secret"`
	ExpireTime int    `json:"expire_time"`
}

type WechatConfig struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type COSConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
}

type EncryptionConfig struct {
	SecretKey string `json:"secret_key"`
}

type SecurityConfig struct {
	AdminIPWhitelist []string `json:"admin_ip_whitelist"`
	MaxRequestSize   int64    `json:"max_request_size"`
	EnableHTTPS      bool     `json:"enable_https"`
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			MySQL: MySQLConfig{
				Host:     getEnv("MYSQL_HOST", "localhost"),
				Port:     getEnv("MYSQL_PORT", "3306"),
				Username: getEnv("MYSQL_USERNAME", "root"),
				Password: getEnv("MYSQL_PASSWORD", ""),
				Database: getEnv("MYSQL_DATABASE", "yun_nian_memorial"),
				LogLevel: getEnv("MYSQL_LOG_LEVEL", "error"), // 默认只显示错误
			},
			Redis: RedisConfig{
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     getEnv("REDIS_PORT", "6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       0,
			},
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "yun-nian-memorial-secret"),
			ExpireTime: 7 * 24 * 3600, // 7天
		},
		Wechat: WechatConfig{
			AppID:     getEnv("WECHAT_APP_ID", ""),
			AppSecret: getEnv("WECHAT_APP_SECRET", ""),
		},
		COS: COSConfig{
			SecretID:  getEnv("COS_SECRET_ID", ""),
			SecretKey: getEnv("COS_SECRET_KEY", ""),
			Region:    getEnv("COS_REGION", "ap-beijing"),
			Bucket:    getEnv("COS_BUCKET", ""),
		},
		Encryption: EncryptionConfig{
			SecretKey: getEnv("ENCRYPTION_SECRET_KEY", "yun-nian-memorial-encryption-key-2024"),
		},
		Security: SecurityConfig{
			AdminIPWhitelist: []string{"127.0.0.1", "::1"}, // 默认只允许本地访问
			MaxRequestSize:   10 * 1024 * 1024,             // 10MB
			EnableHTTPS:      getEnv("ENABLE_HTTPS", "false") == "true",
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
