package services

import (
	"errors"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumService struct {
	db *gorm.DB
}

func NewAlbumService(db *gorm.DB) *AlbumService {
	return &AlbumService{
		db: db,
	}
}

// 创建相册请求
type CreateAlbumRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	IsPublic    bool   `json:"is_public"`
}

// 更新相册请求
type UpdateAlbumRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	IsPublic    *bool  `json:"is_public"`
}

// 添加照片请求
type AddPhotoRequest struct {
	PhotoURL     string     `json:"photo_url" binding:"required"`
	ThumbnailURL string     `json:"thumbnail_url"`
	Caption      string     `json:"caption"`
	TakenDate    *time.Time `json:"taken_date"`
	Location     string     `json:"location"`
}

// 创建相册
func (s *AlbumService) CreateAlbum(userID, memorialID string, req *CreateAlbumRequest) (*models.Album, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, err
	}

	album := &models.Album{
		ID:          uuid.New().String(),
		MemorialID:  memorialID,
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		IsPublic:    req.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(album).Error; err != nil {
		return nil, err
	}

	return album, nil
}

// 获取相册列表
func (s *AlbumService) GetAlbums(userID, memorialID string, page, pageSize int) ([]*models.Album, int64, error) {
	// 验证纪念馆权限
	if err := s.validateMemorialAccess(userID, memorialID); err != nil {
		return nil, 0, err
	}

	var albums []*models.Album
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	s.db.Model(&models.Album{}).Where("memorial_id = ?", memorialID).Count(&total)

	// 查询相册列表
	err := s.db.Where("memorial_id = ?", memorialID).
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&albums).Error

	return albums, total, err
}

// 获取相册详情
func (s *AlbumService) GetAlbum(userID, albumID string) (*models.Album, error) {
	var album models.Album
	err := s.db.Preload("Photos", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, created_at ASC")
	}).First(&album, "id = ?", albumID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("相册不存在")
		}
		return nil, err
	}

	// 验证访问权限
	if err := s.validateMemorialAccess(userID, album.MemorialID); err != nil {
		return nil, err
	}

	return &album, nil
}

// 更新相册
func (s *AlbumService) UpdateAlbum(userID, albumID string, req *UpdateAlbumRequest) error {
	var album models.Album
	err := s.db.First(&album, "id = ?", albumID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("相册不存在")
		}
		return err
	}

	// 验证权限
	if err := s.validateMemorialAccess(userID, album.MemorialID); err != nil {
		return err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CoverURL != "" {
		updates["cover_url"] = req.CoverURL
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&album).Updates(updates).Error
}

// 删除相册
func (s *AlbumService) DeleteAlbum(userID, albumID string) error {
	var album models.Album
	err := s.db.First(&album, "id = ?", albumID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("相册不存在")
		}
		return err
	}

	// 验证权限
	if err := s.validateMemorialAccess(userID, album.MemorialID); err != nil {
		return err
	}

	// 删除相册及其照片
	tx := s.db.Begin()
	
	// 删除相册中的照片
	if err := tx.Where("album_id = ?", albumID).Delete(&models.AlbumPhoto{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除相册
	if err := tx.Delete(&album).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 添加照片到相册
func (s *AlbumService) AddPhoto(userID, albumID string, req *AddPhotoRequest) (*models.AlbumPhoto, error) {
	var album models.Album
	err := s.db.First(&album, "id = ?", albumID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("相册不存在")
		}
		return nil, err
	}

	// 验证权限
	if err := s.validateMemorialAccess(userID, album.MemorialID); err != nil {
		return nil, err
	}

	photo := &models.AlbumPhoto{
		ID:           uuid.New().String(),
		AlbumID:      albumID,
		PhotoURL:     req.PhotoURL,
		ThumbnailURL: req.ThumbnailURL,
		Caption:      req.Caption,
		TakenDate:    req.TakenDate,
		Location:     req.Location,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(photo).Error; err != nil {
		return nil, err
	}

	return photo, nil
}

// 删除照片
func (s *AlbumService) DeletePhoto(userID, photoID string) error {
	var photo models.AlbumPhoto
	err := s.db.Preload("Album").First(&photo, "id = ?", photoID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("照片不存在")
		}
		return err
	}

	// 验证权限
	if err := s.validateMemorialAccess(userID, photo.Album.MemorialID); err != nil {
		return err
	}

	return s.db.Delete(&photo).Error
}

// 更新照片信息
func (s *AlbumService) UpdatePhoto(userID, photoID string, req *AddPhotoRequest) error {
	var photo models.AlbumPhoto
	err := s.db.Preload("Album").First(&photo, "id = ?", photoID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("照片不存在")
		}
		return err
	}

	// 验证权限
	if err := s.validateMemorialAccess(userID, photo.Album.MemorialID); err != nil {
		return err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Caption != "" {
		updates["caption"] = req.Caption
	}
	if req.TakenDate != nil {
		updates["taken_date"] = req.TakenDate
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.ThumbnailURL != "" {
		updates["thumbnail_url"] = req.ThumbnailURL
	}
	updates["updated_at"] = time.Now()

	return s.db.Model(&photo).Updates(updates).Error
}

// 验证纪念馆访问权限
func (s *AlbumService) validateMemorialAccess(userID, memorialID string) error {
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