package model

type Article struct {
	CommonModel
	Island       string `gorm:"column:island;type:varchar(64);not null;comment:'岛屿'"`
	Cover        string `gorm:"column:cover;type:varchar(255);not null;comment:'封面'"`
	Content      string `gorm:"column:content;type:longtext;not null;comment:'内容'"`
	Abstract     string `gorm:"column:abstract;type:varchar(255);not null;comment:'摘要'"`
	Title        string `gorm:"column:title;type:varchar(64);not null;comment:'标题'"`
	UserID       int64  `gorm:"column:user_id;type:bigint;not null;comment:'用户id'"`
	User         User   `gorm:"foreignKey:UserID;references:ID"`
	DiggCount    int    `gorm:"column:digg_count;type:bigint;not null;comment:'点赞数'"`
	CollectCount int    `gorm:"column:collect_count;type:bigint;not null;comment:'收藏数'"`
}

func (a *Article) TableName() string {
	return "article"
}
