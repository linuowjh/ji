package services

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupService struct {
	db         *gorm.DB
	backupPath string
}

func NewBackupService(db *gorm.DB, backupPath string) *BackupService {
	// 确保备份目录存在
	if backupPath == "" {
		backupPath = "./backups"
	}
	os.MkdirAll(backupPath, 0755)
	
	return &BackupService{
		db:         db,
		backupPath: backupPath,
	}
}

// CreateBackup 创建数据备份
func (s *BackupService) CreateBackup(backupType, createdBy string) (*models.DataBackup, error) {
	backup := &models.DataBackup{
		ID:         uuid.New().String(),
		BackupType: backupType,
		Status:     "pending",
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}
	
	// 保存备份记录
	if err := s.db.Create(backup).Error; err != nil {
		return nil, fmt.Errorf("创建备份记录失败: %v", err)
	}
	
	// 异步执行备份
	go s.performBackup(backup)
	
	return backup, nil
}

// performBackup 执行备份操作
func (s *BackupService) performBackup(backup *models.DataBackup) {
	// 更新状态为处理中
	s.db.Model(backup).Update("status", "processing")
	
	var err error
	var backupFilePath string
	
	switch backup.BackupType {
	case "full":
		backupFilePath, err = s.createFullBackup()
	case "incremental":
		backupFilePath, err = s.createIncrementalBackup()
	case "user":
		backupFilePath, err = s.createUserDataBackup(backup.CreatedBy)
	default:
		err = errors.New("不支持的备份类型")
	}
	
	if err != nil {
		// 备份失败
		s.db.Model(backup).Updates(map[string]interface{}{
			"status":        "failed",
			"error_message": err.Error(),
		})
		return
	}
	
	// 获取文件大小
	fileInfo, _ := os.Stat(backupFilePath)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	
	// 备份成功
	completedAt := time.Now()
	s.db.Model(backup).Updates(map[string]interface{}{
		"status":       "completed",
		"backup_path":  backupFilePath,
		"file_size":    fileSize,
		"completed_at": completedAt,
	})
}

// createFullBackup 创建完整备份
func (s *BackupService) createFullBackup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("full_backup_%s.zip", timestamp)
	backupFilePath := filepath.Join(s.backupPath, filename)
	
	// 创建ZIP文件
	zipFile, err := os.Create(backupFilePath)
	if err != nil {
		return "", fmt.Errorf("创建备份文件失败: %v", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 备份所有表数据
	tables := []string{
		"users", "memorials", "worship_records", "families", "family_members",
		"media_files", "prayers", "messages", "memorial_reminders",
		"memorial_families", "visitor_records",
	}
	
	for _, table := range tables {
		if err := s.backupTable(zipWriter, table); err != nil {
			return "", fmt.Errorf("备份表 %s 失败: %v", table, err)
		}
	}
	
	return backupFilePath, nil
}

// createIncrementalBackup 创建增量备份（最近24小时的数据）
func (s *BackupService) createIncrementalBackup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("incremental_backup_%s.zip", timestamp)
	backupFilePath := filepath.Join(s.backupPath, filename)
	
	// 创建ZIP文件
	zipFile, err := os.Create(backupFilePath)
	if err != nil {
		return "", fmt.Errorf("创建备份文件失败: %v", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 获取24小时前的时间
	since := time.Now().Add(-24 * time.Hour)
	
	// 备份最近更新的数据
	tables := []string{
		"users", "memorials", "worship_records", "messages", "prayers",
	}
	
	for _, table := range tables {
		if err := s.backupTableSince(zipWriter, table, since); err != nil {
			return "", fmt.Errorf("备份表 %s 失败: %v", table, err)
		}
	}
	
	return backupFilePath, nil
}

// createUserDataBackup 创建用户数据备份
func (s *BackupService) createUserDataBackup(userID string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("user_backup_%s_%s.zip", userID, timestamp)
	backupFilePath := filepath.Join(s.backupPath, filename)
	
	// 创建ZIP文件
	zipFile, err := os.Create(backupFilePath)
	if err != nil {
		return "", fmt.Errorf("创建备份文件失败: %v", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 备份用户信息
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", fmt.Errorf("获取用户信息失败: %v", err)
	}
	
	if err := s.addJSONToZip(zipWriter, "user.json", user); err != nil {
		return "", err
	}
	
	// 备份用户创建的纪念馆
	var memorials []models.Memorial
	s.db.Where("creator_id = ?", userID).Find(&memorials)
	if err := s.addJSONToZip(zipWriter, "memorials.json", memorials); err != nil {
		return "", err
	}
	
	// 备份用户的祭扫记录
	var worshipRecords []models.WorshipRecord
	s.db.Where("user_id = ?", userID).Find(&worshipRecords)
	if err := s.addJSONToZip(zipWriter, "worship_records.json", worshipRecords); err != nil {
		return "", err
	}
	
	// 备份用户的家族信息
	var familyMembers []models.FamilyMember
	s.db.Where("user_id = ?", userID).Find(&familyMembers)
	if err := s.addJSONToZip(zipWriter, "family_members.json", familyMembers); err != nil {
		return "", err
	}
	
	return backupFilePath, nil
}

// backupTable 备份整个表
func (s *BackupService) backupTable(zipWriter *zip.Writer, tableName string) error {
	var data []map[string]interface{}
	
	// 查询表数据
	if err := s.db.Table(tableName).Find(&data).Error; err != nil {
		return err
	}
	
	// 添加到ZIP
	return s.addJSONToZip(zipWriter, fmt.Sprintf("%s.json", tableName), data)
}

// backupTableSince 备份表中指定时间之后的数据
func (s *BackupService) backupTableSince(zipWriter *zip.Writer, tableName string, since time.Time) error {
	var data []map[string]interface{}
	
	// 查询指定时间之后的数据
	if err := s.db.Table(tableName).Where("created_at >= ? OR updated_at >= ?", since, since).Find(&data).Error; err != nil {
		return err
	}
	
	// 如果没有数据，跳过
	if len(data) == 0 {
		return nil
	}
	
	// 添加到ZIP
	return s.addJSONToZip(zipWriter, fmt.Sprintf("%s.json", tableName), data)
}

// addJSONToZip 将数据以JSON格式添加到ZIP文件
func (s *BackupService) addJSONToZip(zipWriter *zip.Writer, filename string, data interface{}) error {
	// 创建ZIP文件条目
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("创建ZIP条目失败: %v", err)
	}
	
	// 将数据转换为JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}
	
	// 写入ZIP
	if _, err := writer.Write(jsonData); err != nil {
		return fmt.Errorf("写入ZIP失败: %v", err)
	}
	
	return nil
}

// GetBackupList 获取备份列表
func (s *BackupService) GetBackupList(page, pageSize int, backupType string) ([]models.DataBackup, int64, error) {
	var backups []models.DataBackup
	var total int64
	
	query := s.db.Model(&models.DataBackup{})
	
	if backupType != "" {
		query = query.Where("backup_type = ?", backupType)
	}
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&backups).Error
	
	return backups, total, err
}

// GetBackup 获取备份详情
func (s *BackupService) GetBackup(backupID string) (*models.DataBackup, error) {
	var backup models.DataBackup
	err := s.db.Where("id = ?", backupID).First(&backup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("备份不存在")
		}
		return nil, err
	}
	return &backup, nil
}

// DeleteBackup 删除备份
func (s *BackupService) DeleteBackup(backupID string) error {
	var backup models.DataBackup
	err := s.db.Where("id = ?", backupID).First(&backup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("备份不存在")
		}
		return err
	}
	
	// 删除备份文件
	if backup.BackupPath != "" {
		if err := os.Remove(backup.BackupPath); err != nil {
			// 文件删除失败，记录日志但继续删除数据库记录
			fmt.Printf("删除备份文件失败: %v\n", err)
		}
	}
	
	// 删除数据库记录
	return s.db.Delete(&backup).Error
}

// DownloadBackup 下载备份文件
func (s *BackupService) DownloadBackup(backupID string) (string, error) {
	backup, err := s.GetBackup(backupID)
	if err != nil {
		return "", err
	}
	
	if backup.Status != "completed" {
		return "", errors.New("备份未完成")
	}
	
	if backup.BackupPath == "" {
		return "", errors.New("备份文件路径不存在")
	}
	
	// 检查文件是否存在
	if _, err := os.Stat(backup.BackupPath); os.IsNotExist(err) {
		return "", errors.New("备份文件不存在")
	}
	
	return backup.BackupPath, nil
}

// RestoreBackup 恢复备份（简化版本，实际应用中需要更复杂的逻辑）
func (s *BackupService) RestoreBackup(backupID string) error {
	backup, err := s.GetBackup(backupID)
	if err != nil {
		return err
	}
	
	if backup.Status != "completed" {
		return errors.New("备份未完成，无法恢复")
	}
	
	// 打开ZIP文件
	zipReader, err := zip.OpenReader(backup.BackupPath)
	if err != nil {
		return fmt.Errorf("打开备份文件失败: %v", err)
	}
	defer zipReader.Close()
	
	// 遍历ZIP文件中的所有文件
	for _, file := range zipReader.File {
		if err := s.restoreTableFromZip(file); err != nil {
			return fmt.Errorf("恢复表数据失败: %v", err)
		}
	}
	
	return nil
}

// restoreTableFromZip 从ZIP文件中恢复表数据
func (s *BackupService) restoreTableFromZip(file *zip.File) error {
	// 打开文件
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()
	
	// 读取JSON数据
	jsonData, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	
	// 解析JSON
	var data []map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}
	
	// 获取表名（去掉.json后缀）
	tableName := file.Name[:len(file.Name)-5]
	
	// 批量插入数据（这里使用简单的插入，实际应用中可能需要更复杂的逻辑）
	for _, record := range data {
		s.db.Table(tableName).Create(record)
	}
	
	return nil
}

// CleanOldBackups 清理旧备份（保留最近N个备份）
func (s *BackupService) CleanOldBackups(keepCount int) error {
	var backups []models.DataBackup
	
	// 获取所有已完成的备份，按时间倒序
	err := s.db.Where("status = ?", "completed").
		Order("created_at DESC").
		Find(&backups).Error
	if err != nil {
		return err
	}
	
	// 如果备份数量超过保留数量，删除旧备份
	if len(backups) > keepCount {
		for i := keepCount; i < len(backups); i++ {
			if err := s.DeleteBackup(backups[i].ID); err != nil {
				fmt.Printf("删除旧备份失败: %v\n", err)
			}
		}
	}
	
	return nil
}

// GetBackupStats 获取备份统计信息
func (s *BackupService) GetBackupStats() (map[string]interface{}, error) {
	var totalBackups int64
	var completedBackups int64
	var failedBackups int64
	var totalSize int64
	
	s.db.Model(&models.DataBackup{}).Count(&totalBackups)
	s.db.Model(&models.DataBackup{}).Where("status = ?", "completed").Count(&completedBackups)
	s.db.Model(&models.DataBackup{}).Where("status = ?", "failed").Count(&failedBackups)
	
	// 计算总大小
	var backups []models.DataBackup
	s.db.Where("status = ?", "completed").Find(&backups)
	for _, backup := range backups {
		totalSize += backup.FileSize
	}
	
	// 获取最近一次备份时间
	var lastBackup models.DataBackup
	s.db.Where("status = ?", "completed").Order("completed_at DESC").First(&lastBackup)
	
	stats := map[string]interface{}{
		"total_backups":     totalBackups,
		"completed_backups": completedBackups,
		"failed_backups":    failedBackups,
		"total_size":        totalSize,
		"total_size_mb":     float64(totalSize) / 1024 / 1024,
		"last_backup_time":  lastBackup.CompletedAt,
	}
	
	return stats, nil
}
