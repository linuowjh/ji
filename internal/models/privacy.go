package models

import (
	"time"

	"gorm.io/gorm"
)

// 访客权限设置
type VisitorPermissionSetting struct {
	ID           string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID   string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	UserID       string         `json:"user_id" gorm:"type:varchar(36);index"`
	FamilyID     string         `json:"family_id" gorm:"type:varchar(36);index"`
	PermissionType string       `json:"permission_type" gorm:"type:varchar(20);not null"`
	IsAllowed    bool           `json:"is_allowed" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
	Family   Family   `json:"family" gorm:"foreignKey:FamilyID"`
}

func (VisitorPermissionSetting) TableName() string {
	return "visitor_permission_settings"
}

// 访客黑名单
type VisitorBlacklist struct {
	ID         string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	UserID     string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Reason     string         `json:"reason" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (VisitorBlacklist) TableName() string {
	return "visitor_blacklists"
}

// 访问申请
type AccessRequest struct {
	ID         string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	UserID     string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Message    string         `json:"message" gorm:"type:text"`
	Status     string         `json:"status" gorm:"type:varchar(20);default:pending;comment:pending|approved|rejected"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (AccessRequest) TableName() string {
	return "access_requests"
}