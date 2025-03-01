package logic

import (
	"context"
	"tgwp/log/zlog"
	"tgwp/repo"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"time"
)

type ArticleLogic struct {
}

func NewArticleLogic() *ArticleLogic {
	return &ArticleLogic{}
}
func (l *ArticleLogic) ArticleCreate(ctx context.Context, req types.ArticleCreateReq) (resp types.ArticleCreateResp, err error) {
	defer utils.RecordTime(time.Now())()
	//TODO: 拿用户id
	resp, err = repo.NewArticleRepo().ArticleCreate(ctx, req)
	if err != nil {
		zlog.CtxInfof(ctx, "创建文章失败:%v", err)
		return types.ArticleCreateResp{}, response.ErrResp(err, response.ARTICLE_CREATE_ERROR)
	}
	return
}
