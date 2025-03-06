package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/pkg/qqLogin"
	"tgwp/repo"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/jwtUtils"
	"tgwp/utils/snowflake"
	"time"
)

type QQLoginLogic struct {
}

func NewQQLoginLogic() *QQLoginLogic {
	return &QQLoginLogic{}
}
func (r *QQLoginLogic) QQLogin(ctx context.Context, req types.QQLoginReq) (types.QQLoginResp, error) {
	defer utils.RecordTime(time.Now())()
	info, err := qqLogin.GetUserInfo(req.Code)
	if err != nil {
		zlog.CtxInfof(ctx, "获取QQ用户信息失败:%v", err)
		return types.QQLoginResp{}, response.ErrResp(err, response.QQ_LOGIN_ERROR)
	}
	// 入库
	t := repo.NewQQLoginRepo(global.DB)
	user_id, err := t.IsExist(info.OpenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//生成雪花id
			user_id = snowflake.GetIntId(global.Node)
			user := model.User{
				ID:       user_id,
				OpenID:   info.OpenID,
				Avatar:   info.Avatar,
				Username: info.Nickname,
			}
			err = t.CreateUser(user)
			if err != nil {
				zlog.CtxInfof(ctx, "入库失败:%v", err)
				return types.QQLoginResp{}, response.ErrResp(err, response.DATABASE_ERROR)
			}
		} else {
			zlog.CtxInfof(ctx, "入库失败:%v", err)
			return types.QQLoginResp{}, response.ErrResp(err, response.DATABASE_ERROR)
		}
	}
	// 颁发token
	atoken, err := jwtUtils.GenToken(jwtUtils.FullToken(global.AUTH_ENUMS_ATOKEN, user_id))
	if err != nil {
		zlog.CtxErrorf(ctx, "GenToken err: %v", err)
		return types.QQLoginResp{}, response.ErrResp(err, response.GENERATE_TOKEN_ERROR)
	}
	rtoken, err := jwtUtils.GenToken(jwtUtils.FullToken(global.AUTH_ENUMS_RTOKEN, user_id))
	if err != nil {
		zlog.CtxErrorf(ctx, "GenToken err: %v", err)
		return types.QQLoginResp{}, response.ErrResp(err, response.GENERATE_TOKEN_ERROR)
	}
	return types.QQLoginResp{
		Atoken:   atoken,
		Avatar:   info.Avatar,
		Rtoken:   rtoken,
		Username: info.Nickname,
	}, nil
}
