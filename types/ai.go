package types

type AiReq struct {
	Content string `form:"content" json:"content"`
}
type AiAnalysisResp struct {
	Abstract string `json:"abstract"`
}
