package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/jwtUtils"
)

func PKMatching(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.PKMatchingReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "PKMatching request: %v", req)
	req.UserID = jwtUtils.GetUserId(c)
	resp, err := logic.NewPKLogic().PKMatching(ctx, req)
	response.Response(c, resp, err)
}

func GetRoomInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetRoomInfoReq](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "PKMatching request: %v", req)
	resp, err := logic.NewPKLogic().GetRoomInfo(ctx, req)
	response.Response(c, resp, err)
}

func SubmitQuestion(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.SubmitQuestionReq](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "PKMatching request: %v", req)
	resp, err := logic.NewPKLogic().SubmitQuestion(ctx, req)
	response.Response(c, resp, err)
}
