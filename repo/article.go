package repo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
)

type ArticleRepo struct {
	DB *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{
		DB: db,
	}
}
func (r *ArticleRepo) ArticleCreate(ctx context.Context, req types.ArticleCreateReq) (resp types.ArticleCreateResp, err error) {
	article := model.Article{
		Abstract: req.Abstract,
		Content:  req.Content,
		Title:    req.Title,
		UserID:   int64(793478004095),
	}
	err = r.DB.Create(&article).Error
	if err != nil {
		zlog.CtxErrorf(ctx, "创建文章失败:%v", err)
		return types.ArticleCreateResp{}, response.ErrResp(err, response.ARTICLE_CREATE_ERROR)
	}
	resp.ID = article.ID
	return
}
func (r *ArticleRepo) GetArticleByIds(req list.RemoveReq) (resp []model.Article, err error) {
	err = r.DB.Debug().Where("id in ?", req.Ids).Find(&resp).Error
	return
}
func (r *ArticleRepo) ArticleDelete(req []model.Article) (resp string, err error) {
	var successCount int64
	if len(req) > 0 {
		successCount = global.DB.Debug().Delete(&req).RowsAffected
	}
	failCount := int64(len(req)) - successCount
	resp = fmt.Sprintf("操作成功，成功删除了%d篇文章，失败删除了%d篇文章", successCount, failCount)
	return
}
