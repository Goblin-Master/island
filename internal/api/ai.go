package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

func AiGenerateAbstract(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.AiReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, " AiGenerateAbstract request: %s", req)
	resp, err := logic.NewAILogic().GenerateAbstract(ctx, req)
	response.Response(c, resp, err)
}
