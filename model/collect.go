package model

import "time"

// 文章收藏表
type Collect struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time
	UserId    int64 `gorm:"column:user_id;type:bigint;not null;uniqueIndex:idx_user_collect"`
	ArticleId int64 `gorm:"column:article_id;type:bigint;not null;uniqueIndex:idx_user_collect"`
}

func (c *Collect) TableName() string {
	return "collect"
}
