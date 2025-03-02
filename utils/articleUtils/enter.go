package articleUtils

import (
	"context"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"tgwp/global"
	"tgwp/log/zlog"
)

type articleCacheType string

const (
	articleCacheDigg articleCacheType = "article_digg_key"
)

func set(ctx context.Context, key articleCacheType, articleID int64, n int) {
	num, err := global.Rdb.HGet(ctx, string(key), strconv.Itoa(int(articleID))).Int()
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		zlog.CtxErrorf(ctx, "redis文章处理缓存错误: %s\n", err)
		return
	}
	num += n
	err = global.Rdb.HSet(ctx, string(key), strconv.Itoa(int(articleID)), num).Err()
	if err != nil {
		zlog.CtxErrorf(ctx, "redis文章处理缓存错误: %s\n", err)
		return
	}
}
func SetCacheDigg(ctx context.Context, articleID int64, n int) {
	set(ctx, articleCacheDigg, articleID, n)
}
func get(ctx context.Context, key articleCacheType, articleID int64) int {
	num, err := global.Rdb.HGet(ctx, string(key), strconv.Itoa(int(articleID))).Int()
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		zlog.CtxErrorf(ctx, "redis文章处理缓存错误: %s\n", err)
		return 0
	}
	return num
}
func GetCacheDigg(ctx context.Context, articleID int64) int {
	return get(ctx, articleCacheDigg, articleID)
}
func getAll(ctx context.Context, t articleCacheType) (mps map[int64]int) {
	res, err := global.Rdb.HGetAll(ctx, string(t)).Result()
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		logrus.Errorf("redis文章处理缓存错误: %s\n", err)
		return
	}
	mps = make(map[int64]int)
	for key, value := range res {
		num, e := strconv.Atoi(key)
		if e != nil {
			zlog.CtxWarnf(ctx, "类型转换失败: %s\n", err)
			continue
		}
		v, e := strconv.Atoi(value)
		if e != nil {
			zlog.CtxWarnf(ctx, "类型转换失败: %s\n", err)
			continue
		}
		mps[int64(num)] = v
	}
	return
}
func GetCacheDiggList(ctx context.Context) (mps map[int64]int) {
	return getAll(ctx, articleCacheDigg)
}
func ClearCache(ctx context.Context) {
	err := global.Rdb.Del(ctx, string(articleCacheDigg)).Err()
	if err != nil {
		zlog.CtxErrorf(ctx, "redis文章处理缓存错误: %s\n", err)
	}
}
