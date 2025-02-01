package model

import (
	"gorm.io/gorm"
	"time"
)

// 不用CommonModel是因为我想自己定义ID，而不是自动生成
type User struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"column:username;type:varchar(32);not null;comment:'用户名'"`
	Avatar    string         `gorm:"column:avatar;type:varchar(255);not null;comment:'头像'"`
	OpenID    string         `gorm:"column:openid;type:varchar(32);not null;comment:'QQ登录的唯一id'"`
}

func (u *User) TableName() string {
	return "user"
}
