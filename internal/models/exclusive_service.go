package models

import (
	"time"
)

// ExclusiveService 专属服务
type ExclusiveService struct {
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	ServiceName     string    `gorm:"column:service_name;not null" json:"service_name"`
	ServiceType     string    `gorm:"column:service_type;not null" json:"service_type"` // memorial_service|data_backup|photo_restore|custom_design
	Description     string    `gorm:"column:description;type:text" json:"description"`
	BasePrice       float64   `gorm:"column:base_price;not null" json:"base_price"`
	PriceUnit       string    `gorm:"column:price_unit" json:"price_unit"` // per_hour|per_service|per_gb
	Features        string    `gorm:"column:features;type:json" json:"features"`
	RequireBooking  bool      `gorm:"column:require_booking;default:false" json:"require_booking"`
	MaxDuration     int       `gorm:"column:max_duration" json:"max_duration"` // 最大时长（分钟）
	IsActive        bool      `gorm:"column:is_active;default:true" json:"is_active"`
	SortOrder       int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ExclusiveService) TableName() string {
	return "exclusive_services"
}

// ServiceBooking 服务预订
type ServiceBooking struct {
	ID                string     `gorm:"column:id;primaryKey" json:"id"`
	UserID            string     `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	ServiceID         string     `gorm:"column:service_id;type:varchar(36);not null;index" json:"service_id"`
	MemorialID        string     `gorm:"column:memorial_id;type:varchar(36);index" json:"memorial_id"`
	BookingType       string     `gorm:"column:booking_type;not null" json:"booking_type"` // memorial_service|data_backup|custom_design
	ScheduledTime     time.Time  `gorm:"column:scheduled_time" json:"scheduled_time"`
	Duration          int        `gorm:"column:duration" json:"duration"` // 时长（分钟）
	Status            string     `gorm:"column:status;default:pending" json:"status"` // pending|confirmed|in_progress|completed|cancelled
	TotalPrice        float64    `gorm:"column:total_price" json:"total_price"`
	Requirements      string     `gorm:"column:requirements;type:text" json:"requirements"`
	ServiceData       string     `gorm:"column:service_data;type:json" json:"service_data"`
	StaffID           string     `gorm:"column:staff_id" json:"staff_id"`
	CompletedAt       *time.Time `gorm:"column:completed_at" json:"completed_at"`
	CancelledAt       *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
	CancellationReason string    `gorm:"column:cancellation_reason;type:text" json:"cancellation_reason"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	
	// 关联
	User    *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Service *ExclusiveService  `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
}

func (ServiceBooking) TableName() string {
	return "service_bookings"
}

// DataExportRequest 数据导出请求
type DataExportRequest struct {
	ID            string     `gorm:"column:id;primaryKey" json:"id"`
	UserID        string     `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	ExportType    string     `gorm:"column:export_type;not null" json:"export_type"` // full|memorial|family
	TargetID      string     `gorm:"column:target_id" json:"target_id"` // 纪念馆ID或家族ID
	ExportFormat  string     `gorm:"column:export_format;default:zip" json:"export_format"` // zip|pdf|json
	IncludeMedia  bool       `gorm:"column:include_media;default:true" json:"include_media"`
	Encrypted     bool       `gorm:"column:encrypted;default:false" json:"encrypted"`
	EncryptionKey string     `gorm:"column:encryption_key" json:"encryption_key"`
	Status        string     `gorm:"column:status;default:pending" json:"status"` // pending|processing|completed|failed
	FileSize      int64      `gorm:"column:file_size" json:"file_size"`
	FilePath      string     `gorm:"column:file_path" json:"file_path"`
	DownloadURL   string     `gorm:"column:download_url" json:"download_url"`
	ExpiresAt     *time.Time `gorm:"column:expires_at" json:"expires_at"`
	ErrorMessage  string     `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CompletedAt   *time.Time `gorm:"column:completed_at" json:"completed_at"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (DataExportRequest) TableName() string {
	return "data_export_requests"
}

// PhotoRestoreRequest 老照片修复请求
type PhotoRestoreRequest struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	UserID          string     `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	MemorialID      string     `gorm:"column:memorial_id;type:varchar(36);index" json:"memorial_id"`
	OriginalPhotoURL string    `gorm:"column:original_photo_url;not null" json:"original_photo_url"`
	RestoredPhotoURL string    `gorm:"column:restored_photo_url" json:"restored_photo_url"`
	RestoreType     string     `gorm:"column:restore_type;not null" json:"restore_type"` // colorize|enhance|repair|all
	Status          string     `gorm:"column:status;default:pending" json:"status"` // pending|processing|completed|failed
	ProcessingTime  int        `gorm:"column:processing_time" json:"processing_time"` // 处理时长（秒）
	ErrorMessage    string     `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at" json:"completed_at"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PhotoRestoreRequest) TableName() string {
	return "photo_restore_requests"
}

// CustomDesignRequest 定制设计请求
type CustomDesignRequest struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	UserID          string     `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	MemorialID      string     `gorm:"column:memorial_id;type:varchar(36);index" json:"memorial_id"`
	DesignType      string     `gorm:"column:design_type;not null" json:"design_type"` // theme|tombstone|layout|complete
	Requirements    string     `gorm:"column:requirements;type:text" json:"requirements"`
	ReferenceImages string     `gorm:"column:reference_images;type:json" json:"reference_images"` // 参考图片URL列表
	Budget          float64    `gorm:"column:budget" json:"budget"`
	Status          string     `gorm:"column:status;default:pending" json:"status"` // pending|in_design|review|completed|cancelled
	DesignerID      string     `gorm:"column:designer_id" json:"designer_id"`
	DesignFiles     string     `gorm:"column:design_files;type:json" json:"design_files"`
	FeedbackCount   int        `gorm:"column:feedback_count;default:0" json:"feedback_count"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at" json:"completed_at"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (CustomDesignRequest) TableName() string {
	return "custom_design_requests"
}

// ServiceReview 服务评价
type ServiceReview struct {
	ID         string    `gorm:"primaryKey;size:36;comment:评价ID" json:"id"`
	UserID     string    `gorm:"size:36;not null;index;comment:用户ID" json:"user_id"`
	BookingID  string    `gorm:"size:36;not null;uniqueIndex;comment:预订ID" json:"booking_id"`
	Rating     int       `gorm:"not null;comment:评分(1-5星)" json:"rating"`
	Comment    string    `gorm:"type:text;comment:评价内容" json:"comment"`
	Tags       string    `gorm:"type:json;comment:评价标签" json:"tags"`
	IsAnonymous bool     `gorm:"default:false;comment:是否匿名评价" json:"is_anonymous"`
	CreatedAt  time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	
	// 关联
	User    *User            `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Booking *ServiceBooking  `gorm:"foreignKey:BookingID;references:ID;constraint:OnDelete:CASCADE" json:"booking,omitempty"`
}

func (ServiceReview) TableName() string {
	return "service_reviews"
}

// ServiceStaff 服务人员
type ServiceStaff struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	Role         string    `gorm:"column:role;not null" json:"role"` // designer|coordinator|technician
	Specialties  string    `gorm:"column:specialties;type:json" json:"specialties"`
	AvatarURL    string    `gorm:"column:avatar_url" json:"avatar_url"`
	Bio          string    `gorm:"column:bio;type:text" json:"bio"`
	Rating       float64   `gorm:"column:rating;default:5.0" json:"rating"`
	ReviewCount  int       `gorm:"column:review_count;default:0" json:"review_count"`
	IsAvailable  bool      `gorm:"column:is_available;default:true" json:"is_available"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceStaff) TableName() string {
	return "service_staff"
}
