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

type MemorialServiceService struct {
	db *gorm.DB
}

func NewMemorialServiceService(db *gorm.DB) *MemorialServiceService {
	return &MemorialServiceService{
		db: db,
	}
}

// 创建追思会请求
type CreateMemorialServiceRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	MaxParticipants int       `json:"max_participants"`
	IsPublic        bool      `json:"is_public"`
}

// 更新追思会请求
type UpdateMemorialServiceRequest struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	MaxParticipants int        `json:"max_participants"`
	IsPublic        *bool      `json:"is_public"`
}

// 邀请参与者请求
type InviteParticipantRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`
	Message string   `json:"message"`
}

// 发送聊天消息请求
type SendChatMessageRequest struct {
	MessageType string `json:"message_type" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

// 创建追思会
func (s *MemorialServiceService) CreateMemorialService(userID, memorialID string, req *CreateMemorialServiceRequest) (*models.MemorialService, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	// 验证时间
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("结束时间不能早于开始时间")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("开始时间不能早于当前时间")
	}

	// 生成邀请码
	inviteCode := s.generateInviteCode()

	service := &models.MemorialService{
		ID:              uuid.New().String(),
		MemorialID:      memorialID,
		Title:           req.Title,
		Description:     req.Description,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Status:          "scheduled",
		MaxParticipants: req.MaxParticipants,
		IsPublic:        req.IsPublic,
		InviteCode:      inviteCode,
		HostID:          userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if service.MaxParticipants <= 0 {
		service.MaxParticipants = 50
	}

	tx := s.db.Begin()

	if err := tx.Create(service).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加主持人为参与者
	participant := &models.MemorialServiceParticipant{
		ID:        uuid.New().String(),
		ServiceID: service.ID,
		UserID:    userID,
		Role:      "host",
		Status:    "joined",
		JoinedAt:  &service.CreatedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(participant).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新查询包含关联数据的追思会
	s.db.Preload("Host").Preload("Memorial").First(service, "id = ?", service.ID)

	return service, nil
}

// 获取追思会列表
func (s *MemorialServiceService) GetMemorialServices(userID, memorialID string, page, pageSize int) ([]*models.MemorialService, int64, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, 0, err
	}

	var services []*models.MemorialService
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.MemorialService{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 查询追思会列表
	err := s.db.Preload("Host").
		Where("memorial_id = ?", memorialID).
		Order("start_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error

	return services, total, err
}

// 获取追思会详情
func (s *MemorialServiceService) GetMemorialService(userID, serviceID string) (*models.MemorialService, error) {
	var service models.MemorialService
	err := s.db.Preload("Host").
		Preload("Memorial").
		Preload("Participants.User").
		First(&service, "id = ?", serviceID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("追思会不存在")
		}
		return nil, err
	}

	// 验证访问权限
	if err := s.validateServiceAccess(userID, &service); err != nil {
		return nil, err
	}

	return &service, nil
}

// 更新追思会
func (s *MemorialServiceService) UpdateMemorialService(userID, serviceID string, req *UpdateMemorialServiceRequest) error {
	var service models.MemorialService
	err := s.db.First(&service, "id = ?", serviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("追思会不存在")
		}
		return err
	}

	// 验证权限（只有主持人可以修改）
	if service.HostID != userID {
		return errors.New("只有主持人可以修改追思会")
	}

	// 检查状态（只有未开始的追思会可以修改）
	if service.Status != "scheduled" {
		return errors.New("只有未开始的追思会可以修改")
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartTime != nil {
		if req.StartTime.Before(time.Now()) {
			return errors.New("开始时间不能早于当前时间")
		}
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != nil {
		updates["end_time"] = req.EndTime
	}
	if req.MaxParticipants > 0 {
		updates["max_participants"] = req.MaxParticipants
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&service).Updates(updates).Error
}

// 删除追思会
func (s *MemorialServiceService) DeleteMemorialService(userID, serviceID string) error {
	var service models.MemorialService
	err := s.db.First(&service, "id = ?", serviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("追思会不存在")
		}
		return err
	}

	// 验证权限（只有主持人可以删除）
	if service.HostID != userID {
		return errors.New("只有主持人可以删除追思会")
	}

	// 检查状态（进行中的追思会不能删除）
	if service.Status == "ongoing" {
		return errors.New("进行中的追思会不能删除")
	}

	return s.db.Delete(&service).Error
}

// 邀请参与者
func (s *MemorialServiceService) InviteParticipants(userID, serviceID string, req *InviteParticipantRequest) error {
	var service models.MemorialService
	err := s.db.First(&service, "id = ?", serviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("追思会不存在")
		}
		return err
	}

	// 验证权限（主持人和协助主持人可以邀请）
	if !s.canInviteParticipants(userID, serviceID) {
		return errors.New("无权邀请参与者")
	}

	tx := s.db.Begin()

	for _, inviteeID := range req.UserIDs {
		// 检查是否已经邀请过
		var existingInvitation models.ServiceInvitation
		err := tx.Where("service_id = ? AND invitee_id = ? AND status IN (?)", 
			serviceID, inviteeID, []string{"pending", "accepted"}).
			First(&existingInvitation).Error

		if err == nil {
			continue // 已经邀请过，跳过
		}

		// 创建邀请
		invitation := &models.ServiceInvitation{
			ID:        uuid.New().String(),
			ServiceID: serviceID,
			InviterID: userID,
			InviteeID: inviteeID,
			Status:    "pending",
			Message:   req.Message,
			ExpiresAt: service.StartTime.Add(-1 * time.Hour), // 开始前1小时过期
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

// 响应邀请
func (s *MemorialServiceService) RespondToInvitation(userID, invitationID string, accept bool) error {
	var invitation models.ServiceInvitation
	err := s.db.Preload("Service").First(&invitation, "id = ? AND invitee_id = ?", invitationID, userID).Error
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

	// 如果接受邀请，添加为参与者
	if accept {
		participant := &models.MemorialServiceParticipant{
			ID:        uuid.New().String(),
			ServiceID: invitation.ServiceID,
			UserID:    userID,
			Role:      "participant",
			Status:    "invited",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(participant).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// 加入追思会
func (s *MemorialServiceService) JoinService(userID, serviceID string) error {
	var service models.MemorialService
	err := s.db.First(&service, "id = ?", serviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("追思会不存在")
		}
		return err
	}

	// 检查追思会状态
	if service.Status != "ongoing" {
		return errors.New("追思会未开始或已结束")
	}

	// 检查是否是参与者
	var participant models.MemorialServiceParticipant
	err = s.db.Where("service_id = ? AND user_id = ?", serviceID, userID).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您未被邀请参加此追思会")
		}
		return err
	}

	// 更新参与者状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":     "joined",
		"joined_at":  &now,
		"updated_at": now,
	}

	if err := s.db.Model(&participant).Updates(updates).Error; err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(serviceID, userID, "join", nil)

	return nil
}

// 离开追思会
func (s *MemorialServiceService) LeaveService(userID, serviceID string) error {
	var participant models.MemorialServiceParticipant
	err := s.db.Where("service_id = ? AND user_id = ?", serviceID, userID).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您未参加此追思会")
		}
		return err
	}

	// 主持人不能离开
	if participant.Role == "host" {
		return errors.New("主持人不能离开追思会")
	}

	// 更新参与者状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":     "left",
		"left_at":    &now,
		"updated_at": now,
	}

	if err := s.db.Model(&participant).Updates(updates).Error; err != nil {
		return err
	}

	// 记录活动
	s.recordActivity(serviceID, userID, "leave", nil)

	return nil
}

// 发送聊天消息
func (s *MemorialServiceService) SendChatMessage(userID, serviceID string, req *SendChatMessageRequest) (*models.ServiceChat, error) {
	// 验证参与者身份
	if !s.isParticipant(userID, serviceID) {
		return nil, errors.New("您不是此追思会的参与者")
	}

	message := &models.ServiceChat{
		ID:          uuid.New().String(),
		ServiceID:   serviceID,
		UserID:      userID,
		MessageType: req.MessageType,
		Content:     req.Content,
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	// 重新查询包含用户信息的消息
	s.db.Preload("User").First(message, "id = ?", message.ID)

	return message, nil
}

// 获取聊天消息
func (s *MemorialServiceService) GetChatMessages(userID, serviceID string, page, pageSize int) ([]*models.ServiceChat, int64, error) {
	// 验证参与者身份
	if !s.isParticipant(userID, serviceID) {
		return nil, 0, errors.New("您不是此追思会的参与者")
	}

	var messages []*models.ServiceChat
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.ServiceChat{}).Where("service_id = ?", serviceID).Count(&total)

	// 查询消息列表
	err := s.db.Preload("User").
		Where("service_id = ?", serviceID).
		Order("timestamp ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// 开始追思会
func (s *MemorialServiceService) StartService(userID, serviceID string) error {
	var service models.MemorialService
	err := s.db.First(&service, "id = ?", serviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("追思会不存在")
		}
		return err
	}

	// 验证权限（只有主持人可以开始）
	if service.HostID != userID {
		return errors.New("只有主持人可以开始追思会")
	}

	// 检查状态
	if service.Status != "scheduled" {
		return errors.New("追思会已开始或已结束")
	}

	// 更新状态
	updates := map[string]interface{}{
		"status":     "ongoing",
		"updated_at": time.Now(),
	}

	return s.db.Model(&service).Updates(updates).Error
}

// 结束追思会
func (s *MemorialServiceService) EndService(userID, serviceID string) error {
	var service models.MemorialService
	err := s.db.First(&service, "id = ?", serviceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("追思会不存在")
		}
		return err
	}

	// 验证权限（只有主持人可以结束）
	if service.HostID != userID {
		return errors.New("只有主持人可以结束追思会")
	}

	// 检查状态
	if service.Status != "ongoing" {
		return errors.New("追思会未开始或已结束")
	}

	// 更新状态
	updates := map[string]interface{}{
		"status":     "completed",
		"updated_at": time.Now(),
	}

	if err := s.db.Model(&service).Updates(updates).Error; err != nil {
		return err
	}

	// 生成录制视频（异步处理）
	go s.generateRecording(serviceID)

	return nil
}

// 辅助方法

// 生成邀请码
func (s *MemorialServiceService) generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 验证纪念馆访问权限
func (s *MemorialServiceService) validateMemorialAccess(userID, memorialID string) error {
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

// 验证追思会访问权限
func (s *MemorialServiceService) validateServiceAccess(userID string, service *models.MemorialService) error {
	// 公开的追思会任何人都可以查看
	if service.IsPublic {
		return nil
	}

	// 主持人可以访问
	if service.HostID == userID {
		return nil
	}

	// 检查是否是参与者
	if s.isParticipant(userID, service.ID) {
		return nil
	}

	return errors.New("无权访问此追思会")
}

// 检查是否可以邀请参与者
func (s *MemorialServiceService) canInviteParticipants(userID, serviceID string) bool {
	var participant models.MemorialServiceParticipant
	err := s.db.Where("service_id = ? AND user_id = ? AND role IN (?)", 
		serviceID, userID, []string{"host", "co_host"}).First(&participant).Error
	return err == nil
}

// 检查是否是参与者
func (s *MemorialServiceService) isParticipant(userID, serviceID string) bool {
	var participant models.MemorialServiceParticipant
	err := s.db.Where("service_id = ? AND user_id = ?", serviceID, userID).First(&participant).Error
	return err == nil
}

// 记录活动
func (s *MemorialServiceService) recordActivity(serviceID, userID, activityType string, content interface{}) {
	var contentJSON string
	if content != nil {
		if jsonBytes, err := json.Marshal(content); err == nil {
			contentJSON = string(jsonBytes)
		}
	}

	activity := &models.ServiceActivity{
		ID:           uuid.New().String(),
		ServiceID:    serviceID,
		UserID:       userID,
		ActivityType: activityType,
		Content:      contentJSON,
		Timestamp:    time.Now(),
		CreatedAt:    time.Now(),
	}

	s.db.Create(activity)
}

// 生成录制视频（模拟实现）
func (s *MemorialServiceService) generateRecording(serviceID string) {
	// 这里应该调用视频处理服务来生成录制视频
	// 目前只是模拟实现
	time.Sleep(5 * time.Second) // 模拟处理时间

	recordingURL := fmt.Sprintf("/recordings/%s.mp4", serviceID)
	
	recording := &models.ServiceRecording{
		ID:           uuid.New().String(),
		ServiceID:    serviceID,
		RecordingURL: recordingURL,
		ThumbnailURL: fmt.Sprintf("/recordings/%s_thumb.jpg", serviceID),
		Duration:     3600, // 模拟1小时
		FileSize:     1024 * 1024 * 100, // 模拟100MB
		Status:       "completed",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	s.db.Create(recording)

	// 更新追思会的录制URL
	s.db.Model(&models.MemorialService{}).
		Where("id = ?", serviceID).
		Update("recording_url", recordingURL)
}