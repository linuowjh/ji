package utils

import (
	"fmt"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// PermissionManager 权限管理器
type PermissionManager struct {
	db *gorm.DB
}

// NewPermissionManager 创建权限管理器
func NewPermissionManager(db *gorm.DB) *PermissionManager {
	return &PermissionManager{db: db}
}

// CanAccessMemorial 检查用户是否可以访问纪念馆
func (pm *PermissionManager) CanAccessMemorial(userID, memorialID string) (bool, string, error) {
	var memorial models.Memorial
	if err := pm.db.Where("id = ? AND status = ?", memorialID, 1).First(&memorial).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, "", fmt.Errorf("纪念馆不存在")
		}
		return false, "", fmt.Errorf("查询纪念馆失败: %v", err)
	}

	// 1. 创建者有完全访问权限
	if memorial.CreatorID == userID {
		return true, "owner", nil
	}

	// 2. 私密纪念馆只有创建者可以访问
	if memorial.PrivacyLevel == 2 {
		return false, "", fmt.Errorf("无权访问私密纪念馆")
	}

	// 3. 家族可见的纪念馆，检查是否为家族成员
	if memorial.PrivacyLevel == 1 {
		var count int64
		pm.db.Table("memorial_families mf").
			Joins("JOIN family_members fm ON mf.family_id = fm.family_id").
			Where("mf.memorial_id = ? AND fm.user_id = ?", memorialID, userID).
			Count(&count)

		if count > 0 {
			return true, "family", nil
		}
	}

	return false, "", fmt.Errorf("无权访问此纪念馆")
}

// CanModifyMemorial 检查用户是否可以修改纪念馆
func (pm *PermissionManager) CanModifyMemorial(userID, memorialID string) (bool, error) {
	var memorial models.Memorial
	if err := pm.db.Where("id = ? AND status = ?", memorialID, 1).First(&memorial).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("纪念馆不存在")
		}
		return false, fmt.Errorf("查询纪念馆失败: %v", err)
	}

	// 只有创建者可以修改纪念馆
	if memorial.CreatorID != userID {
		return false, fmt.Errorf("只有创建者可以修改纪念馆")
	}

	return true, nil
}

// CanAccessFamily 检查用户是否可以访问家族
func (pm *PermissionManager) CanAccessFamily(userID, familyID string) (bool, string, error) {
	var member models.FamilyMember
	err := pm.db.Where("family_id = ? AND user_id = ?", familyID, userID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, "", fmt.Errorf("您不是该家族成员")
		}
		return false, "", fmt.Errorf("查询家族成员失败: %v", err)
	}

	return true, member.Role, nil
}

// CanManageFamily 检查用户是否可以管理家族
func (pm *PermissionManager) CanManageFamily(userID, familyID string) (bool, error) {
	canAccess, role, err := pm.CanAccessFamily(userID, familyID)
	if err != nil {
		return false, err
	}

	if !canAccess || role != "admin" {
		return false, fmt.Errorf("需要管理员权限")
	}

	return true, nil
}

// IsMemorialOwner 检查用户是否为纪念馆创建者
func (pm *PermissionManager) IsMemorialOwner(userID, memorialID string) (bool, error) {
	var memorial models.Memorial
	if err := pm.db.Where("id = ? AND creator_id = ?", memorialID, userID).First(&memorial).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("查询纪念馆失败: %v", err)
	}
	return true, nil
}

// IsFamilyAdmin 检查用户是否为家族管理员
func (pm *PermissionManager) IsFamilyAdmin(userID, familyID string) (bool, error) {
	var member models.FamilyMember
	err := pm.db.Where("family_id = ? AND user_id = ? AND role = ?", familyID, userID, "admin").First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("查询家族成员失败: %v", err)
	}
	return true, nil
}

// GetUserFamilies 获取用户所属的家族列表
func (pm *PermissionManager) GetUserFamilies(userID string) ([]models.Family, error) {
	var families []models.Family
	err := pm.db.Table("families f").
		Joins("JOIN family_members fm ON f.id = fm.family_id").
		Where("fm.user_id = ?", userID).
		Find(&families).Error

	if err != nil {
		return nil, fmt.Errorf("查询用户家族失败: %v", err)
	}

	return families, nil
}

// GetUserAccessibleMemorials 获取用户可访问的纪念馆列表
func (pm *PermissionManager) GetUserAccessibleMemorials(userID string, page, pageSize int) ([]models.Memorial, int64, error) {
	var memorials []models.Memorial
	var total int64

	// 构建查询条件：用户创建的纪念馆 + 用户所属家族的纪念馆
	query := pm.db.Model(&models.Memorial{}).Where("status = ?", 1)

	// 子查询：用户所属家族的纪念馆
	familyMemorialsQuery := pm.db.Table("memorial_families mf").
		Select("mf.memorial_id").
		Joins("JOIN family_members fm ON mf.family_id = fm.family_id").
		Where("fm.user_id = ?", userID)

	// 最终查询条件
	query = query.Where("creator_id = ? OR id IN (?)", userID, familyMemorialsQuery)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&memorials).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询可访问纪念馆失败: %v", err)
	}

	return memorials, total, nil
}

// RecordVisit 记录访客访问
func (pm *PermissionManager) RecordVisit(memorialID, visitorID, ipAddress string) error {
	// 检查是否已有今日访问记录
	var count int64
	pm.db.Model(&models.VisitorRecord{}).
		Where("memorial_id = ? AND visitor_id = ? AND DATE(visit_time) = CURDATE()", memorialID, visitorID).
		Count(&count)

	// 如果今日已访问，则不重复记录
	if count > 0 {
		return nil
	}

	// 创建访问记录
	record := models.VisitorRecord{
		ID:         GenerateUUID(),
		MemorialID: memorialID,
		VisitorID:  visitorID,
		IPAddress:  ipAddress,
	}

	if err := pm.db.Create(&record).Error; err != nil {
		return fmt.Errorf("记录访问失败: %v", err)
	}

	return nil
}

// GetMemorialVisitors 获取纪念馆访客记录
func (pm *PermissionManager) GetMemorialVisitors(memorialID string, page, pageSize int) ([]models.VisitorRecord, int64, error) {
	var records []models.VisitorRecord
	var total int64

	// 计算总数
	pm.db.Model(&models.VisitorRecord{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 分页查询，预加载访客信息
	offset := (page - 1) * pageSize
	err := pm.db.Where("memorial_id = ?", memorialID).
		Preload("Visitor").
		Order("visit_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询访客记录失败: %v", err)
	}

	return records, total, nil
}