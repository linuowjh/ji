package services

import (
	"errors"
	"fmt"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivacyService struct {
	db *gorm.DB
}

func NewPrivacyService(db *gorm.DB) *PrivacyService {
	return &PrivacyService{
		db: db,
	}
}

// 隐私级别常量
const (
	PrivacyLevelPublic = 0 // 公开
	PrivacyLevelFamily = 1 // 家族可见
	PrivacyLevelPrivate = 2 // 私密
)

// 访客权限类型
const (
	VisitorPermissionView    = "view"    // 查看
	VisitorPermissionWorship = "worship" // 祭扫
	VisitorPermissionComment = "comment" // 留言
	VisitorPermissionShare   = "share"   // 分享
)

// 隐私设置请求
type PrivacySettingsRequest struct {
	MemorialID          string   `json:"memorial_id" binding:"required"`
	PrivacyLevel        int      `json:"privacy_level" binding:"required,oneof=0 1 2"`
	AllowedFamilyIDs    []string `json:"allowed_family_ids"`
	AllowedUserIDs      []string `json:"allowed_user_ids"`
	VisitorPermissions  []string `json:"visitor_permissions"`
	RequireApproval     bool     `json:"require_approval"`
	BlockedUserIDs      []string `json:"blocked_user_ids"`
	AutoDeleteVisitors  bool     `json:"auto_delete_visitors"`
	VisitorRetentionDays int     `json:"visitor_retention_days"`
}



// 设置纪念馆隐私
func (s *PrivacyService) SetMemorialPrivacy(userID string, req *PrivacySettingsRequest) error {
	// 验证用户是否为纪念馆创建者
	var memorial models.Memorial
	err := s.db.Where("id = ? AND creator_id = ?", req.MemorialID, userID).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆不存在或无权操作")
		}
		return err
	}

	// 更新纪念馆隐私级别
	err = s.db.Model(&memorial).Update("privacy_level", req.PrivacyLevel).Error
	if err != nil {
		return fmt.Errorf("更新隐私级别失败: %v", err)
	}

	// 清除现有的权限设置
	s.db.Where("memorial_id = ?", req.MemorialID).Delete(&models.VisitorPermissionSetting{})

	// 设置家族权限
	for _, familyID := range req.AllowedFamilyIDs {
		for _, permission := range req.VisitorPermissions {
			setting := &models.VisitorPermissionSetting{
				ID:             uuid.New().String(),
				MemorialID:     req.MemorialID,
				FamilyID:       familyID,
				PermissionType: permission,
				IsAllowed:      true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			s.db.Create(setting)
		}
	}

	// 设置用户权限
	for _, userIDAllowed := range req.AllowedUserIDs {
		for _, permission := range req.VisitorPermissions {
			setting := &models.VisitorPermissionSetting{
				ID:             uuid.New().String(),
				MemorialID:     req.MemorialID,
				UserID:         userIDAllowed,
				PermissionType: permission,
				IsAllowed:      true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			s.db.Create(setting)
		}
	}

	// 设置黑名单
	s.db.Where("memorial_id = ?", req.MemorialID).Delete(&models.VisitorBlacklist{})
	for _, blockedUserID := range req.BlockedUserIDs {
		blacklist := &models.VisitorBlacklist{
			ID:         uuid.New().String(),
			MemorialID: req.MemorialID,
			UserID:     blockedUserID,
			Reason:     "用户设置",
			CreatedAt:  time.Now(),
		}
		s.db.Create(blacklist)
	}

	return nil
}

// 获取纪念馆隐私设置
func (s *PrivacyService) GetMemorialPrivacySettings(userID, memorialID string) (map[string]interface{}, error) {
	// 验证用户是否为纪念馆创建者
	var memorial models.Memorial
	err := s.db.Where("id = ? AND creator_id = ?", memorialID, userID).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("纪念馆不存在或无权操作")
		}
		return nil, err
	}

	// 获取权限设置
	var permissions []models.VisitorPermissionSetting
	s.db.Where("memorial_id = ?", memorialID).
		Preload("User").
		Preload("Family").
		Find(&permissions)

	// 获取黑名单
	var blacklist []models.VisitorBlacklist
	s.db.Where("memorial_id = ?", memorialID).
		Preload("User").
		Find(&blacklist)

	// 获取访问申请
	var accessRequests []models.AccessRequest
	s.db.Where("memorial_id = ? AND status = ?", memorialID, "pending").
		Preload("User").
		Find(&accessRequests)

	settings := map[string]interface{}{
		"memorial":        memorial,
		"permissions":     permissions,
		"blacklist":       blacklist,
		"access_requests": accessRequests,
	}

	return settings, nil
}

// 检查用户访问权限
func (s *PrivacyService) CheckUserAccess(userID, memorialID string, permissionType string) (bool, error) {
	// 获取纪念馆信息
	var memorial models.Memorial
	err := s.db.Where("id = ? AND status = ?", memorialID, 1).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("纪念馆不存在")
		}
		return false, err
	}

	// 创建者有所有权限
	if memorial.CreatorID == userID {
		return true, nil
	}

	// 检查是否在黑名单中
	var blacklistCount int64
	s.db.Model(&models.VisitorBlacklist{}).
		Where("memorial_id = ? AND user_id = ?", memorialID, userID).
		Count(&blacklistCount)
	if blacklistCount > 0 {
		return false, errors.New("用户已被拉黑")
	}

	// 根据隐私级别检查权限
	switch memorial.PrivacyLevel {
	case PrivacyLevelPublic:
		// 公开纪念馆，所有人都可以访问
		return true, nil

	case PrivacyLevelFamily:
		// 家族可见，检查是否为家族成员或有特殊权限
		return s.checkFamilyOrSpecialAccess(userID, memorialID, permissionType)

	case PrivacyLevelPrivate:
		// 私密纪念馆，只检查特殊权限
		return s.checkSpecialAccess(userID, memorialID, permissionType)

	default:
		return false, errors.New("无效的隐私级别")
	}
}

// 检查家族或特殊访问权限
func (s *PrivacyService) checkFamilyOrSpecialAccess(userID, memorialID, permissionType string) (bool, error) {
	// 检查是否为家族成员
	var familyCount int64
	s.db.Table("memorial_families mf").
		Joins("JOIN family_members fm ON mf.family_id = fm.family_id").
		Where("mf.memorial_id = ? AND fm.user_id = ?", memorialID, userID).
		Count(&familyCount)

	if familyCount > 0 {
		return true, nil
	}

	// 检查特殊权限
	return s.checkSpecialAccess(userID, memorialID, permissionType)
}

// 检查特殊访问权限
func (s *PrivacyService) checkSpecialAccess(userID, memorialID, permissionType string) (bool, error) {
	var permission models.VisitorPermissionSetting
	err := s.db.Where("memorial_id = ? AND user_id = ? AND permission_type = ? AND is_allowed = ?",
		memorialID, userID, permissionType, true).First(&permission).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// 申请访问权限
func (s *PrivacyService) RequestAccess(userID, memorialID, message string) error {
	// 检查是否已有待处理的申请
	var existingRequest models.AccessRequest
	err := s.db.Where("memorial_id = ? AND user_id = ? AND status = ?",
		memorialID, userID, "pending").First(&existingRequest).Error

	if err == nil {
		return errors.New("已有待处理的访问申请")
	}

	// 创建访问申请
	request := &models.AccessRequest{
		ID:         uuid.New().String(),
		MemorialID: memorialID,
		UserID:     userID,
		Message:    message,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.db.Create(request).Error
}

// 处理访问申请
func (s *PrivacyService) HandleAccessRequest(ownerID, requestID string, approve bool, permissions []string) error {
	// 获取访问申请
	var request models.AccessRequest
	err := s.db.Preload("Memorial").Where("id = ?", requestID).First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("访问申请不存在")
		}
		return err
	}

	// 验证权限
	if request.Memorial.CreatorID != ownerID {
		return errors.New("无权处理此申请")
	}

	// 更新申请状态
	status := "rejected"
	if approve {
		status = "approved"
	}

	err = s.db.Model(&request).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return err
	}

	// 如果批准，添加权限
	if approve {
		for _, permission := range permissions {
			setting := &models.VisitorPermissionSetting{
				ID:             uuid.New().String(),
				MemorialID:     request.MemorialID,
				UserID:         request.UserID,
				PermissionType: permission,
				IsAllowed:      true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			s.db.Create(setting)
		}
	}

	return nil
}

// 添加用户到黑名单
func (s *PrivacyService) AddToBlacklist(ownerID, memorialID, userID, reason string) error {
	// 验证权限
	var memorial models.Memorial
	err := s.db.Where("id = ? AND creator_id = ?", memorialID, ownerID).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆不存在或无权操作")
		}
		return err
	}

	// 检查是否已在黑名单中
	var existingBlacklist models.VisitorBlacklist
	err = s.db.Where("memorial_id = ? AND user_id = ?", memorialID, userID).First(&existingBlacklist).Error
	if err == nil {
		return errors.New("用户已在黑名单中")
	}

	// 添加到黑名单
	blacklist := &models.VisitorBlacklist{
		ID:         uuid.New().String(),
		MemorialID: memorialID,
		UserID:     userID,
		Reason:     reason,
		CreatedAt:  time.Now(),
	}

	err = s.db.Create(blacklist).Error
	if err != nil {
		return err
	}

	// 移除该用户的所有权限
	s.db.Where("memorial_id = ? AND user_id = ?", memorialID, userID).Delete(&models.VisitorPermissionSetting{})

	return nil
}

// 从黑名单移除用户
func (s *PrivacyService) RemoveFromBlacklist(ownerID, memorialID, userID string) error {
	// 验证权限
	var memorial models.Memorial
	err := s.db.Where("id = ? AND creator_id = ?", memorialID, ownerID).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆不存在或无权操作")
		}
		return err
	}

	// 从黑名单移除
	return s.db.Where("memorial_id = ? AND user_id = ?", memorialID, userID).Delete(&models.VisitorBlacklist{}).Error
}

// 获取访问申请列表
func (s *PrivacyService) GetAccessRequests(ownerID, memorialID string, page, pageSize int) ([]models.AccessRequest, int64, error) {
	// 验证权限
	var memorial models.Memorial
	err := s.db.Where("id = ? AND creator_id = ?", memorialID, ownerID).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("纪念馆不存在或无权操作")
		}
		return nil, 0, err
	}

	var requests []models.AccessRequest
	var total int64

	// 计算总数
	s.db.Model(&models.AccessRequest{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err = s.db.Where("memorial_id = ?", memorialID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// 清理过期访客记录
func (s *PrivacyService) CleanupExpiredVisitors(memorialID string, retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return s.db.Where("memorial_id = ? AND visit_time < ?", memorialID, cutoffDate).
		Delete(&models.VisitorRecord{}).Error
}