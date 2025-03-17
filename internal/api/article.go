package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/jwtUtils"
)

func ArticleCreate(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	//BindReq里面用泛型进行了处理绑定
	req, err := types.BindReq[types.ArticleCreateReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "CreateArticle request: %v", req.Title)
	req.UserID = jwtUtils.GetUserId(c)
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
	var req types.ArticleRemoveReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "ArticleDelete request: %v", req)
	req.UserID = jwtUtils.GetUserId(c)
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
	if req.UserID == "" {
		req.UserID = strconv.FormatInt(jwtUtils.GetUserId(c), 10)
	}
	resp, err := logic.NewArticleLogic().ArticleDigg(ctx, req)
	response.Response(c, resp, err)
}
func ArticleCollect(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.ArticleCollectReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ArticleCollect request: %v", req)
	if req.UserID == "" {
		req.UserID = strconv.FormatInt(jwtUtils.GetUserId(c), 10)
	}
	resp, err := logic.NewArticleLogic().ArticleCollect(ctx, req)
	response.Response(c, resp, err)
}
func ArticleListByUerID(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	//BindReq里面用泛型进行了处理绑定
	req, err := types.BindReq[types.ArticleListReq](c)
	if err != nil {
		return
	}
	if req.UserID == "" {
		req.UserID = strconv.FormatInt(jwtUtils.GetUserId(c), 10)
	}
	zlog.CtxInfof(ctx, "ArticleListByUerID request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleListByUserID(ctx, req)
	response.Response(c, resp, err)
}

func ArticleCollectList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[list.PageInfo](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ArticleCollectList request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleCollectList(ctx, req, jwtUtils.GetUserId(c))
	response.Response(c, resp, err)
}

func ArticleDiggList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[list.PageInfo](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ArticleDiggList request: %v", req)
	resp, err := logic.NewArticleLogic().ArticleDiggList(ctx, req, jwtUtils.GetUserId(c))
	response.Response(c, resp, err)
}
