package logic

import (
	"context"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo/list"
	"tgwp/types"
	"tgwp/utils"
	"time"
)

type IslandLogic struct {
}

func NewIslandLogic() *IslandLogic {
	return &IslandLogic{}
}
func (l *IslandLogic) GetIsland(ctx context.Context, req list.PageInfo) (resp types.IslandListResp, err error) {
	defer utils.RecordTime(time.Now())()
	_list, count, err := list.ListQuery(model.Island{}, list.Options{
		PageInfo: req,
	})
	if err != nil {
		zlog.CtxInfof(ctx, "获取岛屿列表失败:%v", err)
		return
	}
	var list = make([]types.IslandReq, 0)
	for _, v := range _list {
		list = append(list, types.IslandReq{
			ID:     v.ID,
			Name:   v.Name,
			Path:   v.Path,
			Height: v.Height,
			Width:  v.Width,
			XPoint: v.XPoint,
			YPoint: v.YPoint,
		})
	}
	resp = types.IslandListResp{
		List:  list,
		Count: count,
	}
	return
}
