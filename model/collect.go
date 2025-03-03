package model

import "time"

// 文章收藏表
type Collect struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time
	UserID    int64 `gorm:"column:user_id;type:bigint;not null"`
	ArticleID int64 `gorm:"column:article_id;type:bigint;not null"`
}

func (c *Collect) TableName() string {
	return "collect"
}
