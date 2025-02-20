package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
)

func GetIsland(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[list.PageInfo](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "GetIsland request: %v", req)
	resp, err := logic.NewIslandLogic().GetIsland(ctx, req)
	response.Response(c, resp, err)
}
