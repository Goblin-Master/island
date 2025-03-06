package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/jwtUtils"
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
func CreateIsland(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.IslandReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "CreateIsland request: %v", req)
	req.UserID = jwtUtils.GetUserId(c)
	resp, err := logic.NewIslandLogic().CreateIsland(ctx, req)
	response.Response(c, resp, err)
}
func ModifyIsland(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.IslandReq](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "ModifyIsland request: %v", req)
	resp, err := logic.NewIslandLogic().ModifyIsland(ctx, req)
	response.Response(c, resp, err)
}

func DeleteIsland(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	var req types.IslandDeleteReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "DeleteIsland request: %v", req)
	req.UserID = jwtUtils.GetUserId(c)
	err = logic.NewIslandLogic().DeleteIsland(ctx, req)
	response.Response(c, nil, err)
}
