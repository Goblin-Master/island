package cacheUtils

import (
	"context"
	"tgwp/global"
	"tgwp/log/zlog"
	"time"
)

// SaveContent
//
//	@Description: 用于ai对话保存上下文
//	@param ctx
//	@param key
//	@param history
//	@param reply
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

// GetContent
//
//	@Description: 获取ai对话的上下文
//	@param ctx
//	@param key
//	@return string
//	@return error
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
