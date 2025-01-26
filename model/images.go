package model

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"tgwp/configs"
	"tgwp/log/zlog"
)

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
func (i *Image) WebPath() string {
	return fmt.Sprintf("http://127.0.0.1:%d/%s", configs.Conf.App.Port, i.Path)
}

func (i *Image) BeforeDelete(tx *gorm.DB) (err error) {
	//删除文件
	err = os.Remove(i.Path)
	if err != nil {
		zlog.Warnf("删除图片失败 %s", i.Path)
	}
	return
}
