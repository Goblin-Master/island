package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

func UserDetail(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UserDetailReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "UserDetail request: %v", req)
	resp, err := logic.NewUserLogic().UserDetail(ctx, req)
	response.Response(c, resp, err)
}
