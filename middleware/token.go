package middleware

import (
	"github.com/gin-gonic/gin"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/response"
	"tgwp/utils/jwtUtils"
)

func Authentication(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	token := c.GetHeader("Authorization")
	if token == "" {
		zlog.CtxErrorf(ctx, "token为空")
		response.NewResponse(c).Error(response.TOKEN_IS_BLANK)
		c.Abort()
		return
	}
	//解析token是否有效，并取出上一次的值
	data, err := jwtUtils.IdentifyToken(ctx, token)
	if err != nil {
		zlog.CtxErrorf(ctx, "token验证失败:%v", err)
		response.NewResponse(c).Error(response.TOKEN_IS_EXPIRED)
		//对应token无效，直接让他返回
		c.Abort()
		return
	}
	//判断其是否为atoken
	if data.Class != global.AUTH_ENUMS_ATOKEN {
		zlog.CtxErrorf(ctx, "token类型错误")
		response.NewResponse(c).Error(response.TOKEN_TYPE_ERROR)
		c.Abort()
		return
	}
	//将token内部数据传下去,在logic.token内有对应方法获取userid
	c.Set(global.TOKEN_USER_ID, data.Userid)
	c.Next()
}
