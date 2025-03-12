package model

type Focus struct {
	CommonModel
	UserID  int64 `json:"user_id"`
	FocusID int64 `json:"focus_id"`
	User    User  `gorm:"foreignKey:FocusID;references:ID" json:"-"`
}

func (f *Focus) TableName() string {
	return "focus"
}
