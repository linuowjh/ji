package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FamilyService struct {
	db *gorm.DB
}

func NewFamilyService(db *gorm.DB) *FamilyService {
	return &FamilyService{
		db: db,
	}
}

// 创建家族圈请求
type CreateFamilyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// 更新家族圈请求
type UpdateFamilyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// 邀请成员请求
type InviteFamilyMemberRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`
	Message string   `json:"message"`
}

// 创建家族圈
func (s *FamilyService) CreateFamily(userID string, req *CreateFamilyRequest) (*models.Family, error) {
	// 生成邀请码
	inviteCode := s.generateInviteCode()

	family := &models.Family{
		ID:          uuid.New().String(),
		Name:        req.Name,
		CreatorID:   userID,
		Description: req.Description,
		InviteCode:  inviteCode,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tx := s.db.Begin()

	if err := tx.Create(family).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加创建者为管理员
	member := &models.FamilyMember{
		ID:       uuid.New().String(),
		FamilyID: family.ID,
		UserID:   userID,
		Role:     "admin",
		JoinedAt: time.Now(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新查询包含关联数据的家族
	s.db.Preload("Creator").Preload("Members.User").First(family, "id = ?", family.ID)

	return family, nil
}

// 获取家族圈列表
func (s *FamilyService) GetFamilies(userID string, page, pageSize int) ([]*models.Family, int64, error) {
	var families []*models.Family
	var total int64

	offset := (page - 1) * pageSize

	// 查询用户参与的家族圈
	subQuery := s.db.Model(&models.FamilyMember{}).
		Select("family_id").
		Where("user_id = ?", userID)

	// 查询总数
	s.db.Model(&models.Family{}).
		Where("id IN (?)", subQuery).
		Count(&total)

	// 查询家族列表
	err := s.db.Preload("Creator").
		Preload("Members.User").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&families).Error

	return families, total, err
}

// 获取家族圈详情
func (s *FamilyService) GetFamily(userID, familyID string) (*models.Family, error) {
	var family models.Family
	err := s.db.Preload("Creator").
		Preload("Members.User").
		First(&family, "id = ?", familyID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("家族圈不存在")
		}
		return nil, err
	}

	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, errors.New("您不是此家族圈的成员")
	}

	return &family, nil
}

// 更新家族圈
func (s *FamilyService) UpdateFamily(userID, familyID string, req *UpdateFamilyRequest) error {
	var family models.Family
	err := s.db.First(&family, "id = ?", familyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("家族圈不存在")
		}
		return err
	}

	// 验证权限（只有管理员可以修改）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以修改家族圈")
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&family).Updates(updates).Error
}

// 删除家族圈
func (s *FamilyService) DeleteFamily(userID, familyID string) error {
	var family models.Family
	err := s.db.First(&family, "id = ?", familyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("家族圈不存在")
		}
		return err
	}

	// 验证权限（只有创建者可以删除）
	if family.CreatorID != userID {
		return errors.New("只有创建者可以删除家族圈")
	}

	// 删除家族圈及其相关数据
	tx := s.db.Begin()
	
	// 删除家族成员
	if err := tx.Where("family_id = ?", familyID).Delete(&models.FamilyMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除家族邀请
	if err := tx.Where("family_id = ?", familyID).Delete(&models.FamilyInvitation{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除家族活动
	if err := tx.Where("family_id = ?", familyID).Delete(&models.FamilyActivity{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除纪念馆关联
	if err := tx.Where("family_id = ?", familyID).Delete(&models.MemorialFamily{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除家族圈
	if err := tx.Delete(&family).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 邀请成员
func (s *FamilyService) InviteMembers(userID, familyID string, req *InviteFamilyMemberRequest) error {
	// 验证权限（管理员可以邀请）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以邀请成员")
	}

	tx := s.db.Begin()

	for _, inviteeID := range req.UserIDs {
		// 检查是否已经是成员
		if s.isFamilyMember(inviteeID, familyID) {
			continue // 已经是成员，跳过
		}

		// 检查是否已经邀请过
		var existingInvitation models.FamilyInvitation
		err := tx.Where("family_id = ? AND invitee_id = ? AND status IN (?)", 
			familyID, inviteeID, []string{"pending", "accepted"}).
			First(&existingInvitation).Error

		if err == nil {
			continue // 已经邀请过，跳过
		}

		// 创建邀请
		invitation := &models.FamilyInvitation{
			ID:        uuid.New().String(),
			FamilyID:  familyID,
			InviterID: userID,
			InviteeID: inviteeID,
			Status:    "pending",
			Message:   req.Message,
			ExpiresAt: time.Now().AddDate(0, 0, 30), // 30天后过期
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(invitation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// 通过邀请码加入家族圈
func (s *FamilyService) JoinFamilyByCode(userID, inviteCode string) error {
	var family models.Family
	err := s.db.Where("invite_code = ?", inviteCode).First(&family).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("邀请码无效")
		}
		return err
	}

	// 检查是否已经是成员
	if s.isFamilyMember(userID, family.ID) {
		return errors.New("您已经是此家族圈的成员")
	}

	// 添加为成员
	member := &models.FamilyMember{
		ID:       uuid.New().String(),
		FamilyID: family.ID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	if err := s.db.Create(member).Error; err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(family.ID, userID, "", "join", map[string]interface{}{
		"method": "invite_code",
	})

	return nil
}

// 响应邀请
func (s *FamilyService) RespondToInvitation(userID, invitationID string, accept bool) error {
	var invitation models.FamilyInvitation
	err := s.db.Preload("Family").First(&invitation, "id = ? AND invitee_id = ?", invitationID, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("邀请不存在")
		}
		return err
	}

	// 检查邀请状态
	if invitation.Status != "pending" {
		return errors.New("邀请已处理")
	}

	// 检查是否过期
	if time.Now().After(invitation.ExpiresAt) {
		return errors.New("邀请已过期")
	}

	tx := s.db.Begin()

	// 更新邀请状态
	status := "declined"
	if accept {
		status = "accepted"
	}

	if err := tx.Model(&invitation).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果接受邀请，添加为成员
	if accept {
		// 检查是否已经是成员
		if !s.isFamilyMember(userID, invitation.FamilyID) {
			member := &models.FamilyMember{
				ID:       uuid.New().String(),
				FamilyID: invitation.FamilyID,
				UserID:   userID,
				Role:     "member",
				JoinedAt: time.Now(),
			}

			if err := tx.Create(member).Error; err != nil {
				tx.Rollback()
				return err
			}

			// 记录活动
			s.recordActivity(invitation.FamilyID, userID, "", "join", map[string]interface{}{
				"method": "invitation",
			})
		}
	}

	return tx.Commit().Error
}

// 移除成员
func (s *FamilyService) RemoveMember(userID, familyID, memberID string) error {
	// 验证权限（管理员可以移除成员）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以移除成员")
	}

	// 不能移除创建者
	var family models.Family
	s.db.First(&family, "id = ?", familyID)
	if family.CreatorID == memberID {
		return errors.New("不能移除家族圈创建者")
	}

	// 移除成员
	err := s.db.Where("family_id = ? AND user_id = ?", familyID, memberID).Delete(&models.FamilyMember{}).Error
	if err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, "", "remove_member", map[string]interface{}{
		"removed_user_id": memberID,
	})

	return nil
}

// 离开家族圈
func (s *FamilyService) LeaveFamily(userID, familyID string) error {
	// 创建者不能离开
	var family models.Family
	s.db.First(&family, "id = ?", familyID)
	if family.CreatorID == userID {
		return errors.New("创建者不能离开家族圈")
	}

	// 移除成员身份
	err := s.db.Where("family_id = ? AND user_id = ?", familyID, userID).Delete(&models.FamilyMember{}).Error
	if err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, "", "leave", nil)

	return nil
}

// 设置成员角色
func (s *FamilyService) SetMemberRole(userID, familyID, memberID, role string) error {
	// 验证权限（只有创建者可以设置管理员）
	var family models.Family
	s.db.First(&family, "id = ?", familyID)
	if family.CreatorID != userID {
		return errors.New("只有创建者可以设置成员角色")
	}

	// 不能修改创建者的角色
	if family.CreatorID == memberID {
		return errors.New("不能修改创建者的角色")
	}

	// 验证角色
	if role != "admin" && role != "member" {
		return errors.New("无效的角色")
	}

	// 更新成员角色
	err := s.db.Model(&models.FamilyMember{}).
		Where("family_id = ? AND user_id = ?", familyID, memberID).
		Update("role", role).Error

	if err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, "", "set_role", map[string]interface{}{
		"target_user_id": memberID,
		"new_role":       role,
	})

	return nil
}

// 获取家族成员列表
func (s *FamilyService) GetFamilyMembers(userID, familyID string, page, pageSize int) ([]*models.FamilyMember, int64, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, 0, errors.New("您不是此家族圈的成员")
	}

	var members []*models.FamilyMember
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.FamilyMember{}).Where("family_id = ?", familyID).Count(&total)

	// 查询成员列表
	err := s.db.Preload("User").
		Where("family_id = ?", familyID).
		Order("role DESC, joined_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&members).Error

	return members, total, err
}

// 获取家族活动
func (s *FamilyService) GetFamilyActivities(userID, familyID string, page, pageSize int) ([]*models.FamilyActivity, int64, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, 0, errors.New("您不是此家族圈的成员")
	}

	var activities []*models.FamilyActivity
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.FamilyActivity{}).Where("family_id = ?", familyID).Count(&total)

	// 查询活动列表
	err := s.db.Preload("User").
		Preload("Memorial").
		Where("family_id = ?", familyID).
		Order("timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&activities).Error

	return activities, total, err
}

// 关联纪念馆到家族圈
func (s *FamilyService) AddMemorialToFamily(userID, familyID, memorialID string) error {
	// 验证权限（管理员可以关联纪念馆）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以关联纪念馆")
	}

	// 验证纪念馆权限
	var memorial models.Memorial
	err := s.db.First(&memorial, "id = ? AND creator_id = ?", memorialID, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆不存在或无权操作")
		}
		return err
	}

	// 检查是否已经关联
	var existingRelation models.MemorialFamily
	err = s.db.Where("memorial_id = ? AND family_id = ?", memorialID, familyID).First(&existingRelation).Error
	if err == nil {
		return errors.New("纪念馆已经关联到此家族圈")
	}

	// 创建关联
	relation := &models.MemorialFamily{
		ID:         uuid.New().String(),
		MemorialID: memorialID,
		FamilyID:   familyID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(relation).Error; err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, memorialID, "add_memorial", nil)

	return nil
}

// 移除纪念馆关联
func (s *FamilyService) RemoveMemorialFromFamily(userID, familyID, memorialID string) error {
	// 验证权限（管理员可以移除关联）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以移除纪念馆关联")
	}

	// 移除关联
	err := s.db.Where("memorial_id = ? AND family_id = ?", memorialID, familyID).Delete(&models.MemorialFamily{}).Error
	if err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, memorialID, "remove_memorial", nil)

	return nil
}

// 辅助方法

// 生成邀请码
func (s *FamilyService) generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 检查是否是家族成员
func (s *FamilyService) isFamilyMember(userID, familyID string) bool {
	var member models.FamilyMember
	err := s.db.Where("family_id = ? AND user_id = ?", familyID, userID).First(&member).Error
	return err == nil
}

// 检查是否是家族管理员
func (s *FamilyService) isFamilyAdmin(userID, familyID string) bool {
	var member models.FamilyMember
	err := s.db.Where("family_id = ? AND user_id = ? AND role = ?", familyID, userID, "admin").First(&member).Error
	if err == nil {
		return true
	}

	// 检查是否是创建者
	var family models.Family
	err = s.db.Where("id = ? AND creator_id = ?", familyID, userID).First(&family).Error
	return err == nil
}

// 设置纪念日提醒请求
type SetReminderRequest struct {
	MemorialID   string    `json:"memorial_id" binding:"required"`
	ReminderType string    `json:"reminder_type" binding:"required"` // birthday|death_anniversary|festival
	ReminderDate time.Time `json:"reminder_date" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Content      string    `json:"content"`
}

// 集体祭扫请求
type CollectiveWorshipRequest struct {
	MemorialID   string                 `json:"memorial_id" binding:"required"`
	WorshipType  string                 `json:"worship_type" binding:"required"` // flower|candle|incense|tribute|prayer
	Content      map[string]interface{} `json:"content"`
	ScheduleTime *time.Time             `json:"schedule_time,omitempty"`
}

// 设置纪念日提醒
func (s *FamilyService) SetMemorialReminder(userID, familyID string, req *SetReminderRequest) error {
	// 验证权限（管理员可以设置提醒）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以设置纪念日提醒")
	}

	// 验证纪念馆是否属于家族圈
	var relation models.MemorialFamily
	err := s.db.Where("memorial_id = ? AND family_id = ?", req.MemorialID, familyID).First(&relation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆未关联到此家族圈")
		}
		return err
	}

	// 验证提醒类型
	validTypes := []string{"birthday", "death_anniversary", "festival"}
	isValidType := false
	for _, validType := range validTypes {
		if req.ReminderType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return errors.New("无效的提醒类型")
	}

	// 创建提醒
	reminder := &models.MemorialReminder{
		ID:           uuid.New().String(),
		MemorialID:   req.MemorialID,
		ReminderType: req.ReminderType,
		ReminderDate: req.ReminderDate,
		Title:        req.Title,
		Content:      req.Content,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(reminder).Error; err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, req.MemorialID, "set_reminder", map[string]interface{}{
		"reminder_type": req.ReminderType,
		"reminder_date": req.ReminderDate.Format("2006-01-02"),
		"title":         req.Title,
	})

	return nil
}

// 获取家族纪念日提醒
func (s *FamilyService) GetFamilyReminders(userID, familyID string, page, pageSize int) ([]*models.MemorialReminder, int64, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, 0, errors.New("您不是此家族圈的成员")
	}

	var reminders []*models.MemorialReminder
	var total int64

	offset := (page - 1) * pageSize

	// 获取家族圈关联的纪念馆ID列表
	var memorialIDs []string
	s.db.Model(&models.MemorialFamily{}).
		Where("family_id = ?", familyID).
		Pluck("memorial_id", &memorialIDs)

	if len(memorialIDs) == 0 {
		return reminders, 0, nil
	}

	// 查询总数
	s.db.Model(&models.MemorialReminder{}).
		Where("memorial_id IN (?) AND is_active = ?", memorialIDs, true).
		Count(&total)

	// 查询提醒列表
	err := s.db.Preload("Memorial").
		Where("memorial_id IN (?) AND is_active = ?", memorialIDs, true).
		Order("reminder_date ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&reminders).Error

	return reminders, total, err
}

// 获取即将到来的提醒（3天内）
func (s *FamilyService) GetUpcomingReminders(userID, familyID string) ([]*models.MemorialReminder, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, errors.New("您不是此家族圈的成员")
	}

	// 获取家族圈关联的纪念馆ID列表
	var memorialIDs []string
	s.db.Model(&models.MemorialFamily{}).
		Where("family_id = ?", familyID).
		Pluck("memorial_id", &memorialIDs)

	if len(memorialIDs) == 0 {
		return []*models.MemorialReminder{}, nil
	}

	// 查询3天内的提醒
	now := time.Now()
	threeDaysLater := now.AddDate(0, 0, 3)

	var reminders []*models.MemorialReminder
	err := s.db.Preload("Memorial").
		Where("memorial_id IN (?) AND is_active = ? AND reminder_date BETWEEN ? AND ?", 
			memorialIDs, true, now.Format("2006-01-02"), threeDaysLater.Format("2006-01-02")).
		Order("reminder_date ASC").
		Find(&reminders).Error

	return reminders, err
}

// 删除纪念日提醒
func (s *FamilyService) DeleteReminder(userID, familyID, reminderID string) error {
	// 验证权限（管理员可以删除提醒）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以删除纪念日提醒")
	}

	var reminder models.MemorialReminder
	err := s.db.First(&reminder, "id = ?", reminderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("提醒不存在")
		}
		return err
	}

	// 验证纪念馆是否属于家族圈
	var relation models.MemorialFamily
	err = s.db.Where("memorial_id = ? AND family_id = ?", reminder.MemorialID, familyID).First(&relation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权删除此提醒")
		}
		return err
	}

	// 软删除提醒
	if err := s.db.Model(&reminder).Update("is_active", false).Error; err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(familyID, userID, reminder.MemorialID, "delete_reminder", map[string]interface{}{
		"reminder_type": reminder.ReminderType,
		"title":         reminder.Title,
	})

	return nil
}

// 发起集体祭扫
func (s *FamilyService) InitiateCollectiveWorship(userID, familyID string, req *CollectiveWorshipRequest) error {
	// 验证权限（管理员可以发起集体祭扫）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以发起集体祭扫")
	}

	// 验证纪念馆是否属于家族圈
	var relation models.MemorialFamily
	err := s.db.Where("memorial_id = ? AND family_id = ?", req.MemorialID, familyID).First(&relation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆未关联到此家族圈")
		}
		return err
	}

	// 验证祭扫类型
	validTypes := []string{"flower", "candle", "incense", "tribute", "prayer"}
	isValidType := false
	for _, validType := range validTypes {
		if req.WorshipType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return errors.New("无效的祭扫类型")
	}

	// 创建集体祭扫活动
	activityContent := map[string]interface{}{
		"worship_type":   req.WorshipType,
		"content":        req.Content,
		"schedule_time":  req.ScheduleTime,
		"initiator_id":   userID,
		"participants":   []string{userID}, // 发起者自动参与
		"status":         "active",
	}

	// 记录活动
	s.recordActivity(familyID, userID, req.MemorialID, "collective_worship", activityContent)

	return nil
}

// 参与集体祭扫
func (s *FamilyService) JoinCollectiveWorship(userID, familyID, activityID string) error {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return errors.New("您不是此家族圈的成员")
	}

	var activity models.FamilyActivity
	err := s.db.First(&activity, "id = ? AND family_id = ? AND activity_type = ?", 
		activityID, familyID, "collective_worship").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("集体祭扫活动不存在")
		}
		return err
	}

	// 解析活动内容
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(activity.Content), &content); err != nil {
		return errors.New("活动数据格式错误")
	}

	// 检查活动状态
	if status, ok := content["status"].(string); !ok || status != "active" {
		return errors.New("活动已结束")
	}

	// 检查是否已经参与
	participants, ok := content["participants"].([]interface{})
	if !ok {
		participants = []interface{}{}
	}

	for _, participant := range participants {
		if participantID, ok := participant.(string); ok && participantID == userID {
			return errors.New("您已经参与了此活动")
		}
	}

	// 添加参与者
	participants = append(participants, userID)
	content["participants"] = participants

	// 更新活动内容
	contentJSON, _ := json.Marshal(content)
	if err := s.db.Model(&activity).Update("content", string(contentJSON)).Error; err != nil {
		return err
	}

	// 记录参与活动
	s.recordActivity(familyID, userID, activity.MemorialID, "join_collective_worship", map[string]interface{}{
		"activity_id":  activityID,
		"worship_type": content["worship_type"],
	})

	return nil
}

// 同步祭扫动态到家族圈
func (s *FamilyService) SyncWorshipActivity(userID, memorialID string, worshipType string, content interface{}) error {
	// 查找纪念馆关联的家族圈
	var relations []models.MemorialFamily
	err := s.db.Where("memorial_id = ?", memorialID).Find(&relations).Error
	if err != nil {
		return err
	}

	// 为每个关联的家族圈记录活动
	for _, relation := range relations {
		// 检查用户是否是家族成员
		if s.isFamilyMember(userID, relation.FamilyID) {
			s.recordActivity(relation.FamilyID, userID, memorialID, "worship", map[string]interface{}{
				"worship_type": worshipType,
				"content":      content,
			})
		}
	}

	return nil
}

// 家族谱系相关请求结构
type CreateGenealogyRequest struct {
	PersonName   string     `json:"person_name" binding:"required"`
	Generation   int        `json:"generation" binding:"required"`
	ParentID     string     `json:"parent_id"`
	Gender       string     `json:"gender" binding:"required,oneof=male female"`
	BirthDate    *time.Time `json:"birth_date"`
	DeathDate    *time.Time `json:"death_date"`
	Biography    string     `json:"biography"`
	AvatarURL    string     `json:"avatar_url"`
	MemorialID   string     `json:"memorial_id"`
	Position     string     `json:"position"`
	Achievements string     `json:"achievements"`
}

type UpdateGenealogyRequest struct {
	PersonName   string     `json:"person_name"`
	Generation   int        `json:"generation"`
	ParentID     string     `json:"parent_id"`
	Gender       string     `json:"gender"`
	BirthDate    *time.Time `json:"birth_date"`
	DeathDate    *time.Time `json:"death_date"`
	Biography    string     `json:"biography"`
	AvatarURL    string     `json:"avatar_url"`
	MemorialID   string     `json:"memorial_id"`
	Position     string     `json:"position"`
	Achievements string     `json:"achievements"`
}

// 家族故事相关请求结构
type CreateFamilyStoryRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Category   string   `json:"category" binding:"required"`
	Period     string   `json:"period"`
	Characters []string `json:"characters"`
	Location   string   `json:"location"`
	MediaFiles []string `json:"media_files"`
	Tags       []string `json:"tags"`
	IsPublic   bool     `json:"is_public"`
}

type UpdateFamilyStoryRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Category   string   `json:"category"`
	Period     string   `json:"period"`
	Characters []string `json:"characters"`
	Location   string   `json:"location"`
	MediaFiles []string `json:"media_files"`
	Tags       []string `json:"tags"`
	IsPublic   *bool    `json:"is_public"`
}

// 家族传统相关请求结构
type CreateFamilyTraditionRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Category    string   `json:"category" binding:"required"`
	Origin      string   `json:"origin"`
	Practice    string   `json:"practice"`
	Meaning     string   `json:"meaning"`
	MediaFiles  []string `json:"media_files"`
	IsActive    bool     `json:"is_active"`
}

type UpdateFamilyTraditionRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Origin      string   `json:"origin"`
	Practice    string   `json:"practice"`
	Meaning     string   `json:"meaning"`
	MediaFiles  []string `json:"media_files"`
	IsActive    *bool    `json:"is_active"`
}

// 创建家族谱系成员
func (s *FamilyService) CreateGenealogy(userID, familyID string, req *CreateGenealogyRequest) (*models.FamilyGenealogy, error) {
	// 验证权限（管理员可以创建谱系）
	if !s.isFamilyAdmin(userID, familyID) {
		return nil, errors.New("只有管理员可以创建家族谱系")
	}

	// 验证性别
	if req.Gender != "male" && req.Gender != "female" {
		return nil, errors.New("无效的性别")
	}

	// 如果指定了父辈，验证父辈是否存在
	if req.ParentID != "" {
		var parent models.FamilyGenealogy
		err := s.db.Where("id = ? AND family_id = ?", req.ParentID, familyID).First(&parent).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("指定的父辈不存在")
			}
			return nil, err
		}
	}

	// 如果指定了纪念馆，验证纪念馆是否关联到家族圈
	if req.MemorialID != "" {
		var relation models.MemorialFamily
		err := s.db.Where("memorial_id = ? AND family_id = ?", req.MemorialID, familyID).First(&relation).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("纪念馆未关联到此家族圈")
			}
			return nil, err
		}
	}

	genealogy := &models.FamilyGenealogy{
		ID:           uuid.New().String(),
		FamilyID:     familyID,
		PersonName:   req.PersonName,
		Generation:   req.Generation,
		ParentID:     req.ParentID,
		Gender:       req.Gender,
		BirthDate:    req.BirthDate,
		DeathDate:    req.DeathDate,
		Biography:    req.Biography,
		AvatarURL:    req.AvatarURL,
		MemorialID:   req.MemorialID,
		Position:     req.Position,
		Achievements: req.Achievements,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(genealogy).Error; err != nil {
		return nil, err
	}

	// 重新查询包含关联数据
	s.db.Preload("Parent").Preload("Children").Preload("Memorial").First(genealogy, "id = ?", genealogy.ID)

	// 记录活动
	s.recordActivity(familyID, userID, req.MemorialID, "create_genealogy", map[string]interface{}{
		"person_name": req.PersonName,
		"generation":  req.Generation,
	})

	return genealogy, nil
}

// 获取家族谱系
func (s *FamilyService) GetFamilyGenealogy(userID, familyID string) ([]*models.FamilyGenealogy, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, errors.New("您不是此家族圈的成员")
	}

	var genealogies []*models.FamilyGenealogy
	err := s.db.Preload("Parent").
		Preload("Children").
		Preload("Memorial").
		Where("family_id = ?", familyID).
		Order("generation ASC, person_name ASC").
		Find(&genealogies).Error

	return genealogies, err
}

// 更新家族谱系成员
func (s *FamilyService) UpdateGenealogy(userID, familyID, genealogyID string, req *UpdateGenealogyRequest) error {
	// 验证权限（管理员可以更新谱系）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以更新家族谱系")
	}

	var genealogy models.FamilyGenealogy
	err := s.db.Where("id = ? AND family_id = ?", genealogyID, familyID).First(&genealogy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("谱系成员不存在")
		}
		return err
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.PersonName != "" {
		updates["person_name"] = req.PersonName
	}
	if req.Generation != 0 {
		updates["generation"] = req.Generation
	}
	if req.ParentID != "" {
		updates["parent_id"] = req.ParentID
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.BirthDate != nil {
		updates["birth_date"] = req.BirthDate
	}
	if req.DeathDate != nil {
		updates["death_date"] = req.DeathDate
	}
	if req.Biography != "" {
		updates["biography"] = req.Biography
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.MemorialID != "" {
		updates["memorial_id"] = req.MemorialID
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.Achievements != "" {
		updates["achievements"] = req.Achievements
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&genealogy).Updates(updates).Error
}

// 删除家族谱系成员
func (s *FamilyService) DeleteGenealogy(userID, familyID, genealogyID string) error {
	// 验证权限（管理员可以删除谱系）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以删除家族谱系")
	}

	var genealogy models.FamilyGenealogy
	err := s.db.Where("id = ? AND family_id = ?", genealogyID, familyID).First(&genealogy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("谱系成员不存在")
		}
		return err
	}

	// 检查是否有子代，如果有则不能删除
	var childCount int64
	s.db.Model(&models.FamilyGenealogy{}).Where("parent_id = ?", genealogyID).Count(&childCount)
	if childCount > 0 {
		return errors.New("该成员有子代记录，无法删除")
	}

	return s.db.Delete(&genealogy).Error
}

// 创建家族故事
func (s *FamilyService) CreateFamilyStory(userID, familyID string, req *CreateFamilyStoryRequest) (*models.FamilyStory, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, errors.New("您不是此家族圈的成员")
	}

	// 验证故事分类
	validCategories := []string{"tradition", "achievement", "migration", "business", "education", "war", "love"}
	isValidCategory := false
	for _, category := range validCategories {
		if req.Category == category {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		return nil, errors.New("无效的故事分类")
	}

	// 序列化数组字段
	charactersJSON, _ := json.Marshal(req.Characters)
	mediaFilesJSON, _ := json.Marshal(req.MediaFiles)
	tagsStr := ""
	if len(req.Tags) > 0 {
		tagsStr = fmt.Sprintf(",%s,", fmt.Sprintf("%s", req.Tags))
	}

	story := &models.FamilyStory{
		ID:         uuid.New().String(),
		FamilyID:   familyID,
		AuthorID:   userID,
		Title:      req.Title,
		Content:    req.Content,
		Category:   req.Category,
		Period:     req.Period,
		Characters: string(charactersJSON),
		Location:   req.Location,
		MediaFiles: string(mediaFilesJSON),
		Tags:       tagsStr,
		IsPublic:   req.IsPublic,
		ViewCount:  0,
		LikeCount:  0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(story).Error; err != nil {
		return nil, err
	}

	// 重新查询包含关联数据
	s.db.Preload("Author").First(story, "id = ?", story.ID)

	// 记录活动
	s.recordActivity(familyID, userID, "", "create_story", map[string]interface{}{
		"title":    req.Title,
		"category": req.Category,
	})

	return story, nil
}

// 获取家族故事列表
func (s *FamilyService) GetFamilyStories(userID, familyID string, category string, page, pageSize int) ([]*models.FamilyStory, int64, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, 0, errors.New("您不是此家族圈的成员")
	}

	var stories []*models.FamilyStory
	var total int64

	offset := (page - 1) * pageSize

	query := s.db.Model(&models.FamilyStory{}).Where("family_id = ? AND is_public = ?", familyID, true)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 查询总数
	query.Count(&total)

	// 查询故事列表
	err := query.Preload("Author").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&stories).Error

	return stories, total, err
}

// 获取家族故事详情
func (s *FamilyService) GetFamilyStory(userID, familyID, storyID string) (*models.FamilyStory, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, errors.New("您不是此家族圈的成员")
	}

	var story models.FamilyStory
	err := s.db.Preload("Author").
		Where("id = ? AND family_id = ?", storyID, familyID).
		First(&story).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("故事不存在")
		}
		return nil, err
	}

	// 增加浏览次数
	s.db.Model(&story).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return &story, nil
}

// 更新家族故事
func (s *FamilyService) UpdateFamilyStory(userID, familyID, storyID string, req *UpdateFamilyStoryRequest) error {
	var story models.FamilyStory
	err := s.db.Where("id = ? AND family_id = ?", storyID, familyID).First(&story).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("故事不存在")
		}
		return err
	}

	// 验证权限（作者或管理员可以更新）
	if story.AuthorID != userID && !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有作者或管理员可以更新故事")
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Period != "" {
		updates["period"] = req.Period
	}
	if req.Characters != nil {
		charactersJSON, _ := json.Marshal(req.Characters)
		updates["characters"] = string(charactersJSON)
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.MediaFiles != nil {
		mediaFilesJSON, _ := json.Marshal(req.MediaFiles)
		updates["media_files"] = string(mediaFilesJSON)
	}
	if req.Tags != nil {
		tagsStr := ""
		if len(req.Tags) > 0 {
			tagsStr = fmt.Sprintf(",%s,", fmt.Sprintf("%s", req.Tags))
		}
		updates["tags"] = tagsStr
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&story).Updates(updates).Error
}

// 删除家族故事
func (s *FamilyService) DeleteFamilyStory(userID, familyID, storyID string) error {
	var story models.FamilyStory
	err := s.db.Where("id = ? AND family_id = ?", storyID, familyID).First(&story).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("故事不存在")
		}
		return err
	}

	// 验证权限（作者或管理员可以删除）
	if story.AuthorID != userID && !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有作者或管理员可以删除故事")
	}

	return s.db.Delete(&story).Error
}

// 创建家族传统
func (s *FamilyService) CreateFamilyTradition(userID, familyID string, req *CreateFamilyTraditionRequest) (*models.FamilyTradition, error) {
	// 验证权限（管理员可以创建传统）
	if !s.isFamilyAdmin(userID, familyID) {
		return nil, errors.New("只有管理员可以创建家族传统")
	}

	// 验证传统分类
	validCategories := []string{"festival", "ceremony", "custom", "rule", "recipe"}
	isValidCategory := false
	for _, category := range validCategories {
		if req.Category == category {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		return nil, errors.New("无效的传统分类")
	}

	// 序列化媒体文件
	mediaFilesJSON, _ := json.Marshal(req.MediaFiles)

	tradition := &models.FamilyTradition{
		ID:          uuid.New().String(),
		FamilyID:    familyID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Origin:      req.Origin,
		Practice:    req.Practice,
		Meaning:     req.Meaning,
		MediaFiles:  string(mediaFilesJSON),
		IsActive:    req.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(tradition).Error; err != nil {
		return nil, err
	}

	// 记录活动
	s.recordActivity(familyID, userID, "", "create_tradition", map[string]interface{}{
		"name":     req.Name,
		"category": req.Category,
	})

	return tradition, nil
}

// 获取家族传统列表
func (s *FamilyService) GetFamilyTraditions(userID, familyID string, category string, page, pageSize int) ([]*models.FamilyTradition, int64, error) {
	// 验证访问权限
	if !s.isFamilyMember(userID, familyID) {
		return nil, 0, errors.New("您不是此家族圈的成员")
	}

	var traditions []*models.FamilyTradition
	var total int64

	offset := (page - 1) * pageSize

	query := s.db.Model(&models.FamilyTradition{}).Where("family_id = ?", familyID)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 查询总数
	query.Count(&total)

	// 查询传统列表
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&traditions).Error

	return traditions, total, err
}

// 更新家族传统
func (s *FamilyService) UpdateFamilyTradition(userID, familyID, traditionID string, req *UpdateFamilyTraditionRequest) error {
	// 验证权限（管理员可以更新传统）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以更新家族传统")
	}

	var tradition models.FamilyTradition
	err := s.db.Where("id = ? AND family_id = ?", traditionID, familyID).First(&tradition).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("传统不存在")
		}
		return err
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Origin != "" {
		updates["origin"] = req.Origin
	}
	if req.Practice != "" {
		updates["practice"] = req.Practice
	}
	if req.Meaning != "" {
		updates["meaning"] = req.Meaning
	}
	if req.MediaFiles != nil {
		mediaFilesJSON, _ := json.Marshal(req.MediaFiles)
		updates["media_files"] = string(mediaFilesJSON)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&tradition).Updates(updates).Error
}

// 删除家族传统
func (s *FamilyService) DeleteFamilyTradition(userID, familyID, traditionID string) error {
	// 验证权限（管理员可以删除传统）
	if !s.isFamilyAdmin(userID, familyID) {
		return errors.New("只有管理员可以删除家族传统")
	}

	var tradition models.FamilyTradition
	err := s.db.Where("id = ? AND family_id = ?", traditionID, familyID).First(&tradition).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("传统不存在")
		}
		return err
	}

	return s.db.Delete(&tradition).Error
}

// 记录活动
func (s *FamilyService) recordActivity(familyID, userID, memorialID, activityType string, content interface{}) {
	var contentJSON string
	if content != nil {
		if jsonBytes, err := json.Marshal(content); err == nil {
			contentJSON = string(jsonBytes)
		}
	}

	activity := &models.FamilyActivity{
		ID:           uuid.New().String(),
		FamilyID:     familyID,
		UserID:       userID,
		MemorialID:   memorialID,
		ActivityType: activityType,
		Content:      contentJSON,
		Timestamp:    time.Now(),
		CreatedAt:    time.Now(),
	}

	s.db.Create(activity)
}