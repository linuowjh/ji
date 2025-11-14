package services

import (
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/utils"
)

// 数据加密服务
type EncryptionService struct {
	sensitiveDataManager *utils.SensitiveDataManager
}

// 创建加密服务
func NewEncryptionService(cfg *config.Config) *EncryptionService {
	// 从配置中获取加密密钥，如果没有则使用默认值
	secretKey := cfg.Encryption.SecretKey
	if secretKey == "" {
		secretKey = "yun-nian-memorial-default-secret-key-2024"
	}

	return &EncryptionService{
		sensitiveDataManager: utils.NewSensitiveDataManager(secretKey),
	}
}

// 加密用户敏感信息
func (s *EncryptionService) EncryptUserSensitiveData(data string) (string, error) {
	return s.sensitiveDataManager.EncryptSensitiveField(data)
}

// 解密用户敏感信息
func (s *EncryptionService) DecryptUserSensitiveData(encryptedData string) (string, error) {
	return s.sensitiveDataManager.DecryptSensitiveField(encryptedData)
}

// 加密手机号
func (s *EncryptionService) EncryptPhone(phone string) (string, error) {
	if phone == "" {
		return "", nil
	}
	return s.EncryptUserSensitiveData(phone)
}

// 解密手机号
func (s *EncryptionService) DecryptPhone(encryptedPhone string) (string, error) {
	if encryptedPhone == "" {
		return "", nil
	}
	return s.DecryptUserSensitiveData(encryptedPhone)
}

// 加密邮箱
func (s *EncryptionService) EncryptEmail(email string) (string, error) {
	if email == "" {
		return "", nil
	}
	return s.EncryptUserSensitiveData(email)
}

// 解密邮箱
func (s *EncryptionService) DecryptEmail(encryptedEmail string) (string, error) {
	if encryptedEmail == "" {
		return "", nil
	}
	return s.DecryptUserSensitiveData(encryptedEmail)
}

// 加密身份证号
func (s *EncryptionService) EncryptIDCard(idCard string) (string, error) {
	if idCard == "" {
		return "", nil
	}
	return s.EncryptUserSensitiveData(idCard)
}

// 解密身份证号
func (s *EncryptionService) DecryptIDCard(encryptedIDCard string) (string, error) {
	if encryptedIDCard == "" {
		return "", nil
	}
	return s.DecryptUserSensitiveData(encryptedIDCard)
}

// 加密地址信息
func (s *EncryptionService) EncryptAddress(address string) (string, error) {
	if address == "" {
		return "", nil
	}
	return s.EncryptUserSensitiveData(address)
}

// 解密地址信息
func (s *EncryptionService) DecryptAddress(encryptedAddress string) (string, error) {
	if encryptedAddress == "" {
		return "", nil
	}
	return s.DecryptUserSensitiveData(encryptedAddress)
}

// 脱敏显示手机号
func (s *EncryptionService) MaskPhone(phone string) string {
	if len(phone) != 11 {
		return utils.MaskSensitiveData(phone)
	}
	return phone[:3] + "****" + phone[7:]
}

// 脱敏显示邮箱
func (s *EncryptionService) MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}
	
	if atIndex <= 0 {
		return utils.MaskSensitiveData(email)
	}
	
	username := email[:atIndex]
	domain := email[atIndex:]
	
	if len(username) <= 2 {
		return username + domain
	} else if len(username) <= 4 {
		return username[:1] + "**" + username[len(username)-1:] + domain
	} else {
		return username[:2] + "****" + username[len(username)-2:] + domain
	}
}

// 脱敏显示身份证号
func (s *EncryptionService) MaskIDCard(idCard string) string {
	if len(idCard) != 18 {
		return utils.MaskSensitiveData(idCard)
	}
	return idCard[:6] + "********" + idCard[14:]
}

// 脱敏显示地址
func (s *EncryptionService) MaskAddress(address string) string {
	if len(address) <= 10 {
		return utils.MaskSensitiveData(address)
	}
	return address[:4] + "****" + address[len(address)-4:]
}