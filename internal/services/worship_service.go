package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorshipService struct {
	db            *gorm.DB
	familyService *FamilyService
}

func NewWorshipService(db *gorm.DB) *WorshipService {
	return &WorshipService{
		db: db,
	}
}

// SetFamilyService 设置家族服务依赖（避免循环依赖）
func (s *WorshipService) SetFamilyService(familyService *FamilyService) {
	s.familyService = familyService
}

// 献花请求结构
type OfferFlowersRequest struct {
	FlowerType   string `json:"flowerType" binding:"required"`     // 花卉类型：chrysanthemum|carnation|lily|rose
	Quantity     int    `json:"quantity" binding:"required,min=1"` // 数量
	Message      string `json:"message"`                           // 献花留言
	IsScheduled  bool   `json:"isScheduled"`                       // 是否定时送花
	ScheduleTime string `json:"scheduleTime"`                      // 定时时间 (格式: "2006-01-02 15:04:05")
}

// 点烛请求结构
type LightCandleRequest struct {
	CandleType string `json:"candleType" binding:"required"` // 蜡烛类型：red|white|yellow
	Duration   int    `json:"duration" binding:"required"`   // 燃烧时长(分钟)
	Message    string `json:"message"`                       // 点烛留言
}

// 上香请求结构
type OfferIncenseRequest struct {
	IncenseCount int    `json:"incenseCount" binding:"required,oneof=3 9"` // 香柱数量：3或9
	IncenseType  string `json:"incenseType" binding:"required"`            // 香的类型：sandalwood|agarwood|traditional
	Message      string `json:"message"`                                   // 上香留言
}

// 供奉供品请求结构
type OfferTributeRequest struct {
	TributeType string   `json:"tributeType" binding:"required"` // 供品类型：fruit|pastry|wine|tea
	Items       []string `json:"items" binding:"required"`       // 具体供品项目
	Message     string   `json:"message"`                        // 供奉留言
}

// 祈福请求结构
type CreatePrayerRequest struct {
	Content  string `json:"content" binding:"required"` // 祈福内容
	IsPublic bool   `json:"is_public"`                  // 是否公开显示
}

// 留言请求结构
type CreateMessageRequest struct {
	MessageType string `json:"message_type" binding:"required,oneof=text audio video"` // 留言类型
	Content     string `json:"content"`                                                // 文字内容
	MediaURL    string `json:"media_url"`                                              // 音频/视频URL
	Duration    int    `json:"duration"`                                               // 音频/视频时长
}

// 献花内容结构
type FlowerContent struct {
	FlowerType   string    `json:"flower_type"`
	Quantity     int       `json:"quantity"`
	Message      string    `json:"message"`
	IsScheduled  bool      `json:"is_scheduled"`
	ScheduleTime time.Time `json:"schedule_time,omitempty"`
}

// 点烛内容结构
type CandleContent struct {
	CandleType string `json:"candle_type"`
	Duration   int    `json:"duration"`
	Message    string `json:"message"`
	ExpireTime string `json:"expire_time"` // 蜡烛熄灭时间
}

// 上香内容结构
type IncenseContent struct {
	IncenseCount int    `json:"incense_count"`
	IncenseType  string `json:"incense_type"`
	Message      string `json:"message"`
}

// 供品内容结构
type TributeContent struct {
	TributeType string   `json:"tribute_type"`
	Items       []string `json:"items"`
	Message     string   `json:"message"`
}

// 献花
func (s *WorshipService) OfferFlowers(userID, memorialID string, req *OfferFlowersRequest) error {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return err
	}

	// 构建献花内容
	content := FlowerContent{
		FlowerType:  req.FlowerType,
		Quantity:    req.Quantity,
		Message:     req.Message,
		IsScheduled: req.IsScheduled,
	}

	// 处理定时送花
	if req.IsScheduled && req.ScheduleTime != "" {
		scheduleTime, err := time.Parse("2006-01-02 15:04:05", req.ScheduleTime)
		if err != nil {
			return errors.New("定时时间格式错误")
		}
		if scheduleTime.Before(time.Now()) {
			return errors.New("定时时间不能早于当前时间")
		}
		content.ScheduleTime = scheduleTime
	}

	contentJSON, _ := json.Marshal(content)

	// 创建祭扫记录
	record := &models.WorshipRecord{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		WorshipType: "flower",
		Content:     string(contentJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(record).Error; err != nil {
		return err
	}

	// 同步到家族圈
	if s.familyService != nil {
		s.familyService.SyncWorshipActivity(userID, memorialID, "flower", content)
	}

	return nil
}

// 点烛
func (s *WorshipService) LightCandle(userID, memorialID string, req *LightCandleRequest) error {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return err
	}

	// 计算蜡烛熄灭时间
	expireTime := time.Now().Add(time.Duration(req.Duration) * time.Minute)

	// 构建点烛内容
	content := CandleContent{
		CandleType: req.CandleType,
		Duration:   req.Duration,
		Message:    req.Message,
		ExpireTime: expireTime.Format("2006-01-02 15:04:05"),
	}

	contentJSON, _ := json.Marshal(content)

	// 创建祭扫记录
	record := &models.WorshipRecord{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		WorshipType: "candle",
		Content:     string(contentJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.db.Create(record).Error
}

// 上香
func (s *WorshipService) OfferIncense(userID, memorialID string, req *OfferIncenseRequest) error {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return err
	}

	// 构建上香内容
	content := IncenseContent{
		IncenseCount: req.IncenseCount,
		IncenseType:  req.IncenseType,
		Message:      req.Message,
	}

	contentJSON, _ := json.Marshal(content)

	// 创建祭扫记录
	record := &models.WorshipRecord{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		WorshipType: "incense",
		Content:     string(contentJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.db.Create(record).Error
}

// 供奉供品
func (s *WorshipService) OfferTribute(userID, memorialID string, req *OfferTributeRequest) error {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return err
	}

	// 构建供品内容
	content := TributeContent{
		TributeType: req.TributeType,
		Items:       req.Items,
		Message:     req.Message,
	}

	contentJSON, _ := json.Marshal(content)

	// 创建祭扫记录
	record := &models.WorshipRecord{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		WorshipType: "tribute",
		Content:     string(contentJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.db.Create(record).Error
}

// 创建祈福
func (s *WorshipService) CreatePrayer(userID, memorialID string, req *CreatePrayerRequest) (*models.Prayer, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	// 创建祈福记录
	prayer := &models.Prayer{
		ID:         uuid.New().String(),
		MemorialID: memorialID,
		UserID:     userID,
		Content:    req.Content,
		IsPublic:   req.IsPublic,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(prayer).Error; err != nil {
		return nil, err
	}

	// 同时创建祭扫记录
	contentJSON, _ := json.Marshal(map[string]interface{}{
		"content":   req.Content,
		"is_public": req.IsPublic,
	})

	record := &models.WorshipRecord{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		WorshipType: "prayer",
		Content:     string(contentJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.db.Create(record)

	return prayer, nil
}

// 创建留言
func (s *WorshipService) CreateMessage(userID, memorialID string, req *CreateMessageRequest) (*models.Message, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	// 验证留言内容
	if req.MessageType == "text" && req.Content == "" {
		return nil, errors.New("文字留言内容不能为空")
	}
	if (req.MessageType == "audio" || req.MessageType == "video") && req.MediaURL == "" {
		return nil, errors.New("音频/视频留言必须提供媒体文件")
	}

	// 创建留言记录
	message := &models.Message{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		MessageType: req.MessageType,
		Content:     req.Content,
		MediaURL:    req.MediaURL,
		Duration:    req.Duration,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	// 同时创建祭扫记录
	contentJSON, _ := json.Marshal(map[string]interface{}{
		"message_type": req.MessageType,
		"content":      req.Content,
		"media_url":    req.MediaURL,
		"duration":     req.Duration,
	})

	record := &models.WorshipRecord{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		UserID:      userID,
		WorshipType: "message",
		Content:     string(contentJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.db.Create(record)

	return message, nil
}

// GetWorshipRecords 获取祭扫记录
func (s *WorshipService) GetWorshipRecords(userID, memorialID string, page, pageSize int) ([]*models.WorshipRecord, int64, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, 0, err
	}

	var records []*models.WorshipRecord
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 查询记录，包含用户信息
	err := s.db.Preload("User").
		Where("memorial_id = ?", memorialID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

// 获取祈福墙
func (s *WorshipService) GetPrayerWall(userID, memorialID string, page, pageSize int) ([]*models.Prayer, int64, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, 0, err
	}

	var prayers []*models.Prayer
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询总数（只查询公开的祈福）
	s.db.Model(&models.Prayer{}).Where("memorial_id = ? AND is_public = ?", memorialID, true).Count(&total)

	// 查询记录，包含用户信息
	err := s.db.Preload("User").
		Where("memorial_id = ? AND is_public = ?", memorialID, true).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&prayers).Error

	return prayers, total, err
}

// 获取时光信箱
func (s *WorshipService) GetTimeMessages(userID, memorialID string, page, pageSize int) ([]*models.Message, int64, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, 0, err
	}

	var messages []*models.Message
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.Message{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 查询记录，包含用户信息
	err := s.db.Preload("User").
		Where("memorial_id = ?", memorialID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// 续烛功能
func (s *WorshipService) RenewCandle(userID, memorialID string, additionalMinutes int) error {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return err
	}

	// 查找用户最近的点烛记录
	var record models.WorshipRecord
	err := s.db.Where("memorial_id = ? AND user_id = ? AND worship_type = ?",
		memorialID, userID, "candle").
		Order("created_at DESC").
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到可续的蜡烛")
		}
		return err
	}

	// 解析原有内容
	var content CandleContent
	if err := json.Unmarshal([]byte(record.Content), &content); err != nil {
		return errors.New("蜡烛记录数据异常")
	}

	// 检查蜡烛是否已熄灭
	expireTime, err := time.Parse("2006-01-02 15:04:05", content.ExpireTime)
	if err != nil {
		return errors.New("蜡烛时间数据异常")
	}

	if time.Now().After(expireTime) {
		return errors.New("蜡烛已熄灭，无法续烛")
	}

	// 更新蜡烛熄灭时间
	newExpireTime := expireTime.Add(time.Duration(additionalMinutes) * time.Minute)
	content.ExpireTime = newExpireTime.Format("2006-01-02 15:04:05")
	content.Duration += additionalMinutes

	contentJSON, _ := json.Marshal(content)

	// 更新记录
	return s.db.Model(&record).Update("content", string(contentJSON)).Error
}

// 获取当前燃烧的蜡烛状态
func (s *WorshipService) GetActiveCandleStatus(userID, memorialID string) (map[string]interface{}, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	// 查找所有用户的点烛记录
	var records []models.WorshipRecord
	err := s.db.Where("memorial_id = ? AND worship_type = ?", memorialID, "candle").
		Order("created_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	activeCandlesCount := 0
	var activeCandlesByUser []map[string]interface{}

	for _, record := range records {
		var content CandleContent
		if err := json.Unmarshal([]byte(record.Content), &content); err != nil {
			continue
		}

		expireTime, err := time.Parse("2006-01-02 15:04:05", content.ExpireTime)
		if err != nil {
			continue
		}

		// 检查蜡烛是否还在燃烧
		if time.Now().Before(expireTime) {
			activeCandlesCount++

			// 获取用户信息
			var user models.User
			s.db.First(&user, "id = ?", record.UserID)

			activeCandlesByUser = append(activeCandlesByUser, map[string]interface{}{
				"user_id":     record.UserID,
				"user_name":   user.Nickname,
				"candle_type": content.CandleType,
				"expire_time": content.ExpireTime,
				"message":     content.Message,
				"lit_at":      record.CreatedAt,
			})
		}
	}

	return map[string]interface{}{
		"active_candles_count": activeCandlesCount,
		"active_candles":       activeCandlesByUser,
	}, nil
}

// 验证纪念馆访问权限
func (s *WorshipService) validateMemorialAccess(userID, memorialID string) error {
	var memorial models.Memorial
	err := s.db.First(&memorial, "id = ? AND status = ?", memorialID, 1).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆不存在")
		}
		return err
	}

	// 检查隐私设置
	if memorial.PrivacyLevel == 2 { // 私密
		if memorial.CreatorID != userID {
			// 检查是否是家族成员
			var count int64
			s.db.Table("family_members fm").
				Joins("JOIN memorial_families mf ON fm.family_id = mf.family_id").
				Where("mf.memorial_id = ? AND fm.user_id = ?", memorialID, userID).
				Count(&count)

			if count == 0 {
				return errors.New("无权访问此纪念馆")
			}
		}
	}

	return nil
}

// 获取祭扫统计
func (s *WorshipService) GetWorshipStatistics(userID, memorialID string) (map[string]interface{}, error) {
	// 验证纪念馆是否存在且用户有权限访问
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})

	// 统计各类祭扫行为的数量
	worshipTypes := []string{"flower", "candle", "incense", "tribute", "prayer", "message"}

	for _, worshipType := range worshipTypes {
		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND worship_type = ?", memorialID, worshipType).
			Count(&count)
		stats[worshipType+"_count"] = count
	}

	// 统计总访问次数
	var totalVisits int64
	s.db.Model(&models.WorshipRecord{}).
		Where("memorial_id = ?", memorialID).
		Count(&totalVisits)
	stats["total_visits"] = totalVisits

	// 统计独立访客数
	var uniqueVisitors int64
	s.db.Model(&models.WorshipRecord{}).
		Where("memorial_id = ?", memorialID).
		Distinct("user_id").
		Count(&uniqueVisitors)
	stats["unique_visitors"] = uniqueVisitors

	// 统计最近7天的访问情况
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var recentVisits int64
	s.db.Model(&models.WorshipRecord{}).
		Where("memorial_id = ? AND created_at >= ?", memorialID, sevenDaysAgo).
		Count(&recentVisits)
	stats["recent_visits"] = recentVisits

	return stats, nil
}

// PrayerCardTemplate 祈福卡样式模板
type PrayerCardTemplate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"` // traditional|modern|festival
}

// 获取祈福卡模板
func (s *WorshipService) GetPrayerCardTemplates() []PrayerCardTemplate {
	return []PrayerCardTemplate{
		{
			ID:          "template-1",
			Name:        "传统祈福卡",
			Description: "古典雅致的传统祈福卡样式",
			ImageURL:    "/static/templates/traditional-prayer-card.jpg",
			Category:    "traditional",
		},
		{
			ID:          "template-2",
			Name:        "莲花祈福卡",
			Description: "清净庄严的莲花主题祈福卡",
			ImageURL:    "/static/templates/lotus-prayer-card.jpg",
			Category:    "traditional",
		},
		{
			ID:          "template-3",
			Name:        "现代简约卡",
			Description: "简洁现代的祈福卡设计",
			ImageURL:    "/static/templates/modern-prayer-card.jpg",
			Category:    "modern",
		},
		{
			ID:          "template-4",
			Name:        "节日祈福卡",
			Description: "适合特殊节日的祈福卡",
			ImageURL:    "/static/templates/festival-prayer-card.jpg",
			Category:    "festival",
		},
	}
}

// 生成祈福卡图片
func (s *WorshipService) GeneratePrayerCard(templateID, content, userName string) (string, error) {
	// 这里应该调用图片生成服务，将祈福内容渲染到模板上
	// 目前返回模拟的生成结果
	cardURL := fmt.Sprintf("/generated/prayer-cards/%s-%d.jpg", templateID, time.Now().Unix())

	// 实际实现中，这里会：
	// 1. 根据templateID获取模板
	// 2. 将content和userName渲染到模板上
	// 3. 生成图片并上传到对象存储
	// 4. 返回生成的图片URL

	return cardURL, nil
}

// 留言审核状态
type MessageModerationStatus struct {
	IsApproved bool   `json:"is_approved"`
	Reason     string `json:"reason,omitempty"`
}

// 审核留言内容
func (s *WorshipService) ModerateMessage(messageID string) (*MessageModerationStatus, error) {
	var message models.Message
	err := s.db.First(&message, "id = ?", messageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("留言不存在")
		}
		return nil, err
	}

	// 这里应该调用内容审核服务
	// 目前返回模拟的审核结果
	status := &MessageModerationStatus{
		IsApproved: true,
		Reason:     "",
	}

	// 实际实现中，这里会：
	// 1. 调用腾讯云内容安全API审核文本/音频/视频
	// 2. 检查敏感词
	// 3. 更新留言的审核状态
	// 4. 如果不通过，记录原因

	return status, nil
}

// 获取用户的祭扫历史统计
func (s *WorshipService) GetUserWorshipHistory(userID string, page, pageSize int) (map[string]interface{}, error) {
	offset := (page - 1) * pageSize

	// 获取用户的祭扫记录
	var records []models.WorshipRecord
	err := s.db.Preload("Memorial").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	// 统计各类祭扫行为
	var stats map[string]int64 = make(map[string]int64)
	worshipTypes := []string{"flower", "candle", "incense", "tribute", "prayer", "message"}

	for _, worshipType := range worshipTypes {
		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("user_id = ? AND worship_type = ?", userID, worshipType).
			Count(&count)
		stats[worshipType+"_count"] = count
	}

	// 统计总记录数
	var totalCount int64
	s.db.Model(&models.WorshipRecord{}).Where("user_id = ?", userID).Count(&totalCount)

	// 统计参与的纪念馆数量
	var memorialCount int64
	s.db.Model(&models.WorshipRecord{}).
		Where("user_id = ?", userID).
		Distinct("memorial_id").
		Count(&memorialCount)

	return map[string]interface{}{
		"records":        records,
		"total_records":  totalCount,
		"memorial_count": memorialCount,
		"statistics":     stats,
		"page":           page,
		"page_size":      pageSize,
	}, nil
}

// 创建定时祈福提醒
type ScheduledPrayerRequest struct {
	MemorialID    string    `json:"memorial_id" binding:"required"`
	Content       string    `json:"content" binding:"required"`
	ScheduleTime  time.Time `json:"schedule_time" binding:"required"`
	IsRecurring   bool      `json:"is_recurring"`
	RecurringType string    `json:"recurring_type"` // daily|weekly|monthly|yearly
}

func (s *WorshipService) CreateScheduledPrayer(userID string, req *ScheduledPrayerRequest) error {
	// 验证纪念馆访问权限
	if err := s.validateMemorialAccess(userID, req.MemorialID); err != nil {
		return err
	}

	// 验证定时时间
	if req.ScheduleTime.Before(time.Now()) {
		return errors.New("定时时间不能早于当前时间")
	}

	// 创建定时祈福记录
	reminder := &models.MemorialReminder{
		ID:           uuid.New().String(),
		MemorialID:   req.MemorialID,
		ReminderType: "scheduled_prayer",
		ReminderDate: req.ScheduleTime,
		Title:        "定时祈福提醒",
		Content:      req.Content,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.db.Create(reminder).Error
}

// 获取热门祈福内容（用于推荐）
func (s *WorshipService) GetPopularPrayerContents(limit int) ([]string, error) {
	// 这里应该基于祈福内容的点赞数、使用频率等来推荐
	// 目前返回一些预设的热门祈福内容
	popularContents := []string{
		"愿您在天堂安好，保佑家人平安健康",
		"思念如潮水，愿您安息",
		"您的音容笑貌永远在我心中",
		"愿您在另一个世界快乐无忧",
		"感谢您给我们的爱，永远怀念您",
		"愿您的灵魂得到安息，家人得到慰藉",
		"您的教诲我们永远铭记在心",
		"愿天堂没有病痛，只有快乐",
		"您永远是我们心中最亮的星",
		"愿您在天堂与亲人团聚",
	}

	if limit > len(popularContents) {
		limit = len(popularContents)
	}

	return popularContents[:limit], nil
}

// EmotionAnalysisResult 情感分析结果
type EmotionAnalysisResult struct {
	Emotion    string   `json:"emotion"`    // happy|sad|nostalgic|grateful|peaceful
	Confidence float64  `json:"confidence"` // 置信度 0-1
	Keywords   []string `json:"keywords"`   // 关键词
	Suggestion string   `json:"suggestion"` // 回复建议
}

// 分析留言情感
func (s *WorshipService) AnalyzeMessageEmotion(content string) (*EmotionAnalysisResult, error) {
	// 这里应该调用情感分析API，目前返回模拟结果
	// 实际实现中可以使用腾讯云文本分析API或其他NLP服务

	// 简单的关键词匹配来模拟情感分析
	sadKeywords := []string{"想念", "思念", "难过", "伤心", "离别", "痛苦"}
	happyKeywords := []string{"快乐", "开心", "幸福", "美好", "温暖", "笑容"}
	nostalgicKeywords := []string{"回忆", "往昔", "从前", "过去", "曾经", "那时"}
	gratefulKeywords := []string{"感谢", "感恩", "谢谢", "感激", "恩情", "教诲"}
	peacefulKeywords := []string{"安息", "安好", "平静", "宁静", "安详", "祝福"}

	content = strings.ToLower(content)

	// 统计各种情感关键词出现次数
	emotionScores := map[string]int{
		"sad":       0,
		"happy":     0,
		"nostalgic": 0,
		"grateful":  0,
		"peaceful":  0,
	}

	var foundKeywords []string

	for _, keyword := range sadKeywords {
		if strings.Contains(content, keyword) {
			emotionScores["sad"]++
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	for _, keyword := range happyKeywords {
		if strings.Contains(content, keyword) {
			emotionScores["happy"]++
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	for _, keyword := range nostalgicKeywords {
		if strings.Contains(content, keyword) {
			emotionScores["nostalgic"]++
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	for _, keyword := range gratefulKeywords {
		if strings.Contains(content, keyword) {
			emotionScores["grateful"]++
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	for _, keyword := range peacefulKeywords {
		if strings.Contains(content, keyword) {
			emotionScores["peaceful"]++
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	// 找出得分最高的情感
	maxScore := 0
	dominantEmotion := "peaceful" // 默认为平静

	for emotion, score := range emotionScores {
		if score > maxScore {
			maxScore = score
			dominantEmotion = emotion
		}
	}

	// 计算置信度
	totalKeywords := len(foundKeywords)
	confidence := 0.5 // 默认置信度
	if totalKeywords > 0 {
		confidence = float64(maxScore) / float64(totalKeywords)
		if confidence > 1.0 {
			confidence = 1.0
		}
	}

	// 生成回复建议
	suggestions := map[string]string{
		"sad":       "您的思念之情让人动容，相信逝者能感受到您深深的爱意。时间会慢慢抚平伤痛，但美好的回忆会永远陪伴着您。",
		"happy":     "感谢您分享这些美好的回忆，逝者一定也希望看到您如此快乐。让我们一起珍藏这些温暖的时光。",
		"nostalgic": "往昔的美好时光值得永远珍藏，这些回忆是您与逝者之间最珍贵的纽带。",
		"grateful":  "您的感恩之心令人敬佩，逝者的教诲和恩情将永远指引着您前行的道路。",
		"peaceful":  "愿逝者安息，愿您内心平静。在这个特殊的空间里，让爱与思念得到最好的表达。",
	}

	return &EmotionAnalysisResult{
		Emotion:    dominantEmotion,
		Confidence: confidence,
		Keywords:   foundKeywords,
		Suggestion: suggestions[dominantEmotion],
	}, nil
}

// 获取留言回复建议
func (s *WorshipService) GetMessageReplySuggestions(messageType, content string) ([]string, error) {
	suggestions := []string{}

	switch messageType {
	case "text":
		// 基于内容分析生成建议
		emotion, err := s.AnalyzeMessageEmotion(content)
		if err == nil {
			suggestions = append(suggestions, emotion.Suggestion)
		}

		// 添加通用回复建议
		suggestions = append(suggestions, []string{
			"您的话语充满了爱与思念",
			"相信逝者能感受到您的真挚情感",
			"这份深情让人动容",
			"愿这份爱能给您带来慰藉",
		}...)

	case "audio":
		suggestions = []string{
			"您的声音传达了深深的思念",
			"声音是最温暖的陪伴",
			"相信逝者能听到您的呼唤",
			"这份用心的表达很珍贵",
		}

	case "video":
		suggestions = []string{
			"影像记录了最真挚的情感",
			"这是最珍贵的纪念方式",
			"您的用心让人感动",
			"愿这份美好永远保存",
		}
	}

	return suggestions, nil
}

// 创建留言时的智能提醒
type MessageCreationTip struct {
	Type    string `json:"type"`    // suggestion|reminder|encouragement
	Content string `json:"content"` // 提醒内容
}

func (s *WorshipService) GetMessageCreationTips(memorialID, messageType string) ([]MessageCreationTip, error) {
	tips := []MessageCreationTip{}

	// 基于留言类型给出建议
	switch messageType {
	case "text":
		tips = append(tips, []MessageCreationTip{
			{
				Type:    "suggestion",
				Content: "可以分享一些与逝者的美好回忆，或者表达您的思念之情",
			},
			{
				Type:    "reminder",
				Content: "文字会永久保存，成为珍贵的纪念",
			},
		}...)

	case "audio":
		tips = append(tips, []MessageCreationTip{
			{
				Type:    "suggestion",
				Content: "可以说说您想对逝者说的话，或者分享近况",
			},
			{
				Type:    "reminder",
				Content: "建议在安静的环境中录制，时长控制在3分钟内",
			},
		}...)

	case "video":
		tips = append(tips, []MessageCreationTip{
			{
				Type:    "suggestion",
				Content: "可以录制一段对逝者的话，或者展示相关的纪念物品",
			},
			{
				Type:    "reminder",
				Content: "请注意光线充足，建议时长控制在5分钟内",
			},
		}...)
	}

	// 添加通用鼓励
	tips = append(tips, MessageCreationTip{
		Type:    "encouragement",
		Content: "每一份真挚的表达都是对逝者最好的纪念",
	})

	return tips, nil
}

// 获取纪念馆的留言统计分析
func (s *WorshipService) GetMemorialMessageAnalytics(memorialID string) (map[string]interface{}, error) {
	// 统计各类型留言数量
	var textCount, audioCount, videoCount int64

	s.db.Model(&models.Message{}).Where("memorial_id = ? AND message_type = ?", memorialID, "text").Count(&textCount)
	s.db.Model(&models.Message{}).Where("memorial_id = ? AND message_type = ?", memorialID, "audio").Count(&audioCount)
	s.db.Model(&models.Message{}).Where("memorial_id = ? AND message_type = ?", memorialID, "video").Count(&videoCount)

	// 统计最近30天的留言趋势
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var recentMessages []models.Message
	s.db.Where("memorial_id = ? AND created_at >= ?", memorialID, thirtyDaysAgo).
		Order("created_at ASC").
		Find(&recentMessages)

	// 按天统计留言数量
	dailyStats := make(map[string]int)
	for _, message := range recentMessages {
		day := message.CreatedAt.Format("2006-01-02")
		dailyStats[day]++
	}

	// 统计活跃用户
	var activeUsers int64
	s.db.Model(&models.Message{}).
		Where("memorial_id = ? AND created_at >= ?", memorialID, thirtyDaysAgo).
		Distinct("user_id").
		Count(&activeUsers)

	return map[string]interface{}{
		"message_types": map[string]int64{
			"text":  textCount,
			"audio": audioCount,
			"video": videoCount,
		},
		"daily_stats":    dailyStats,
		"active_users":   activeUsers,
		"total_messages": textCount + audioCount + videoCount,
	}, nil
}

// WorshipRecordStats 祭扫记录详细统计
type WorshipRecordStats struct {
	TotalRecords   int64                 `json:"total_records"`
	UniqueVisitors int64                 `json:"unique_visitors"`
	TypeStatistics map[string]int64      `json:"type_statistics"`
	MonthlyTrend   []MonthlyStats        `json:"monthly_trend"`
	HourlyPattern  []HourlyStats         `json:"hourly_pattern"`
	TopVisitors    []VisitorStats        `json:"top_visitors"`
	RecentActivity []RecentActivityStats `json:"recent_activity"`
}

type MonthlyStats struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type HourlyStats struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

type VisitorStats struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	AvatarURL string `json:"avatar_url"`
	Count     int64  `json:"count"`
	LastVisit string `json:"last_visit"`
}

type RecentActivityStats struct {
	Date         string `json:"date"`
	WorshipCount int64  `json:"worship_count"`
	VisitorCount int64  `json:"visitor_count"`
}

// 获取详细的祭扫统计
func (s *WorshipService) GetDetailedWorshipStatistics(userID, memorialID string) (*WorshipRecordStats, error) {
	// 验证纪念馆访问权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	stats := &WorshipRecordStats{}

	// 总记录数
	s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorialID).Count(&stats.TotalRecords)

	// 独立访客数
	s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorialID).Distinct("user_id").Count(&stats.UniqueVisitors)

	// 各类型统计
	stats.TypeStatistics = make(map[string]int64)
	worshipTypes := []string{"flower", "candle", "incense", "tribute", "prayer", "message"}
	for _, worshipType := range worshipTypes {
		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND worship_type = ?", memorialID, worshipType).
			Count(&count)
		stats.TypeStatistics[worshipType] = count
	}

	// 月度趋势（最近12个月）
	stats.MonthlyTrend = []MonthlyStats{}
	for i := 11; i >= 0; i-- {
		month := time.Now().AddDate(0, -i, 0)
		monthStr := month.Format("2006-01")

		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND DATE_FORMAT(created_at, '%Y-%m') = ?", memorialID, monthStr).
			Count(&count)

		stats.MonthlyTrend = append(stats.MonthlyTrend, MonthlyStats{
			Month: monthStr,
			Count: count,
		})
	}

	// 小时分布模式
	stats.HourlyPattern = []HourlyStats{}
	for hour := 0; hour < 24; hour++ {
		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND HOUR(created_at) = ?", memorialID, hour).
			Count(&count)

		stats.HourlyPattern = append(stats.HourlyPattern, HourlyStats{
			Hour:  hour,
			Count: count,
		})
	}

	// 访客排行榜（前10名）
	type visitorCount struct {
		UserID string
		Count  int64
	}

	var visitorCounts []visitorCount
	s.db.Model(&models.WorshipRecord{}).
		Select("user_id, COUNT(*) as count").
		Where("memorial_id = ?", memorialID).
		Group("user_id").
		Order("count DESC").
		Limit(10).
		Scan(&visitorCounts)

	stats.TopVisitors = []VisitorStats{}
	for _, vc := range visitorCounts {
		var user models.User
		s.db.First(&user, "id = ?", vc.UserID)

		var lastRecord models.WorshipRecord
		s.db.Where("memorial_id = ? AND user_id = ?", memorialID, vc.UserID).
			Order("created_at DESC").
			First(&lastRecord)

		stats.TopVisitors = append(stats.TopVisitors, VisitorStats{
			UserID:    user.ID,
			UserName:  user.Nickname,
			AvatarURL: user.AvatarURL,
			Count:     vc.Count,
			LastVisit: lastRecord.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 最近30天活动统计
	stats.RecentActivity = []RecentActivityStats{}
	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		var worshipCount int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND DATE(created_at) = ?", memorialID, dateStr).
			Count(&worshipCount)

		var visitorCount int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND DATE(created_at) = ?", memorialID, dateStr).
			Distinct("user_id").
			Count(&visitorCount)

		stats.RecentActivity = append(stats.RecentActivity, RecentActivityStats{
			Date:         dateStr,
			WorshipCount: worshipCount,
			VisitorCount: visitorCount,
		})
	}

	return stats, nil
}

// 用户祭扫行为分析
type UserWorshipBehavior struct {
	UserID           string           `json:"user_id"`
	TotalWorships    int64            `json:"total_worships"`
	FavoriteType     string           `json:"favorite_type"`
	ActiveHours      []int            `json:"active_hours"`
	WorshipFrequency map[string]int64 `json:"worship_frequency"` // 按类型统计
	MemorialCount    int64            `json:"memorial_count"`    // 参与的纪念馆数量
	FirstWorship     string           `json:"first_worship"`
	LastWorship      string           `json:"last_worship"`
}

// 分析用户祭扫行为
func (s *WorshipService) AnalyzeUserWorshipBehavior(userID string) (*UserWorshipBehavior, error) {
	behavior := &UserWorshipBehavior{
		UserID: userID,
	}

	// 总祭扫次数
	s.db.Model(&models.WorshipRecord{}).Where("user_id = ?", userID).Count(&behavior.TotalWorships)

	// 参与的纪念馆数量
	s.db.Model(&models.WorshipRecord{}).Where("user_id = ?", userID).Distinct("memorial_id").Count(&behavior.MemorialCount)

	// 各类型祭扫频率
	behavior.WorshipFrequency = make(map[string]int64)
	worshipTypes := []string{"flower", "candle", "incense", "tribute", "prayer", "message"}
	maxCount := int64(0)

	for _, worshipType := range worshipTypes {
		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("user_id = ? AND worship_type = ?", userID, worshipType).
			Count(&count)
		behavior.WorshipFrequency[worshipType] = count

		if count > maxCount {
			maxCount = count
			behavior.FavoriteType = worshipType
		}
	}

	// 活跃时段分析
	hourCounts := make(map[int]int64)
	var records []models.WorshipRecord
	s.db.Where("user_id = ?", userID).Find(&records)

	for _, record := range records {
		hour := record.CreatedAt.Hour()
		hourCounts[hour]++
	}

	// 找出最活跃的时段（前3个）
	type hourCount struct {
		hour  int
		count int64
	}

	var sortedHours []hourCount
	for hour, count := range hourCounts {
		sortedHours = append(sortedHours, hourCount{hour, count})
	}

	// 简单排序（实际应该使用sort包）
	for i := 0; i < len(sortedHours)-1; i++ {
		for j := i + 1; j < len(sortedHours); j++ {
			if sortedHours[i].count < sortedHours[j].count {
				sortedHours[i], sortedHours[j] = sortedHours[j], sortedHours[i]
			}
		}
	}

	behavior.ActiveHours = []int{}
	for i := 0; i < len(sortedHours) && i < 3; i++ {
		behavior.ActiveHours = append(behavior.ActiveHours, sortedHours[i].hour)
	}

	// 首次和最近祭扫时间
	var firstRecord, lastRecord models.WorshipRecord
	s.db.Where("user_id = ?", userID).Order("created_at ASC").First(&firstRecord)
	s.db.Where("user_id = ?", userID).Order("created_at DESC").First(&lastRecord)

	if firstRecord.ID != "" {
		behavior.FirstWorship = firstRecord.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if lastRecord.ID != "" {
		behavior.LastWorship = lastRecord.CreatedAt.Format("2006-01-02 15:04:05")
	}

	return behavior, nil
}

// 生成祭扫报告
type WorshipReport struct {
	MemorialID      string                 `json:"memorial_id"`
	MemorialName    string                 `json:"memorial_name"`
	ReportPeriod    string                 `json:"report_period"`
	Summary         map[string]interface{} `json:"summary"`
	Highlights      []string               `json:"highlights"`
	Recommendations []string               `json:"recommendations"`
	GeneratedAt     string                 `json:"generated_at"`
}

// 生成祭扫报告
func (s *WorshipService) GenerateWorshipReport(userID, memorialID string, period string) (*WorshipReport, error) {
	// 验证纪念馆访问权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	// 获取纪念馆信息
	var memorial models.Memorial
	err := s.db.First(&memorial, "id = ?", memorialID).Error
	if err != nil {
		return nil, err
	}

	// 根据周期计算时间范围
	var startTime time.Time
	switch period {
	case "week":
		startTime = time.Now().AddDate(0, 0, -7)
	case "month":
		startTime = time.Now().AddDate(0, -1, 0)
	case "quarter":
		startTime = time.Now().AddDate(0, -3, 0)
	case "year":
		startTime = time.Now().AddDate(-1, 0, 0)
	default:
		startTime = time.Now().AddDate(0, -1, 0) // 默认一个月
		period = "month"
	}

	report := &WorshipReport{
		MemorialID:   memorialID,
		MemorialName: memorial.DeceasedName,
		ReportPeriod: period,
		GeneratedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 统计摘要
	var totalRecords, uniqueVisitors int64
	s.db.Model(&models.WorshipRecord{}).
		Where("memorial_id = ? AND created_at >= ?", memorialID, startTime).
		Count(&totalRecords)

	s.db.Model(&models.WorshipRecord{}).
		Where("memorial_id = ? AND created_at >= ?", memorialID, startTime).
		Distinct("user_id").
		Count(&uniqueVisitors)

	// 各类型统计
	typeStats := make(map[string]int64)
	worshipTypes := []string{"flower", "candle", "incense", "tribute", "prayer", "message"}
	for _, worshipType := range worshipTypes {
		var count int64
		s.db.Model(&models.WorshipRecord{}).
			Where("memorial_id = ? AND worship_type = ? AND created_at >= ?", memorialID, worshipType, startTime).
			Count(&count)
		typeStats[worshipType] = count
	}

	report.Summary = map[string]interface{}{
		"total_records":   totalRecords,
		"unique_visitors": uniqueVisitors,
		"type_statistics": typeStats,
		"period_start":    startTime.Format("2006-01-02"),
		"period_end":      time.Now().Format("2006-01-02"),
	}

	// 生成亮点
	report.Highlights = []string{}
	if totalRecords > 0 {
		report.Highlights = append(report.Highlights, fmt.Sprintf("本%s共收到%d次祭扫", getPeriodName(period), totalRecords))
		report.Highlights = append(report.Highlights, fmt.Sprintf("有%d位访客表达了思念", uniqueVisitors))

		// 找出最受欢迎的祭扫方式
		maxType := ""
		maxCount := int64(0)
		for worshipType, count := range typeStats {
			if count > maxCount {
				maxCount = count
				maxType = worshipType
			}
		}
		if maxType != "" {
			report.Highlights = append(report.Highlights, fmt.Sprintf("最受欢迎的祭扫方式是%s", getWorshipTypeName(maxType)))
		}
	}

	// 生成建议
	report.Recommendations = []string{}
	if totalRecords == 0 {
		report.Recommendations = append(report.Recommendations, "可以邀请更多亲友参与纪念活动")
	} else {
		report.Recommendations = append(report.Recommendations, "感谢大家的参与，让我们继续传承这份美好的纪念")
	}

	return report, nil
}

// 辅助函数
func getPeriodName(period string) string {
	switch period {
	case "week":
		return "周"
	case "month":
		return "月"
	case "quarter":
		return "季度"
	case "year":
		return "年"
	default:
		return "期间"
	}
}

func getWorshipTypeName(worshipType string) string {
	switch worshipType {
	case "flower":
		return "献花"
	case "candle":
		return "点烛"
	case "incense":
		return "上香"
	case "tribute":
		return "供奉"
	case "prayer":
		return "祈福"
	case "message":
		return "留言"
	default:
		return "祭扫"
	}
}
