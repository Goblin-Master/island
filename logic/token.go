package logic

import (
	"context"
	"github.com/gin-gonic/gin"
	"tgwp/global"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/jwtUtils"
)

type TokenLogic struct {
}

func NewTokenLogic() *TokenLogic {
	return &TokenLogic{}
}

// RefreshToken
//
//	@Description: 用于rtoken刷新atoken
//	@receiver l
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (l *TokenLogic) RefreshToken(ctx context.Context, req types.TokenReq) (resp types.TokenResp, err error) {
	//解析token是否有效，并取出上一次的值
	data, err := jwtUtils.IdentifyToken(ctx, req.Token)
	if err != nil {
		//对应token无效，直接让他返回
		return resp, err
	}
	//判断其是否为rtoken
	if data.Class != global.AUTH_ENUMS_RTOKEN {
		return resp, response.ErrResp(err, response.TOKEN_TYPE_ERROR)
	}
	//生成新的token
	resp.Atoken, err = jwtUtils.GenToken(jwtUtils.FullToken(global.AUTH_ENUMS_ATOKEN, data.Userid))
	if err != nil {
		return resp, response.ErrResp(err, response.GENERATE_TOKEN_ERROR)
	}
	resp.Rtoken, err = jwtUtils.GenToken(jwtUtils.TokenData{
		Class:  global.AUTH_ENUMS_RTOKEN,
		Time:   data.Time,
		Userid: data.Userid,
	})
	if err != nil {
		return resp, response.ErrResp(err, response.GENERATE_TOKEN_ERROR)
	}
	return resp, nil
}

// GetUserId
//
//	@Description: 这个方法只能是使用了鉴权中间件才能用
//	@param c
//	@return int64
func GetUserId(c *gin.Context) int64 {
	if data, exists := c.Get(global.TOKEN_USER_ID); exists {
		userid, ok := data.(int64)
		if ok {
			return userid
		}
	}
	return 0
}
