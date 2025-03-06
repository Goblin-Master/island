package types

type AiReq struct {
	UserID  int64  `form:"-"`
	Content string `form:"content" json:"content"`
	Type    int    `form:"type" json:"type"`
}
type AiAnalysisResp struct {
	Abstract string `json:"abstract"`
}
