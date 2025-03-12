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

type UserLogic struct {
}

func NewUserLogic() *UserLogic {
	return &UserLogic{}
}
func (l *UserLogic) UserDetail(ctx context.Context, req types.UserDetailReq) (resp types.UserDetailResp, err error) {
	defer utils.RecordTime(time.Now())()
	user_id, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "类型转换:%v", err)
		return types.UserDetailResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	resp, err = repo.NewUserRepo(global.DB).UserDetail(user_id)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户详情失败:%v", err)
		return types.UserDetailResp{}, response.ErrResp(err, response.GET_USER_ERROR)
	}
	return
}

func (l *UserLogic) FocusUser(ctx context.Context, req types.FocusUserReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	focus_id, err := strconv.ParseInt(req.FocusID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "类型转换:%v", err)
		return "", response.ErrResp(err, response.COMMON_FAIL)
	}
	if focus_id == req.UserID {
		return "", response.ErrResp(err, response.USER_FOCUS_SELF)
	}
	r := repo.NewUserRepo(global.DB)
	err = r.IsExist(focus_id)
	if err != nil {
		zlog.CtxErrorf(ctx, "用户不存在:%v", err)
		return "", response.ErrResp(err, response.USER_NOT_EXIST)
	}
	resp, err = r.FocusUser(req.UserID, focus_id)
	if err != nil {
		zlog.CtxErrorf(ctx, "用户关注相关操作失败:%v", err)
	}
	return
}

func (l *UserLogic) FocusList(ctx context.Context, req types.FocusListReq) (resp types.FocusListResp, err error) {
	defer utils.RecordTime(time.Now())()
	user_id, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "类型转换:%v", err)
		return types.FocusListResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	focusIDList, err := repo.NewUserRepo(global.DB).FocusIDList(user_id)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户关注列表失败:%v", err)
		return types.FocusListResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	if len(focusIDList) == 0 {
		return types.FocusListResp{}, nil
	}
	_list, count, err := list.ListQuery(model.User{}, list.Options{
		PageInfo: req.PageInfo,
		Order:    "created_at desc",
		Where:    global.DB.Where("id in ?", focusIDList),
	})
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户关注列表失败:%v", err)
		return types.FocusListResp{}, response.ErrResp(err, response.COMMON_FAIL)
	}
	var focusList []types.FocusList
	for _, v := range _list {
		focusList = append(focusList, types.FocusList{
			Avatar:    v.Avatar,
			CreatedAt: v.CreatedAt,
			UserID:    v.ID,
			Username:  v.Username,
		})
	}
	return types.FocusListResp{
		Count: count,
		List:  focusList,
	}, nil
}
