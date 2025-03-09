package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/repo"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/runcodeUtils"
	"tgwp/utils/snowflake"
	"time"
)

type QuestionLogic struct {
}

func NewQuestionLogic() *QuestionLogic {
	return &QuestionLogic{}
}

func (l *QuestionLogic) CreateQuestion(ctx context.Context, req types.CreateQuestionReq) (resp types.CreateQuestionResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 如果指定了题库ID，先判断题库是否存在
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.QuestionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	if len(req.QuestionBankID) > 0 {
		if !repo.NewQuestionRepo(global.DB).CheckQuestionBankExist(questionBankID) {
			zlog.CtxErrorf(ctx, "题库 %d 不存在", req.QuestionBankID)
			return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
		}
	}
	// 将Options和Answer转换为json字符串
	options, _ := json.Marshal(req.Options)
	answers, _ := json.Marshal(req.Answers)
	// 填入数据
	id := snowflake.GetIntId(global.Node)
	question := model.Question{
		ID:         id,
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		Options:    string(options),
		Answer:     string(answers),
		Difficulty: req.Difficulty,
	}
	// 存入数据库
	err = repo.NewQuestionRepo(global.DB).CreateQuestion(question)
	if err != nil {
		zlog.CtxErrorf(ctx, "create question error: %v", err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	// 如果指定了题库ID，则添加到题库中
	if len(req.QuestionBankID) > 0 {
		AddQuestionReq := types.AddQuestionReq{
			QuestionID:     fmt.Sprintf("%d", id),
			QuestionBankID: req.QuestionBankID,
		}
		_, err = l.AddQuestion(ctx, AddQuestionReq)
		if err != nil {
			return resp, err
		}
	}
	// 返回响应
	resp.QuestionID = id
	return
}

func (l *QuestionLogic) CreateQuestionBank(ctx context.Context, req types.CreateQuestionBankReq) (resp types.CreateQuestionBankResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 如果指定了岛屿ID，先判断岛屿是否存在
	if len(req.IslandID) > 0 {
		IslandID, err := strconv.ParseInt(req.IslandID, 10, 64)
		if err != nil {
			zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.IslandID, err)
			return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
		}
		if !repo.NewQuestionRepo(global.DB).CheckQuestionBankExist(IslandID) {
			zlog.CtxErrorf(ctx, "题库 %d 不存在", req.IslandID)
			return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
		}
	}
	// 填入数据
	id := snowflake.GetIntId(global.Node)
	questionBank := model.QuestionBank{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
	}
	// 存入数据库
	err = repo.NewQuestionRepo(global.DB).CreateQuestionBank(questionBank)
	if err != nil {
		zlog.CtxErrorf(ctx, "create question error: %v", err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	// 如果指定了岛屿ID，则添加到岛屿中
	if len(req.IslandID) > 0 {
		AddQuestionReq := types.AddQuestionReq{
			QuestionID:     fmt.Sprintf("%d", id),
			QuestionBankID: req.IslandID,
		}
		_, err = l.AddQuestion(ctx, AddQuestionReq)
		if err != nil {
			return resp, err
		}
	}
	// 返回响应
	resp.QuestionBankID = id
	return
}

func (l *QuestionLogic) AddQuestion(ctx context.Context, req types.AddQuestionReq) (resp types.AddQuestionResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转换为 int64
	questionID, err := strconv.ParseInt(req.QuestionID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 创建表
	questionBankQuestion := model.QuestionBankQuestion{
		QuestionBankID: questionBankID,
		QuestionID:     questionID,
	}
	// 判断问题是否存在
	if !repo.NewQuestionRepo(global.DB).CheckQuestionExist(questionID) {
		zlog.Debugf("题目 %d 不存在", req.QuestionID)
		return resp, response.ErrResp(err, response.QUESTION_NOT_EXIST)
	}
	// 判断题库是否存在
	if !repo.NewQuestionRepo(global.DB).CheckQuestionBankExist(questionBankID) {
		zlog.Debugf("题库 %d 不存在", req.QuestionBankID)
		return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
	}
	// 数据库操作
	err = repo.NewQuestionRepo(global.DB).AddQuestion(questionBankQuestion)
	return resp, nil
}

// AddQuestionBank 向岛屿添加题库
func (l *QuestionLogic) AddQuestionBank(ctx context.Context, req types.AddQuestionBankReq) (resp types.AddQuestionBankResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转换为 int64
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	islandID, err := strconv.ParseInt(req.IslandID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", islandID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 创建表
	islandQuestionBank := model.IslandQuestionBank{
		QuestionBankID: questionBankID,
		IslandID:       islandID,
	}
	// 判断题库是否存在
	if !repo.NewQuestionRepo(global.DB).CheckQuestionBankExist(questionBankID) {
		zlog.Debugf("题库 %d 不存在", req.QuestionBankID)
		return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
	}
	// 判断岛屿是否存在
	if !repo.NewQuestionRepo(global.DB).CheckIslandExist(islandID) {
		zlog.Debugf("岛屿 %d 不存在", req.IslandID)
		return resp, response.ErrResp(err, response.ISLAND_NOT_EXIST)
	}
	// 数据库操作
	err = repo.NewQuestionRepo(global.DB).AddQuestionBank(islandQuestionBank)
	return resp, nil
}

// GetQuestion 获取题目详情
func (l *QuestionLogic) GetQuestion(ctx context.Context, req types.GetQuestionReq) (resp types.GetQuestionResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转换为 int64
	questionID, err := strconv.ParseInt(req.QuestionID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 查询数据库
	question, err := repo.NewQuestionRepo(global.DB).GetQuestion(questionID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "题目 %d 不存在", questionID)
		return resp, response.ErrResp(err, response.QUESTION_NOT_EXIST)
	} else if err != nil {
		zlog.CtxErrorf(ctx, "获取题目 %d 失败: %v", questionID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	// json 反序列化
	options := make([]string, 0)
	err = json.Unmarshal([]byte(question.Options), &options)
	if err != nil {
		zlog.CtxErrorf(ctx, "json 反序列化失败: %v", err)
		return resp, response.ErrResp(err, response.JSON_UNMARSHAL_ERROR)
	}
	answers := make([]string, 0)
	err = json.Unmarshal([]byte(question.Answer), &answers)
	if err != nil {
		zlog.CtxErrorf(ctx, "json 反序列化失败: %v", err)
		return resp, response.ErrResp(err, response.JSON_UNMARSHAL_ERROR)
	}
	// 返回响应
	resp = types.GetQuestionResp{
		Type:       question.Type,
		Title:      question.Title,
		Content:    question.Content,
		Options:    options,
		Answers:    answers,
		Difficulty: question.Difficulty,
	}
	return
}

func (l *QuestionLogic) GetQuestionBank(ctx context.Context, req types.GetQuestionBankReq) (resp types.GetQuestionBankResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转换为 int64
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 查询数据库
	question, err := repo.NewQuestionRepo(global.DB).GetQuestionBank(questionBankID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "题库 %d 不存在", questionBankID)
		return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
	} else if err != nil {
		zlog.CtxErrorf(ctx, "获取题库 %d 失败: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	// 获取题库题目数量
	count, err := repo.NewQuestionRepo(global.DB).GetQuestionBankQuestionCount(questionBankID)
	// 返回响应
	resp = types.GetQuestionBankResp{
		Title:       question.Title,
		Description: question.Description,
		Count:       count,
	}
	return
}

func (l *QuestionLogic) GetQuestionList(ctx context.Context, req types.GetQuestionListReq) (resp types.GetQuestionListResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转换为 int64
	questionBankID, err := strconv.ParseInt(req.QuestionBankID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 计算分页
	offset, limit := (req.Page-1)*req.PageSize, req.PageSize
	// 查询数据库
	list, err := repo.NewQuestionRepo(global.DB).GetQuestionList(questionBankID, offset, limit)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "题库 %d 不存在或没有记录", questionBankID)
		return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
	} else if err != nil {
		zlog.CtxErrorf(ctx, "获取题库 %d 题目列表失败: %v", questionBankID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	fmt.Println(list)
	// 赋值
	resp.Length = 0
	resp.Questions = make([]types.LiteQuestion, 0)
	for _, item := range list {
		question := types.LiteQuestion{
			QuestionID: item.ID,
			Type:       item.Type,
			Title:      item.Title,
			Difficulty: item.Difficulty,
		}
		resp.Questions = append(resp.Questions, question)
		resp.Length++
	}
	return resp, nil
}

func (l *QuestionLogic) GetQuestionBankList(ctx context.Context, req types.GetQuestionBankListReq) (resp types.GetQuestionBankListResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 转换为 int64
	IslandID, err := strconv.ParseInt(req.IslandID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", IslandID, err)
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	// 查询数据库
	list, err := repo.NewQuestionRepo(global.DB).GetQuestionBankList(IslandID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "岛屿 %d 不存在或没有记录", IslandID)
		return resp, response.ErrResp(err, response.QUESTION_BANK_NOT_EXIST)
	} else if err != nil {
		zlog.CtxErrorf(ctx, "获取题库 %d 题目列表失败: %v", IslandID, err)
		return resp, response.ErrResp(err, response.DATABASE_ERROR)
	}
	// 赋值
	resp.Length = 0
	resp.QuestionBanks = make([]types.LiteQuestionBank, 0)
	for _, item := range list {
		var count int
		count, err = repo.NewQuestionRepo(global.DB).GetQuestionBankQuestionCount(item.ID)
		if err != nil {
			zlog.CtxErrorf(ctx, "获取题库 %d 题目数量失败: %v", item.ID, err)
			return resp, response.ErrResp(err, response.DATABASE_ERROR)
		}
		questionBank := types.LiteQuestionBank{
			QuestionBankID: item.ID,
			Title:          item.Title,
			Description:    item.Description,
			Count:          count,
		}
		resp.QuestionBanks = append(resp.QuestionBanks, questionBank)
		resp.Length++
	}
	return resp, nil
}

func (l *QuestionLogic) RunCode(ctx context.Context, req types.RunCodeReq) (resp types.RunCodeResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 默认值处理
	if req.TimeLimit == 0 {
		req.TimeLimit = 1000
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 512 * 1024 * 1024
	}
	// 运行代码
	output, is_limit_out, err := runcodeUtils.RunCode(req.Language, req.Code, req.Stdin, req.TimeLimit, req.MemoryLimit)
	if err != nil {
		zlog.CtxErrorf(ctx, "run code error: %v", err)
		return types.RunCodeResp{}, err
	}

	// 超出限制
	if is_limit_out {
		resp.Output = ""
		resp.StatusCode = 1
		resp.StatusMsg = "运行超时或超出内存限制"
	} else {
		resp.Output = output
		resp.StatusCode = 0
		resp.StatusMsg = "运行成功"
	}

	// 需要判断是否与预期答案相同
	if resp.StatusCode == 0 && len(req.Ans) > 0 {
		if runcodeUtils.CompareOutput(output, req.Ans) {
			resp.StatusCode = 0
			resp.StatusMsg = "答案正确"
		} else {
			resp.StatusCode = 2
			resp.StatusMsg = "答案错误"
		}
	}

	return
}
