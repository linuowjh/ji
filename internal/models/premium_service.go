package models

import (
	"time"
)

// PremiumPackage 高级套餐
type PremiumPackage struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	PackageName string    `gorm:"column:package_name;not null" json:"package_name"`
	PackageType string    `gorm:"column:package_type;not null" json:"package_type"` // memorial|service|storage
	Description string    `gorm:"column:description;type:text" json:"description"`
	Features    string    `gorm:"column:features;type:json" json:"features"` // JSON格式的功能列表
	Price       float64   `gorm:"column:price;not null" json:"price"`
	Duration    int       `gorm:"column:duration;default:365" json:"duration"` // 有效期（天）
	StorageSize int64     `gorm:"column:storage_size" json:"storage_size"`     // 存储空间（字节）
	IsActive    bool      `gorm:"column:is_active;default:true" json:"is_active"`
	SortOrder   int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PremiumPackage) TableName() string {
	return "premium_packages"
}

// UserSubscription 用户订阅
type UserSubscription struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	UserID          string     `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	PackageID       string     `gorm:"column:package_id;type:varchar(36);not null;index" json:"package_id"`
	MemorialID      string     `gorm:"column:memorial_id;type:varchar(36);index" json:"memorial_id"` // 关联的纪念馆（可选）
	Status          string     `gorm:"column:status;default:active" json:"status"`  // active|expired|cancelled
	StartDate       time.Time  `gorm:"column:start_date;not null" json:"start_date"`
	EndDate         time.Time  `gorm:"column:end_date;not null" json:"end_date"`
	AutoRenew       bool       `gorm:"column:auto_renew;default:false" json:"auto_renew"`
	PaymentAmount   float64    `gorm:"column:payment_amount" json:"payment_amount"`
	PaymentMethod   string     `gorm:"column:payment_method" json:"payment_method"` // wechat|alipay
	TransactionID   string     `gorm:"column:transaction_id" json:"transaction_id"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CancelledAt     *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
	
	// 关联
	Package *PremiumPackage `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	User    *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

// MemorialUpgrade 纪念馆升级记录
type MemorialUpgrade struct {
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	MemorialID      string    `gorm:"column:memorial_id;type:varchar(36);not null;index" json:"memorial_id"`
	SubscriptionID  string    `gorm:"column:subscription_id;type:varchar(36);not null;index" json:"subscription_id"`
	UpgradeType     string    `gorm:"column:upgrade_type;not null" json:"upgrade_type"` // theme|tombstone|storage|feature
	UpgradeData     string    `gorm:"column:upgrade_data;type:json" json:"upgrade_data"`
	IsActive        bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	
	// 关联
	Memorial     *Memorial         `gorm:"foreignKey:MemorialID" json:"memorial,omitempty"`
	Subscription *UserSubscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
}

func (MemorialUpgrade) TableName() string {
	return "memorial_upgrades"
}

// CustomTemplate 定制模板
type CustomTemplate struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	UserID       string    `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	MemorialID   string    `gorm:"column:memorial_id;type:varchar(36);index" json:"memorial_id"`
	TemplateType string    `gorm:"column:template_type;not null" json:"template_type"` // theme|tombstone|layout
	TemplateName string    `gorm:"column:template_name;not null" json:"template_name"`
	TemplateData string    `gorm:"column:template_data;type:json" json:"template_data"`
	PreviewURL   string    `gorm:"column:preview_url" json:"preview_url"`
	Status       string    `gorm:"column:status;default:draft" json:"status"` // draft|active|archived
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	
	// 关联
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Memorial *Memorial `gorm:"foreignKey:MemorialID" json:"memorial,omitempty"`
}

func (CustomTemplate) TableName() string {
	return "custom_templates"
}

// StorageUsage 存储使用情况
type StorageUsage struct {
	ID            string    `gorm:"column:id;primaryKey" json:"id"`
	UserID        string    `gorm:"column:user_id;type:varchar(36);not null;uniqueIndex" json:"user_id"`
	UsedSpace     int64     `gorm:"column:used_space;default:0" json:"used_space"`         // 已使用空间（字节）
	TotalSpace    int64     `gorm:"column:total_space;default:104857600" json:"total_space"` // 总空间（字节），默认100MB
	FileCount     int       `gorm:"column:file_count;default:0" json:"file_count"`
	LastUpdated   time.Time `gorm:"column:last_updated;autoUpdateTime" json:"last_updated"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (StorageUsage) TableName() string {
	return "storage_usages"
}

// PaymentOrder 支付订单
type PaymentOrder struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	OrderNo         string     `gorm:"column:order_no;type:varchar(100);uniqueIndex;not null" json:"order_no"`
	UserID          string     `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	PackageID       string     `gorm:"column:package_id;type:varchar(36);not null;index" json:"package_id"`
	OrderType       string     `gorm:"column:order_type;not null" json:"order_type"` // subscription|upgrade|renewal
	Amount          float64    `gorm:"column:amount;not null" json:"amount"`
	PaymentMethod   string     `gorm:"column:payment_method" json:"payment_method"` // wechat|alipay
	PaymentStatus   string     `gorm:"column:payment_status;default:pending" json:"payment_status"` // pending|paid|failed|refunded
	TransactionID   string     `gorm:"column:transaction_id" json:"transaction_id"`
	PaymentTime     *time.Time `gorm:"column:payment_time" json:"payment_time"`
	RefundTime      *time.Time `gorm:"column:refund_time" json:"refund_time"`
	RefundAmount    float64    `gorm:"column:refund_amount" json:"refund_amount"`
	RefundReason    string     `gorm:"column:refund_reason;type:text" json:"refund_reason"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	
	// 关联
	User    *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Package *PremiumPackage  `gorm:"foreignKey:PackageID" json:"package,omitempty"`
}

func (PaymentOrder) TableName() string {
	return "payment_orders"
}

// ServiceUsageLog 服务使用日志
type ServiceUsageLog struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	UserID      string    `gorm:"column:user_id;type:varchar(36);not null;index" json:"user_id"`
	ServiceType string    `gorm:"column:service_type;not null" json:"service_type"` // photo_restore|custom_template|premium_service
	ServiceData string    `gorm:"column:service_data;type:json" json:"service_data"`
	UsageCount  int       `gorm:"column:usage_count;default:1" json:"usage_count"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;index" json:"created_at"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ServiceUsageLog) TableName() string {
	return "service_usage_logs"
}
