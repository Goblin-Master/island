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
	var list []model.Article
	global.DB.Find(&list)
	for _, t := range list {
		digg := diggMap[t.ID]
		if digg == 0 {
			continue
		}
		err := global.DB.Model(&t).Updates(map[string]any{
			"digg_count": t.DiggCount + digg,
		}).Error
		if err != nil {
			zlog.CtxErrorf(ctx, "更新文章失败 %s", err)
			continue
		}
		zlog.CtxInfof(ctx, "id为%d的文章更新成功", t.ID)
	}
	// 回填增量
	_diggMap := articleUtils.GetCacheDiggList(ctx)
	//同步回去
	for _, t := range list {
		digg := _diggMap[t.ID] - diggMap[t.ID]
		if digg == 0 {
			continue
		}
		articleUtils.SetCacheDigg(ctx, t.ID, digg)
	}
	articleUtils.ClearCache(ctx)
}
