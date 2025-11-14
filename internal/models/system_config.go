package models

import (
	"time"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	ConfigKey   string    `gorm:"column:config_key;type:varchar(255);uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"column:config_value;type:text" json:"config_value"`
	ConfigType  string    `gorm:"column:config_type;not null" json:"config_type"` // festival|template|system
	Description string    `gorm:"column:description;type:text" json:"description"`
	IsActive    bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

// FestivalConfig 祭扫节日配置
type FestivalConfig struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	FestivalDate string    `gorm:"column:festival_date;not null" json:"festival_date"` // MM-DD 格式
	Description  string    `gorm:"column:description;type:text" json:"description"`
	ReminderDays int       `gorm:"column:reminder_days;default:3" json:"reminder_days"` // 提前几天提醒
	IsActive     bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (FestivalConfig) TableName() string {
	return "festival_configs"
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	TemplateType string    `gorm:"column:template_type;not null" json:"template_type"` // theme|tombstone|prayer
	TemplateName string    `gorm:"column:template_name;not null" json:"template_name"`
	TemplateData string    `gorm:"column:template_data;type:json" json:"template_data"`
	PreviewURL   string    `gorm:"column:preview_url" json:"preview_url"`
	IsPremium    bool      `gorm:"column:is_premium;default:false" json:"is_premium"`
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsActive     bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (TemplateConfig) TableName() string {
	return "template_configs"
}

// DataBackup 数据备份记录
type DataBackup struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	BackupType   string    `gorm:"column:backup_type;not null" json:"backup_type"` // full|incremental|user
	BackupPath   string    `gorm:"column:backup_path;not null" json:"backup_path"`
	FileSize     int64     `gorm:"column:file_size" json:"file_size"`
	Status       string    `gorm:"column:status;default:pending" json:"status"` // pending|processing|completed|failed
	ErrorMessage string    `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completed_at"`
}

func (DataBackup) TableName() string {
	return "data_backups"
}

// SystemLog 系统日志
type SystemLog struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	LogLevel  string    `gorm:"column:log_level;not null" json:"log_level"` // info|warning|error|critical
	LogType   string    `gorm:"column:log_type;not null" json:"log_type"`   // admin|system|security|api
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	Action    string    `gorm:"column:action" json:"action"`
	Details   string    `gorm:"column:details;type:json" json:"details"`
	IPAddress string    `gorm:"column:ip_address" json:"ip_address"`
	UserAgent string    `gorm:"column:user_agent" json:"user_agent"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index" json:"created_at"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}

// SystemMonitor 系统监控指标
type SystemMonitor struct {
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	MetricType      string    `gorm:"column:metric_type;not null" json:"metric_type"` // cpu|memory|disk|api|database
	MetricValue     float64   `gorm:"column:metric_value" json:"metric_value"`
	MetricUnit      string    `gorm:"column:metric_unit" json:"metric_unit"`
	AdditionalInfo  string    `gorm:"column:additional_info;type:json" json:"additional_info"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime;index" json:"created_at"`
}

func (SystemMonitor) TableName() string {
	return "system_monitors"
}
