package types

type AiReq struct {
	Content string `form:"content" json:"content"`
	Type    int    `form:"type" json:"type"`
}
type AiAnalysisResp struct {
	Abstract string `json:"abstract"`
}
