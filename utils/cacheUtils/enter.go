package cacheUtils

import (
	"context"
	"tgwp/global"
	"tgwp/log/zlog"
	"time"
)

func SaveContent(ctx context.Context, key string, history, reply string) {
	if len(history)+len(reply) > 4000 {
		err := global.Rdb.Set(ctx, key, reply, time.Second*3600).Err()
		if err != nil {
			zlog.Errorf("redis存储错误:%v", err)
		}
	}
	err := global.Rdb.Set(ctx, key, history+reply, time.Second*3600).Err()
	if err != nil {
		zlog.Errorf("redis存储错误:%v", err)
	}
}
func GetContent(ctx context.Context, key string) (string, error) {
	value, err := global.Rdb.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", nil
		}
		return "", err
	}
	return value, nil
}
