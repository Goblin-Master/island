package model

type Image struct {
	CommonModel
	Filename string `gorm:"column:filename;type:varchar(64);not null;comment:'文件名'"`
	Path     string `gorm:"column:path;type:varchar(255);not null;comment:'文件路径'"`
	Size     int64  `gorm:"column:size;type:bigint;not null;comment:'文件大小'"`
	Hash     string `gorm:"column:hash;type:varchar(32);not null;comment:'文件哈希'"`
}

func (i *Image) TableName() string {
	return "image"
}
