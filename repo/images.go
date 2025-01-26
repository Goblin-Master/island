package repo

import (
	"fmt"
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/model"
	"tgwp/repo/list"
)

type ImagesRepo struct {
	DB *gorm.DB
}

func NewImagesRepo(db *gorm.DB) *ImagesRepo {
	return &ImagesRepo{
		DB: db,
	}
}
func (r *ImagesRepo) GetImagesByIds(req list.RemoveReq) (resp []model.Image, err error) {
	err = r.DB.Debug().Where("id in ?", req.Ids).Find(&resp).Error
	return
}
func (r *ImagesRepo) DeleteImages(req []model.Image) (resp string, err error) {
	var successCount int64
	if len(req) > 0 {
		successCount = global.DB.Debug().Delete(&req).RowsAffected
	}
	failCount := int64(len(req)) - successCount
	resp = fmt.Sprintf("操作成功，成功删除了%d张图片，失败删除了%d张图片", successCount, failCount)
	return
}
