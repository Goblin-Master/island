package logic

import (
	"context"
	"tgwp/log/zlog"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/runcodeUtils"
	"time"
)

type QuestionLogic struct {
}

func NewQuestionLogic() *QuestionLogic {
	return &QuestionLogic{}
}

func (l *QuestionLogic) RunCode(ctx context.Context, req types.RunCodeReq) (resp types.RunCodeResp, err error) {
	defer utils.RecordTime(time.Now())()
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
