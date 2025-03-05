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
	"tgwp/utils/snowflake"
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
		return types.IslandListResp{}, response.ErrResp(err, response.ISLAND_GET_ERROR)
	}
	var islandList = make([]types.ListResp, 0)
	for _, v := range _list {
		islandList = append(islandList, types.ListResp{
			ID:     v.ID,
			Name:   v.Name,
			Path:   v.Path,
			Height: v.Height,
			Width:  v.Width,
			XPoint: v.XPoint,
			YPoint: v.YPoint,
			UserID: v.UserID,
		})
	}
	resp = types.IslandListResp{
		List:  islandList,
		Count: count,
	}
	return
}

func (l *IslandLogic) CreateIsland(ctx context.Context, req types.IslandReq) (resp types.IslandResp, err error) {
	defer utils.RecordTime(time.Now())()
	//生成雪花id
	id := snowflake.GetIntId(global.Node)
	// TODO:创建岛屿得登录，拿用户的id
	userid := int64(793478004095)
	r := repo.NewIslandRepo(global.DB)
	if r.IdentifyIslandName(req.Name) {
		zlog.CtxInfof(ctx, "岛屿名字重复:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_EXIST)
	}
	var island = model.Island{
		ID:     id,
		Name:   req.Name,
		Path:   req.Path,
		Height: req.Height,
		Width:  req.Width,
		XPoint: req.XPoint,
		YPoint: req.YPoint,
		UserID: userid,
	}
	err = r.CreateIsland(island)
	if err != nil {
		zlog.CtxInfof(ctx, "创建岛屿失败:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	resp = types.IslandResp{
		ID:     id,
		UserID: userid,
		Name:   req.Name,
	}
	return
}

func (l *IslandLogic) ModifyIsland(ctx context.Context, req types.IslandReq) (resp types.IslandResp, err error) {
	defer utils.RecordTime(time.Now())()
	// TODO:创建岛屿得登录，拿用户的id
	user_id := int64(793478004095)
	r := repo.NewIslandRepo(global.DB)
	island_id, e := strconv.ParseInt(req.ID, 10, 64)
	if e != nil {
		zlog.CtxErrorf(ctx, "类型转换:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	if !r.IdentifyIslandById(island_id, user_id) {
		zlog.CtxInfof(ctx, "修改岛屿权限不足:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_NOT_UPDATE)
	}
	if r.IdentifyIslandNameAndId(req.Name, island_id) {
		zlog.CtxInfof(ctx, "岛屿名字重复:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_EXIST)
	}
	var island = model.Island{
		ID:     island_id,
		Name:   req.Name,
		Path:   req.Path,
		Height: req.Height,
		Width:  req.Width,
		XPoint: req.XPoint,
		YPoint: req.YPoint,
		UserID: user_id,
	}
	err = r.UpdatesIsland(island)
	if err != nil {
		zlog.CtxInfof(ctx, "修改岛屿失败:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_UPDATE_ERROR)
	}
	resp = types.IslandResp{
		ID:     island_id,
		UserID: user_id,
		Name:   req.Name,
	}
	return
}
func (l *IslandLogic) DeleteIsland(ctx context.Context, req types.IslandDeleteReq) (err error) {
	defer utils.RecordTime(time.Now())()
	// TODO:创建岛屿得登录，拿用户的id
	userid := int64(793478004095)
	island_id, e := strconv.ParseInt(req.ID, 10, 64)
	if e != nil {
		zlog.CtxErrorf(ctx, "类型转换:%v", err)
		return response.ErrResp(err, response.COMMON_FAIL)
	}
	err = repo.NewIslandRepo(global.DB).DeleteIsland(island_id, userid)
	if err != nil {
		zlog.CtxInfof(ctx, "删除岛屿失败:%v", err)
		return response.ErrResp(err, response.ISLAND_DELETE_ERROR)
	}
	return
}
