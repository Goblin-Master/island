package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/pkg/ai"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/cacheUtils"
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
	userid := int64(793478004095)
	history, err := cacheUtils.GetContent(ctx, strconv.FormatInt(userid, 10))
	if err != nil {
		zlog.CtxErrorf(ctx, "获取历史记录失败 %s", err)
		return
	}
	msgChan, err := ai.ChatStream(ctx, req.Content, history)
	if err != nil {
		response.SSEFail(err.Error(), c)
		return
	}
	var reply string
	for msg := range msgChan {
		reply += msg
		response.SSESuccess(msg, c)
	}
	// 保存历史记录
	cacheUtils.SaveContent(ctx, strconv.FormatInt(userid, 10), history, fmt.Sprintf("user:%s,assistant:%s ", req.Content, reply))
}
