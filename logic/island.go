package logic

import (
	"context"
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
	var list = make([]types.ListReq, 0)
	for _, v := range _list {
		list = append(list, types.ListReq{
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
		List:  list,
		Count: count,
	}
	return
}

func (l *IslandLogic) CreateIsland(ctx context.Context, req types.IslandReq) (resp types.IslandResp, err error) {
	defer utils.RecordTime(time.Now())()
	//雪花id的生成格式
	node, err := snowflake.NewNode(global.DEFAULT_NODE_ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "NewNode err: %v", err)
		return types.IslandResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	//一般是生成12位的int64id，也可以生成string的，看snowflakes包
	id := snowflake.GetInt12Id(node)
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
	userid := int64(793478004095)
	r := repo.NewIslandRepo(global.DB)
	if !r.IdentifyIslandById(req.ID, userid) {
		zlog.CtxInfof(ctx, "修改岛屿权限不足:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_NOT_UPDATE)
	}
	if r.IdentifyIslandName(req.Name) {
		zlog.CtxInfof(ctx, "岛屿名字重复:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_EXIST)
	}
	var island = model.Island{
		Name:   req.Name,
		Path:   req.Path,
		Height: req.Height,
		Width:  req.Width,
		XPoint: req.XPoint,
		YPoint: req.YPoint,
		UserID: userid,
	}
	err = r.UpdatesIsland(island)
	if err != nil {
		zlog.CtxInfof(ctx, "修改岛屿失败:%v", err)
		return types.IslandResp{}, response.ErrResp(err, response.ISLAND_UPDATE_ERROR)
	}
	resp = types.IslandResp{
		ID:     req.ID,
		UserID: userid,
		Name:   req.Name,
	}
	return
}
func (l *IslandLogic) DeleteIsland(ctx context.Context, res types.IslandDeleteReq) (err error) {
	defer utils.RecordTime(time.Now())()
	// TODO:创建岛屿得登录，拿用户的id
	userid := int64(793478004095)
	err = repo.NewIslandRepo(global.DB).DeleteIsland(res.ID, userid)
	if err != nil {
		zlog.CtxInfof(ctx, "删除岛屿失败:%v", err)
		return response.ErrResp(err, response.ISLAND_DELETE_ERROR)
	}
	return
}
