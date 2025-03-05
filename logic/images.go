package logic

import (
	"context"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo"
	"tgwp/repo/list"
	"tgwp/response"
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
	var imageList = make([]types.ImageResp, 0)
	for _, v := range _list {
		imageList = append(imageList, types.ImageResp{
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
		List:  imageList,
	}
	return
}
func (r *ImagesLogic) DeleteImages(ctx context.Context, req types.ImageRemoveReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	db := repo.NewImagesRepo(global.DB)
	var _req list.RemoveReq
	for _, v := range req.Ids {
		id, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			zlog.CtxErrorf(ctx, "类型转换失败 %s", e)
			return "", response.ErrResp(err, response.COMMON_FAIL)
		}
		_req.Ids = append(_req.Ids, id)
	}
	imageList, err := db.GetImagesByIds(_req)
	if err != nil {
		zlog.CtxInfof(ctx, "获取图片列表失败:%v", err)
		return
	}
	resp, err = db.DeleteImages(imageList)
	return
}
