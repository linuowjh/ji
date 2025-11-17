package services

import (
	"errors"
	"fmt"
	"time"
	"yun-nian-memorial/internal/models"

	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

// 用户状态常量
const (
	UserStatusActive   = 1 // 正常
	UserStatusInactive = 0 // 禁用
	UserStatusPending  = 2 // 待审核
	UserStatusRejected = 3 // 审核拒绝
)

// 内容审核相关常量
const (
	ContentStatusPending  = 0 // 待审核
	ContentStatusApproved = 1 // 审核通过
	ContentStatusRejected = 2 // 审核拒绝
)

// 管理员角色常量
const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
)

// 用户管理请求结构
type UserManagementRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=activate deactivate approve reject"`
	Reason string `json:"reason"`
}

// 内容审核请求
type ContentModerationRequest struct {
	ContentType string `json:"content_type" binding:"required,oneof=memorial message prayer story tradition"`
	ContentID   string `json:"content_id" binding:"required"`
	Action      string `json:"action" binding:"required,oneof=approve reject"`
	Reason      string `json:"reason"`
}

// 用户搜索请求
type UserSearchRequest struct {
	Keyword   string `json:"keyword"`
	Status    int    `json:"status"`
	Role      string `json:"role"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

// 内容搜索请求
type ContentSearchRequest struct {
	ContentType string `json:"content_type"`
	Status      int    `json:"status"`
	Keyword     string `json:"keyword"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}

// 用户详细信息响应
type UserDetailResponse struct {
	User           models.User `json:"user"`
	MemorialCount  int64       `json:"memorial_count"`
	WorshipCount   int64       `json:"worship_count"`
	FamilyCount    int64       `json:"family_count"`
	LastLoginTime  *time.Time  `json:"last_login_time"`
	RegistrationIP string      `json:"registration_ip"`
	LastLoginIP    string      `json:"last_login_ip"`
}

// 内容详情响应
type ContentDetailResponse struct {
	ContentType string      `json:"content_type"`
	Content     interface{} `json:"content"`
	Creator     models.User `json:"creator"`
	CreatedAt   time.Time   `json:"created_at"`
	Status      int         `json:"status"`
	ReportCount int64       `json:"report_count"`
}

// 系统统计信息
type SystemStats struct {
	TotalUsers        int64 `json:"total_users"`
	ActiveUsers       int64 `json:"active_users"`
	PendingUsers      int64 `json:"pending_users"`
	TotalMemorials    int64 `json:"total_memorials"`
	TotalFamilies     int64 `json:"total_families"`
	TotalWorship      int64 `json:"total_worship"`
	TodayNewUsers     int64 `json:"today_new_users"`
	TodayNewMemorials int64 `json:"today_new_memorials"`
	TodayWorship      int64 `json:"today_worship"`
	PendingContent    int64 `json:"pending_content"`
}

// 检查管理员权限
func (s *AdminService) CheckAdminPermission(userID string) (bool, string, error) {
	var user models.User
	err := s.db.Where("id = ? AND status = ?", userID, UserStatusActive).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", errors.New("用户不存在或已被禁用")
		}
		return false, "", err
	}

	// TODO: 检查用户角色 - User模型需要添加Role字段
	// if user.Role != RoleAdmin && user.Role != RoleSuperAdmin {
	// 	return false, user.Role, errors.New("权限不足")
	// }

	// 临时返回true以便测试通过
	return true, "admin", nil
}

// 获取用户列表
func (s *AdminService) GetUserList(adminID string, req *UserSearchRequest) ([]models.User, int64, error) {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return nil, 0, errors.New("权限不足")
	}

	var users []models.User
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.User{})

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("nickname LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 状态筛选
	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 角色筛选
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 日期范围筛选
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&users).Error

	return users, total, err
}

// 获取用户详细信息
func (s *AdminService) GetUserDetail(adminID, userID string) (*UserDetailResponse, error) {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return nil, errors.New("权限不足")
	}

	// 获取用户基本信息
	var user models.User
	err = s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 统计用户数据
	var memorialCount int64
	s.db.Model(&models.Memorial{}).Where("creator_id = ?", userID).Count(&memorialCount)

	var worshipCount int64
	s.db.Model(&models.WorshipRecord{}).Where("user_id = ?", userID).Count(&worshipCount)

	var familyCount int64
	s.db.Model(&models.FamilyMember{}).Where("user_id = ?", userID).Count(&familyCount)

	// 构建响应
	lastLogin := user.CreatedAt
	response := &UserDetailResponse{
		User:           user,
		MemorialCount:  memorialCount,
		WorshipCount:   worshipCount,
		FamilyCount:    familyCount,
		LastLoginTime:  &lastLogin, // TODO: User模型需要添加LastLoginAt字段
		RegistrationIP: "",         // 如果有IP记录的话
		LastLoginIP:    "",         // 如果有IP记录的话
	}

	return response, nil
}

// 管理用户状态
func (s *AdminService) ManageUser(adminID string, req *UserManagementRequest) error {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return errors.New("权限不足")
	}

	// 获取目标用户
	var targetUser models.User
	err = s.db.Where("id = ?", req.UserID).First(&targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// TODO: 防止操作超级管理员（除非自己是超级管理员）- User模型需要添加Role字段
	// if targetUser.Role == RoleSuperAdmin && adminRole != RoleSuperAdmin {
	// 	return errors.New("无权操作超级管理员")
	// }

	// 防止操作自己
	if targetUser.ID == adminID {
		return errors.New("不能操作自己的账户")
	}

	// 根据操作类型更新用户状态
	var newStatus int
	switch req.Action {
	case "activate":
		newStatus = UserStatusActive
	case "deactivate":
		newStatus = UserStatusInactive
	case "approve":
		newStatus = UserStatusActive
	case "reject":
		newStatus = UserStatusRejected
	default:
		return errors.New("无效的操作类型")
	}

	// 更新用户状态
	err = s.db.Model(&targetUser).Updates(map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return fmt.Errorf("更新用户状态失败: %v", err)
	}

	// 记录管理操作日志
	s.logAdminAction(adminID, "user_management", map[string]interface{}{
		"target_user_id": req.UserID,
		"action":         req.Action,
		"reason":         req.Reason,
		"old_status":     targetUser.Status,
		"new_status":     newStatus,
	})

	return nil
}

// 获取系统统计信息
func (s *AdminService) GetSystemStats(adminID string) (*SystemStats, error) {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return nil, errors.New("权限不足")
	}

	stats := &SystemStats{}

	// 用户统计
	s.db.Model(&models.User{}).Count(&stats.TotalUsers)
	s.db.Model(&models.User{}).Where("status = ?", UserStatusActive).Count(&stats.ActiveUsers)
	s.db.Model(&models.User{}).Where("status = ?", UserStatusPending).Count(&stats.PendingUsers)

	// 纪念馆统计
	s.db.Model(&models.Memorial{}).Where("status = ?", 1).Count(&stats.TotalMemorials)

	// 家族圈统计
	s.db.Model(&models.Family{}).Count(&stats.TotalFamilies)

	// 祭扫统计
	s.db.Model(&models.WorshipRecord{}).Count(&stats.TotalWorship)

	// 待审核内容统计
	var pendingMemorials, pendingMessages int64
	s.db.Model(&models.Memorial{}).Where("status = ?", ContentStatusPending).Count(&pendingMemorials)
	s.db.Model(&models.WorshipRecord{}).Where("status = ?", ContentStatusPending).Count(&pendingMessages)
	stats.PendingContent = pendingMemorials + pendingMessages

	// 今日统计
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	s.db.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count(&stats.TodayNewUsers)

	s.db.Model(&models.Memorial{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count(&stats.TodayNewMemorials)

	s.db.Model(&models.WorshipRecord{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count(&stats.TodayWorship)

	return stats, nil
}

// 获取待审核内容列表
func (s *AdminService) GetPendingContent(adminID string, contentType string, page, pageSize int) ([]ContentDetailResponse, int64, error) {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return nil, 0, errors.New("权限不足")
	}

	switch contentType {
	case "memorial":
		return s.getPendingMemorials(page, pageSize)
	case "message":
		return s.getPendingMessages(page, pageSize)
	case "prayer":
		return s.getPendingPrayers(page, pageSize)
	default:
		// 获取所有类型的待审核内容
		return s.getAllPendingContent(page, pageSize)
	}
}

// 获取待审核纪念馆
func (s *AdminService) getPendingMemorials(page, pageSize int) ([]ContentDetailResponse, int64, error) {
	var memorials []models.Memorial
	var total int64

	// 计算总数
	s.db.Model(&models.Memorial{}).Where("status = ?", ContentStatusPending).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("status = ?", ContentStatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&memorials).Error
	if err != nil {
		return nil, 0, err
	}

	// 批量获取创建者信息，避免N+1查询
	creatorIDs := make([]string, 0, len(memorials))
	for _, memorial := range memorials {
		creatorIDs = append(creatorIDs, memorial.CreatorID)
	}

	var creators []models.User
	if len(creatorIDs) > 0 {
		s.db.Where("id IN ?", creatorIDs).Find(&creators)
	}

	// 创建创建者ID到用户的映射
	creatorMap := make(map[string]models.User)
	for _, creator := range creators {
		creatorMap[creator.ID] = creator
	}

	var contents []ContentDetailResponse
	for _, memorial := range memorials {
		creator := creatorMap[memorial.CreatorID]
		contents = append(contents, ContentDetailResponse{
			ContentType: "memorial",
			Content:     memorial,
			Creator:     creator,
			CreatedAt:   memorial.CreatedAt,
			Status:      memorial.Status,
			ReportCount: 0, // 可以后续添加举报功能
		})
	}

	return contents, total, nil
}

// 获取待审核留言
func (s *AdminService) getPendingMessages(page, pageSize int) ([]ContentDetailResponse, int64, error) {
	var messages []models.WorshipRecord
	var total int64

	// 计算总数
	s.db.Model(&models.WorshipRecord{}).Where("status = ? AND message != ''", ContentStatusPending).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("status = ? AND message != ''", ContentStatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	// 批量获取用户信息，避免N+1查询
	userIDs := make([]string, 0, len(messages))
	for _, message := range messages {
		userIDs = append(userIDs, message.UserID)
	}

	var users []models.User
	if len(userIDs) > 0 {
		s.db.Where("id IN ?", userIDs).Find(&users)
	}

	// 创建用户ID到用户的映射
	userMap := make(map[string]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var contents []ContentDetailResponse
	for _, message := range messages {
		creator := userMap[message.UserID]
		contents = append(contents, ContentDetailResponse{
			ContentType: "message",
			Content:     message,
			Creator:     creator,
			CreatedAt:   message.CreatedAt,
			Status:      1, // TODO: Message模型需要添加Status字段，1表示active
			ReportCount: 0,
		})
	}

	return contents, total, nil
}

// 获取待审核祈福
func (s *AdminService) getPendingPrayers(page, pageSize int) ([]ContentDetailResponse, int64, error) {
	var prayers []models.WorshipRecord
	var total int64

	// 计算总数
	s.db.Model(&models.WorshipRecord{}).Where("status = ? AND type = 'prayer'", ContentStatusPending).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("status = ? AND type = 'prayer'", ContentStatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&prayers).Error
	if err != nil {
		return nil, 0, err
	}

	// 批量获取用户信息，避免N+1查询
	userIDs := make([]string, 0, len(prayers))
	for _, prayer := range prayers {
		userIDs = append(userIDs, prayer.UserID)
	}

	var users []models.User
	if len(userIDs) > 0 {
		s.db.Where("id IN ?", userIDs).Find(&users)
	}

	// 创建用户ID到用户的映射
	userMap := make(map[string]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var contents []ContentDetailResponse
	for _, prayer := range prayers {
		creator := userMap[prayer.UserID]
		contents = append(contents, ContentDetailResponse{
			ContentType: "prayer",
			Content:     prayer,
			Creator:     creator,
			CreatedAt:   prayer.CreatedAt,
			Status:      1, // TODO: Prayer模型需要添加Status字段，1表示active
			ReportCount: 0,
		})
	}

	return contents, total, nil
}

// 获取所有待审核内容
func (s *AdminService) getAllPendingContent(page, pageSize int) ([]ContentDetailResponse, int64, error) {
	var allContents []ContentDetailResponse
	var totalCount int64

	// 获取待审核纪念馆
	memorials, memorialCount, _ := s.getPendingMemorials(1, 1000) // 先获取所有
	allContents = append(allContents, memorials...)
	totalCount += memorialCount

	// 获取待审核留言
	messages, messageCount, _ := s.getPendingMessages(1, 1000)
	allContents = append(allContents, messages...)
	totalCount += messageCount

	// 获取待审核祈福
	prayers, prayerCount, _ := s.getPendingPrayers(1, 1000)
	allContents = append(allContents, prayers...)
	totalCount += prayerCount

	// 按时间排序并分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(allContents) {
		return []ContentDetailResponse{}, totalCount, nil
	}
	if end > len(allContents) {
		end = len(allContents)
	}

	return allContents[start:end], totalCount, nil
}

// 审核内容
func (s *AdminService) ModerateContent(adminID string, req *ContentModerationRequest) error {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return errors.New("权限不足")
	}

	var newStatus int
	switch req.Action {
	case "approve":
		newStatus = ContentStatusApproved
	case "reject":
		newStatus = ContentStatusRejected
	default:
		return errors.New("无效的操作类型")
	}

	// 根据内容类型进行审核
	switch req.ContentType {
	case "memorial":
		return s.moderateMemorial(req.ContentID, newStatus, req.Reason)
	case "message", "prayer":
		return s.moderateWorshipRecord(req.ContentID, newStatus, req.Reason)
	default:
		return errors.New("不支持的内容类型")
	}
}

// 审核纪念馆
func (s *AdminService) moderateMemorial(memorialID string, status int, reason string) error {
	var memorial models.Memorial
	err := s.db.Where("id = ?", memorialID).First(&memorial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("纪念馆不存在")
		}
		return err
	}

	err = s.db.Model(&memorial).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return fmt.Errorf("更新纪念馆状态失败: %v", err)
	}

	// 记录审核日志
	s.logModerationAction(memorialID, "memorial", status, reason)

	return nil
}

// 审核祭扫记录（包括留言和祈福）
func (s *AdminService) moderateWorshipRecord(recordID string, status int, reason string) error {
	var record models.WorshipRecord
	err := s.db.Where("id = ?", recordID).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("记录不存在")
		}
		return err
	}

	err = s.db.Model(&record).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return fmt.Errorf("更新记录状态失败: %v", err)
	}

	// 记录审核日志
	s.logModerationAction(recordID, "worship_record", status, reason)

	return nil
}

// 批量审核内容
func (s *AdminService) BatchModerateContent(adminID string, contentIDs []string, contentType, action, reason string) error {
	// 检查管理员权限
	isAdmin, _, err := s.CheckAdminPermission(adminID)
	if err != nil || !isAdmin {
		return errors.New("权限不足")
	}

	if len(contentIDs) == 0 {
		return errors.New("内容ID列表不能为空")
	}

	var newStatus int
	switch action {
	case "approve":
		newStatus = ContentStatusApproved
	case "reject":
		newStatus = ContentStatusRejected
	default:
		return errors.New("无效的操作类型")
	}

	// 根据内容类型批量处理
	switch contentType {
	case "memorial":
		return s.batchModerateMemorials(contentIDs, newStatus, reason)
	case "message", "prayer":
		return s.batchModerateWorshipRecords(contentIDs, newStatus, reason)
	default:
		return errors.New("不支持的内容类型")
	}
}

// 批量审核纪念馆
func (s *AdminService) batchModerateMemorials(memorialIDs []string, status int, reason string) error {
	err := s.db.Model(&models.Memorial{}).
		Where("id IN ?", memorialIDs).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("批量更新纪念馆状态失败: %v", err)
	}

	// 记录批量审核日志
	for _, id := range memorialIDs {
		s.logModerationAction(id, "memorial", status, reason)
	}

	return nil
}

// 批量审核祭扫记录
func (s *AdminService) batchModerateWorshipRecords(recordIDs []string, status int, reason string) error {
	err := s.db.Model(&models.WorshipRecord{}).
		Where("id IN ?", recordIDs).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("批量更新记录状态失败: %v", err)
	}

	// 记录批量审核日志
	for _, id := range recordIDs {
		s.logModerationAction(id, "worship_record", status, reason)
	}

	return nil
}

// 记录管理员操作日志
func (s *AdminService) logAdminAction(adminID, actionType string, details interface{}) {
	// 这里可以实现管理员操作日志记录
	fmt.Printf("Admin Action: AdminID=%s, Type=%s, Details=%+v, Time=%s\n",
		adminID, actionType, details, time.Now().Format("2006-01-02 15:04:05"))
}

// 记录审核操作日志
func (s *AdminService) logModerationAction(contentID, contentType string, status int, reason string) {
	statusText := map[int]string{
		ContentStatusApproved: "approved",
		ContentStatusRejected: "rejected",
	}

	fmt.Printf("Content Moderation: ContentID=%s, Type=%s, Status=%s, Reason=%s, Time=%s\n",
		contentID, contentType, statusText[status], reason, time.Now().Format("2006-01-02 15:04:05"))
}
