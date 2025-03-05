package aiUtils

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
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
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func InsertHistory(ctx context.Context, key, question, reply string) {
	// 只保留三轮对话
	val, err := global.Rdb.LLen(ctx, key).Result()
	if err != nil {
		zlog.Errorf("redis获取对话历史错误:%v", err)
	}
	if val > 6 {
		// 删除最老的轮对话
		global.Rdb.LPop(ctx, key)
		global.Rdb.LPop(ctx, key)
	}
	// 插入一轮对话
	err = global.Rdb.RPush(ctx, key, question, reply).Err()
	if err != nil {
		zlog.Errorf("redis存储对话历史错误:%v", err)
	}
}

func GetHistory(ctx context.Context, key string) ([]string, error) {
	// 获取所有对话
	value, err := global.Rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return value, nil
}

func ClearHistory(ctx context.Context, key string) (err error) {
	return global.Rdb.Del(ctx, key).Err()
}
