package logic

import (
	"context"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo"
	"tgwp/repo/list"
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
	resp, err = repo.NewArticleRepo(global.DB).ArticleCreate(ctx, req)
	if err != nil {
		zlog.CtxInfof(ctx, "创建文章失败:%v", err)
		return types.ArticleCreateResp{}, response.ErrResp(err, response.ARTICLE_CREATE_ERROR)
	}
	return
}
func (l *ArticleLogic) ArticleList(ctx context.Context, req list.PageInfo) (resp types.ArticleList, err error) {
	defer utils.RecordTime(time.Now())()
	_list, count, err := list.ListQuery(model.Article{}, list.Options{
		PageInfo: req,
		Preloads: []string{"User"},
		Order:    "created_at desc",
		Likes:    []string{"title"},
	})
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return
	}
	var articleList = make([]types.Article, 0)
	for _, v := range _list {
		articleList = append(articleList, types.Article{
			CreatedAt: v.CreatedAt,
			DiggCount: v.DiggCount,
			Content:   v.Content,
			ID:        v.ID,
			Title:     v.Title,
			Cover:     v.Cover,
			Abstract:  v.Abstract,
			Avatar:    v.User.Avatar,
			Username:  v.User.Username,
			UserID:    v.UserID,
		})
	}
	resp = types.ArticleList{
		Count: count,
		List:  articleList,
	}
	return
}
func (l *ArticleLogic) ArticleDelete(ctx context.Context, req list.RemoveReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	db := repo.NewArticleRepo(global.DB)
	articleList, err := db.GetArticleByIds(req)
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return
	}
	resp, err = db.ArticleDelete(articleList)
	return
}

func (l *ArticleLogic) ArticleDigg(ctx context.Context, req types.ArticleDiggReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	db := repo.NewArticleRepo(global.DB)
	resp, err = db.ArticleDigg(ctx, req)
	return
}
