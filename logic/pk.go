package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sort"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/snowflake"
	"time"
)

type PKLogic struct {
}

func NewPKLogic() *PKLogic {
	return &PKLogic{}
}

const (
	REDIS_PK_WAITING_LIST = "pk:question_bank:%d:waiting_list"
	REDIS_PK_WAITING_SUB  = "pk:question_bank:%d:waiting_sub"
	REDIS_ROOM_INFO       = "pk:room:%d:info"
)

// PKSetRule 设置 PK 规则
func (l *PKLogic) PKSetRule(ctx context.Context, req types.PKSetRuleReq) (resp types.PKSetRuleResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转 int64
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 判断题目数量是否合法
	var cnt int
	cnt, err = repo.NewQuestionRepo(global.DB).GetQuestionBankQuestionCount(questionBankID)
	if req.QuestionCount > cnt {
		zlog.CtxErrorf(ctx, "题库 %d 题目数量少于 %d", questionBankID, req.QuestionCount)
		return resp, response.ErrResp(err, response.QUESTION_COUNT_NOT_ENOUGH)
	}
	// 修改数据库规则
	err = repo.NewQuestionRepo(global.DB).PKSetRule(questionBankID, req.QuestionCount, req.Duration)
	if err != nil {
		zlog.CtxErrorf(ctx, "修改数据库规则错误: %v", err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	return resp, nil
}

// PKGetRule 获取 PK 规则
func (l *PKLogic) PKGetRule(ctx context.Context, req types.PKGetRuleReq) (resp types.PKGetRuleResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转 int64
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 获取规则
	var rule model.QuestionBankPKRule
	rule, err = repo.NewQuestionRepo(global.DB).PKGetRule(questionBankID)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取题库 %d PK 规则错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	// 转换为返回结构
	resp.QuestionCount = rule.QuestionCount
	resp.Duration = rule.Duration
	return resp, nil
}

func (l *PKLogic) PKMatching(ctx context.Context, req types.PKMatchingReq) (resp types.PKMatchingResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转 int64
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 如果规则不存在，说明此题库不支持 PK
	if !repo.NewQuestionRepo(global.DB).CheckPKRuleExist(questionBankID) {
		zlog.CtxErrorf(ctx, "题库 %d 未设置 PK 规则", questionBankID)
	}
	// 题库中题目数量至少要大于规则要求的数量
	var cnt int
	cnt, err = repo.NewQuestionRepo(global.DB).GetQuestionBankQuestionCount(questionBankID)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取题库 %d 题目数量错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	var rule model.QuestionBankPKRule
	rule, err = repo.NewQuestionRepo(global.DB).PKGetRule(questionBankID)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取题库 %d PK 规则错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	if cnt < rule.QuestionCount {
		zlog.CtxErrorf(ctx, "题库 %d 题目数量不足", questionBankID)
		return resp, response.ErrResp(err, response.QUESTION_COUNT_NOT_ENOUGH)
	}
	// 判断等待队列是否已经有人
	var val string
	val, err = global.Rdb.Get(ctx, fmt.Sprintf(REDIS_PK_WAITING_LIST, questionBankID)).Result()
	if errors.Is(err, redis.Nil) {
		// 没有人在队列中，加入匹配队列长轮询
		// 先创建一个房间号
		id := snowflake.GetIntId(global.Node)
		val = fmt.Sprintf("%d", id)
		err = global.Rdb.Set(ctx, fmt.Sprintf(REDIS_PK_WAITING_LIST, questionBankID), id, 5*time.Minute).Err()
		if err != nil {
			zlog.CtxErrorf(ctx, "redis 设置等待队列错误: %v", err)
			return resp, response.ErrResp(err, response.REDIS_ERROR)
		}
		// 订阅匹配队列，等待有人匹配并通知自己
		var user2ID int64
		resp.Code, user2ID, err = WaitPKMatching(ctx, questionBankID)
		// 先清除等待队列
		zlog.CtxDebugf(ctx, "清除等待队列")
		err = global.Rdb.Del(context.Background(), fmt.Sprintf(REDIS_PK_WAITING_LIST, questionBankID)).Err()
		if err != nil {
			zlog.CtxErrorf(ctx, "redis 清除等待队列错误: %v", err)
			return resp, response.ErrResp(err, response.REDIS_ERROR)
		}
		// 如果超时和取消
		if resp.Code == -1 {
			resp.Msg = "匹配超时"
			return resp, nil
		} else if resp.Code == -2 {
			resp.Msg = "取消匹配"
			return resp, nil
		}
		// 匹配成功，处理房间逻辑
		InitRoom(id, questionBankID, req.UserID, user2ID, rule)
	} else if err == nil {
		// 队列中有人，直接匹配成功，向对方发送订阅信息
		err = global.Rdb.Publish(ctx, fmt.Sprintf(REDIS_PK_WAITING_SUB, questionBankID), req.UserID).Err()
		if err != nil {
			zlog.CtxErrorf(ctx, "redis 发布订阅信息错误: %v", err)
			return resp, response.ErrResp(err, response.REDIS_ERROR)
		}
	} else {
		zlog.CtxErrorf(ctx, "redis 获取等待队列错误: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	resp.Code = 0
	resp.Msg = "匹配成功"
	// 房间号转 int64
	resp.RoomID, err = strconv.ParseInt(val, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", val, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	return
}

func (l *PKLogic) GetRoomInfo(ctx context.Context, req types.GetRoomInfoReq) (resp types.GetRoomInfoResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转 int64
	var roomID int64
	roomID, err = strconv.ParseInt(req.RoomID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", roomID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 从 redis 中获取房间信息
	redisRoomInfoJSON, err := global.Rdb.Get(ctx, fmt.Sprintf(REDIS_ROOM_INFO, roomID)).Result()
	if errors.Is(err, redis.Nil) {
		zlog.CtxErrorf(ctx, "房间 %d 不存在", roomID)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 解析 json 字符串
	var redisRoomInfo types.RedisRoomInfo
	err = json.Unmarshal([]byte(redisRoomInfoJSON), &redisRoomInfo)
	if err != nil {
		zlog.CtxErrorf(ctx, "json 解析房间信息错误: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 统计数据
	allSubmit := true
	redisRoomInfo.User1ScoreTotal = 0
	redisRoomInfo.User2ScoreTotal = 0
	for i, item := range redisRoomInfo.Questions {
		if !item.User1Submit || !item.User2Submit {
			allSubmit = false
		}
		redisRoomInfo.User1ScoreTotal += item.User1Score
		redisRoomInfo.User2ScoreTotal += item.User2Score
		// 处理发生请求用户的分数
		if req.UserID == redisRoomInfo.User1ID {
			redisRoomInfo.Questions[i].YourSubmit = item.User1Submit
			redisRoomInfo.Questions[i].YourScore = item.User1Score
		} else {
			redisRoomInfo.Questions[i].YourSubmit = item.User2Submit
			redisRoomInfo.Questions[i].YourScore = item.User2Score
		}
	}
	// 判断是否已结束游戏(时间结束或双方均已全部交卷)
	if time.Now().UnixMilli() > redisRoomInfo.EndTimestamp || allSubmit {
		if redisRoomInfo.User1ScoreTotal == redisRoomInfo.User2ScoreTotal {
			// 同分，比较谁先达到这个分数
			if redisRoomInfo.User1FinalSubmitTimestamp < redisRoomInfo.User2FinalSubmitTimestamp {
				redisRoomInfo.WinnerID = redisRoomInfo.User1ID
			} else {
				redisRoomInfo.WinnerID = redisRoomInfo.User2ID
			}
		} else {
			if redisRoomInfo.User1ScoreTotal > redisRoomInfo.User2ScoreTotal {
				redisRoomInfo.WinnerID = redisRoomInfo.User1ID
			} else {
				redisRoomInfo.WinnerID = redisRoomInfo.User2ID
			}
		}
	}
	// 转换成返回结构
	resp = types.GetRoomInfoResp{
		User1ID:                   redisRoomInfo.User1ID,
		User2ID:                   redisRoomInfo.User2ID,
		Questions:                 redisRoomInfo.Questions,
		User1ScoreTotal:           redisRoomInfo.User1ScoreTotal,
		User2ScoreTotal:           redisRoomInfo.User2ScoreTotal,
		User1FinalSubmitTimestamp: redisRoomInfo.User1FinalSubmitTimestamp,
		User2FinalSubmitTimestamp: redisRoomInfo.User2FinalSubmitTimestamp,
		StartTimestamp:            redisRoomInfo.StartTimestamp,
		EndTimestamp:              redisRoomInfo.EndTimestamp,
		WinnerID:                  redisRoomInfo.WinnerID,
	}
	return
}

func (l *PKLogic) SubmitQuestion(ctx context.Context, req types.SubmitQuestionReq) (resp types.SubmitQuestionResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 如果答案为空，则不提交
	if len(req.Answer) == 0 {
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 转 int64
	roomID, err := strconv.ParseInt(req.RoomID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", roomID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	questionID, err := strconv.ParseInt(req.QuestionID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}

	// 获取分布式锁
	lockKey := fmt.Sprintf("lock:room:%d", roomID)
	lockAcquired, err := global.Rdb.SetNX(ctx, lockKey, "locked", 10*time.Second).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "获取锁失败: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	if !lockAcquired {
		zlog.CtxErrorf(ctx, "房间 %d 操作冲突", roomID)
		return resp, response.ErrResp(errors.New("系统繁忙，请稍后再试"), response.REDIS_ERROR)
	}
	defer func() {
		if _, err := global.Rdb.Del(ctx, lockKey).Result(); err != nil {
			zlog.CtxErrorf(ctx, "释放锁失败: %v", err)
		}
	}()

	// 从 redis 中获取房间信息
	redisRoomInfoJSON, err := global.Rdb.Get(ctx, fmt.Sprintf(REDIS_ROOM_INFO, roomID)).Result()
	if errors.Is(err, redis.Nil) {
		zlog.CtxErrorf(ctx, "房间 %d 不存在", roomID)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}

	// 解析 json 字符串
	var redisRoomInfo types.RedisRoomInfo
	err = json.Unmarshal([]byte(redisRoomInfoJSON), &redisRoomInfo)
	if err != nil {
		zlog.CtxErrorf(ctx, "json 解析房间信息错误: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 寻找提交的题目的索引
	questionIndex := -1
	for i, item := range redisRoomInfo.Questions {
		if item.QuestionID == questionID {
			questionIndex = i
			break
		}
	}
	if questionIndex == -1 {
		zlog.CtxErrorf(ctx, "题目 %d 不在房间题目列表中", questionID)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 判断用户是否是房间的参与者
	if req.UserID != redisRoomInfo.User1ID && req.UserID != redisRoomInfo.User2ID {
		zlog.CtxErrorf(ctx, "用户 %d 不在房间中", req.UserID)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 判断是否已经交过
	if (redisRoomInfo.Questions[questionIndex].User1Submit && req.UserID == redisRoomInfo.User1ID) || (redisRoomInfo.Questions[questionIndex].User2Submit && req.UserID == redisRoomInfo.User2ID) {
		zlog.CtxErrorf(ctx, "用户 %d 已经交过题目 %d", req.UserID, questionID)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 判断答案得分
	question, err := repo.NewQuestionRepo(global.DB).GetQuestion(questionID)
	score := int64(0)
	// 答案 json 反序列化
	var answer []string
	err = json.Unmarshal([]byte(question.Answer), &answer)
	if err != nil {
		zlog.CtxErrorf(ctx, "json 反序列化答案错误: %v", err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	if question.Type == global.QUESTION_TYPE_CHOICE {
		// 单选题
		if req.Answer[0] == answer[0] {
			score = 100
		}
	} else if question.Type == global.QUESTION_TYPE_MULTIPLE {
		// 多选题
		allRight := true
		if len(req.Answer) != len(answer) {
			allRight = false
		} else {
			zlog.Debugf("刚开始: req.Answer: %v, answer: %v", req.Answer, answer)
			// 对两个答案切片进行排序
			sort.Strings(req.Answer)
			sort.Strings(answer)
			zlog.Debugf("排序后: req.Answer: %v, answer: %v", req.Answer, answer)
			for i := 0; i < len(answer); i++ {
				if req.Answer[i] != answer[i] {
					allRight = false
					break
				}
			}
		}
		// 两个答案都进行排序

		if allRight {
			score = 100
		}
	} else if question.Type == global.QUESTION_TYPE_BLANK {
		// 填空题
		for i := 0; i < len(answer); i++ {
			if req.Answer[0] == answer[i] {
				score = 100
				break
			}
		}
	}
	// 记录得分和提交时间
	if req.UserID == redisRoomInfo.User1ID {
		redisRoomInfo.Questions[questionIndex].User1Submit = true
		redisRoomInfo.Questions[questionIndex].User1Score = score
		redisRoomInfo.User1FinalSubmitTimestamp = time.Now().UnixMilli()
		if req.UserID == redisRoomInfo.User1ID {
			redisRoomInfo.Questions[questionIndex].YourSubmit = true
			redisRoomInfo.Questions[questionIndex].YourScore = score
		}
	} else {
		redisRoomInfo.Questions[questionIndex].User2Score = score
		redisRoomInfo.Questions[questionIndex].User2Submit = true
		redisRoomInfo.User2FinalSubmitTimestamp = time.Now().UnixMilli()
		if req.UserID == redisRoomInfo.User2ID {
			redisRoomInfo.Questions[questionIndex].YourSubmit = true
			redisRoomInfo.Questions[questionIndex].YourScore = score
		}
	}
	// 重新计算总分
	redisRoomInfo.User1ScoreTotal = 0
	redisRoomInfo.User2ScoreTotal = 0
	for _, item := range redisRoomInfo.Questions {
		redisRoomInfo.User1ScoreTotal += item.User1Score
		redisRoomInfo.User2ScoreTotal += item.User2Score
	}
	// 计算缓存时间
	cacheTime := time.Duration((redisRoomInfo.EndTimestamp-time.Now().UnixMilli())/1000)*time.Second + 5*time.Minute
	zlog.Debugf("房间 %d 缓存时间: %v", roomID, cacheTime)
	// 转 json 字符串并存入 redis
	newRedisRoomInfoJSON, _ := json.Marshal(redisRoomInfo)
	err = global.Rdb.Set(ctx, fmt.Sprintf(REDIS_ROOM_INFO, roomID), newRedisRoomInfoJSON, cacheTime).Err()
	if err != nil {
		zlog.CtxErrorf(ctx, "redis 设置房间信息错误: %v", err)
		return resp, response.ErrResp(err, response.REDIS_ERROR)
	}
	// 发送提交结果
	resp.Score = score
	return
}

func WaitPKMatching(ctx context.Context, questionBankID int64) (code int, user2ID int64, err error) {
	// 订阅匹配队列，等待有人匹配并通知自己
	pubsub := global.Rdb.Subscribe(ctx, fmt.Sprintf(REDIS_PK_WAITING_SUB, questionBankID))
	defer pubsub.Close()
	// 长轮询等待
	timeout := time.After(5 * time.Minute)
	for {
		select {
		case temp := <-pubsub.Channel():
			// 有人加入队列，匹配成功
			user2ID, err = strconv.ParseInt(temp.Payload, 10, 64)
			return 0, user2ID, nil
		case <-timeout:
			// 超时
			return -1, 0, nil
		case <-ctx.Done():
			// 取消订阅
			zlog.CtxDebugf(ctx, "取消订阅")
			return -2, 0, nil
		}
	}
}

// InitRoom 初始化房间
func InitRoom(roomID int64, questionBankID int64, user1ID int64, user2ID int64, rule model.QuestionBankPKRule) {
	// 搭建房间信息
	zlog.Debugf("房间 %d 开始搭建中，选手为: %d 与 %d", roomID, user1ID, user2ID)
	// 随机选取题目
	questions, err := repo.NewQuestionRepo(global.DB).GetRandomQuestions(questionBankID, rule.QuestionCount)
	if err != nil {
		zlog.Errorf("获取题目错误: %v", err)
		return
	}
	// 编辑初始房间信息
	nowTimestamp := time.Now().UnixMilli()
	redisRoomInfo := types.RedisRoomInfo{
		RoomID:                    roomID,
		User1ID:                   user1ID,
		User2ID:                   user2ID,
		Questions:                 make([]types.QuestionsInfo, 0),
		User1ScoreTotal:           0,
		User2ScoreTotal:           0,
		User1FinalSubmitTimestamp: nowTimestamp,
		User2FinalSubmitTimestamp: nowTimestamp,
		StartTimestamp:            nowTimestamp,
		EndTimestamp:              nowTimestamp + int64(rule.Duration)*1000,
		WinnerID:                  0,
	}
	for _, item := range questions {
		redisRoomInfo.Questions = append(redisRoomInfo.Questions, types.QuestionsInfo{
			QuestionID:  item.ID,
			User1Submit: false,
			User2Submit: false,
			User1Score:  0,
			User2Score:  0,
			YourSubmit:  false,
			YourScore:   0,
		})
	}
	// 转 json 字符串并存入 redis
	redisRoomInfoJSON, _ := json.Marshal(redisRoomInfo)
	zlog.CtxDebugf(context.Background(), "房间信息: %s", redisRoomInfoJSON)
	err = global.Rdb.Set(context.Background(), fmt.Sprintf(REDIS_ROOM_INFO, roomID), redisRoomInfoJSON, time.Duration(rule.Duration)*time.Second+5*time.Minute).Err()
	if err != nil {
		zlog.Errorf("redis 设置房间信息错误: %v", err)
		return
	}
}
