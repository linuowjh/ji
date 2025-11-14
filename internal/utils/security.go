package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// 数据加密工具
type Encryptor struct {
	key []byte
}

// 创建新的加密器
func NewEncryptor(secretKey string) *Encryptor {
	// 使用SHA256生成32字节的密钥
	hash := sha256.Sum256([]byte(secretKey))
	return &Encryptor{
		key: hash[:],
	}
}

// AES加密
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// 返回base64编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES解密
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// 哈希密码
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// 验证密码
func VerifyPassword(password, hashedPassword string) bool {
	return HashPassword(password) == hashedPassword
}

// SQL注入防护
type SQLSanitizer struct {
	// 危险的SQL关键词
	dangerousKeywords []string
	// 危险的字符模式
	dangerousPatterns []*regexp.Regexp
}

// 创建SQL清理器
func NewSQLSanitizer() *SQLSanitizer {
	dangerousKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "UNION", "SCRIPT", "JAVASCRIPT", "VBSCRIPT",
		"ONLOAD", "ONERROR", "ONCLICK", "ALERT", "CONFIRM", "PROMPT",
		"EVAL", "EXPRESSION", "APPLET", "OBJECT", "EMBED", "FORM",
		"INPUT", "TEXTAREA", "IFRAME", "FRAME", "FRAMESET", "META",
		"LINK", "STYLE", "TITLE", "BASE", "BGSOUND", "BLINK",
	}

	dangerousPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(\s|^)(union|select|insert|update|delete|drop|create|alter|exec|execute)(\s|$)`),
		regexp.MustCompile(`(?i)(\s|^)(script|javascript|vbscript)(\s|$)`),
		regexp.MustCompile(`(?i)(<|&lt;)\s*script`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)vbscript:`),
		regexp.MustCompile(`(?i)on\w+\s*=`),
		regexp.MustCompile(`(?i)expression\s*\(`),
		regexp.MustCompile(`(?i)url\s*\(`),
		regexp.MustCompile(`(?i)@import`),
		regexp.MustCompile(`(?i)<!--.*-->`),
		regexp.MustCompile(`(?i)<\s*\w+.*>`),
		regexp.MustCompile(`(?i)</\s*\w+.*>`),
		regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`), // 控制字符
	}

	return &SQLSanitizer{
		dangerousKeywords: dangerousKeywords,
		dangerousPatterns: dangerousPatterns,
	}
}

// 检查是否包含危险内容
func (s *SQLSanitizer) ContainsDangerousContent(input string) bool {
	if input == "" {
		return false
	}

	// 转换为大写进行关键词检查
	upperInput := strings.ToUpper(input)
	
	// 检查危险关键词
	for _, keyword := range s.dangerousKeywords {
		if strings.Contains(upperInput, keyword) {
			return true
		}
	}

	// 检查危险模式
	for _, pattern := range s.dangerousPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}

	return false
}

// 清理输入内容
func (s *SQLSanitizer) SanitizeInput(input string) string {
	if input == "" {
		return input
	}

	// 移除危险字符
	sanitized := input
	
	// 替换HTML特殊字符
	sanitized = strings.ReplaceAll(sanitized, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
	sanitized = strings.ReplaceAll(sanitized, "'", "&#x27;")
	sanitized = strings.ReplaceAll(sanitized, "&", "&amp;")
	
	// 移除控制字符
	controlCharPattern := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	sanitized = controlCharPattern.ReplaceAllString(sanitized, "")
	
	// 移除多余的空白字符
	spacePattern := regexp.MustCompile(`\s+`)
	sanitized = spacePattern.ReplaceAllString(sanitized, " ")
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// 验证输入安全性
func (s *SQLSanitizer) ValidateInput(input string) error {
	if s.ContainsDangerousContent(input) {
		return errors.New("输入内容包含不安全的字符或关键词")
	}
	return nil
}

// 全局SQL清理器实例
var GlobalSQLSanitizer = NewSQLSanitizer()

// 便捷函数
func SanitizeString(input string) string {
	return GlobalSQLSanitizer.SanitizeInput(input)
}

func ValidateString(input string) error {
	return GlobalSQLSanitizer.ValidateInput(input)
}

func ContainsDangerousContent(input string) bool {
	return GlobalSQLSanitizer.ContainsDangerousContent(input)
}

// 敏感数据字段加密
type SensitiveDataManager struct {
	encryptor *Encryptor
}

func NewSensitiveDataManager(secretKey string) *SensitiveDataManager {
	return &SensitiveDataManager{
		encryptor: NewEncryptor(secretKey),
	}
}

// 加密敏感字段
func (sdm *SensitiveDataManager) EncryptSensitiveField(data string) (string, error) {
	if data == "" {
		return "", nil
	}
	return sdm.encryptor.Encrypt(data)
}

// 解密敏感字段
func (sdm *SensitiveDataManager) DecryptSensitiveField(encryptedData string) (string, error) {
	if encryptedData == "" {
		return "", nil
	}
	return sdm.encryptor.Decrypt(encryptedData)
}

// 脱敏显示（用于日志和调试）
func MaskSensitiveData(data string) string {
	if len(data) <= 4 {
		return strings.Repeat("*", len(data))
	}
	
	if len(data) <= 8 {
		return data[:2] + strings.Repeat("*", len(data)-4) + data[len(data)-2:]
	}
	
	return data[:3] + strings.Repeat("*", len(data)-6) + data[len(data)-3:]
}

// 生成安全的随机字符串
func GenerateSecureRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	
	return string(b), nil
}

// IP地址验证
func IsValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		
		num := 0
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
			num = num*10 + int(char-'0')
		}
		
		if num > 255 {
			return false
		}
	}
	
	return true
}

// 检查是否为内网IP
func IsPrivateIP(ip string) bool {
	if !IsValidIP(ip) {
		return false
	}
	
	parts := strings.Split(ip, ".")
	first := 0
	second := 0
	
	for _, char := range parts[0] {
		first = first*10 + int(char-'0')
	}
	for _, char := range parts[1] {
		second = second*10 + int(char-'0')
	}
	
	// 10.0.0.0/8
	if first == 10 {
		return true
	}
	
	// 172.16.0.0/12
	if first == 172 && second >= 16 && second <= 31 {
		return true
	}
	
	// 192.168.0.0/16
	if first == 192 && second == 168 {
		return true
	}
	
	// 127.0.0.0/8 (localhost)
	if first == 127 {
		return true
	}
	
	return false
}

// 安全日志记录（脱敏）
func SecureLog(format string, args ...interface{}) string {
	// 对参数进行脱敏处理
	maskedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		if str, ok := arg.(string); ok {
			// 如果是字符串，检查是否可能是敏感信息
			if len(str) > 10 && (strings.Contains(strings.ToLower(str), "password") || 
				strings.Contains(strings.ToLower(str), "token") ||
				strings.Contains(strings.ToLower(str), "secret")) {
				maskedArgs[i] = MaskSensitiveData(str)
			} else {
				maskedArgs[i] = arg
			}
		} else {
			maskedArgs[i] = arg
		}
	}
	
	return fmt.Sprintf(format, maskedArgs...)
}