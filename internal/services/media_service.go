package services

import (
	"fmt"
	"mime/multipart"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/utils"

	"gorm.io/gorm"
)

type MediaService struct {
	db                *gorm.DB
	fileUploadManager *utils.FileUploadManager
	permissionManager *utils.PermissionManager
}

type UploadMediaRequest struct {
	MemorialID  string `form:"memorial_id"` // 可选，创建纪念馆时上传头像不需要
	Description string `form:"description"`
}

type MediaFileResponse struct {
	ID          string `json:"id"`
	MemorialID  string `json:"memorial_id"`
	FileType    string `json:"file_type"`
	FileURL     string `json:"file_url"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func NewMediaService(db *gorm.DB, uploadDir string) *MediaService {
	return &MediaService{
		db:                db,
		fileUploadManager: utils.NewFileUploadManager(uploadDir, 100*1024*1024), // 100MB
		permissionManager: utils.NewPermissionManager(db),
	}
}

// UploadImage 上传图片
func (s *MediaService) UploadImage(userID string, req *UploadMediaRequest, file *multipart.FileHeader) (*models.MediaFile, error) {
	// 如果提供了 memorial_id，检查纪念馆访问权限
	if req.MemorialID != "" {
		canAccess, _, err := s.permissionManager.CanAccessMemorial(userID, req.MemorialID)
		if err != nil {
			return nil, err
		}
		if !canAccess {
			return nil, fmt.Errorf("无权访问此纪念馆")
		}
	}

	// 验证图片文件
	if err := utils.ValidateImageFile(file); err != nil {
		return nil, err
	}

	// 确定上传路径
	var uploadPath string
	if req.MemorialID != "" {
		uploadPath = fmt.Sprintf("memorials/%s/images", req.MemorialID)
	} else {
		// 临时上传（如创建纪念馆时的头像）
		uploadPath = fmt.Sprintf("temp/%s/images", userID)
	}

	// 上传文件
	relativePath, err := s.fileUploadManager.UploadFile(file, uploadPath)
	if err != nil {
		return nil, err
	}

	// 构建返回的文件信息（不保存到数据库，因为纪念馆还未创建）
	mediaFile := &models.MediaFile{
		ID:          utils.GenerateUUID(),
		MemorialID:  req.MemorialID,
		FileType:    "image",
		FileURL:     s.fileUploadManager.GetFileURL(relativePath),
		FileName:    file.Filename,
		FileSize:    file.Size,
		Description: req.Description,
	}

	// 只有在提供了 memorial_id 时才保存到数据库
	if req.MemorialID != "" {
		if err := s.db.Create(mediaFile).Error; err != nil {
			// 如果数据库保存失败，删除已上传的文件
			s.fileUploadManager.DeleteFile(relativePath)
			return nil, fmt.Errorf("保存文件信息失败: %v", err)
		}
	}

	return mediaFile, nil
}

// UploadVideo 上传视频
func (s *MediaService) UploadVideo(userID string, req *UploadMediaRequest, file *multipart.FileHeader) (*models.MediaFile, error) {
	// 如果提供了 memorial_id，检查纪念馆访问权限
	if req.MemorialID != "" {
		canAccess, _, err := s.permissionManager.CanAccessMemorial(userID, req.MemorialID)
		if err != nil {
			return nil, err
		}
		if !canAccess {
			return nil, fmt.Errorf("无权访问此纪念馆")
		}
	}

	// 验证视频文件
	if err := utils.ValidateVideoFile(file); err != nil {
		return nil, err
	}

	// 确定上传路径
	var uploadPath string
	if req.MemorialID != "" {
		uploadPath = fmt.Sprintf("memorials/%s/videos", req.MemorialID)
	} else {
		uploadPath = fmt.Sprintf("temp/%s/videos", userID)
	}

	// 上传文件
	relativePath, err := s.fileUploadManager.UploadFile(file, uploadPath)
	if err != nil {
		return nil, err
	}

	// 构建返回的文件信息
	mediaFile := &models.MediaFile{
		ID:          utils.GenerateUUID(),
		MemorialID:  req.MemorialID,
		FileType:    "video",
		FileURL:     s.fileUploadManager.GetFileURL(relativePath),
		FileName:    file.Filename,
		FileSize:    file.Size,
		Description: req.Description,
	}

	// 只有在提供了 memorial_id 时才保存到数据库
	if req.MemorialID != "" {
		if err := s.db.Create(mediaFile).Error; err != nil {
			// 如果数据库保存失败，删除已上传的文件
			s.fileUploadManager.DeleteFile(relativePath)
			return nil, fmt.Errorf("保存文件信息失败: %v", err)
		}
	}

	return mediaFile, nil
}

// UploadAudio 上传音频
func (s *MediaService) UploadAudio(userID string, req *UploadMediaRequest, file *multipart.FileHeader) (*models.MediaFile, error) {
	// 如果提供了 memorial_id，检查纪念馆访问权限
	if req.MemorialID != "" {
		canAccess, _, err := s.permissionManager.CanAccessMemorial(userID, req.MemorialID)
		if err != nil {
			return nil, err
		}
		if !canAccess {
			return nil, fmt.Errorf("无权访问此纪念馆")
		}
	}

	// 验证音频文件
	if err := utils.ValidateAudioFile(file); err != nil {
		return nil, err
	}

	// 确定上传路径
	var uploadPath string
	if req.MemorialID != "" {
		uploadPath = fmt.Sprintf("memorials/%s/audios", req.MemorialID)
	} else {
		uploadPath = fmt.Sprintf("temp/%s/audios", userID)
	}

	// 上传文件
	relativePath, err := s.fileUploadManager.UploadFile(file, uploadPath)
	if err != nil {
		return nil, err
	}

	// 构建返回的文件信息
	mediaFile := &models.MediaFile{
		ID:          utils.GenerateUUID(),
		MemorialID:  req.MemorialID,
		FileType:    "audio",
		FileURL:     s.fileUploadManager.GetFileURL(relativePath),
		FileName:    file.Filename,
		FileSize:    file.Size,
		Description: req.Description,
	}

	// 只有在提供了 memorial_id 时才保存到数据库
	if req.MemorialID != "" {
		if err := s.db.Create(mediaFile).Error; err != nil {
			// 如果数据库保存失败，删除已上传的文件
			s.fileUploadManager.DeleteFile(relativePath)
			return nil, fmt.Errorf("保存文件信息失败: %v", err)
		}
	}

	return mediaFile, nil
}

// GetMediaFiles 获取纪念馆媒体文件列表
func (s *MediaService) GetMediaFiles(userID, memorialID string, fileType string, page, pageSize int) ([]MediaFileResponse, int64, error) {
	// 检查纪念馆访问权限
	canAccess, _, err := s.permissionManager.CanAccessMemorial(userID, memorialID)
	if err != nil {
		return nil, 0, err
	}
	if !canAccess {
		return nil, 0, fmt.Errorf("无权访问此纪念馆")
	}

	var mediaFiles []models.MediaFile
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.MediaFile{}).Where("memorial_id = ?", memorialID)
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&mediaFiles).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询媒体文件失败: %v", err)
	}

	// 转换为响应格式
	var response []MediaFileResponse
	for _, file := range mediaFiles {
		response = append(response, MediaFileResponse{
			ID:          file.ID,
			MemorialID:  file.MemorialID,
			FileType:    file.FileType,
			FileURL:     file.FileURL,
			FileName:    file.FileName,
			FileSize:    file.FileSize,
			Description: file.Description,
			CreatedAt:   file.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, total, nil
}

// DeleteMediaFile 删除媒体文件
func (s *MediaService) DeleteMediaFile(userID, fileID string) error {
	// 查询文件信息
	var mediaFile models.MediaFile
	if err := s.db.Where("id = ?", fileID).First(&mediaFile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("文件不存在")
		}
		return fmt.Errorf("查询文件失败: %v", err)
	}

	// 检查纪念馆修改权限
	canModify, err := s.permissionManager.CanModifyMemorial(userID, mediaFile.MemorialID)
	if err != nil {
		return err
	}
	if !canModify {
		return fmt.Errorf("无权删除此文件")
	}

	// 从数据库删除记录
	if err := s.db.Delete(&mediaFile).Error; err != nil {
		return fmt.Errorf("删除文件记录失败: %v", err)
	}

	// 删除物理文件（从URL中提取相对路径）
	// 这里需要根据实际的URL格式来提取相对路径
	// 暂时跳过物理文件删除，避免路径解析错误

	return nil
}

// UpdateMediaFile 更新媒体文件信息
func (s *MediaService) UpdateMediaFile(userID, fileID, description string) error {
	// 查询文件信息
	var mediaFile models.MediaFile
	if err := s.db.Where("id = ?", fileID).First(&mediaFile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("文件不存在")
		}
		return fmt.Errorf("查询文件失败: %v", err)
	}

	// 检查纪念馆修改权限
	canModify, err := s.permissionManager.CanModifyMemorial(userID, mediaFile.MemorialID)
	if err != nil {
		return err
	}
	if !canModify {
		return fmt.Errorf("无权修改此文件")
	}

	// 更新描述
	if err := s.db.Model(&mediaFile).Update("description", description).Error; err != nil {
		return fmt.Errorf("更新文件信息失败: %v", err)
	}

	return nil
}

// GetMediaFileStats 获取纪念馆媒体文件统计
func (s *MediaService) GetMediaFileStats(userID, memorialID string) (map[string]int64, error) {
	// 检查纪念馆访问权限
	canAccess, _, err := s.permissionManager.CanAccessMemorial(userID, memorialID)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, fmt.Errorf("无权访问此纪念馆")
	}

	stats := make(map[string]int64)

	// 统计各类型文件数量
	fileTypes := []string{"image", "video", "audio"}
	for _, fileType := range fileTypes {
		var count int64
		s.db.Model(&models.MediaFile{}).
			Where("memorial_id = ? AND file_type = ?", memorialID, fileType).
			Count(&count)
		stats[fileType] = count
	}

	// 统计总数
	var total int64
	s.db.Model(&models.MediaFile{}).Where("memorial_id = ?", memorialID).Count(&total)
	stats["total"] = total

	return stats, nil
}
