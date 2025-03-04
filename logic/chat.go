package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/snowflake"
	"time"
)

type ChatLogic struct {
}

func NewChatLogic() *ChatLogic {
	return &ChatLogic{}
}

const (
	REDIS_CHAT = "chat:island:%d"
)

func (l *ChatLogic) SendMessage(ctx context.Context, req types.SendMessageReq) (resp types.SendMessageResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 生成雪花ID
	var id int64
	id = snowflake.GetIntId(global.Node)
	// 记录发送时间
	timestamp := time.Now().UnixMilli()
	resp.Timestamp = timestamp
	// 序列化为 JSON
	messageJSON, err := json.Marshal(types.Message{
		ID:      id,
		UserID:  req.UserID,
		Message: req.Message,
	})
	if err != nil {
		zlog.Errorf("序列化消息失败: %v", err)
		return resp, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	// 往redis中写入消息
	err = global.Rdb.ZAdd(ctx,
		fmt.Sprintf(REDIS_CHAT, req.IslandID),
		&redis.Z{
			Score:  float64(timestamp),
			Member: messageJSON,
		}).Err()
	// 清理redis中一分钟前的消息
	minScore := float64(time.Now().Add(-time.Minute).UnixMilli())
	err = global.Rdb.ZRemRangeByScore(ctx,
		fmt.Sprintf(REDIS_CHAT, req.IslandID),
		"-inf",
		fmt.Sprintf("%f", minScore)).Err()
	if err != nil {
		zlog.Errorf("清理redis中一分钟前的消息失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	return
}
