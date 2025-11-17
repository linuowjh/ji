package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/utils"

	"gorm.io/gorm"
)

// FlexibleDate 支持多种日期格式的自定义类型
type FlexibleDate struct {
	time.Time
}

// UnmarshalJSON 自定义 JSON 解析，支持多种日期格式
func (fd *FlexibleDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// 尝试多种日期格式
	formats := []string{
		"2006-01-02",                // YYYY-MM-DD
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05Z",      // RFC3339 without timezone
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
	}

	var err error
	for _, format := range formats {
		fd.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("无法解析日期格式: %s", s)
}

// MarshalJSON 自定义 JSON 序列化
func (fd FlexibleDate) MarshalJSON() ([]byte, error) {
	if fd.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(fd.Time.Format("2006-01-02"))
}

// convertFlexibleDateToTimePtr 将 FlexibleDate 转换为 *time.Time
func convertFlexibleDateToTimePtr(fd *FlexibleDate) *time.Time {
	if fd == nil || fd.Time.IsZero() {
		return nil
	}
	return &fd.Time
}

type MemorialService struct {
	db                *gorm.DB
	permissionManager *utils.PermissionManager
}

type CreateMemorialRequest struct {
	DeceasedName   string        `json:"deceasedName" binding:"required"` // 支持驼峰命名
	BirthDate      *FlexibleDate `json:"birthDate"`                       // 支持多种日期格式
	DeathDate      *FlexibleDate `json:"deathDate"`                       // 支持多种日期格式
	Biography      string        `json:"biography"`
	AvatarURL      string        `json:"avatarUrl"`
	ThemeStyle     string        `json:"themeStyle"`
	TombstoneStyle string        `json:"tombstoneStyle"`
	Epitaph        string        `json:"epitaph"`
	PrivacyLevel   int           `json:"privacyLevel"`
}

type UpdateMemorialRequest struct {
	DeceasedName   string        `json:"deceasedName"` // 支持驼峰命名
	BirthDate      *FlexibleDate `json:"birthDate"`    // 支持多种日期格式
	DeathDate      *FlexibleDate `json:"deathDate"`    // 支持多种日期格式
	Biography      string        `json:"biography"`
	AvatarURL      string        `json:"avatarUrl"`
	ThemeStyle     string        `json:"themeStyle"`
	TombstoneStyle string        `json:"tombstoneStyle"`
	Epitaph        string        `json:"epitaph"`
	PrivacyLevel   int           `json:"privacyLevel"`
}

type MemorialListResponse struct {
	ID           string     `json:"id"`
	DeceasedName string     `json:"deceasedName"` // 驼峰命名
	BirthDate    *time.Time `json:"birthDate"`    // 驼峰命名
	DeathDate    *time.Time `json:"deathDate"`    // 驼峰命名
	AvatarURL    string     `json:"avatarUrl"`    // 驼峰命名
	ThemeStyle   string     `json:"themeStyle"`   // 驼峰命名
	PrivacyLevel int        `json:"privacyLevel"` // 驼峰命名
	CreatedAt    time.Time  `json:"createdAt"`    // 驼峰命名
	CreatorName  string     `json:"creatorName"`  // 驼峰命名
	VisitorCount int64      `json:"visitorCount"` // 驼峰命名
	WorshipCount int64      `json:"worshipCount"` // 驼峰命名
}

func NewMemorialService(db *gorm.DB) *MemorialService {
	return &MemorialService{
		db:                db,
		permissionManager: utils.NewPermissionManager(db),
	}
}

// CreateMemorial 创建纪念馆
func (s *MemorialService) CreateMemorial(userID string, req *CreateMemorialRequest) (*models.Memorial, error) {
	// 验证输入参数
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 创建纪念馆
	memorial := &models.Memorial{
		ID:             utils.GenerateUUID(),
		CreatorID:      userID,
		DeceasedName:   req.DeceasedName,
		BirthDate:      convertFlexibleDateToTimePtr(req.BirthDate),
		DeathDate:      convertFlexibleDateToTimePtr(req.DeathDate),
		Biography:      req.Biography,
		AvatarURL:      req.AvatarURL,
		ThemeStyle:     s.getDefaultThemeStyle(req.ThemeStyle),
		TombstoneStyle: s.getDefaultTombstoneStyle(req.TombstoneStyle),
		Epitaph:        req.Epitaph,
		PrivacyLevel:   s.getDefaultPrivacyLevel(req.PrivacyLevel),
		Status:         1,
	}

	if err := s.db.Create(memorial).Error; err != nil {
		return nil, fmt.Errorf("创建纪念馆失败: %v", err)
	}

	// 预加载创建者信息
	if err := s.db.Preload("Creator").Where("id = ?", memorial.ID).First(memorial).Error; err != nil {
		return nil, fmt.Errorf("获取纪念馆信息失败: %v", err)
	}

	return memorial, nil
}

// GetMemorial 获取纪念馆详情
func (s *MemorialService) GetMemorial(userID, memorialID string) (*models.Memorial, error) {
	// 检查访问权限
	canAccess, accessLevel, err := s.permissionManager.CanAccessMemorial(userID, memorialID)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, fmt.Errorf("无权访问此纪念馆")
	}

	// 获取纪念馆详情
	var memorial models.Memorial
	if err := s.db.Preload("Creator").Where("id = ? AND status = ?", memorialID, 1).First(&memorial).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("纪念馆不存在")
		}
		return nil, fmt.Errorf("查询纪念馆失败: %v", err)
	}

	// 记录访问
	if accessLevel != "owner" {
		s.permissionManager.RecordVisit(memorialID, userID, "")
	}

	return &memorial, nil
}

// UpdateMemorial 更新纪念馆信息
func (s *MemorialService) UpdateMemorial(userID, memorialID string, req *UpdateMemorialRequest) error {
	// 检查修改权限
	canModify, err := s.permissionManager.CanModifyMemorial(userID, memorialID)
	if err != nil {
		return err
	}
	if !canModify {
		return fmt.Errorf("无权修改此纪念馆")
	}

	// 验证输入参数
	if err := s.validateUpdateRequest(req); err != nil {
		return err
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.DeceasedName != "" {
		updates["deceased_name"] = req.DeceasedName
	}
	if req.BirthDate != nil {
		updates["birth_date"] = convertFlexibleDateToTimePtr(req.BirthDate)
	}
	if req.DeathDate != nil {
		updates["death_date"] = convertFlexibleDateToTimePtr(req.DeathDate)
	}
	if req.Biography != "" {
		updates["biography"] = req.Biography
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.ThemeStyle != "" {
		updates["theme_style"] = req.ThemeStyle
	}
	if req.TombstoneStyle != "" {
		updates["tombstone_style"] = req.TombstoneStyle
	}
	if req.Epitaph != "" {
		updates["epitaph"] = req.Epitaph
	}
	if req.PrivacyLevel > 0 {
		updates["privacy_level"] = req.PrivacyLevel
	}

	if len(updates) == 0 {
		return fmt.Errorf("没有需要更新的信息")
	}

	// 执行更新
	if err := s.db.Model(&models.Memorial{}).Where("id = ?", memorialID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新纪念馆失败: %v", err)
	}

	return nil
}

// DeleteMemorial 删除纪念馆（软删除）
func (s *MemorialService) DeleteMemorial(userID, memorialID string) error {
	// 检查修改权限
	canModify, err := s.permissionManager.CanModifyMemorial(userID, memorialID)
	if err != nil {
		return err
	}
	if !canModify {
		return fmt.Errorf("无权删除此纪念馆")
	}

	// 软删除纪念馆
	if err := s.db.Where("id = ?", memorialID).Delete(&models.Memorial{}).Error; err != nil {
		return fmt.Errorf("删除纪念馆失败: %v", err)
	}

	return nil
}

// GetMemorialList 获取纪念馆列表
func (s *MemorialService) GetMemorialList(userID string, page, pageSize int) ([]MemorialListResponse, int64, error) {
	// 获取用户可访问的纪念馆
	memorials, total, err := s.permissionManager.GetUserAccessibleMemorials(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var response []MemorialListResponse
	for _, memorial := range memorials {
		// 获取创建者信息
		var creator models.User
		s.db.Where("id = ?", memorial.CreatorID).First(&creator)

		// 获取访客数量
		var visitorCount int64
		s.db.Model(&models.VisitorRecord{}).Where("memorial_id = ?", memorial.ID).Count(&visitorCount)

		// 获取祭扫次数
		var worshipCount int64
		s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorial.ID).Count(&worshipCount)

		response = append(response, MemorialListResponse{
			ID:           memorial.ID,
			DeceasedName: memorial.DeceasedName,
			BirthDate:    memorial.BirthDate,
			DeathDate:    memorial.DeathDate,
			AvatarURL:    memorial.AvatarURL,
			ThemeStyle:   memorial.ThemeStyle,
			PrivacyLevel: memorial.PrivacyLevel,
			CreatedAt:    memorial.CreatedAt,
			CreatorName:  creator.Nickname,
			VisitorCount: visitorCount,
			WorshipCount: worshipCount,
		})
	}

	return response, total, nil
}

// GetMemorialVisitors 获取纪念馆访客记录
func (s *MemorialService) GetMemorialVisitors(userID, memorialID string, page, pageSize int) ([]models.VisitorRecord, int64, error) {
	// 检查访问权限（只有创建者可以查看访客记录）
	isOwner, err := s.permissionManager.IsMemorialOwner(userID, memorialID)
	if err != nil {
		return nil, 0, err
	}
	if !isOwner {
		return nil, 0, fmt.Errorf("只有创建者可以查看访客记录")
	}

	return s.permissionManager.GetMemorialVisitors(memorialID, page, pageSize)
}

// validateCreateRequest 验证创建请求
func (s *MemorialService) validateCreateRequest(req *CreateMemorialRequest) error {
	if req.DeceasedName == "" {
		return fmt.Errorf("逝者姓名不能为空")
	}

	if len(req.DeceasedName) > 50 {
		return fmt.Errorf("逝者姓名不能超过50个字符")
	}

	if req.Biography != "" && len(req.Biography) > 2000 {
		return fmt.Errorf("生平简介不能超过2000个字符")
	}

	if req.Epitaph != "" && len(req.Epitaph) > 500 {
		return fmt.Errorf("墓志铭不能超过500个字符")
	}

	if req.PrivacyLevel != 0 && req.PrivacyLevel != 1 && req.PrivacyLevel != 2 {
		return fmt.Errorf("隐私级别必须为1（家族可见）或2（私密）")
	}

	// 验证日期逻辑
	if req.BirthDate != nil && req.DeathDate != nil {
		if req.BirthDate.Time.After(req.DeathDate.Time) {
			return fmt.Errorf("出生日期不能晚于逝世日期")
		}
	}

	return nil
}

// validateUpdateRequest 验证更新请求
func (s *MemorialService) validateUpdateRequest(req *UpdateMemorialRequest) error {
	if req.DeceasedName != "" && len(req.DeceasedName) > 50 {
		return fmt.Errorf("逝者姓名不能超过50个字符")
	}

	if req.Biography != "" && len(req.Biography) > 2000 {
		return fmt.Errorf("生平简介不能超过2000个字符")
	}

	if req.Epitaph != "" && len(req.Epitaph) > 500 {
		return fmt.Errorf("墓志铭不能超过500个字符")
	}

	if req.PrivacyLevel != 0 && req.PrivacyLevel != 1 && req.PrivacyLevel != 2 {
		return fmt.Errorf("隐私级别必须为1（家族可见）或2（私密）")
	}

	// 验证日期逻辑
	if req.BirthDate != nil && req.DeathDate != nil {
		if req.BirthDate.Time.After(req.DeathDate.Time) {
			return fmt.Errorf("出生日期不能晚于逝世日期")
		}
	}

	return nil
}

// getDefaultThemeStyle 获取默认主题风格
func (s *MemorialService) getDefaultThemeStyle(style string) string {
	validStyles := []string{"traditional", "elegant", "natural", "modern"}
	for _, validStyle := range validStyles {
		if style == validStyle {
			return style
		}
	}
	return "traditional"
}

// getDefaultTombstoneStyle 获取默认墓碑样式
func (s *MemorialService) getDefaultTombstoneStyle(style string) string {
	validStyles := []string{"marble", "granite", "jade", "wood"}
	for _, validStyle := range validStyles {
		if style == validStyle {
			return style
		}
	}
	return "marble"
}

// getDefaultPrivacyLevel 获取默认隐私级别
func (s *MemorialService) getDefaultPrivacyLevel(level int) int {
	if level == 1 || level == 2 {
		return level
	}
	return 1 // 默认家族可见
}

// GetTombstoneStyles 获取可用的墓碑样式列表
func (s *MemorialService) GetTombstoneStyles() []TombstoneStyle {
	return []TombstoneStyle{
		{
			ID:          "marble",
			Name:        "汉白玉",
			Description: "经典汉白玉材质，庄重典雅",
			PreviewURL:  "/static/tombstones/marble.jpg",
			Price:       0, // 免费
		},
		{
			ID:          "granite",
			Name:        "花岗岩",
			Description: "坚固耐用的花岗岩材质",
			PreviewURL:  "/static/tombstones/granite.jpg",
			Price:       0, // 免费
		},
		{
			ID:          "jade",
			Name:        "青玉",
			Description: "温润如玉，寓意美好",
			PreviewURL:  "/static/tombstones/jade.jpg",
			Price:       99, // 付费样式
		},
		{
			ID:          "wood",
			Name:        "檀木",
			Description: "天然檀木，古朴自然",
			PreviewURL:  "/static/tombstones/wood.jpg",
			Price:       199, // 付费样式
		},
		{
			ID:          "crystal",
			Name:        "水晶",
			Description: "透明水晶，现代简约",
			PreviewURL:  "/static/tombstones/crystal.jpg",
			Price:       299, // 高级付费样式
		},
	}
}

// GetThemeStyles 获取可用的主题风格列表
func (s *MemorialService) GetThemeStyles() []ThemeStyle {
	return []ThemeStyle{
		{
			ID:          "traditional",
			Name:        "中式传统",
			Description: "传统中式风格，庄重肃穆",
			PreviewURL:  "/static/themes/traditional.jpg",
			Price:       0,
		},
		{
			ID:          "elegant",
			Name:        "简约素雅",
			Description: "简约现代风格，素雅清新",
			PreviewURL:  "/static/themes/elegant.jpg",
			Price:       0,
		},
		{
			ID:          "natural",
			Name:        "自然清新",
			Description: "自然风光主题，清新怡人",
			PreviewURL:  "/static/themes/natural.jpg",
			Price:       0,
		},
		{
			ID:          "modern",
			Name:        "现代简约",
			Description: "现代简约风格，时尚大方",
			PreviewURL:  "/static/themes/modern.jpg",
			Price:       99,
		},
		{
			ID:          "luxury",
			Name:        "奢华典雅",
			Description: "奢华典雅风格，尊贵大气",
			PreviewURL:  "/static/themes/luxury.jpg",
			Price:       299,
		},
	}
}

// UpdateTombstoneStyle 更新墓碑样式
func (s *MemorialService) UpdateTombstoneStyle(userID, memorialID, styleID string) error {
	// 检查修改权限
	canModify, err := s.permissionManager.CanModifyMemorial(userID, memorialID)
	if err != nil {
		return err
	}
	if !canModify {
		return fmt.Errorf("无权修改此纪念馆")
	}

	// 验证样式是否存在
	styles := s.GetTombstoneStyles()
	var selectedStyle *TombstoneStyle
	for _, style := range styles {
		if style.ID == styleID {
			selectedStyle = &style
			break
		}
	}

	if selectedStyle == nil {
		return fmt.Errorf("墓碑样式不存在")
	}

	// 如果是付费样式，这里可以添加付费验证逻辑
	if selectedStyle.Price > 0 {
		// TODO: 验证用户是否已购买此样式
		// 暂时允许所有用户使用付费样式
	}

	// 更新墓碑样式
	if err := s.db.Model(&models.Memorial{}).
		Where("id = ?", memorialID).
		Update("tombstone_style", styleID).Error; err != nil {
		return fmt.Errorf("更新墓碑样式失败: %v", err)
	}

	return nil
}

// UpdateEpitaph 更新墓志铭
func (s *MemorialService) UpdateEpitaph(userID, memorialID, epitaph string) error {
	// 检查修改权限
	canModify, err := s.permissionManager.CanModifyMemorial(userID, memorialID)
	if err != nil {
		return err
	}
	if !canModify {
		return fmt.Errorf("无权修改此纪念馆")
	}

	// 验证墓志铭长度
	if len(epitaph) > 500 {
		return fmt.Errorf("墓志铭不能超过500个字符")
	}

	// 更新墓志铭
	if err := s.db.Model(&models.Memorial{}).
		Where("id = ?", memorialID).
		Update("epitaph", epitaph).Error; err != nil {
		return fmt.Errorf("更新墓志铭失败: %v", err)
	}

	return nil
}

// GenerateCalligraphy 生成书法字体墓志铭（模拟功能）
func (s *MemorialService) GenerateCalligraphy(text, fontStyle string) (*CalligraphyResult, error) {
	// 验证输入
	if text == "" {
		return nil, fmt.Errorf("文本不能为空")
	}
	if len(text) > 200 {
		return nil, fmt.Errorf("文本长度不能超过200个字符")
	}

	// 支持的字体样式
	validFonts := map[string]string{
		"kaishu":   "楷书",
		"xingshu":  "行书",
		"lishu":    "隶书",
		"caoshu":   "草书",
		"zhuanshu": "篆书",
	}

	fontName, exists := validFonts[fontStyle]
	if !exists {
		fontStyle = "kaishu"
		fontName = "楷书"
	}

	// 模拟生成书法图片
	// 在实际实现中，这里会调用书法生成服务或AI接口
	result := &CalligraphyResult{
		Text:      text,
		FontStyle: fontStyle,
		FontName:  fontName,
		ImageURL:  fmt.Sprintf("/api/v1/calligraphy/generate?text=%s&font=%s", text, fontStyle),
		CreatedAt: time.Now(),
	}

	return result, nil
}

// ProcessHandwritingImage 处理手写照片转换（模拟功能）
func (s *MemorialService) ProcessHandwritingImage(imageURL string) (*HandwritingResult, error) {
	if imageURL == "" {
		return nil, fmt.Errorf("图片URL不能为空")
	}

	// 模拟OCR识别和字体转换
	// 在实际实现中，这里会调用OCR服务和字体转换API
	result := &HandwritingResult{
		OriginalImageURL:  imageURL,
		RecognizedText:    "模拟识别的文字内容", // 实际应该是OCR识别结果
		ProcessedImageURL: "/api/v1/handwriting/processed/" + utils.GenerateUUID() + ".jpg",
		Confidence:        0.95, // 识别置信度
		CreatedAt:         time.Now(),
	}

	return result, nil
}

// 定义相关的数据结构
type TombstoneStyle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PreviewURL  string `json:"preview_url"`
	Price       int    `json:"price"` // 价格，0表示免费
}

type ThemeStyle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PreviewURL  string `json:"preview_url"`
	Price       int    `json:"price"`
}

type CalligraphyResult struct {
	Text      string    `json:"text"`
	FontStyle string    `json:"font_style"`
	FontName  string    `json:"font_name"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type HandwritingResult struct {
	OriginalImageURL  string    `json:"original_image_url"`
	RecognizedText    string    `json:"recognized_text"`
	ProcessedImageURL string    `json:"processed_image_url"`
	Confidence        float64   `json:"confidence"`
	CreatedAt         time.Time `json:"created_at"`
}

// GetRecentMemorials 获取用户最近访问的纪念馆
func (s *MemorialService) GetRecentMemorials(userID string, limit int) ([]MemorialListResponse, error) {
	var response []MemorialListResponse

	// 通过访客记录获取最近访问的纪念馆
	// 按访问时间倒序，去重
	var visitorRecords []models.VisitorRecord
	err := s.db.Where("visitor_id = ?", userID).
		Order("visit_time DESC").
		Limit(limit * 3). // 多查询一些，因为可能有重复和已删除的
		Find(&visitorRecords).Error

	if err != nil {
		return nil, fmt.Errorf("查询访客记录失败: %v", err)
	}

	// 去重并获取纪念馆信息
	memorialIDs := make([]string, 0)
	memorialIDSet := make(map[string]bool)

	for _, record := range visitorRecords {
		if !memorialIDSet[record.MemorialID] {
			memorialIDSet[record.MemorialID] = true
			memorialIDs = append(memorialIDs, record.MemorialID)
			if len(memorialIDs) >= limit {
				break
			}
		}
	}

	// 如果没有访问记录，返回用户创建的纪念馆
	if len(memorialIDs) == 0 {
		var memorials []models.Memorial
		err = s.db.Where("creator_id = ? AND status = ?", userID, 1).
			Order("created_at DESC").
			Limit(limit).
			Find(&memorials).Error

		if err != nil {
			return nil, fmt.Errorf("查询纪念馆失败: %v", err)
		}

		for _, memorial := range memorials {
			// 获取统计信息
			var visitorCount, worshipCount int64
			s.db.Model(&models.VisitorRecord{}).Where("memorial_id = ?", memorial.ID).Count(&visitorCount)
			s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorial.ID).Count(&worshipCount)

			response = append(response, MemorialListResponse{
				ID:           memorial.ID,
				DeceasedName: memorial.DeceasedName,
				BirthDate:    memorial.BirthDate,
				DeathDate:    memorial.DeathDate,
				AvatarURL:    memorial.AvatarURL,
				ThemeStyle:   memorial.ThemeStyle,
				PrivacyLevel: memorial.PrivacyLevel,
				CreatedAt:    memorial.CreatedAt,
				VisitorCount: visitorCount,
				WorshipCount: worshipCount,
			})
		}

		return response, nil
	}

	// 批量获取纪念馆信息
	var memorials []models.Memorial
	err = s.db.Where("id IN ? AND status = ?", memorialIDs, 1).
		Find(&memorials).Error

	if err != nil {
		return nil, fmt.Errorf("查询纪念馆失败: %v", err)
	}

	// 按访问顺序排序
	memorialMap := make(map[string]models.Memorial)
	for _, memorial := range memorials {
		memorialMap[memorial.ID] = memorial
	}

	for _, memorialID := range memorialIDs {
		memorial, exists := memorialMap[memorialID]
		if !exists {
			continue // 跳过已删除的纪念馆
		}

		// 获取统计信息
		var visitorCount, worshipCount int64
		s.db.Model(&models.VisitorRecord{}).Where("memorial_id = ?", memorial.ID).Count(&visitorCount)
		s.db.Model(&models.WorshipRecord{}).Where("memorial_id = ?", memorial.ID).Count(&worshipCount)

		response = append(response, MemorialListResponse{
			ID:           memorial.ID,
			DeceasedName: memorial.DeceasedName,
			BirthDate:    memorial.BirthDate,
			DeathDate:    memorial.DeathDate,
			AvatarURL:    memorial.AvatarURL,
			ThemeStyle:   memorial.ThemeStyle,
			PrivacyLevel: memorial.PrivacyLevel,
			CreatedAt:    memorial.CreatedAt,
			VisitorCount: visitorCount,
			WorshipCount: worshipCount,
		})
	}

	return response, nil
}
