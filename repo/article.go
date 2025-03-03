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
		Island:   req.Island,
		Cover:    req.Cover,
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
	err = r.DB.Where("user_id = ? and article_id = ?", req.UserID, req.ArticleID).Take(&digg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			digg = model.Digg{
				ArticleID: req.ArticleID,
				UserID:    req.UserID,
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
func (r *ArticleRepo) ArticleCollect(ctx context.Context, req types.ArticleCollectReq) (resp string, err error) {
	var collect model.Collect
	err = r.DB.Where("user_id = ? and article_id = ?", req.UserID, req.ArticleID).Take(&collect).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collect = model.Collect{
				ArticleID: req.ArticleID,
				UserID:    req.UserID,
			}
			err = r.DB.Create(&collect).Error
			if err != nil {
				zlog.CtxErrorf(ctx, "收藏失败:%v", err)
				return "", response.ErrResp(err, response.ARTICLE_COLLECT_ERROR)
			}
			// 收藏成功，设置缓存
			articleUtils.SetCacheCollect(ctx, req.ArticleID, 1)
			return "收藏成功", nil
		}
		zlog.CtxErrorf(ctx, "收藏时数据库出现错误:%v", err)
		return "", response.ErrResp(err, response.ARTICLE_COLLECT_ERROR)
	}
	// 已经点赞过了
	if err = r.DB.Delete(&collect).Error; err != nil {
		zlog.CtxErrorf(ctx, "取消收藏失败:%v", err)
		return "", response.ErrResp(err, response.ARTICLE_COLLECT_ERROR)
	}
	// 取消点赞成功，设置缓存
	articleUtils.SetCacheCollect(ctx, req.ArticleID, -1)
	return "取消收藏成功", nil
}

func (r *ArticleRepo) CollectList(ctx context.Context, userid int64) (articleIDList []int64, err error) {
	err = r.DB.Model(&model.Collect{}).Where("user_id = ?", userid).Pluck("article_id", &articleIDList).Error
	if err != nil {
		zlog.CtxErrorf(ctx, "获取收藏列表失败:%v", err)
		return nil, response.ErrResp(err, response.GET_COLLECT_ARTICLE_ERROR)
	}
	return
}
func (r *ArticleRepo) DiggList(ctx context.Context, userid int64) (articleIDList []int64, err error) {
	err = r.DB.Model(&model.Digg{}).Where("user_id = ?", userid).Pluck("article_id", &articleIDList).Error
	if err != nil {
		zlog.CtxErrorf(ctx, "获取点赞列表失败:%v", err)
		return nil, response.ErrResp(err, response.GET_DIGG_ARTICLE_ERROR)
	}
	return
}
