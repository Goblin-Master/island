package model

import (
	"time"
)

type Island struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint;uniqueIndex:idx_name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string  `gorm:"column:name;type:varchar(64);not null;unique;comment:'名称'"`
	Path      string  `gorm:"column:path;type:varchar(255);not null;comment:'岛屿图片'"`
	Width     float64 `gorm:"column:width;type:float;not null;comment:'宽度'"`
	Height    float64 `gorm:"column:height;type:float;not null;comment:'高度'"`
	XPoint    float64 `gorm:"column:xPoint;type:float;not null;comment:'x坐标'"`
	YPoint    float64 `gorm:"column:yPoint;type:float;not null;comment:'y坐标'"`
	UserID    int64   `gorm:"column:user_id;type:bigint;not null;comment:'用户id';uniqueIndex:idx_name"`
}

func (i *Island) TableName() string {
	return "island"
}
