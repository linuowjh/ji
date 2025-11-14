package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PremiumService struct {
	db *gorm.DB
}

func NewPremiumService(db *gorm.DB) *PremiumService {
	return &PremiumService{
		db: db,
	}
}

// 套餐管理

// GetPremiumPackages 获取高级套餐列表
func (s *PremiumService) GetPremiumPackages(packageType string, activeOnly bool) ([]models.PremiumPackage, error) {
	var packages []models.PremiumPackage
	query := s.db.Model(&models.PremiumPackage{})
	
	if packageType != "" {
		query = query.Where("package_type = ?", packageType)
	}
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Order("sort_order ASC, price ASC").Find(&packages).Error
	return packages, err
}

// GetPremiumPackage 获取单个套餐详情
func (s *PremiumService) GetPremiumPackage(packageID string) (*models.PremiumPackage, error) {
	var pkg models.PremiumPackage
	err := s.db.Where("id = ?", packageID).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("套餐不存在")
		}
		return nil, err
	}
	return &pkg, nil
}

// CreatePremiumPackage 创建高级套餐
func (s *PremiumService) CreatePremiumPackage(pkg *models.PremiumPackage) error {
	pkg.ID = uuid.New().String()
	pkg.CreatedAt = time.Now()
	pkg.UpdatedAt = time.Now()
	
	return s.db.Create(pkg).Error
}

// UpdatePremiumPackage 更新套餐信息
func (s *PremiumService) UpdatePremiumPackage(packageID string, updates map[string]interface{}) error {
	var pkg models.PremiumPackage
	err := s.db.Where("id = ?", packageID).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("套餐不存在")
		}
		return err
	}
	
	updates["updated_at"] = time.Now()
	return s.db.Model(&pkg).Updates(updates).Error
}

// 用户订阅管理

// Subscribe 用户订阅套餐
func (s *PremiumService) Subscribe(userID, packageID, memorialID string) (*models.UserSubscription, error) {
	// 获取套餐信息
	pkg, err := s.GetPremiumPackage(packageID)
	if err != nil {
		return nil, err
	}
	
	if !pkg.IsActive {
		return nil, errors.New("该套餐已下架")
	}
	
	// 检查用户是否已有相同类型的有效订阅
	var existingSub models.UserSubscription
	err = s.db.Where("user_id = ? AND package_id = ? AND status = ? AND end_date > ?", 
		userID, packageID, "active", time.Now()).First(&existingSub).Error
	
	if err == nil {
		return nil, errors.New("您已订阅该套餐")
	}
	
	// 创建订阅
	subscription := &models.UserSubscription{
		ID:            uuid.New().String(),
		UserID:        userID,
		PackageID:     packageID,
		MemorialID:    memorialID,
		Status:        "active",
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, pkg.Duration),
		AutoRenew:     false,
		PaymentAmount: pkg.Price,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	if err := s.db.Create(subscription).Error; err != nil {
		return nil, err
	}
	
	// 如果是存储套餐，更新用户存储空间
	if pkg.PackageType == "storage" && pkg.StorageSize > 0 {
		if err := s.UpdateUserStorage(userID, pkg.StorageSize); err != nil {
			return nil, err
		}
	}
	
	// 记录服务使用日志
	s.LogServiceUsage(userID, "subscription", map[string]interface{}{
		"package_id":      packageID,
		"subscription_id": subscription.ID,
	})
	
	return subscription, nil
}

// GetUserSubscriptions 获取用户订阅列表
func (s *PremiumService) GetUserSubscriptions(userID string, status string) ([]models.UserSubscription, error) {
	var subscriptions []models.UserSubscription
	query := s.db.Where("user_id = ?", userID).Preload("Package")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Order("created_at DESC").Find(&subscriptions).Error
	return subscriptions, err
}

// GetSubscription 获取订阅详情
func (s *PremiumService) GetSubscription(subscriptionID string) (*models.UserSubscription, error) {
	var subscription models.UserSubscription
	err := s.db.Where("id = ?", subscriptionID).Preload("Package").First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订阅不存在")
		}
		return nil, err
	}
	return &subscription, nil
}

// CancelSubscription 取消订阅
func (s *PremiumService) CancelSubscription(subscriptionID, userID string) error {
	var subscription models.UserSubscription
	err := s.db.Where("id = ? AND user_id = ?", subscriptionID, userID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订阅不存在")
		}
		return err
	}
	
	if subscription.Status != "active" {
		return errors.New("订阅状态不正确")
	}
	
	now := time.Now()
	return s.db.Model(&subscription).Updates(map[string]interface{}{
		"status":       "cancelled",
		"auto_renew":   false,
		"cancelled_at": now,
		"updated_at":   now,
	}).Error
}

// RenewSubscription 续订
func (s *PremiumService) RenewSubscription(subscriptionID, userID string) error {
	subscription, err := s.GetSubscription(subscriptionID)
	if err != nil {
		return err
	}
	
	if subscription.UserID != userID {
		return errors.New("无权操作")
	}
	
	pkg, err := s.GetPremiumPackage(subscription.PackageID)
	if err != nil {
		return err
	}
	
	// 计算新的结束日期
	var newEndDate time.Time
	if subscription.EndDate.After(time.Now()) {
		// 如果还未过期，从原结束日期延长
		newEndDate = subscription.EndDate.AddDate(0, 0, pkg.Duration)
	} else {
		// 如果已过期，从现在开始计算
		newEndDate = time.Now().AddDate(0, 0, pkg.Duration)
	}
	
	return s.db.Model(subscription).Updates(map[string]interface{}{
		"status":     "active",
		"end_date":   newEndDate,
		"updated_at": time.Now(),
	}).Error
}

// CheckSubscriptionExpiry 检查订阅是否过期
func (s *PremiumService) CheckSubscriptionExpiry() error {
	// 查找所有已过期但状态仍为active的订阅
	var expiredSubs []models.UserSubscription
	err := s.db.Where("status = ? AND end_date < ?", "active", time.Now()).Find(&expiredSubs).Error
	if err != nil {
		return err
	}
	
	// 更新状态为expired
	for _, sub := range expiredSubs {
		s.db.Model(&sub).Updates(map[string]interface{}{
			"status":     "expired",
			"updated_at": time.Now(),
		})
	}
	
	return nil
}

// 纪念馆升级管理

// UpgradeMemorial 升级纪念馆
func (s *PremiumService) UpgradeMemorial(memorialID, subscriptionID, upgradeType string, upgradeData map[string]interface{}) error {
	// 验证订阅是否有效
	subscription, err := s.GetSubscription(subscriptionID)
	if err != nil {
		return err
	}
	
	if subscription.Status != "active" {
		return errors.New("订阅已失效")
	}
	
	if subscription.EndDate.Before(time.Now()) {
		return errors.New("订阅已过期")
	}
	
	// 序列化升级数据
	dataJSON, err := json.Marshal(upgradeData)
	if err != nil {
		return fmt.Errorf("序列化升级数据失败: %v", err)
	}
	
	// 创建升级记录
	upgrade := &models.MemorialUpgrade{
		ID:             uuid.New().String(),
		MemorialID:     memorialID,
		SubscriptionID: subscriptionID,
		UpgradeType:    upgradeType,
		UpgradeData:    string(dataJSON),
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	if err := s.db.Create(upgrade).Error; err != nil {
		return err
	}
	
	// 记录服务使用日志
	s.LogServiceUsage(subscription.UserID, "memorial_upgrade", map[string]interface{}{
		"memorial_id": memorialID,
		"upgrade_type": upgradeType,
	})
	
	return nil
}

// GetMemorialUpgrades 获取纪念馆升级记录
func (s *PremiumService) GetMemorialUpgrades(memorialID string) ([]models.MemorialUpgrade, error) {
	var upgrades []models.MemorialUpgrade
	err := s.db.Where("memorial_id = ? AND is_active = ?", memorialID, true).
		Preload("Subscription").
		Preload("Subscription.Package").
		Order("created_at DESC").
		Find(&upgrades).Error
	return upgrades, err
}

// 定制模板管理

// CreateCustomTemplate 创建定制模板
func (s *PremiumService) CreateCustomTemplate(template *models.CustomTemplate) error {
	template.ID = uuid.New().String()
	template.Status = "draft"
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	
	if err := s.db.Create(template).Error; err != nil {
		return err
	}
	
	// 记录服务使用日志
	s.LogServiceUsage(template.UserID, "custom_template", map[string]interface{}{
		"template_id":   template.ID,
		"template_type": template.TemplateType,
	})
	
	return nil
}

// GetUserCustomTemplates 获取用户定制模板列表
func (s *PremiumService) GetUserCustomTemplates(userID, templateType string) ([]models.CustomTemplate, error) {
	var templates []models.CustomTemplate
	query := s.db.Where("user_id = ?", userID)
	
	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
	}
	
	err := query.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

// UpdateCustomTemplate 更新定制模板
func (s *PremiumService) UpdateCustomTemplate(templateID, userID string, updates map[string]interface{}) error {
	var template models.CustomTemplate
	err := s.db.Where("id = ? AND user_id = ?", templateID, userID).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模板不存在")
		}
		return err
	}
	
	updates["updated_at"] = time.Now()
	return s.db.Model(&template).Updates(updates).Error
}

// DeleteCustomTemplate 删除定制模板
func (s *PremiumService) DeleteCustomTemplate(templateID, userID string) error {
	result := s.db.Where("id = ? AND user_id = ?", templateID, userID).Delete(&models.CustomTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("模板不存在")
	}
	return nil
}

// 存储管理

// GetUserStorage 获取用户存储使用情况
func (s *PremiumService) GetUserStorage(userID string) (*models.StorageUsage, error) {
	var storage models.StorageUsage
	err := s.db.Where("user_id = ?", userID).First(&storage).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，创建默认记录
			storage = models.StorageUsage{
				ID:          uuid.New().String(),
				UserID:      userID,
				UsedSpace:   0,
				TotalSpace:  104857600, // 默认100MB
				FileCount:   0,
				LastUpdated: time.Now(),
			}
			if err := s.db.Create(&storage).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	return &storage, nil
}

// UpdateUserStorage 更新用户存储空间
func (s *PremiumService) UpdateUserStorage(userID string, additionalSpace int64) error {
	storage, err := s.GetUserStorage(userID)
	if err != nil {
		return err
	}
	
	return s.db.Model(storage).Updates(map[string]interface{}{
		"total_space":  storage.TotalSpace + additionalSpace,
		"last_updated": time.Now(),
	}).Error
}

// UpdateStorageUsage 更新存储使用量
func (s *PremiumService) UpdateStorageUsage(userID string, spaceChange int64, fileCountChange int) error {
	storage, err := s.GetUserStorage(userID)
	if err != nil {
		return err
	}
	
	newUsedSpace := storage.UsedSpace + spaceChange
	newFileCount := storage.FileCount + fileCountChange
	
	// 检查是否超出限制
	if newUsedSpace > storage.TotalSpace {
		return errors.New("存储空间不足")
	}
	
	return s.db.Model(storage).Updates(map[string]interface{}{
		"used_space":   newUsedSpace,
		"file_count":   newFileCount,
		"last_updated": time.Now(),
	}).Error
}

// CheckStorageLimit 检查存储限制
func (s *PremiumService) CheckStorageLimit(userID string, requiredSpace int64) (bool, error) {
	storage, err := s.GetUserStorage(userID)
	if err != nil {
		return false, err
	}
	
	availableSpace := storage.TotalSpace - storage.UsedSpace
	return availableSpace >= requiredSpace, nil
}

// 服务使用日志

// LogServiceUsage 记录服务使用
func (s *PremiumService) LogServiceUsage(userID, serviceType string, serviceData map[string]interface{}) error {
	dataJSON, _ := json.Marshal(serviceData)
	
	log := &models.ServiceUsageLog{
		ID:          uuid.New().String(),
		UserID:      userID,
		ServiceType: serviceType,
		ServiceData: string(dataJSON),
		UsageCount:  1,
		CreatedAt:   time.Now(),
	}
	
	return s.db.Create(log).Error
}

// GetServiceUsageStats 获取服务使用统计
func (s *PremiumService) GetServiceUsageStats(userID string, startTime, endTime time.Time) (map[string]interface{}, error) {
	var logs []models.ServiceUsageLog
	query := s.db.Where("user_id = ?", userID)
	
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}
	
	err := query.Find(&logs).Error
	if err != nil {
		return nil, err
	}
	
	// 统计各类服务使用次数
	stats := make(map[string]int)
	for _, log := range logs {
		stats[log.ServiceType] += log.UsageCount
	}
	
	return map[string]interface{}{
		"total_usage": len(logs),
		"by_service":  stats,
		"start_time":  startTime,
		"end_time":    endTime,
	}, nil
}

// 初始化默认套餐
func (s *PremiumService) InitDefaultPackages() error {
	defaultPackages := []models.PremiumPackage{
		{
			ID:          uuid.New().String(),
			PackageName: "基础版",
			PackageType: "memorial",
			Description: "适合个人使用的基础纪念馆服务",
			Features:    `["100MB存储空间", "基础主题模板", "标准祭扫功能"]`,
			Price:       0,
			Duration:    365,
			StorageSize: 104857600, // 100MB
			IsActive:    true,
			SortOrder:   1,
		},
		{
			ID:          uuid.New().String(),
			PackageName: "高级版",
			PackageType: "memorial",
			Description: "提供更多定制化功能和存储空间",
			Features:    `["500MB存储空间", "高级主题模板", "定制墓碑样式", "老照片修复", "优先客服支持"]`,
			Price:       99.00,
			Duration:    365,
			StorageSize: 524288000, // 500MB
			IsActive:    true,
			SortOrder:   2,
		},
		{
			ID:          uuid.New().String(),
			PackageName: "尊享版",
			PackageType: "memorial",
			Description: "最完整的纪念馆服务体验",
			Features:    `["2GB存储空间", "所有主题模板", "完全定制化", "老照片修复", "专属追思会", "数据备份服务", "专属客服"]`,
			Price:       299.00,
			Duration:    365,
			StorageSize: 2147483648, // 2GB
			IsActive:    true,
			SortOrder:   3,
		},
		{
			ID:          uuid.New().String(),
			PackageName: "扩展存储包",
			PackageType: "storage",
			Description: "额外增加1GB存储空间",
			Features:    `["1GB额外存储空间", "永久有效"]`,
			Price:       49.00,
			Duration:    36500, // 100年
			StorageSize: 1073741824, // 1GB
			IsActive:    true,
			SortOrder:   4,
		},
	}
	
	for _, pkg := range defaultPackages {
		// 检查是否已存在
		var existing models.PremiumPackage
		err := s.db.Where("package_name = ?", pkg.PackageName).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			if err := s.db.Create(&pkg).Error; err != nil {
				return fmt.Errorf("创建默认套餐失败: %v", err)
			}
		}
	}
	
	return nil
}
