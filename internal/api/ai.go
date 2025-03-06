package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/pkg/ai_eino"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/aiUtils"
	"tgwp/utils/jwtUtils"
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
	zlog.CtxInfof(ctx, " AiGenerateAbstract request: %v", req)
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
	zlog.CtxInfof(ctx, " AiChatStream request: %v", req)
	req.UserID = jwtUtils.GetUserId(c)
	outStream, err := ai_eino.StreamChat(ctx, req)
	if err != nil {
		response.SSEFail(err.Error(), c)
		return
	}
	defer outStream.Close() // 注意要关闭
	var reply string
	for {
		chunk, err := outStream.Recv()
		if err != nil {
			response.SSEFail(err.Error(), c)
			break
		}
		reply += chunk.Content
		response.SSESuccess(chunk.Content, c)
	}
	// 保存历史记录
	aiUtils.InsertHistory(ctx, strconv.FormatInt(req.UserID, 10), req.Content, reply)
}

func ClearHistory(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	zlog.CtxInfof(ctx, "ClearHistory request")
	err := logic.NewAILogic().ClearHistory(ctx, jwtUtils.GetUserId(c))
	response.Response(c, nil, err)
}
