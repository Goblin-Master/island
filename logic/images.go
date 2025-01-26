package logic

import (
	"context"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo"
	"tgwp/repo/list"
	"tgwp/types"
	"tgwp/utils"
	"time"
)

type ImagesLogic struct{}

func NewImagesLogic() *ImagesLogic {
	return &ImagesLogic{}
}
func (r *ImagesLogic) GetImages(ctx context.Context, req list.PageInfo) (resp types.ImageListResp, err error) {
	defer utils.RecordTime(time.Now())()
	_list, count, err := list.ListQuery(model.Image{}, list.Options{
		PageInfo: req,
	})
	if err != nil {
		zlog.CtxInfof(ctx, "获取图片列表失败:%v", err)
		return
	}
	var list = make([]types.ImageResp, 0)
	for _, v := range _list {
		list = append(list, types.ImageResp{
			ID:       v.ID,
			Filename: v.Filename,
			Hash:     v.Hash,
			Path:     v.Path,
			Size:     v.Size,
			WebPath:  v.WebPath(),
		})
	}
	resp = types.ImageListResp{
		Count: count,
		List:  list,
	}
	return
}
func (r *ImagesLogic) DeleteImages(ctx context.Context, req list.RemoveReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	db := repo.NewImagesRepo(global.DB)
	imageList, err := db.GetImagesByIds(req)
	if err != nil {
		zlog.CtxInfof(ctx, "获取图片列表失败:%v", err)
		return
	}
	resp, err = db.DeleteImages(imageList)
	return
}
