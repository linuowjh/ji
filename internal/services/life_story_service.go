package services

import (
	"errors"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LifeStoryService struct {
	db *gorm.DB
}

func NewLifeStoryService(db *gorm.DB) *LifeStoryService {
	return &LifeStoryService{
		db: db,
	}
}

// 创建生平故事请求
type CreateLifeStoryRequest struct {
	Title     string     `json:"title" binding:"required"`
	Content   string     `json:"content" binding:"required"`
	StoryDate *time.Time `json:"story_date"`
	AgeAtTime int        `json:"age_at_time"`
	Location  string     `json:"location"`
	Category  string     `json:"category"`
	IsPublic  bool       `json:"is_public"`
	Media     []MediaItem `json:"media"`
}

type MediaItem struct {
	MediaType    string `json:"media_type" binding:"required"`
	MediaURL     string `json:"media_url" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url"`
	Caption      string `json:"caption"`
}

// 更新生平故事请求
type UpdateLifeStoryRequest struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	StoryDate *time.Time `json:"story_date"`
	AgeAtTime int        `json:"age_at_time"`
	Location  string     `json:"location"`
	Category  string     `json:"category"`
	IsPublic  *bool      `json:"is_public"`
}

// 创建时间轴事件请求
type CreateTimelineRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	EventType   string    `json:"event_type"`
	IconURL     string    `json:"icon_url"`
	IsPublic    bool      `json:"is_public"`
}

// 创建生平故事
func (s *LifeStoryService) CreateLifeStory(userID, memorialID string, req *CreateLifeStoryRequest) (*models.LifeStory, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	tx := s.db.Begin()

	story := &models.LifeStory{
		ID:         uuid.New().String(),
		MemorialID: memorialID,
		Title:      req.Title,
		Content:    req.Content,
		StoryDate:  req.StoryDate,
		AgeAtTime:  req.AgeAtTime,
		Location:   req.Location,
		Category:   req.Category,
		IsPublic:   req.IsPublic,
		AuthorID:   userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := tx.Create(story).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加媒体文件
	for i, media := range req.Media {
		storyMedia := &models.LifeStoryMedia{
			ID:           uuid.New().String(),
			LifeStoryID:  story.ID,
			MediaType:    media.MediaType,
			MediaURL:     media.MediaURL,
			ThumbnailURL: media.ThumbnailURL,
			Caption:      media.Caption,
			SortOrder:    i,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := tx.Create(storyMedia).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新查询包含媒体文件的故事
	s.db.Preload("Media", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("Author").First(story, "id = ?", story.ID)

	return story, nil
}

// 获取生平故事列表
func (s *LifeStoryService) GetLifeStories(userID, memorialID string, page, pageSize int) ([]*models.LifeStory, int64, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, 0, err
	}

	var stories []*models.LifeStory
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.LifeStory{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 查询故事列表
	err := s.db.Preload("Author").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("memorial_id = ?", memorialID).
		Order("story_date DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&stories).Error

	return stories, total, err
}

// 获取生平故事详情
func (s *LifeStoryService) GetLifeStory(userID, storyID string) (*models.LifeStory, error) {
	var story models.LifeStory
	err := s.db.Preload("Author").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&story, "id = ?", storyID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("生平故事不存在")
		}
		return nil, err
	}

	// 验证访问权限
	if err := s.validateMemorialAccess(userID, story.MemorialID); err != nil {
		return nil, err
	}

	return &story, nil
}

// 更新生平故事
func (s *LifeStoryService) UpdateLifeStory(userID, storyID string, req *UpdateLifeStoryRequest) error {
	var story models.LifeStory
	err := s.db.First(&story, "id = ?", storyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("生平故事不存在")
		}
		return err
	}

	// 验证权限（只有作者或纪念馆创建者可以修改）
	if story.AuthorID != userID {
		if err := s.validateMemorialOwnership(userID, story.MemorialID); err != nil {
			return errors.New("无权修改此故事")
		}
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.StoryDate != nil {
		updates["story_date"] = req.StoryDate
	}
	if req.AgeAtTime > 0 {
		updates["age_at_time"] = req.AgeAtTime
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&story).Updates(updates).Error
}

// 删除生平故事
func (s *LifeStoryService) DeleteLifeStory(userID, storyID string) error {
	var story models.LifeStory
	err := s.db.First(&story, "id = ?", storyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("生平故事不存在")
		}
		return err
	}

	// 验证权限（只有作者或纪念馆创建者可以删除）
	if story.AuthorID != userID {
		if err := s.validateMemorialOwnership(userID, story.MemorialID); err != nil {
			return errors.New("无权删除此故事")
		}
	}

	// 删除故事及其媒体文件
	tx := s.db.Begin()
	
	// 删除媒体文件
	if err := tx.Where("life_story_id = ?", storyID).Delete(&models.LifeStoryMedia{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除故事
	if err := tx.Delete(&story).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 创建时间轴事件
func (s *LifeStoryService) CreateTimeline(userID, memorialID string, req *CreateTimelineRequest) (*models.Timeline, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	timeline := &models.Timeline{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
		EventType:   req.EventType,
		IconURL:     req.IconURL,
		IsPublic:    req.IsPublic,
		AuthorID:    userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(timeline).Error; err != nil {
		return nil, err
	}

	// 重新查询包含作者信息的时间轴
	s.db.Preload("Author").First(timeline, "id = ?", timeline.ID)

	return timeline, nil
}

// 获取时间轴
func (s *LifeStoryService) GetTimeline(userID, memorialID string) ([]*models.Timeline, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	var timelines []*models.Timeline

	err := s.db.Preload("Author").
		Where("memorial_id = ?", memorialID).
		Order("event_date ASC").
		Find(&timelines).Error

	return timelines, err
}

// 删除时间轴事件
func (s *LifeStoryService) DeleteTimeline(userID, timelineID string) error {
	var timeline models.Timeline
	err := s.db.First(&timeline, "id = ?", timelineID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("时间轴事件不存在")
		}
		return err
	}

	// 验证权限（只有作者或纪念馆创建者可以删除）
	if timeline.AuthorID != userID {
		if err := s.validateMemorialOwnership(userID, timeline.MemorialID); err != nil {
			return errors.New("无权删除此事件")
		}
	}

	return s.db.Delete(&timeline).Error
}

// 按分类获取生平故事
func (s *LifeStoryService) GetStoriesByCategory(userID, memorialID, category string) ([]*models.LifeStory, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	var stories []*models.LifeStory

	query := s.db.Preload("Author").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("memorial_id = ?", memorialID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Order("story_date DESC, created_at DESC").Find(&stories).Error

	return stories, err
}

// 验证纪念馆访问权限
func (s *LifeStoryService) validateMemorialAccess(userID, memorialID string) error {
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

// 验证纪念馆所有权
func (s *LifeStoryService) validateMemorialOwnership(userID, memorialID string) error {
	var memorial models.Memorial
	err := s.db.First(&memorial, "id = ? AND creator_id = ? AND status = ?", memorialID, userID, 1).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权操作此纪念馆")
		}
		return err
	}
	return nil
}