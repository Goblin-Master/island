package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/pkg/ai"
	"tgwp/response"
	"tgwp/types"
)

// AiGenerateAbstract
//
//	@Description: ai生成简介
//	@param c
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

// AiChatStream
//
//	@Description: 流式ai对话
//	@param c
func AiChatStream(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.AiReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, " AiChatStream request: %s", req)
	msgChan, err := ai.ChatStream(ctx, req.Content)
	if err != nil {
		response.SSEFail(err.Error(), c)
		return
	}
	for msg := range msgChan {
		response.SSESuccess(msg, c)
	}
}
