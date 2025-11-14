package services

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExclusiveServiceService struct {
	db         *gorm.DB
	exportPath string
}

func NewExclusiveServiceService(db *gorm.DB, exportPath string) *ExclusiveServiceService {
	// 确保导出目录存在
	if exportPath == "" {
		exportPath = "./exports"
	}
	os.MkdirAll(exportPath, 0755)
	
	return &ExclusiveServiceService{
		db:         db,
		exportPath: exportPath,
	}
}

// 专属服务管理

// GetExclusiveServices 获取专属服务列表
func (s *ExclusiveServiceService) GetExclusiveServices(serviceType string, activeOnly bool) ([]models.ExclusiveService, error) {
	var services []models.ExclusiveService
	query := s.db.Model(&models.ExclusiveService{})
	
	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Order("sort_order ASC, base_price ASC").Find(&services).Error
	return services, err
}

// GetExclusiveService 获取专属服务详情
func (s *ExclusiveServiceService) GetExclusiveService(serviceID string) (*models.ExclusiveService, error) {
	var service models.ExclusiveService
	err := s.db.Where("id = ?", serviceID).First(&service).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务不存在")
		}
		return nil, err
	}
	return &service, nil
}

// 服务预订管理

// CreateBooking 创建服务预订
func (s *ExclusiveServiceService) CreateBooking(booking *models.ServiceBooking) error {
	// 获取服务信息
	service, err := s.GetExclusiveService(booking.ServiceID)
	if err != nil {
		return err
	}
	
	if !service.IsActive {
		return errors.New("该服务暂不可用")
	}
	
	// 计算价格
	totalPrice := service.BasePrice
	if service.PriceUnit == "per_hour" && booking.Duration > 0 {
		hours := float64(booking.Duration) / 60.0
		totalPrice = service.BasePrice * hours
	}
	
	booking.ID = uuid.New().String()
	booking.Status = "pending"
	booking.TotalPrice = totalPrice
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	
	return s.db.Create(booking).Error
}

// GetUserBookings 获取用户预订列表
func (s *ExclusiveServiceService) GetUserBookings(userID string, status string) ([]models.ServiceBooking, error) {
	var bookings []models.ServiceBooking
	query := s.db.Where("user_id = ?", userID).Preload("Service")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

// GetBooking 获取预订详情
func (s *ExclusiveServiceService) GetBooking(bookingID string) (*models.ServiceBooking, error) {
	var booking models.ServiceBooking
	err := s.db.Where("id = ?", bookingID).Preload("Service").Preload("User").First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预订不存在")
		}
		return nil, err
	}
	return &booking, nil
}

// UpdateBookingStatus 更新预订状态
func (s *ExclusiveServiceService) UpdateBookingStatus(bookingID, status string, staffID string) error {
	var booking models.ServiceBooking
	err := s.db.Where("id = ?", bookingID).First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("预订不存在")
		}
		return err
	}
	
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	
	if staffID != "" {
		updates["staff_id"] = staffID
	}
	
	if status == "completed" {
		now := time.Now()
		updates["completed_at"] = now
	}
	
	return s.db.Model(&booking).Updates(updates).Error
}

// CancelBooking 取消预订
func (s *ExclusiveServiceService) CancelBooking(bookingID, userID, reason string) error {
	var booking models.ServiceBooking
	err := s.db.Where("id = ? AND user_id = ?", bookingID, userID).First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("预订不存在")
		}
		return err
	}
	
	if booking.Status == "completed" || booking.Status == "cancelled" {
		return errors.New("无法取消该预订")
	}
	
	now := time.Now()
	return s.db.Model(&booking).Updates(map[string]interface{}{
		"status":              "cancelled",
		"cancelled_at":        now,
		"cancellation_reason": reason,
		"updated_at":          now,
	}).Error
}

// 数据导出服务

// CreateDataExportRequest 创建数据导出请求
func (s *ExclusiveServiceService) CreateDataExportRequest(req *models.DataExportRequest) error {
	req.ID = uuid.New().String()
	req.Status = "pending"
	req.CreatedAt = time.Now()
	
	if err := s.db.Create(req).Error; err != nil {
		return err
	}
	
	// 异步处理导出
	go s.processDataExport(req)
	
	return nil
}

// processDataExport 处理数据导出
func (s *ExclusiveServiceService) processDataExport(req *models.DataExportRequest) {
	// 更新状态为处理中
	s.db.Model(req).Update("status", "processing")
	
	var err error
	var filePath string
	
	switch req.ExportType {
	case "full":
		filePath, err = s.exportFullUserData(req)
	case "memorial":
		filePath, err = s.exportMemorialData(req)
	case "family":
		filePath, err = s.exportFamilyData(req)
	default:
		err = errors.New("不支持的导出类型")
	}
	
	if err != nil {
		// 导出失败
		s.db.Model(req).Updates(map[string]interface{}{
			"status":        "failed",
			"error_message": err.Error(),
		})
		return
	}
	
	// 获取文件大小
	fileInfo, _ := os.Stat(filePath)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	
	// 设置过期时间（7天后）
	expiresAt := time.Now().AddDate(0, 0, 7)
	completedAt := time.Now()
	
	// 导出成功
	s.db.Model(req).Updates(map[string]interface{}{
		"status":       "completed",
		"file_path":    filePath,
		"file_size":    fileSize,
		"expires_at":   expiresAt,
		"completed_at": completedAt,
	})
}

// exportFullUserData 导出用户完整数据
func (s *ExclusiveServiceService) exportFullUserData(req *models.DataExportRequest) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("user_data_%s_%s.zip", req.UserID, timestamp)
	filePath := filepath.Join(s.exportPath, filename)
	
	// 创建ZIP文件
	zipFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建导出文件失败: %v", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 导出用户信息
	var user models.User
	s.db.Where("id = ?", req.UserID).First(&user)
	s.addJSONToZip(zipWriter, "user.json", user)
	
	// 导出纪念馆
	var memorials []models.Memorial
	s.db.Where("creator_id = ?", req.UserID).Find(&memorials)
	s.addJSONToZip(zipWriter, "memorials.json", memorials)
	
	// 导出祭扫记录
	var worshipRecords []models.WorshipRecord
	s.db.Where("user_id = ?", req.UserID).Find(&worshipRecords)
	s.addJSONToZip(zipWriter, "worship_records.json", worshipRecords)
	
	// 导出家族信息
	var familyMembers []models.FamilyMember
	s.db.Where("user_id = ?", req.UserID).Preload("Family").Find(&familyMembers)
	s.addJSONToZip(zipWriter, "families.json", familyMembers)
	
	// 导出订阅信息
	var subscriptions []models.UserSubscription
	s.db.Where("user_id = ?", req.UserID).Preload("Package").Find(&subscriptions)
	s.addJSONToZip(zipWriter, "subscriptions.json", subscriptions)
	
	return filePath, nil
}

// exportMemorialData 导出纪念馆数据
func (s *ExclusiveServiceService) exportMemorialData(req *models.DataExportRequest) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("memorial_data_%s_%s.zip", req.TargetID, timestamp)
	filePath := filepath.Join(s.exportPath, filename)
	
	// 创建ZIP文件
	zipFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建导出文件失败: %v", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 导出纪念馆信息
	var memorial models.Memorial
	s.db.Where("id = ?", req.TargetID).First(&memorial)
	s.addJSONToZip(zipWriter, "memorial.json", memorial)
	
	// 导出媒体文件列表
	if req.IncludeMedia {
		var mediaFiles []models.MediaFile
		s.db.Where("memorial_id = ?", req.TargetID).Find(&mediaFiles)
		s.addJSONToZip(zipWriter, "media_files.json", mediaFiles)
	}
	
	// 导出祭扫记录
	var worshipRecords []models.WorshipRecord
	s.db.Where("memorial_id = ?", req.TargetID).Find(&worshipRecords)
	s.addJSONToZip(zipWriter, "worship_records.json", worshipRecords)
	
	// 导出留言和祈福
	var messages []models.Message
	s.db.Where("memorial_id = ?", req.TargetID).Find(&messages)
	s.addJSONToZip(zipWriter, "messages.json", messages)
	
	return filePath, nil
}

// exportFamilyData 导出家族数据
func (s *ExclusiveServiceService) exportFamilyData(req *models.DataExportRequest) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("family_data_%s_%s.zip", req.TargetID, timestamp)
	filePath := filepath.Join(s.exportPath, filename)
	
	// 创建ZIP文件
	zipFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建导出文件失败: %v", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 导出家族信息
	var family models.Family
	s.db.Where("id = ?", req.TargetID).First(&family)
	s.addJSONToZip(zipWriter, "family.json", family)
	
	// 导出家族成员
	var members []models.FamilyMember
	s.db.Where("family_id = ?", req.TargetID).Preload("User").Find(&members)
	s.addJSONToZip(zipWriter, "members.json", members)
	
	// 导出家族故事
	var stories []models.FamilyStory
	s.db.Where("family_id = ?", req.TargetID).Find(&stories)
	s.addJSONToZip(zipWriter, "stories.json", stories)
	
	return filePath, nil
}

// addJSONToZip 将数据以JSON格式添加到ZIP文件
func (s *ExclusiveServiceService) addJSONToZip(zipWriter *zip.Writer, filename string, data interface{}) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("创建ZIP条目失败: %v", err)
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}
	
	if _, err := writer.Write(jsonData); err != nil {
		return fmt.Errorf("写入ZIP失败: %v", err)
	}
	
	return nil
}

// GetUserExportRequests 获取用户导出请求列表
func (s *ExclusiveServiceService) GetUserExportRequests(userID string) ([]models.DataExportRequest, error) {
	var requests []models.DataExportRequest
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// GetExportRequest 获取导出请求详情
func (s *ExclusiveServiceService) GetExportRequest(requestID string) (*models.DataExportRequest, error) {
	var request models.DataExportRequest
	err := s.db.Where("id = ?", requestID).First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("导出请求不存在")
		}
		return nil, err
	}
	return &request, nil
}

// DownloadExport 下载导出文件
func (s *ExclusiveServiceService) DownloadExport(requestID, userID string) (string, error) {
	var request models.DataExportRequest
	err := s.db.Where("id = ? AND user_id = ?", requestID, userID).First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("导出请求不存在")
		}
		return "", err
	}
	
	if request.Status != "completed" {
		return "", errors.New("导出尚未完成")
	}
	
	// 检查是否过期
	if request.ExpiresAt != nil && request.ExpiresAt.Before(time.Now()) {
		return "", errors.New("下载链接已过期")
	}
	
	// 检查文件是否存在
	if _, err := os.Stat(request.FilePath); os.IsNotExist(err) {
		return "", errors.New("导出文件不存在")
	}
	
	return request.FilePath, nil
}

// 老照片修复服务

// CreatePhotoRestoreRequest 创建老照片修复请求
func (s *ExclusiveServiceService) CreatePhotoRestoreRequest(req *models.PhotoRestoreRequest) error {
	req.ID = uuid.New().String()
	req.Status = "pending"
	req.CreatedAt = time.Now()
	
	if err := s.db.Create(req).Error; err != nil {
		return err
	}
	
	// 异步处理修复（这里简化处理，实际应该调用AI服务）
	go s.processPhotoRestore(req)
	
	return nil
}

// processPhotoRestore 处理照片修复
func (s *ExclusiveServiceService) processPhotoRestore(req *models.PhotoRestoreRequest) {
	// 更新状态为处理中
	s.db.Model(req).Update("status", "processing")
	
	// 模拟处理时间
	time.Sleep(5 * time.Second)
	
	// 这里应该调用实际的AI修复服务
	// 暂时使用原图URL作为修复后的URL
	restoredURL := req.OriginalPhotoURL + "?restored=true"
	
	completedAt := time.Now()
	processingTime := int(completedAt.Sub(req.CreatedAt).Seconds())
	
	// 更新为完成状态
	s.db.Model(req).Updates(map[string]interface{}{
		"status":            "completed",
		"restored_photo_url": restoredURL,
		"processing_time":   processingTime,
		"completed_at":      completedAt,
	})
}

// GetUserPhotoRestoreRequests 获取用户照片修复请求列表
func (s *ExclusiveServiceService) GetUserPhotoRestoreRequests(userID string) ([]models.PhotoRestoreRequest, error) {
	var requests []models.PhotoRestoreRequest
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// 服务评价

// CreateServiceReview 创建服务评价
func (s *ExclusiveServiceService) CreateServiceReview(review *models.ServiceReview) error {
	// 检查预订是否存在且已完成
	var booking models.ServiceBooking
	err := s.db.Where("id = ? AND user_id = ?", review.BookingID, review.UserID).First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("预订不存在")
		}
		return err
	}
	
	if booking.Status != "completed" {
		return errors.New("只能评价已完成的服务")
	}
	
	// 检查是否已评价
	var existingReview models.ServiceReview
	err = s.db.Where("booking_id = ?", review.BookingID).First(&existingReview).Error
	if err == nil {
		return errors.New("该服务已评价")
	}
	
	review.ID = uuid.New().String()
	review.CreatedAt = time.Now()
	
	return s.db.Create(review).Error
}

// GetServiceReviews 获取服务评价列表
func (s *ExclusiveServiceService) GetServiceReviews(serviceID string, page, pageSize int) ([]models.ServiceReview, int64, error) {
	var reviews []models.ServiceReview
	var total int64
	
	// 获取该服务的所有预订
	var bookingIDs []string
	s.db.Model(&models.ServiceBooking{}).
		Where("service_id = ?", serviceID).
		Pluck("id", &bookingIDs)
	
	query := s.db.Where("booking_id IN ?", bookingIDs).Preload("User")
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error
	
	return reviews, total, err
}

// 初始化默认专属服务
func (s *ExclusiveServiceService) InitDefaultServices() error {
	defaultServices := []models.ExclusiveService{
		{
			ID:             uuid.New().String(),
			ServiceName:    "专属追思会策划",
			ServiceType:    "memorial_service",
			Description:    "专业团队为您策划和主持线上追思会，提供全程技术支持",
			BasePrice:      299.00,
			PriceUnit:      "per_service",
			Features:       `["专业主持人", "定制流程", "技术支持", "录制回放", "最多50人参与"]`,
			RequireBooking: true,
			MaxDuration:    180,
			IsActive:       true,
			SortOrder:      1,
		},
		{
			ID:             uuid.New().String(),
			ServiceName:    "数据备份导出服务",
			ServiceType:    "data_backup",
			Description:    "将您的纪念馆数据导出为加密文件，永久保存",
			BasePrice:      49.00,
			PriceUnit:      "per_service",
			Features:       `["完整数据导出", "加密保护", "多种格式", "7天有效期"]`,
			RequireBooking: false,
			IsActive:       true,
			SortOrder:      2,
		},
		{
			ID:             uuid.New().String(),
			ServiceName:    "老照片AI修复",
			ServiceType:    "photo_restore",
			Description:    "使用AI技术修复老照片，包括上色、增强、修补等",
			BasePrice:      19.00,
			PriceUnit:      "per_service",
			Features:       `["AI智能修复", "色彩还原", "清晰度增强", "划痕修复"]`,
			RequireBooking: false,
			IsActive:       true,
			SortOrder:      3,
		},
		{
			ID:             uuid.New().String(),
			ServiceName:    "定制设计服务",
			ServiceType:    "custom_design",
			Description:    "专业设计师为您定制独一无二的纪念馆主题和墓碑",
			BasePrice:      599.00,
			PriceUnit:      "per_service",
			Features:       `["专业设计师", "一对一沟通", "3次修改机会", "源文件交付"]`,
			RequireBooking: true,
			IsActive:       true,
			SortOrder:      4,
		},
	}
	
	for _, service := range defaultServices {
		// 检查是否已存在
		var existing models.ExclusiveService
		err := s.db.Where("service_name = ?", service.ServiceName).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			if err := s.db.Create(&service).Error; err != nil {
				return fmt.Errorf("创建默认服务失败: %v", err)
			}
		}
	}
	
	return nil
}
