package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/jwtUtils"
)

func ChatSendMessage(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.SendMessageReq](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "ChatSendMessage request: %v", req)
	resp, err := logic.NewChatLogic().SendMessage(ctx, req)
	response.Response(c, resp, err)
}

func ChatGetMessages(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetMessagesReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ChatGetMessages request: %v", req)
	resp, err := logic.NewChatLogic().GetMessages(ctx, req)
	response.Response(c, resp, err)
}
