package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
)

func ArticleCreate(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	//BindReq里面用泛型进行了处理绑定
	req, err := types.BindReq[types.ArticleCreateReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "CreateArticle request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleCreate(ctx, req)
	response.Response(c, resp, err)
}

func ArticleList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	//BindReq里面用泛型进行了处理绑定
	req, err := types.BindReq[list.PageInfo](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ArticleList request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleList(ctx, req)
	response.Response(c, resp, err)
}

func ArticleDelete(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	var req list.RemoveReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "ArticleDelete request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleDelete(ctx, req)
	response.Response(c, resp, err)
}

func ArticleDigg(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.ArticleDiggReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ArticleDigg request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleDigg(ctx, req)
	response.Response(c, resp, err)
}
