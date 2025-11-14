package models

import (
	"time"

	"gorm.io/gorm"
)

// MemorialService 追思会
type MemorialService struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemorialID  string         `json:"memorial_id" gorm:"type:varchar(36);not null;index"`
	Title       string         `json:"title" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	StartTime   time.Time      `json:"start_time" gorm:"not null"`
	EndTime     time.Time      `json:"end_time" gorm:"not null"`
	Status      string         `json:"status" gorm:"type:varchar(20);default:scheduled;comment:scheduled|ongoing|completed|cancelled"`
	MaxParticipants int        `json:"max_participants" gorm:"default:50"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	InviteCode  string         `json:"invite_code" gorm:"uniqueIndex;type:varchar(20)"`
	HostID      string         `json:"host_id" gorm:"type:varchar(36);not null;index"`
	RecordingURL string        `json:"recording_url" gorm:"type:varchar(255)"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Memorial     Memorial                    `json:"memorial" gorm:"foreignKey:MemorialID"`
	Host         User                        `json:"host" gorm:"foreignKey:HostID"`
	Participants []MemorialServiceParticipant `json:"participants" gorm:"foreignKey:ServiceID"`
	Activities   []ServiceActivity           `json:"activities" gorm:"foreignKey:ServiceID"`
}

func (MemorialService) TableName() string {
	return "memorial_services"
}

// MemorialServiceParticipant 追思会参与者
type MemorialServiceParticipant struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ServiceID string         `json:"service_id" gorm:"type:varchar(36);not null;index"`
	UserID    string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Role      string         `json:"role" gorm:"type:varchar(20);default:participant;comment:host|co_host|participant"`
	Status    string         `json:"status" gorm:"type:varchar(20);default:invited;comment:invited|joined|left|removed"`
	JoinedAt  *time.Time     `json:"joined_at"`
	LeftAt    *time.Time     `json:"left_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Service MemorialService `json:"service" gorm:"foreignKey:ServiceID"`
	User    User            `json:"user" gorm:"foreignKey:UserID"`
}

func (MemorialServiceParticipant) TableName() string {
	return "memorial_service_participants"
}

// ServiceActivity 追思会活动记录
type ServiceActivity struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ServiceID   string         `json:"service_id" gorm:"type:varchar(36);not null;index"`
	UserID      string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	ActivityType string        `json:"activity_type" gorm:"type:varchar(30);not null;comment:join|leave|worship|speak|share_screen"`
	Content     string         `json:"content" gorm:"type:json"`
	Timestamp   time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Service MemorialService `json:"service" gorm:"foreignKey:ServiceID"`
	User    User            `json:"user" gorm:"foreignKey:UserID"`
}

func (ServiceActivity) TableName() string {
	return "service_activities"
}

// ServiceInvitation 追思会邀请
type ServiceInvitation struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ServiceID string         `json:"service_id" gorm:"type:varchar(36);not null;index"`
	InviterID string         `json:"inviter_id" gorm:"type:varchar(36);not null;index"`
	InviteeID string         `json:"invitee_id" gorm:"type:varchar(36);not null;index"`
	Status    string         `json:"status" gorm:"type:varchar(20);default:pending;comment:pending|accepted|declined|expired"`
	Message   string         `json:"message" gorm:"type:text"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Service MemorialService `json:"service" gorm:"foreignKey:ServiceID"`
	Inviter User            `json:"inviter" gorm:"foreignKey:InviterID"`
	Invitee User            `json:"invitee" gorm:"foreignKey:InviteeID"`
}

func (ServiceInvitation) TableName() string {
	return "service_invitations"
}

// ServiceRecording 追思会录制
type ServiceRecording struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ServiceID   string         `json:"service_id" gorm:"type:varchar(36);not null;index"`
	RecordingURL string        `json:"recording_url" gorm:"type:varchar(255);not null"`
	ThumbnailURL string        `json:"thumbnail_url" gorm:"type:varchar(255)"`
	Duration    int            `json:"duration" gorm:"comment:录制时长(秒)"`
	FileSize    int64          `json:"file_size" gorm:"comment:文件大小(字节)"`
	Status      string         `json:"status" gorm:"type:varchar(20);default:processing;comment:processing|completed|failed"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Service MemorialService `json:"service" gorm:"foreignKey:ServiceID"`
}

func (ServiceRecording) TableName() string {
	return "service_recordings"
}

// ServiceChat 追思会聊天消息
type ServiceChat struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ServiceID string         `json:"service_id" gorm:"type:varchar(36);not null;index"`
	UserID    string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	MessageType string       `json:"message_type" gorm:"type:varchar(20);not null;comment:text|image|emoji"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Timestamp time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Service MemorialService `json:"service" gorm:"foreignKey:ServiceID"`
	User    User            `json:"user" gorm:"foreignKey:UserID"`
}

func (ServiceChat) TableName() string {
	return "service_chats"
}