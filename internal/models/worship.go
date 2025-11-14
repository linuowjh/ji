package models

import (
	"time"

	"gorm.io/gorm"
)

type WorshipRecord struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:祭扫记录ID"`
	MemorialID  string         `json:"memorial_id" gorm:"type:varchar(36);not null;index;comment:纪念馆ID"`
	UserID      string         `json:"user_id" gorm:"type:varchar(36);not null;index;comment:用户ID"`
	WorshipType string         `json:"worship_type" gorm:"type:varchar(20);not null;comment:祭扫类型:flower献花 candle点烛 incense上香 tribute供品 prayer祈福"`
	Content     string         `json:"content" gorm:"type:json;comment:祭扫内容(JSON格式)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (WorshipRecord) TableName() string {
	return "worship_records"
}

type Prayer struct {
	ID         string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:祈福ID"`
	MemorialID string         `json:"memorial_id" gorm:"type:varchar(36);not null;index;comment:纪念馆ID"`
	UserID     string         `json:"user_id" gorm:"type:varchar(36);not null;index;comment:用户ID"`
	Content    string         `json:"content" gorm:"type:text;not null;comment:祈福内容"`
	IsPublic   bool           `json:"is_public" gorm:"default:true;comment:是否公开显示"`
	CreatedAt  time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (Prayer) TableName() string {
	return "prayers"
}

type Message struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:留言ID"`
	MemorialID  string         `json:"memorial_id" gorm:"type:varchar(36);not null;index;comment:纪念馆ID"`
	UserID      string         `json:"user_id" gorm:"type:varchar(36);not null;index;comment:用户ID"`
	MessageType string         `json:"message_type" gorm:"type:varchar(20);not null;comment:留言类型:text文字 audio音频 video视频"`
	Content     string         `json:"content" gorm:"type:text;comment:留言内容"`
	MediaURL    string         `json:"media_url" gorm:"type:varchar(255);comment:媒体文件URL"`
	Duration    int            `json:"duration" gorm:"comment:音频/视频时长(秒)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (Message) TableName() string {
	return "messages"
}