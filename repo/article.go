package repo

import (
	"context"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/response"
	"tgwp/types"
)

type ArticleRepo struct {
}

func NewArticleRepo() *ArticleRepo {
	return &ArticleRepo{}
}
func (r *ArticleRepo) ArticleCreate(ctx context.Context, req types.ArticleCreateReq) (resp types.ArticleCreateResp, err error) {
	article := model.Article{
		Abstract: req.Abstract,
		Content:  req.Content,
		Title:    req.Title,
		UserID:   int64(793478004095),
	}
	err = global.DB.Create(&article).Error
	if err != nil {
		zlog.CtxErrorf(ctx, "创建文章失败:%v", err)
		return types.ArticleCreateResp{}, response.ErrResp(err, response.ARTICLE_CREATE_ERROR)
	}
	resp.ID = article.ID
	return
}
