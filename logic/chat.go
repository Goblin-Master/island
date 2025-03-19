package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
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
	REDIS_CHAT_MESSAGES    = "chat:island:%d:messages"
	REDIS_CHAT_UPDATE_TIME = "chat:island:%d:update_time"
	REDIS_CHAT_UPDATE_CHAN = "chat:island:%d:update_chan"
)

func (l *ChatLogic) SendMessage(ctx context.Context, req types.SendMessageReq) (resp types.SendMessageResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 生成雪花ID
	var id int64
	id = snowflake.GetIntId(global.Node)
	// IslandID 转化成 int64
	islandID, err := strconv.ParseInt(req.IslandID, 10, 64)
	if err != nil {
		zlog.Errorf("IslandID 转化成 int64 失败: %v", err)
		return resp, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	// 记录发送时间
	timestamp := time.Now().UnixMilli()
	resp.Timestamp = timestamp
	// 序列化为 JSON
	messageJSON, err := json.Marshal(types.Message{
		ID:        id,
		UserID:    req.UserID,
		Message:   req.Message,
		Timestamp: timestamp,
	})
	if err != nil {
		zlog.Errorf("序列化消息失败: %v", err)
		return resp, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	// 往redis中写入消息
	err = global.Rdb.ZAdd(ctx,
		fmt.Sprintf(REDIS_CHAT_MESSAGES, islandID),
		&redis.Z{
			Score:  float64(timestamp),
			Member: messageJSON,
		}).Err()
	// 记录最近一次更新时间为当前时间
	err = global.Rdb.Set(ctx,
		fmt.Sprintf(REDIS_CHAT_UPDATE_TIME, islandID),
		timestamp,
		0).Err()
	if err != nil {
		zlog.Errorf("记录最近一次更新时间失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 清理redis中一分钟前的消息
	minScore := float64(time.Now().Add(-time.Minute).UnixMilli())
	err = global.Rdb.ZRemRangeByScore(ctx,
		fmt.Sprintf(REDIS_CHAT_MESSAGES, islandID),
		"-inf",
		fmt.Sprintf("%f", minScore)).Err()
	if err != nil {
		zlog.Errorf("清理redis中一分钟前的消息失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 发布新消息通知
	err = global.Rdb.Publish(ctx, fmt.Sprintf(REDIS_CHAT_UPDATE_CHAN, islandID), "have_new_message").Err()
	if err != nil {
		zlog.Errorf("发布新消息通知失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	return
}

func (l *ChatLogic) GetMessages(ctx context.Context, req types.GetMessagesReq) (resp types.GetMessagesResp, err error) {
	defer utils.RecordTime(time.Now())()
	// IslandID 转化成 int64
	islandID, err := strconv.ParseInt(req.IslandID, 10, 64)
	if err != nil {
		zlog.Errorf("IslandID 转化成 int64 失败: %v", err)
		return resp, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	// 清理redis中一分钟前的消息
	minScore := float64(time.Now().Add(-time.Minute).UnixMilli())
	err = global.Rdb.ZRemRangeByScore(ctx,
		fmt.Sprintf(REDIS_CHAT_MESSAGES, islandID),
		"-inf",
		fmt.Sprintf("%f", minScore)).Err()
	if err != nil {
		zlog.Errorf("清理redis中一分钟前的消息失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 进入判断是否更新，如果需要更新，则进行长轮询等待
	var yes bool
	yes, err = checkUpdate(ctx, islandID, req.Timestamp)
	if !yes {
		// 长轮询等待超时，依旧没有更新，返回空数据
		return resp, nil
	}
	// 可以更新，读取redis从timestamp开始到现在的消息
	zlog.Debugf("开始读取redis中的消息")

	messages, err := global.Rdb.ZRangeByScore(ctx,
		fmt.Sprintf(REDIS_CHAT_MESSAGES, islandID),
		&redis.ZRangeBy{
			Min: fmt.Sprintf("%d", req.Timestamp),
			Max: "+inf",
		},
	).Result()

	if err != nil {
		zlog.Errorf("读取redis中的消息失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	zlog.Debugf("读取redis中的消息成功: %v", messages)
	// 反序列化 JSON
	for _, messageJSON := range messages {
		var message types.Message
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			zlog.Errorf("反序列化 JSON 失败: %v", err)
		}
		resp.Messages = append(resp.Messages, message)
		resp.Length += 1
	}
	return
}

func checkUpdate(ctx context.Context, islandID int64, timestamp int64) (bool, error) {
	// 读取最近一次更新时间
	yes, err := checkUpdateByTimestamp(ctx, islandID, timestamp)
	if err != nil {
		return false, err
	}
	if yes {
		return true, nil
	}
	// 没有更新，建立Redis订阅准备进行长轮询
	pubsub := global.Rdb.Subscribe(ctx, fmt.Sprintf(REDIS_CHAT_UPDATE_CHAN, islandID))
	defer pubsub.Close()
	// 期间再次判断有无消息更新，避免无故等待
	yes, err = checkUpdateByTimestamp(ctx, islandID, timestamp)
	if err != nil {
		return false, err
	}
	if yes {
		return true, nil
	}
	// 长轮询等待
	zlog.Debugf("岛屿 %d 未找到新消息，开始长轮询等待新消息通知", islandID)
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-pubsub.Channel():
			// 收到新消息通知，检查是否有新数据
			zlog.Debugf("岛屿 %d 收到新消息通知", islandID)
			yes, err = checkUpdateByTimestamp(ctx, islandID, timestamp)
			if err != nil {
				return false, err
			}
			if yes {
				return true, nil
			}
		case <-timeout:
			// 超时返回
			zlog.Debugf("岛屿 %d 长轮询等待超时", islandID)
			return false, nil
		case <-ctx.Done():
			// 客户端取消请求
			zlog.Debugf("岛屿 %d 长轮询等待取消", islandID)
			return false, ctx.Err()
		}
	}
}

func checkUpdateByTimestamp(ctx context.Context, islandID int64, timestamp int64) (bool, error) {
	// 读取最近一次更新时间
	key := fmt.Sprintf(REDIS_CHAT_UPDATE_TIME, islandID)
	lastUpdateTime, err := global.Rdb.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return true, nil // 键不存在视为需要更新
	} else if err != nil {
		zlog.Errorf("读取最近一次更新时间失败: %v", err)
		return false, response.ErrResp(err, response.REDIS_ERROR)
	}
	zlog.Debugf("最近一次更新时间: %v", lastUpdateTime)
	// 比较时间戳
	if lastUpdateTime < timestamp {
		return false, nil
	}
	// 有新消息，直接更新
	return true, nil
}
