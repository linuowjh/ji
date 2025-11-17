package models

import (
	"time"

	"gorm.io/gorm"
)

type Family struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:家族ID"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:家族名称"`
	CreatorID   string         `json:"creatorId" gorm:"type:varchar(36);not null;index;comment:创建者ID"`
	Description string         `json:"description" gorm:"type:text;comment:家族描述"`
	InviteCode  string         `json:"inviteCode" gorm:"uniqueIndex;type:varchar(20);comment:邀请码"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Creator User           `json:"creator" gorm:"foreignKey:CreatorID"`
	Members []FamilyMember `json:"members" gorm:"foreignKey:FamilyID"`
}

func (Family) TableName() string {
	return "families"
}

type FamilyMember struct {
	ID       string    `json:"id" gorm:"primaryKey;type:varchar(36);comment:成员ID"`
	FamilyID string    `json:"familyId" gorm:"column:family_id;type:varchar(36);not null;index;comment:家族ID"`
	UserID   string    `json:"userId" gorm:"column:user_id;type:varchar(36);not null;index;comment:用户ID"`
	Role     string    `json:"role" gorm:"type:varchar(20);default:member;comment:角色:admin管理员 member成员"`
	JoinedAt time.Time `json:"joinedAt" gorm:"column:joined_at;comment:加入时间"`

	// 关联关系
	Family Family `json:"family" gorm:"foreignKey:FamilyID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

func (FamilyMember) TableName() string {
	return "family_members"
}

type MemorialReminder struct {
	ID           string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID   string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	ReminderType string         `json:"reminder_type" gorm:"type:varchar(20);not null;comment:birthday|death_anniversary|festival"`
	ReminderDate time.Time      `json:"reminder_date"`
	Title        string         `json:"title" gorm:"type:varchar(100)"`
	Content      string         `json:"content" gorm:"type:text"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
}

func (MemorialReminder) TableName() string {
	return "memorial_reminders"
}

// FamilyInvitation 家族邀请
type FamilyInvitation struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FamilyID  string         `json:"family_id" gorm:"type:varchar(36);not null;index"`
	InviterID string         `json:"inviter_id" gorm:"type:varchar(36);not null;index"`
	InviteeID string         `json:"invitee_id" gorm:"type:varchar(36);not null;index"`
	Status    string         `json:"status" gorm:"type:varchar(20);default:pending;comment:pending|accepted|declined|expired"`
	Message   string         `json:"message" gorm:"type:text"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Family  Family `json:"family" gorm:"foreignKey:FamilyID"`
	Inviter User   `json:"inviter" gorm:"foreignKey:InviterID"`
	Invitee User   `json:"invitee" gorm:"foreignKey:InviteeID"`
}

func (FamilyInvitation) TableName() string {
	return "family_invitations"
}

// FamilyActivity 家族活动记录
type FamilyActivity struct {
	ID           string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FamilyID     string         `json:"familyId" gorm:"column:family_id;type:varchar(36);not null;index"`
	UserID       string         `json:"userId" gorm:"column:user_id;type:varchar(36);not null;index"`
	MemorialID   string         `json:"memorialId" gorm:"column:memorial_id;type:varchar(36);index"`
	ActivityType string         `json:"activityType" gorm:"column:activity_type;type:varchar(30);not null;comment:worship|join|create_memorial|create_story"`
	Content      string         `json:"content" gorm:"type:json"`
	Timestamp    time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Family   Family   `json:"family" gorm:"foreignKey:FamilyID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
	Memorial Memorial `json:"memorial" gorm:"foreignKey:MemorialID"`
}

func (FamilyActivity) TableName() string {
	return "family_activities"
}

// FamilyGenealogy 家族谱系
type FamilyGenealogy struct {
	ID           string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FamilyID     string         `json:"family_id" gorm:"type:varchar(36);not null;index"`
	PersonName   string         `json:"person_name" gorm:"type:varchar(100);not null"`
	Generation   int            `json:"generation" gorm:"not null;comment:辈分，数字越小辈分越高"`
	ParentID     string         `json:"parent_id" gorm:"type:varchar(36);index;comment:父辈ID"`
	Gender       string         `json:"gender" gorm:"type:varchar(10);comment:male|female"`
	BirthDate    *time.Time     `json:"birth_date"`
	DeathDate    *time.Time     `json:"death_date"`
	Biography    string         `json:"biography" gorm:"type:text"`
	AvatarURL    string         `json:"avatar_url" gorm:"type:varchar(255)"`
	MemorialID   string         `json:"memorial_id" gorm:"type:varchar(36);index;comment:关联的纪念馆ID"`
	Position     string         `json:"position" gorm:"type:varchar(100);comment:在家族中的地位或职业"`
	Achievements string         `json:"achievements" gorm:"type:text;comment:主要成就"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Family   Family            `json:"family" gorm:"foreignKey:FamilyID"`
	Parent   *FamilyGenealogy  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []FamilyGenealogy `json:"children" gorm:"foreignKey:ParentID"`
	Memorial *Memorial         `json:"memorial" gorm:"foreignKey:MemorialID"`
}

func (FamilyGenealogy) TableName() string {
	return "family_genealogies"
}

// FamilyStory 家族故事
type FamilyStory struct {
	ID         string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FamilyID   string         `json:"family_id" gorm:"type:varchar(36);not null;index"`
	AuthorID   string         `json:"author_id" gorm:"type:varchar(36);not null;index"`
	Title      string         `json:"title" gorm:"type:varchar(200);not null"`
	Content    string         `json:"content" gorm:"type:longtext;not null"`
	Category   string         `json:"category" gorm:"type:varchar(50);comment:tradition|achievement|migration|business|education|war|love"`
	Period     string         `json:"period" gorm:"type:varchar(100);comment:故事发生的时期"`
	Characters string         `json:"characters" gorm:"type:json;comment:故事中的人物"`
	Location   string         `json:"location" gorm:"type:varchar(200);comment:故事发生地点"`
	MediaFiles string         `json:"media_files" gorm:"type:json;comment:相关媒体文件"`
	Tags       string         `json:"tags" gorm:"type:varchar(500);comment:标签，逗号分隔"`
	IsPublic   bool           `json:"is_public" gorm:"default:true;comment:是否对家族成员公开"`
	ViewCount  int            `json:"view_count" gorm:"default:0"`
	LikeCount  int            `json:"like_count" gorm:"default:0"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Family Family `json:"family" gorm:"foreignKey:FamilyID"`
	Author User   `json:"author" gorm:"foreignKey:AuthorID"`
}

func (FamilyStory) TableName() string {
	return "family_stories"
}

// FamilyTradition 家族传统
type FamilyTradition struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FamilyID    string         `json:"family_id" gorm:"type:varchar(36);not null;index"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Category    string         `json:"category" gorm:"type:varchar(50);comment:festival|ceremony|custom|rule|recipe"`
	Origin      string         `json:"origin" gorm:"type:text;comment:传统的起源"`
	Practice    string         `json:"practice" gorm:"type:text;comment:具体做法"`
	Meaning     string         `json:"meaning" gorm:"type:text;comment:意义和价值"`
	MediaFiles  string         `json:"media_files" gorm:"type:json;comment:相关媒体文件"`
	IsActive    bool           `json:"is_active" gorm:"default:true;comment:是否仍在传承"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Family Family `json:"family" gorm:"foreignKey:FamilyID"`
}

func (FamilyTradition) TableName() string {
	return "family_traditions"
}
