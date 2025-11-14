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

type SystemConfigService struct {
	db *gorm.DB
}

func NewSystemConfigService(db *gorm.DB) *SystemConfigService {
	return &SystemConfigService{
		db: db,
	}
}

// 祭扫节日配置管理

// GetFestivalConfigs 获取所有祭扫节日配置
func (s *SystemConfigService) GetFestivalConfigs(activeOnly bool) ([]models.FestivalConfig, error) {
	var festivals []models.FestivalConfig
	query := s.db.Model(&models.FestivalConfig{})
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Order("festival_date ASC").Find(&festivals).Error
	return festivals, err
}

// GetFestivalConfig 获取单个节日配置
func (s *SystemConfigService) GetFestivalConfig(festivalID string) (*models.FestivalConfig, error) {
	var festival models.FestivalConfig
	err := s.db.Where("id = ?", festivalID).First(&festival).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("节日配置不存在")
		}
		return nil, err
	}
	return &festival, nil
}

// CreateFestivalConfig 创建祭扫节日配置
func (s *SystemConfigService) CreateFestivalConfig(festival *models.FestivalConfig) error {
	festival.ID = uuid.New().String()
	festival.CreatedAt = time.Now()
	festival.UpdatedAt = time.Now()
	
	// 验证日期格式 MM-DD
	if len(festival.FestivalDate) != 5 || festival.FestivalDate[2] != '-' {
		return errors.New("节日日期格式错误，应为 MM-DD")
	}
	
	return s.db.Create(festival).Error
}

// UpdateFestivalConfig 更新祭扫节日配置
func (s *SystemConfigService) UpdateFestivalConfig(festivalID string, updates map[string]interface{}) error {
	var festival models.FestivalConfig
	err := s.db.Where("id = ?", festivalID).First(&festival).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("节日配置不存在")
		}
		return err
	}
	
	updates["updated_at"] = time.Now()
	return s.db.Model(&festival).Updates(updates).Error
}

// DeleteFestivalConfig 删除祭扫节日配置
func (s *SystemConfigService) DeleteFestivalConfig(festivalID string) error {
	result := s.db.Where("id = ?", festivalID).Delete(&models.FestivalConfig{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("节日配置不存在")
	}
	return nil
}

// 模板配置管理

// GetTemplateConfigs 获取模板配置列表
func (s *SystemConfigService) GetTemplateConfigs(templateType string, activeOnly bool) ([]models.TemplateConfig, error) {
	var templates []models.TemplateConfig
	query := s.db.Model(&models.TemplateConfig{})
	
	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
	}
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Order("sort_order ASC, created_at DESC").Find(&templates).Error
	return templates, err
}

// GetTemplateConfig 获取单个模板配置
func (s *SystemConfigService) GetTemplateConfig(templateID string) (*models.TemplateConfig, error) {
	var template models.TemplateConfig
	err := s.db.Where("id = ?", templateID).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("模板配置不存在")
		}
		return nil, err
	}
	return &template, nil
}

// CreateTemplateConfig 创建模板配置
func (s *SystemConfigService) CreateTemplateConfig(template *models.TemplateConfig) error {
	template.ID = uuid.New().String()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	
	// 验证模板类型
	validTypes := map[string]bool{
		"theme":     true,
		"tombstone": true,
		"prayer":    true,
	}
	if !validTypes[template.TemplateType] {
		return errors.New("无效的模板类型")
	}
	
	return s.db.Create(template).Error
}

// UpdateTemplateConfig 更新模板配置
func (s *SystemConfigService) UpdateTemplateConfig(templateID string, updates map[string]interface{}) error {
	var template models.TemplateConfig
	err := s.db.Where("id = ?", templateID).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模板配置不存在")
		}
		return err
	}
	
	updates["updated_at"] = time.Now()
	return s.db.Model(&template).Updates(updates).Error
}

// DeleteTemplateConfig 删除模板配置
func (s *SystemConfigService) DeleteTemplateConfig(templateID string) error {
	result := s.db.Where("id = ?", templateID).Delete(&models.TemplateConfig{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("模板配置不存在")
	}
	return nil
}

// 系统配置管理

// GetSystemConfig 获取系统配置
func (s *SystemConfigService) GetSystemConfig(configKey string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := s.db.Where("config_key = ? AND is_active = ?", configKey, true).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("配置不存在")
		}
		return nil, err
	}
	return &config, nil
}

// GetSystemConfigsByType 根据类型获取系统配置列表
func (s *SystemConfigService) GetSystemConfigsByType(configType string) ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	query := s.db.Model(&models.SystemConfig{}).Where("is_active = ?", true)
	
	if configType != "" {
		query = query.Where("config_type = ?", configType)
	}
	
	err := query.Order("config_key ASC").Find(&configs).Error
	return configs, err
}

// SetSystemConfig 设置系统配置
func (s *SystemConfigService) SetSystemConfig(configKey, configValue, configType, description string) error {
	var config models.SystemConfig
	err := s.db.Where("config_key = ?", configKey).First(&config).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新配置
			config = models.SystemConfig{
				ID:          uuid.New().String(),
				ConfigKey:   configKey,
				ConfigValue: configValue,
				ConfigType:  configType,
				Description: description,
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			return s.db.Create(&config).Error
		}
		return err
	}
	
	// 更新现有配置
	return s.db.Model(&config).Updates(map[string]interface{}{
		"config_value": configValue,
		"config_type":  configType,
		"description":  description,
		"updated_at":   time.Now(),
	}).Error
}

// DeleteSystemConfig 删除系统配置
func (s *SystemConfigService) DeleteSystemConfig(configKey string) error {
	result := s.db.Where("config_key = ?", configKey).Delete(&models.SystemConfig{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("配置不存在")
	}
	return nil
}

// InitDefaultConfigs 初始化默认配置
func (s *SystemConfigService) InitDefaultConfigs() error {
	// 初始化默认祭扫节日
	defaultFestivals := []models.FestivalConfig{
		{
			ID:           uuid.New().String(),
			Name:         "清明节",
			FestivalDate: "04-05",
			Description:  "清明节是中国传统的祭祖节日",
			ReminderDays: 3,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			Name:         "中元节",
			FestivalDate: "08-15",
			Description:  "中元节（农历七月十五）是祭祀祖先的重要节日",
			ReminderDays: 3,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			Name:         "寒衣节",
			FestivalDate: "10-01",
			Description:  "寒衣节（农历十月初一）是祭祀祖先、送寒衣的节日",
			ReminderDays: 3,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			Name:         "除夕",
			FestivalDate: "12-31",
			Description:  "除夕是农历年的最后一天，祭祖迎新",
			ReminderDays: 3,
			IsActive:     true,
		},
	}
	
	for _, festival := range defaultFestivals {
		// 检查是否已存在
		var existing models.FestivalConfig
		err := s.db.Where("name = ?", festival.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			if err := s.db.Create(&festival).Error; err != nil {
				return fmt.Errorf("创建默认节日配置失败: %v", err)
			}
		}
	}
	
	// 初始化默认模板配置
	defaultTemplates := []models.TemplateConfig{
		{
			ID:           uuid.New().String(),
			TemplateType: "theme",
			TemplateName: "中式传统",
			TemplateData: `{"background": "traditional-bg.jpg", "color_scheme": "warm", "font": "serif"}`,
			IsPremium:    false,
			SortOrder:    1,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			TemplateType: "theme",
			TemplateName: "简约素雅",
			TemplateData: `{"background": "elegant-bg.jpg", "color_scheme": "neutral", "font": "sans-serif"}`,
			IsPremium:    false,
			SortOrder:    2,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			TemplateType: "theme",
			TemplateName: "自然清新",
			TemplateData: `{"background": "nature-bg.jpg", "color_scheme": "green", "font": "sans-serif"}`,
			IsPremium:    false,
			SortOrder:    3,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			TemplateType: "tombstone",
			TemplateName: "大理石",
			TemplateData: `{"material": "marble", "color": "white", "style": "classic"}`,
			IsPremium:    false,
			SortOrder:    1,
			IsActive:     true,
		},
		{
			ID:           uuid.New().String(),
			TemplateType: "tombstone",
			TemplateName: "花岗岩",
			TemplateData: `{"material": "granite", "color": "gray", "style": "modern"}`,
			IsPremium:    false,
			SortOrder:    2,
			IsActive:     true,
		},
	}
	
	for _, template := range defaultTemplates {
		// 检查是否已存在
		var existing models.TemplateConfig
		err := s.db.Where("template_type = ? AND template_name = ?", template.TemplateType, template.TemplateName).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			if err := s.db.Create(&template).Error; err != nil {
				return fmt.Errorf("创建默认模板配置失败: %v", err)
			}
		}
	}
	
	// 初始化默认系统配置
	defaultConfigs := map[string]map[string]string{
		"max_memorial_per_user": {
			"value":       "10",
			"type":        "system",
			"description": "每个用户最多可创建的纪念馆数量",
		},
		"max_upload_size": {
			"value":       "10485760",
			"type":        "system",
			"description": "最大上传文件大小（字节）",
		},
		"enable_auto_backup": {
			"value":       "true",
			"type":        "system",
			"description": "是否启用自动备份",
		},
		"backup_interval_hours": {
			"value":       "24",
			"type":        "system",
			"description": "自动备份间隔（小时）",
		},
	}
	
	for key, config := range defaultConfigs {
		if err := s.SetSystemConfig(key, config["value"], config["type"], config["description"]); err != nil {
			return fmt.Errorf("创建默认系统配置失败: %v", err)
		}
	}
	
	return nil
}

// GetUpcomingFestivals 获取即将到来的节日（用于提醒）
func (s *SystemConfigService) GetUpcomingFestivals(daysAhead int) ([]models.FestivalConfig, error) {
	var festivals []models.FestivalConfig
	err := s.db.Where("is_active = ?", true).Find(&festivals).Error
	if err != nil {
		return nil, err
	}
	
	now := time.Now()
	var upcomingFestivals []models.FestivalConfig
	
	for _, festival := range festivals {
		// 解析节日日期
		festivalDate := fmt.Sprintf("%d-%s", now.Year(), festival.FestivalDate)
		festivalTime, err := time.Parse("2006-01-02", festivalDate)
		if err != nil {
			continue
		}
		
		// 如果今年的节日已过，检查明年的
		if festivalTime.Before(now) {
			festivalDate = fmt.Sprintf("%d-%s", now.Year()+1, festival.FestivalDate)
			festivalTime, _ = time.Parse("2006-01-02", festivalDate)
		}
		
		// 计算距离节日的天数
		daysUntil := int(festivalTime.Sub(now).Hours() / 24)
		
		// 如果在提醒范围内
		if daysUntil >= 0 && daysUntil <= festival.ReminderDays {
			upcomingFestivals = append(upcomingFestivals, festival)
		}
	}
	
	return upcomingFestivals, nil
}

// ExportConfig 导出配置为JSON
func (s *SystemConfigService) ExportConfig(configType string) (string, error) {
	var result map[string]interface{}
	
	switch configType {
	case "festival":
		festivals, err := s.GetFestivalConfigs(false)
		if err != nil {
			return "", err
		}
		result = map[string]interface{}{"festivals": festivals}
		
	case "template":
		templates, err := s.GetTemplateConfigs("", false)
		if err != nil {
			return "", err
		}
		result = map[string]interface{}{"templates": templates}
		
	case "system":
		configs, err := s.GetSystemConfigsByType("")
		if err != nil {
			return "", err
		}
		result = map[string]interface{}{"system_configs": configs}
		
	case "all":
		festivals, _ := s.GetFestivalConfigs(false)
		templates, _ := s.GetTemplateConfigs("", false)
		configs, _ := s.GetSystemConfigsByType("")
		result = map[string]interface{}{
			"festivals":      festivals,
			"templates":      templates,
			"system_configs": configs,
		}
		
	default:
		return "", errors.New("无效的配置类型")
	}
	
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("导出配置失败: %v", err)
	}
	
	return string(jsonData), nil
}
