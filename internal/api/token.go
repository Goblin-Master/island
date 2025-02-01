package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

// RefreshToken
//
//	@Description: 前端用rtoken刷新token
//	@param c
func RefreshToken(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.TokenReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "RefreshRtoken request: %v", req)
	resp, err := logic.NewTokenLogic().RefreshToken(ctx, req)
	response.Response(c, resp, err)
}
