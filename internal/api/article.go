package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
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
