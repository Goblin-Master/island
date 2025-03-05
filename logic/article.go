package logic

import (
	"context"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/articleUtils"
	"time"
)

type ArticleLogic struct {
}

func NewArticleLogic() *ArticleLogic {
	return &ArticleLogic{}
}
func (l *ArticleLogic) ArticleCreate(ctx context.Context, req types.ArticleCreateReq) (resp types.ArticleCreateResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 判断岛屿名称是否存在
	if !repo.NewIslandRepo(global.DB).IdentifyIslandName(req.Island) {
		zlog.CtxInfof(ctx, "岛屿不存在")
		return types.ArticleCreateResp{}, response.ErrResp(err, response.ISLAND_NOT_EXIST)
	}
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
		Likes:    []string{"title", "island"},
	})
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return
	}
	var articleList = make([]types.Article, 0)
	for _, v := range _list {
		articleList = append(articleList, types.Article{
			CreatedAt: v.CreatedAt,
			//从缓存同步数据
			DiggCount:    v.DiggCount + articleUtils.GetCacheDigg(ctx, v.ID),
			CollectCount: v.CollectCount + articleUtils.GetCacheCollect(ctx, v.ID),
			Content:      v.Content,
			ID:           v.ID,
			Title:        v.Title,
			Cover:        v.Cover,
			Abstract:     v.Abstract,
			Avatar:       v.User.Avatar,
			Username:     v.User.Username,
			UserID:       v.UserID,
			Island:       v.Island,
		})
	}
	resp = types.ArticleList{
		Count: count,
		List:  articleList,
	}
	return
}
func (l *ArticleLogic) ArticleDelete(ctx context.Context, req types.ArticleRemoveReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	var _req list.RemoveReq
	for _, v := range req.Ids {
		id, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			zlog.CtxErrorf(ctx, "类型转换失败 %s", e)
			return "", response.ErrResp(err, response.COMMON_FAIL)
		}
		_req.Ids = append(_req.Ids, id)
	}
	db := repo.NewArticleRepo(global.DB)
	articleList, err := db.GetArticleByIds(_req)
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return
	}
	resp, err = db.ArticleDelete(articleList)
	return
}

func (l *ArticleLogic) ArticleDigg(ctx context.Context, req types.ArticleDiggReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	// TODO:获取用户id填进去
	user_id := int64(793478004095)
	article_id, e := strconv.ParseInt(req.ArticleID, 10, 64)
	if e != nil {
		zlog.CtxErrorf(ctx, "类型转换失败 %s", e)
		return "", response.ErrResp(err, response.COMMON_FAIL)
	}
	db := repo.NewArticleRepo(global.DB)
	resp, err = db.ArticleDigg(ctx, user_id, article_id)
	return
}

func (l *ArticleLogic) ArticleCollect(ctx context.Context, req types.ArticleCollectReq) (resp string, err error) {
	defer utils.RecordTime(time.Now())()
	// TODO:获取用户id填进去
	user_id := int64(793478004095)
	article_id, e := strconv.ParseInt(req.ArticleID, 10, 64)
	if e != nil {
		zlog.CtxErrorf(ctx, "类型转换失败 %s", e)
		return "", response.ErrResp(err, response.COMMON_FAIL)
	}
	db := repo.NewArticleRepo(global.DB)
	resp, err = db.ArticleCollect(ctx, user_id, article_id)
	return
}
func (l *ArticleLogic) ArticleListByUserID(ctx context.Context, req list.PageInfo) (resp types.ArticleList, err error) {
	defer utils.RecordTime(time.Now())()
	_list, count, err := list.ListQuery(model.Article{
		// TODO:获取用户id填进去
		UserID: int64(793478004095),
	}, list.Options{
		PageInfo: req,
		Preloads: []string{"User"},
		Order:    "created_at desc",
		Likes:    []string{"title", "island"},
	})
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return
	}
	var articleList = make([]types.Article, 0)
	for _, v := range _list {
		articleList = append(articleList, types.Article{
			CreatedAt: v.CreatedAt,
			//从缓存同步数据
			DiggCount:    v.DiggCount + articleUtils.GetCacheDigg(ctx, v.ID),
			CollectCount: v.CollectCount + articleUtils.GetCacheCollect(ctx, v.ID),
			Content:      v.Content,
			ID:           v.ID,
			Title:        v.Title,
			Cover:        v.Cover,
			Abstract:     v.Abstract,
			Avatar:       v.User.Avatar,
			Username:     v.User.Username,
			UserID:       v.UserID,
			Island:       v.Island,
		})
	}
	resp = types.ArticleList{
		Count: count,
		List:  articleList,
	}
	return
}
func (l *ArticleLogic) ArticleCollectList(ctx context.Context, req list.PageInfo) (resp types.ArticleList, err error) {
	defer utils.RecordTime(time.Now())()
	//TODO: 获取用户id填进去
	//先查用户收藏表，获取所有文章id
	collectList, err := repo.NewArticleRepo(global.DB).CollectList(ctx, int64(793478004095))
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return types.ArticleList{}, response.ErrResp(err, response.GET_COLLECT_ARTICLE_ERROR)
	}
	if len(collectList) == 0 {
		return types.ArticleList{}, nil
	}
	_list, count, err := list.ListQuery(model.Article{}, list.Options{
		PageInfo: req,
		Preloads: []string{"User"},
		Order:    "created_at desc",
		Likes:    []string{"title", "island"},
		Where:    global.DB.Where("id in ?", collectList),
	})
	var articleList = make([]types.Article, 0)
	for _, v := range _list {
		articleList = append(articleList, types.Article{
			CreatedAt: v.CreatedAt,
			//从缓存同步数据
			DiggCount:    v.DiggCount + articleUtils.GetCacheDigg(ctx, v.ID),
			CollectCount: v.CollectCount + articleUtils.GetCacheCollect(ctx, v.ID),
			Content:      v.Content,
			ID:           v.ID,
			Title:        v.Title,
			Cover:        v.Cover,
			Abstract:     v.Abstract,
			Avatar:       v.User.Avatar,
			Username:     v.User.Username,
			UserID:       v.UserID,
			Island:       v.Island,
		})
	}
	resp = types.ArticleList{
		Count: count,
		List:  articleList,
	}
	return
}
func (l *ArticleLogic) ArticleDiggList(ctx context.Context, req list.PageInfo) (resp types.ArticleList, err error) {
	defer utils.RecordTime(time.Now())()
	//TODO: 获取用户id填进去
	//先查用户点赞表，获取所有文章id
	diggList, err := repo.NewArticleRepo(global.DB).DiggList(ctx, int64(793478004095))
	if err != nil {
		zlog.CtxInfof(ctx, "获取文章列表失败:%v", err)
		return types.ArticleList{}, response.ErrResp(err, response.GET_DIGG_ARTICLE_ERROR)
	}
	if len(diggList) == 0 {
		return types.ArticleList{}, nil
	}
	_list, count, err := list.ListQuery(model.Article{}, list.Options{
		PageInfo: req,
		Preloads: []string{"User"},
		Order:    "created_at desc",
		Likes:    []string{"title", "island"},
		Where:    global.DB.Where("id in ?", diggList),
	})
	var articleList = make([]types.Article, 0)
	for _, v := range _list {
		articleList = append(articleList, types.Article{
			CreatedAt: v.CreatedAt,
			//从缓存同步数据
			DiggCount:    v.DiggCount + articleUtils.GetCacheDigg(ctx, v.ID),
			CollectCount: v.CollectCount + articleUtils.GetCacheCollect(ctx, v.ID),
			Content:      v.Content,
			ID:           v.ID,
			Title:        v.Title,
			Cover:        v.Cover,
			Abstract:     v.Abstract,
			Avatar:       v.User.Avatar,
			Username:     v.User.Username,
			UserID:       v.UserID,
			Island:       v.Island,
		})
	}
	resp = types.ArticleList{
		Count: count,
		List:  articleList,
	}
	return
}
