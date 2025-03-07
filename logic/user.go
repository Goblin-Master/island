package logic

import (
	"context"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/repo"
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
