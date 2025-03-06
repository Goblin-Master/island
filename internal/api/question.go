package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

func CreateQuestion(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateQuestionReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "CreateQuestion request: %v", req)
	resp, err := logic.NewQuestionLogic().CreateQuestion(ctx, req)
	response.Response(c, resp, err)
}

func CreateQuestionBank(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateQuestionBankReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "CreateQuestion request: %v", req)
	resp, err := logic.NewQuestionLogic().CreateQuestionBank(ctx, req)
	response.Response(c, resp, err)
}

func AddQuestion(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.AddQuestionReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "AddQuestion request: %v", req)
	resp, err := logic.NewQuestionLogic().AddQuestion(ctx, req)
	response.Response(c, resp, err)
}

func AddQuestionBank(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.AddQuestionBankReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "AddQuestionBank request: %v", req)
	resp, err := logic.NewQuestionLogic().AddQuestionBank(ctx, req)
	response.Response(c, resp, err)
}

func GetQuestion(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetQuestionReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "AddQuestionBank request: %v", req)
	resp, err := logic.NewQuestionLogic().GetQuestion(ctx, req)
	response.Response(c, resp, err)
}

func GetQuestionBank(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetQuestionBankReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "AddQuestionBank request: %v", req)
	resp, err := logic.NewQuestionLogic().GetQuestionBank(ctx, req)
	response.Response(c, resp, err)
}

func GetQuestionList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetQuestionListReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "AddQuestionBank request: %v", req)
	resp, err := logic.NewQuestionLogic().GetQuestionList(ctx, req)
	response.Response(c, resp, err)
}

func GetQuestionBankList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetQuestionBankListReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "AddQuestionBank request: %v", req)
	resp, err := logic.NewQuestionLogic().GetQuestionBankList(ctx, req)
	response.Response(c, resp, err)
}

func RunCode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.RunCodeReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "RunCode request: %v", req)
	resp, err := logic.NewQuestionLogic().RunCode(ctx, req)
	response.Response(c, resp, err)
}
