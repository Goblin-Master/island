package logic

import (
	"context"
	"strconv"
	"tgwp/log/zlog"
	"tgwp/pkg/ai_eino"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils"
	"tgwp/utils/aiUtils"
	"time"
)

type AILogic struct {
}

func NewAILogic() *AILogic {
	return &AILogic{}
}
func (l *AILogic) GenerateAbstract(ctx context.Context, req types.AiReq) (resp types.AiAnalysisResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp.Abstract, err = ai_eino.Chat(ctx, req)
	return
}

func (l *AILogic) ClearHistory(ctx context.Context, user_id int64) (err error) {
	err = aiUtils.ClearHistory(ctx, strconv.FormatInt(user_id, 10))
	if err != nil {
		zlog.CtxErrorf(ctx, "清除历史记录失败 %s", err)
		return response.ErrResp(err, response.COMMON_FAIL)
	}
	return
}
