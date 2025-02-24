package logic

import (
	"context"
	"tgwp/pkg/ai"
	"tgwp/types"
	"tgwp/utils"
	"time"
)

type AILogic struct {
}

func NewAILogic() *AILogic {
	return &AILogic{}
}
func (l *AILogic) GenerateAbstract(ctx context.Context, res types.AiReq) (resp types.AiAnalysisResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp.Abstract, err = ai.Chat(ctx, res.Content)
	return
}
