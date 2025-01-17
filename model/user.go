package model

type User struct {
	CommonModel
	Username string `gorm:"column:username;type:varchar(32);not null;comment:'用户名'"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);not null;comment:'头像'"`
	OpenID   string `gorm:"column:openid;type:varchar(32);not null;comment:'QQ登录的唯一id'"`
}

func (u *User) TableName() string {
	return "user"
}
