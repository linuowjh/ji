package models

import (
	"time"
)

// VisitorRecord 访客记录
type VisitorRecord struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36);comment:访客记录ID"`
	MemorialID string    `json:"memorial_id" gorm:"type:varchar(36);not null;index;comment:纪念馆ID"`
	VisitorID  string    `json:"visitor_id" gorm:"type:varchar(36);not null;index;comment:访客用户ID"`
	VisitTime  time.Time `json:"visit_time" gorm:"comment:访问时间"`
	IPAddress  string    `json:"ip_address" gorm:"type:varchar(45);comment:访客IP地址"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	Visitor  User     `json:"visitor" gorm:"foreignKey:VisitorID"`
}

func (VisitorRecord) TableName() string {
	return "visitor_records"
}