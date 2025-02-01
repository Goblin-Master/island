package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

func QQLogin(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.QQLoginReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "QQLogin request: %s", req)
	resp, err := logic.NewQQLoginLogic().QQLogin(ctx, req)
	response.Response(c, resp, err)
}
