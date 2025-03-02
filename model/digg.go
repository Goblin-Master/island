package model

import "time"

// 文章点赞表
type Digg struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time
	UserId    int64 `gorm:"column:user_id;type:bigint;not null"`
	ArticleId int64 `gorm:"column:article_id;type:bigint;not null"`
}
