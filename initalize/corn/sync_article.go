package corn

import (
	"context"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/utils/articleUtils"
)

func SyncArticle() {
	// 从缓存获取数据
	ctx := context.Background()
	zlog.CtxInfof(ctx, "开始同步文章点赞数")
	diggMap := articleUtils.GetCacheDiggList(ctx)
	collectMap := articleUtils.GetCacheCollectList(ctx)
	var list []model.Article
	global.DB.Find(&list)
	for _, t := range list {
		digg := diggMap[t.ID]
		collect := collectMap[t.ID]
		if digg == 0 && collect == 0 {
			continue
		}
		err := global.DB.Model(&t).Updates(map[string]any{
			"digg_count":    t.DiggCount + digg,
			"collect_count": t.CollectCount + collect,
		}).Error
		if err != nil {
			zlog.CtxErrorf(ctx, "更新文章失败 %s", err)
			continue
		}
		zlog.CtxInfof(ctx, "id为%d的文章更新成功", t.ID)
	}
	// 回填增量
	_diggMap := articleUtils.GetCacheDiggList(ctx)
	_collectMap := articleUtils.GetCacheCollectList(ctx)
	//同步回去
	for _, t := range list {
		digg := _diggMap[t.ID] - diggMap[t.ID]
		collect := _collectMap[t.ID] - collectMap[t.ID]
		if digg == 0 && collect == 0 {
			continue
		}
		articleUtils.SetCacheDigg(ctx, t.ID, digg)
		articleUtils.SetCacheCollect(ctx, t.ID, collect)
	}
	articleUtils.ClearCache(ctx)
}
