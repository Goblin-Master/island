package repo

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/articleUtils"
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
func (r *ArticleRepo) ArticleDigg(ctx context.Context, req types.ArticleDiggReq) (resp string, err error) {
	var digg model.Digg
	err = r.DB.Where("article_id = ? and user_id = ?", req.ArticleID, req.UserID).Take(&digg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			digg = model.Digg{
				ArticleId: req.ArticleID,
				UserId:    req.UserID,
			}
			err = r.DB.Create(&digg).Error
			if err != nil {
				zlog.CtxErrorf(ctx, "点赞失败:%v", err)
				return "", response.ErrResp(err, response.ARTICLE_DIGG_ERROR)
			}
			// 点赞成功，设置缓存
			articleUtils.SetCacheDigg(ctx, req.ArticleID, 1)
			return "点赞成功", nil
		}
		zlog.CtxErrorf(ctx, "点赞时数据库出现错误:%v", err)
		return "", response.ErrResp(err, response.ARTICLE_DIGG_ERROR)
	}
	// 已经点赞过了
	if err = r.DB.Delete(&digg).Error; err != nil {
		zlog.CtxErrorf(ctx, "取消点赞失败:%v", err)
		return "", response.ErrResp(err, response.ARTICLE_DIGG_ERROR)
	}
	// 取消点赞成功，设置缓存
	articleUtils.SetCacheDigg(ctx, req.ArticleID, -1)
	return "取消点赞成功", nil
}
