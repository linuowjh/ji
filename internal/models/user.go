package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(36);comment:用户ID"`
	WechatOpenID  string         `json:"wechatOpenid" gorm:"column:wechat_open_id;uniqueIndex;type:varchar(100);not null;comment:微信OpenID"`
	WechatUnionID string         `json:"wechatUnionid" gorm:"column:wechat_union_id;type:varchar(100);comment:微信UnionID"`
	Nickname      string         `json:"nickname" gorm:"type:varchar(50);comment:用户昵称"`
	AvatarURL     string         `json:"avatarUrl" gorm:"column:avatar_url;type:varchar(255);comment:头像URL"`
	Phone         string         `json:"phone" gorm:"type:varchar(20);comment:手机号"`
	Status        int            `json:"status" gorm:"default:1;comment:状态:1正常 0禁用"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (User) TableName() string {
	return "users"
}
