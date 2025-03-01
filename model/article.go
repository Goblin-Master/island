package model

type Article struct {
	CommonModel
	Cover     string `gorm:"column:cover;type:varchar(255);not null;comment:'封面'"`
	Content   string `gorm:"column:content;type:longtext;not null;comment:'内容'"`
	Abstract  string `gorm:"column:abstract;type:varchar(255);not null;comment:'摘要'"`
	Title     string `gorm:"column:title;type:varchar(64);not null;comment:'标题'"`
	UserID    int64  `gorm:"column:userid;type:bigint;not null;comment:'用户id'"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	DiggCount int64  `gorm:"column:digg_count;type:bigint;not null;comment:'点赞数'"`
}

func (a *Article) TableName() string {
	return "article"
}
