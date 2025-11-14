package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoggingService struct {
	db *gorm.DB
}

func NewLoggingService(db *gorm.DB) *LoggingService {
	return &LoggingService{
		db: db,
	}
}

// LogLevel 日志级别
const (
	LogLevelInfo     = "info"
	LogLevelWarning  = "warning"
	LogLevelError    = "error"
	LogLevelCritical = "critical"
)

// LogType 日志类型
const (
	LogTypeAdmin    = "admin"
	LogTypeSystem   = "system"
	LogTypeSecurity = "security"
	LogTypeAPI      = "api"
)

// LogEntry 日志条目结构
type LogEntry struct {
	Level     string
	Type      string
	UserID    string
	Action    string
	Details   map[string]interface{}
	IPAddress string
	UserAgent string
}

// Log 记录日志
func (s *LoggingService) Log(entry *LogEntry) error {
	// 将详情转换为JSON字符串
	detailsJSON, err := json.Marshal(entry.Details)
	if err != nil {
		return fmt.Errorf("序列化日志详情失败: %v", err)
	}
	
	log := &models.SystemLog{
		ID:        uuid.New().String(),
		LogLevel:  entry.Level,
		LogType:   entry.Type,
		UserID:    entry.UserID,
		Action:    entry.Action,
		Details:   string(detailsJSON),
		IPAddress: entry.IPAddress,
		UserAgent: entry.UserAgent,
		CreatedAt: time.Now(),
	}
	
	return s.db.Create(log).Error
}

// LogInfo 记录信息日志
func (s *LoggingService) LogInfo(logType, userID, action string, details map[string]interface{}) error {
	return s.Log(&LogEntry{
		Level:   LogLevelInfo,
		Type:    logType,
		UserID:  userID,
		Action:  action,
		Details: details,
	})
}

// LogWarning 记录警告日志
func (s *LoggingService) LogWarning(logType, userID, action string, details map[string]interface{}) error {
	return s.Log(&LogEntry{
		Level:   LogLevelWarning,
		Type:    logType,
		UserID:  userID,
		Action:  action,
		Details: details,
	})
}

// LogError 记录错误日志
func (s *LoggingService) LogError(logType, userID, action string, details map[string]interface{}) error {
	return s.Log(&LogEntry{
		Level:   LogLevelError,
		Type:    logType,
		UserID:  userID,
		Action:  action,
		Details: details,
	})
}

// LogCritical 记录严重错误日志
func (s *LoggingService) LogCritical(logType, userID, action string, details map[string]interface{}) error {
	return s.Log(&LogEntry{
		Level:   LogLevelCritical,
		Type:    logType,
		UserID:  userID,
		Action:  action,
		Details: details,
	})
}

// LogAdminAction 记录管理员操作
func (s *LoggingService) LogAdminAction(adminID, action string, details map[string]interface{}, ipAddress, userAgent string) error {
	return s.Log(&LogEntry{
		Level:     LogLevelInfo,
		Type:      LogTypeAdmin,
		UserID:    adminID,
		Action:    action,
		Details:   details,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// LogSecurityEvent 记录安全事件
func (s *LoggingService) LogSecurityEvent(userID, action string, details map[string]interface{}, ipAddress string) error {
	return s.Log(&LogEntry{
		Level:     LogLevelWarning,
		Type:      LogTypeSecurity,
		UserID:    userID,
		Action:    action,
		Details:   details,
		IPAddress: ipAddress,
	})
}

// LogAPIRequest 记录API请求
func (s *LoggingService) LogAPIRequest(userID, endpoint, method string, statusCode int, duration int64, ipAddress string) error {
	return s.Log(&LogEntry{
		Level:  LogLevelInfo,
		Type:   LogTypeAPI,
		UserID: userID,
		Action: fmt.Sprintf("%s %s", method, endpoint),
		Details: map[string]interface{}{
			"status_code": statusCode,
			"duration_ms": duration,
		},
		IPAddress: ipAddress,
	})
}

// GetLogs 获取日志列表
func (s *LoggingService) GetLogs(filter *LogFilter) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64
	
	query := s.db.Model(&models.SystemLog{})
	
	// 应用过滤条件
	if filter.LogLevel != "" {
		query = query.Where("log_level = ?", filter.LogLevel)
	}
	if filter.LogType != "" {
		query = query.Where("log_type = ?", filter.LogType)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action LIKE ?", "%"+filter.Action+"%")
	}
	if filter.IPAddress != "" {
		query = query.Where("ip_address = ?", filter.IPAddress)
	}
	if filter.StartTime != "" {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if filter.EndTime != "" {
		query = query.Where("created_at <= ?", filter.EndTime)
	}
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&logs).Error
	
	return logs, total, err
}

// LogFilter 日志过滤器
type LogFilter struct {
	LogLevel  string
	LogType   string
	UserID    string
	Action    string
	IPAddress string
	StartTime string
	EndTime   string
	Page      int
	PageSize  int
}

// GetLogStats 获取日志统计信息
func (s *LoggingService) GetLogStats(startTime, endTime string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	query := s.db.Model(&models.SystemLog{})
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}
	
	// 总日志数
	var totalLogs int64
	query.Count(&totalLogs)
	stats["total_logs"] = totalLogs
	
	// 按级别统计
	levelStats := make(map[string]int64)
	levels := []string{LogLevelInfo, LogLevelWarning, LogLevelError, LogLevelCritical}
	for _, level := range levels {
		var count int64
		query.Where("log_level = ?", level).Count(&count)
		levelStats[level] = count
	}
	stats["by_level"] = levelStats
	
	// 按类型统计
	typeStats := make(map[string]int64)
	types := []string{LogTypeAdmin, LogTypeSystem, LogTypeSecurity, LogTypeAPI}
	for _, logType := range types {
		var count int64
		s.db.Model(&models.SystemLog{}).Where("log_type = ?", logType).Count(&count)
		typeStats[logType] = count
	}
	stats["by_type"] = typeStats
	
	// 最近的错误日志
	var recentErrors []models.SystemLog
	s.db.Where("log_level IN ?", []string{LogLevelError, LogLevelCritical}).
		Order("created_at DESC").
		Limit(10).
		Find(&recentErrors)
	stats["recent_errors"] = recentErrors
	
	return stats, nil
}

// CleanOldLogs 清理旧日志
func (s *LoggingService) CleanOldLogs(daysToKeep int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	
	result := s.db.Where("created_at < ?", cutoffDate).Delete(&models.SystemLog{})
	if result.Error != nil {
		return fmt.Errorf("清理旧日志失败: %v", result.Error)
	}
	
	fmt.Printf("已清理 %d 条旧日志\n", result.RowsAffected)
	return nil
}

// ExportLogs 导出日志
func (s *LoggingService) ExportLogs(filter *LogFilter) (string, error) {
	logs, _, err := s.GetLogs(filter)
	if err != nil {
		return "", err
	}
	
	jsonData, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return "", fmt.Errorf("导出日志失败: %v", err)
	}
	
	return string(jsonData), nil
}

// GetLogDetail 获取日志详情
func (s *LoggingService) GetLogDetail(logID string) (*models.SystemLog, error) {
	var log models.SystemLog
	err := s.db.Where("id = ?", logID).First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("日志不存在")
		}
		return nil, err
	}
	return &log, nil
}

// SearchLogs 搜索日志
func (s *LoggingService) SearchLogs(keyword string, page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64
	
	query := s.db.Model(&models.SystemLog{}).
		Where("action LIKE ? OR details LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error
	
	return logs, total, err
}

// GetUserActivityLogs 获取用户活动日志
func (s *LoggingService) GetUserActivityLogs(userID string, page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64
	
	query := s.db.Model(&models.SystemLog{}).Where("user_id = ?", userID)
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error
	
	return logs, total, err
}

// GetSecurityAlerts 获取安全警报
func (s *LoggingService) GetSecurityAlerts(page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64
	
	query := s.db.Model(&models.SystemLog{}).
		Where("log_type = ? AND log_level IN ?", LogTypeSecurity, []string{LogLevelWarning, LogLevelError, LogLevelCritical})
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error
	
	return logs, total, err
}
