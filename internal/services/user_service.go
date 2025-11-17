package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db     *gorm.DB
	config *config.Config
}

type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type WechatLoginResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	IsNewUser bool        `json:"is_new_user"`
}

type WechatSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// UpcomingReminderResponse 即将到来的提醒响应结构
type UpcomingReminderResponse struct {
	ID           string    `json:"id"`
	ReminderType string    `json:"reminder_type"`
	ReminderDate time.Time `json:"reminder_date"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	DaysUntil    int       `json:"days_until"`
	Memorial     struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"memorial"`
	Family struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"family"`
}

func NewUserService(db *gorm.DB, config *config.Config) *UserService {
	return &UserService{
		db:     db,
		config: config,
	}
}

// WechatLogin 微信小程序登录
func (s *UserService) WechatLogin(req *WechatLoginRequest) (*WechatLoginResponse, error) {
	// 1. 通过code获取微信用户信息
	sessionResp, err := s.getWechatSession(req.Code)
	if err != nil {
		return nil, fmt.Errorf("获取微信会话失败: %v", err)
	}

	if sessionResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s", sessionResp.ErrMsg)
	}

	// 2. 查找或创建用户
	var user models.User
	var isNewUser bool

	err = s.db.Where("wechat_open_id = ?", sessionResp.OpenID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新用户
			user = models.User{
				ID:            uuid.New().String(),
				WechatOpenID:  sessionResp.OpenID,
				WechatUnionID: sessionResp.UnionID,
				Nickname:      req.Nickname,
				AvatarURL:     req.Avatar,
				Status:        1,
			}

			if err := s.db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("创建用户失败: %v", err)
			}
			isNewUser = true
		} else {
			return nil, fmt.Errorf("查询用户失败: %v", err)
		}
	} else {
		// 更新用户信息
		if req.Nickname != "" {
			user.Nickname = req.Nickname
		}
		if req.Avatar != "" {
			user.AvatarURL = req.Avatar
		}
		if err := s.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("更新用户信息失败: %v", err)
		}
	}

	// 3. 生成JWT Token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	return &WechatLoginResponse{
		Token:     token,
		User:      user,
		IsNewUser: isNewUser,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID string) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ? AND status = ?", userID, 1).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(userID string, nickname, phone string) error {
	updates := make(map[string]interface{})

	if nickname != "" {
		updates["nickname"] = nickname
	}
	if phone != "" {
		updates["phone"] = phone
	}

	if len(updates) == 0 {
		return fmt.Errorf("没有需要更新的信息")
	}

	err := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}
	return nil
}

// getWechatSession 获取微信会话信息
func (s *UserService) getWechatSession(code string) (*WechatSessionResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.config.Wechat.AppID,
		s.config.Wechat.AppSecret,
		code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sessionResp WechatSessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return nil, err
	}

	return &sessionResp, nil
}

// generateJWT 生成JWT Token
func (s *UserService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.config.JWT.ExpireTime) * time.Second).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// ValidateJWT 验证JWT Token
func (s *UserService) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", fmt.Errorf("invalid user_id in token")
	}

	return "", fmt.Errorf("invalid token")
}

// GenerateInviteCode 生成邀请码
func (s *UserService) GenerateInviteCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetUserMemorials 获取用户创建的纪念馆列表
func (s *UserService) GetUserMemorials(userID string, page, pageSize int) ([]models.Memorial, int64, error) {
	var memorials []models.Memorial
	var total int64

	// 计算总数
	s.db.Model(&models.Memorial{}).Where("creator_id = ? AND status = ?", userID, 1).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("creator_id = ? AND status = ?", userID, 1).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&memorials).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询用户纪念馆失败: %v", err)
	}

	return memorials, total, nil
}

// GetUserWorshipRecords 获取用户祭扫记录
func (s *UserService) GetUserWorshipRecords(userID string, page, pageSize int) ([]models.WorshipRecord, int64, error) {
	var records []models.WorshipRecord
	var total int64

	// 计算总数
	s.db.Model(&models.WorshipRecord{}).Where("user_id = ?", userID).Count(&total)

	// 分页查询，预加载纪念馆信息
	offset := (page - 1) * pageSize
	err := s.db.Where("user_id = ?", userID).
		Preload("Memorial").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询用户祭扫记录失败: %v", err)
	}

	return records, total, nil
}

// GetUserMemorialDetails 获取用户纪念馆详细信息（包含统计数据）
func (s *UserService) GetUserMemorialDetails(userID string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var memorials []models.Memorial
	var total int64

	// 计算总数
	s.db.Model(&models.Memorial{}).Where("creator_id = ? AND status = ?", userID, 1).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("creator_id = ? AND status = ?", userID, 1).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&memorials).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询用户纪念馆失败: %v", err)
	}

	// 为每个纪念馆添加统计信息
	var result []map[string]interface{}
	for _, memorial := range memorials {
		// 统计访客数量
		var visitorCount int64
		s.db.Model(&models.VisitorRecord{}).Where("memorial_id = ?", memorial.ID).Count(&visitorCount)

		// 统计祭扫次数
		var worshipCount int64
		s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorial.ID).Count(&worshipCount)

		// 统计祈福数量
		var prayerCount int64
		s.db.Model(&models.Prayer{}).Where("memorial_id = ?", memorial.ID).Count(&prayerCount)

		// 统计留言数量
		var messageCount int64
		s.db.Model(&models.Message{}).Where("memorial_id = ?", memorial.ID).Count(&messageCount)

		// 获取最近访客
		var recentVisitors []models.VisitorRecord
		s.db.Where("memorial_id = ?", memorial.ID).
			Preload("User").
			Order("visit_time DESC").
			Limit(5).
			Find(&recentVisitors)

		memorialData := map[string]interface{}{
			"memorial":        memorial,
			"visitor_count":   visitorCount,
			"worship_count":   worshipCount,
			"prayer_count":    prayerCount,
			"message_count":   messageCount,
			"recent_visitors": recentVisitors,
		}

		result = append(result, memorialData)
	}

	return result, total, nil
}

// GetMemorialVisitors 获取纪念馆访客记录
func (s *UserService) GetMemorialVisitors(userID, memorialID string, page, pageSize int) ([]models.VisitorRecord, int64, error) {
	// 验证纪念馆所有权
	var memorial models.Memorial
	err := s.db.Where("id = ? AND creator_id = ?", memorialID, userID).First(&memorial).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("纪念馆不存在或无权访问")
		}
		return nil, 0, fmt.Errorf("查询纪念馆失败: %v", err)
	}

	var visitors []models.VisitorRecord
	var total int64

	// 计算总数
	s.db.Model(&models.VisitorRecord{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 分页查询，预加载用户信息
	offset := (page - 1) * pageSize
	err = s.db.Where("memorial_id = ?", memorialID).
		Preload("User").
		Order("visit_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&visitors).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询访客记录失败: %v", err)
	}

	return visitors, total, nil
}

// GetUserFamilies 获取用户参与的家族圈
func (s *UserService) GetUserFamilies(userID string, page, pageSize int) ([]models.Family, int64, error) {
	var families []models.Family
	var total int64

	// 通过家族成员表查询用户参与的家族圈
	subQuery := s.db.Model(&models.FamilyMember{}).
		Select("family_id").
		Where("user_id = ?", userID)

	// 计算总数
	s.db.Model(&models.Family{}).Where("id IN (?)", subQuery).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("id IN (?)", subQuery).
		Preload("Creator").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&families).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询用户家族圈失败: %v", err)
	}

	return families, total, nil
}

// GetUserStatistics 获取用户统计信息
func (s *UserService) GetUserStatistics(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 统计创建的纪念馆数量
	var memorialCount int64
	s.db.Model(&models.Memorial{}).Where("creator_id = ? AND status = ?", userID, 1).Count(&memorialCount)
	stats["memorialCount"] = memorialCount

	// 统计祭扫次数
	var worshipCount int64
	s.db.Model(&models.WorshipRecord{}).Where("user_id = ?", userID).Count(&worshipCount)
	stats["worshipCount"] = worshipCount

	// 统计参与的家族圈数量
	var familyCount int64
	s.db.Model(&models.FamilyMember{}).Where("user_id = ?", userID).Count(&familyCount)
	stats["familyCount"] = familyCount

	// 统计发布的祈福数量
	var prayerCount int64
	s.db.Model(&models.Prayer{}).Where("user_id = ?", userID).Count(&prayerCount)
	stats["prayerCount"] = prayerCount

	// 统计发布的留言数量
	var messageCount int64
	s.db.Model(&models.Message{}).Where("user_id = ?", userID).Count(&messageCount)
	stats["messageCount"] = messageCount

	// 统计最近7天的活动
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var recentWorshipCount int64
	s.db.Model(&models.WorshipRecord{}).
		Where("user_id = ? AND created_at >= ?", userID, sevenDaysAgo).
		Count(&recentWorshipCount)
	stats["recentWorshipCount"] = recentWorshipCount

	// 获取用户创建的纪念馆的总访客数
	var totalVisitors int64
	s.db.Table("visitor_records vr").
		Joins("JOIN memorials m ON vr.memorial_id = m.id").
		Where("m.creator_id = ? AND m.status = ?", userID, 1).
		Count(&totalVisitors)
	stats["totalVisitors"] = totalVisitors

	return stats, nil
}

// GetUserRecentActivities 获取用户最近活动
func (s *UserService) GetUserRecentActivities(userID string, limit int) ([]map[string]interface{}, error) {
	var activities []map[string]interface{}

	// 获取最近的祭扫记录
	var worshipRecords []models.WorshipRecord
	s.db.Where("user_id = ?", userID).
		Preload("Memorial", "deleted_at IS NULL"). // 只加载未删除的纪念馆
		Order("created_at DESC").
		Limit(limit).
		Find(&worshipRecords)

	for _, record := range worshipRecords {
		// 跳过纪念馆已被删除的记录
		if record.Memorial.ID == "" {
			continue
		}
		activity := map[string]interface{}{
			"type":       "worship",
			"action":     record.WorshipType,
			"memorial":   record.Memorial,
			"created_at": record.CreatedAt,
		}
		activities = append(activities, activity)
	}

	// 获取最近的家族活动
	var familyActivities []models.FamilyActivity
	s.db.Where("user_id = ?", userID).
		Preload("Family", "deleted_at IS NULL").   // 只加载未删除的家族
		Preload("Memorial", "deleted_at IS NULL"). // 只加载未删除的纪念馆
		Order("timestamp DESC").
		Limit(limit).
		Find(&familyActivities)

	for _, activity := range familyActivities {
		// 跳过纪念馆或家族已被删除的记录
		if activity.Memorial.ID == "" && activity.Family.ID == "" {
			continue
		}
		activityData := map[string]interface{}{
			"type":       "family",
			"action":     activity.ActivityType,
			"family":     activity.Family,
			"memorial":   activity.Memorial,
			"created_at": activity.Timestamp,
		}
		activities = append(activities, activityData)
	}

	return activities, nil
}

// GetUpcomingReminders 获取用户所有家族的即将到来的纪念日提醒
func (s *UserService) GetUpcomingReminders(userID string) ([]*UpcomingReminderResponse, error) {
	// 1. Get all families the user belongs to
	var familyMembers []models.FamilyMember
	err := s.db.Where("user_id = ?", userID).Find(&familyMembers).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户家族失败: %v", err)
	}

	if len(familyMembers) == 0 {
		return []*UpcomingReminderResponse{}, nil
	}

	// 2. Extract family IDs
	familyIDs := make([]string, len(familyMembers))
	for i, member := range familyMembers {
		familyIDs[i] = member.FamilyID
	}

	// 3. Get all memorial IDs associated with these families
	var memorialFamilies []models.MemorialFamily
	err = s.db.Where("family_id IN ?", familyIDs).Find(&memorialFamilies).Error
	if err != nil {
		return nil, fmt.Errorf("查询家族纪念馆关联失败: %v", err)
	}

	if len(memorialFamilies) == 0 {
		return []*UpcomingReminderResponse{}, nil
	}

	// 4. Create a map of memorial ID to family IDs for later lookup
	memorialToFamilies := make(map[string][]string)
	for _, mf := range memorialFamilies {
		memorialToFamilies[mf.MemorialID] = append(memorialToFamilies[mf.MemorialID], mf.FamilyID)
	}

	memorialIDs := make([]string, 0, len(memorialToFamilies))
	for memorialID := range memorialToFamilies {
		memorialIDs = append(memorialIDs, memorialID)
	}

	// 5. Query reminders within next 30 days
	now := time.Now()
	thirtyDaysLater := now.AddDate(0, 0, 30)

	var reminders []models.MemorialReminder
	err = s.db.Preload("Memorial").
		Where("memorial_id IN ? AND is_active = ? AND reminder_date BETWEEN ? AND ?",
			memorialIDs, true, now.Format("2006-01-02"), thirtyDaysLater.Format("2006-01-02")).
		Order("reminder_date ASC").
		Find(&reminders).Error
	if err != nil {
		return nil, fmt.Errorf("查询纪念日提醒失败: %v", err)
	}

	// 6. Load family information
	var families []models.Family
	err = s.db.Where("id IN ?", familyIDs).Find(&families).Error
	if err != nil {
		return nil, fmt.Errorf("查询家族信息失败: %v", err)
	}

	familyMap := make(map[string]*models.Family)
	for i := range families {
		familyMap[families[i].ID] = &families[i]
	}

	// 7. Build response with enriched data
	responses := make([]*UpcomingReminderResponse, 0, len(reminders))
	for _, reminder := range reminders {
		// Get the first family associated with this memorial
		familyIDsForMemorial := memorialToFamilies[reminder.MemorialID]
		if len(familyIDsForMemorial) == 0 {
			continue
		}

		family := familyMap[familyIDsForMemorial[0]]
		if family == nil {
			continue
		}

		// Calculate days until reminder
		daysUntil := int(reminder.ReminderDate.Sub(now).Hours() / 24)

		response := &UpcomingReminderResponse{
			ID:           reminder.ID,
			ReminderType: reminder.ReminderType,
			ReminderDate: reminder.ReminderDate,
			Title:        reminder.Title,
			Content:      reminder.Content,
			DaysUntil:    daysUntil,
		}

		response.Memorial.ID = reminder.Memorial.ID
		response.Memorial.Name = reminder.Memorial.DeceasedName
		response.Memorial.AvatarURL = reminder.Memorial.AvatarURL

		response.Family.ID = family.ID
		response.Family.Name = family.Name

		responses = append(responses, response)
	}

	return responses, nil
}

// UpdatePhoneRequest 更新手机号请求
type UpdatePhoneRequest struct {
	Code string `json:"code" binding:"required"`
}

// UpdatePhone 更新用户手机号
func (s *UserService) UpdatePhone(userID string, code string) (string, error) {
	// 这里应该调用微信API解密手机号
	// 由于是开发环境，我们暂时模拟返回
	// 生产环境需要调用: https://api.weixin.qq.com/wxa/business/getuserphonenumber

	// TODO: 实际生产环境需要调用微信API
	// url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken)
	// 发送POST请求，body: {"code": code}

	// 开发环境模拟
	phone := "138****8888" // 模拟的手机号

	// 更新数据库
	err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("phone", phone).Error
	if err != nil {
		return "", fmt.Errorf("更新手机号失败: %v", err)
	}

	return phone, nil
}
