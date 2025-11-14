package models

import (
	"time"

	"gorm.io/gorm"
)

// Album 纪念相册
type Album struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID  string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	Title       string         `json:"title" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	CoverURL    string         `json:"cover_url" gorm:"type:varchar(255)"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial    `json:"memorial" gorm:"foreignKey:MemorialID"`
	Photos   []AlbumPhoto `json:"photos" gorm:"foreignKey:AlbumID"`
}

func (Album) TableName() string {
	return "albums"
}

// AlbumPhoto 相册照片
type AlbumPhoto struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AlbumID     string         `json:"album_id" gorm:"type:varchar(36);not null;index"`
	PhotoURL    string         `json:"photo_url" gorm:"type:varchar(255);not null"`
	ThumbnailURL string        `json:"thumbnail_url" gorm:"type:varchar(255)"`
	Caption     string         `json:"caption" gorm:"type:text"`
	TakenDate   *time.Time     `json:"taken_date"`
	Location    string         `json:"location" gorm:"type:varchar(100)"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Album Album `json:"album" gorm:"foreignKey:AlbumID"`
}

func (AlbumPhoto) TableName() string {
	return "album_photos"
}

// LifeStory 生平故事
type LifeStory struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID  string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	Title       string         `json:"title" gorm:"type:varchar(100);not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	StoryDate   *time.Time     `json:"story_date"`
	AgeAtTime   int            `json:"age_at_time"`
	Location    string         `json:"location" gorm:"type:varchar(100)"`
	Category    string         `json:"category" gorm:"type:varchar(50)"` // childhood|youth|career|family|achievement|other
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	AuthorID    string         `json:"author_id" gorm:"type:varchar(36);not null;index"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial           `json:"memorial" gorm:"foreignKey:MemorialID"`
	Author   User               `json:"author" gorm:"foreignKey:AuthorID"`
	Media    []LifeStoryMedia   `json:"media" gorm:"foreignKey:LifeStoryID"`
}

func (LifeStory) TableName() string {
	return "life_stories"
}

// LifeStoryMedia 生平故事媒体文件
type LifeStoryMedia struct {
	ID           string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	LifeStoryID  string         `json:"life_story_id" gorm:"type:varchar(36);not null;index"`
	MediaType    string         `json:"media_type" gorm:"type:varchar(20);not null;comment:image|video|audio"`
	MediaURL     string         `json:"media_url" gorm:"type:varchar(255);not null"`
	ThumbnailURL string         `json:"thumbnail_url" gorm:"type:varchar(255)"`
	Caption      string         `json:"caption" gorm:"type:text"`
	SortOrder    int            `json:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	LifeStory LifeStory `json:"life_story" gorm:"foreignKey:LifeStoryID"`
}

func (LifeStoryMedia) TableName() string {
	return "life_story_media"
}

// Timeline 时间轴事件
type Timeline struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID  string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	Title       string         `json:"title" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	EventDate   time.Time      `json:"event_date" gorm:"not null"`
	EventType   string         `json:"event_type" gorm:"type:varchar(50)"` // birth|education|career|marriage|achievement|death|other
	IconURL     string         `json:"icon_url" gorm:"type:varchar(255)"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	AuthorID    string         `json:"author_id" gorm:"type:varchar(36);not null;index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
	Author   User     `json:"author" gorm:"foreignKey:AuthorID"`
}

func (Timeline) TableName() string {
	return "timelines"
}