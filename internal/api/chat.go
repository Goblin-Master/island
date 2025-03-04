package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

func ChatSendMessage(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.SendMessageReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ChatSendMessage request: %v", req)
	resp, err := logic.NewChatLogic().SendMessage(ctx, req)
	response.Response(c, resp, err)
}
