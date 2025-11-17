package models

import (
	"time"

	"gorm.io/gorm"
)

type Memorial struct {
	ID             string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:纪念馆ID"`
	CreatorID      string         `json:"creatorId" gorm:"type:varchar(36);not null;index;comment:创建者ID"`
	DeceasedName   string         `json:"deceasedName" gorm:"type:varchar(50);not null;comment:逝者姓名"`
	BirthDate      *time.Time     `json:"birthDate" gorm:"comment:出生日期"`
	DeathDate      *time.Time     `json:"deathDate" gorm:"comment:逝世日期"`
	Biography      string         `json:"biography" gorm:"type:text;comment:生平简介"`
	AvatarURL      string         `json:"avatarUrl" gorm:"type:varchar(255);comment:头像URL"`
	ThemeStyle     string         `json:"themeStyle" gorm:"type:varchar(50);default:traditional;comment:主题风格"`
	TombstoneStyle string         `json:"tombstoneStyle" gorm:"type:varchar(50);default:marble;comment:墓碑样式"`
	Epitaph        string         `json:"epitaph" gorm:"type:text;comment:墓志铭"`
	PrivacyLevel   int            `json:"privacyLevel" gorm:"default:1;comment:隐私级别:1家族可见 2私密"`
	Status         int            `json:"status" gorm:"default:1;comment:状态:1正常 0禁用"`
	CreatedAt      time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt      time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Creator User `json:"creator" gorm:"foreignKey:CreatorID"`
}

func (Memorial) TableName() string {
	return "memorials"
}

type MediaFile struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:媒体文件ID"`
	MemorialID  string         `json:"memorialId" gorm:"type:varchar(36);not null;index;comment:纪念馆ID"`
	FileType    string         `json:"fileType" gorm:"type:varchar(20);not null;comment:文件类型:image图片 video视频 audio音频"`
	FileURL     string         `json:"url" gorm:"type:varchar(255);not null;comment:文件URL"` // 简化为 url
	FileName    string         `json:"fileName" gorm:"type:varchar(255);comment:文件名"`
	FileSize    int64          `json:"fileSize" gorm:"comment:文件大小(字节)"`
	Description string         `json:"description" gorm:"type:text;comment:文件描述"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
}

func (MediaFile) TableName() string {
	return "media_files"
}

// MemorialFamily 纪念馆家族关联表
type MemorialFamily struct {
	ID         string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:关联ID"`
	MemorialID string         `json:"memorial_id" gorm:"type:varchar(36);not null;index;comment:纪念馆ID"`
	FamilyID   string         `json:"family_id" gorm:"type:varchar(36);not null;index;comment:家族ID"`
	CreatedAt  time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	Family   Family   `json:"family" gorm:"foreignKey:FamilyID"`
}

func (MemorialFamily) TableName() string {
	return "memorial_families"
}
