package types

type ArticleCreateReq struct {
	UserID   int64  `json:"userid"`
	Abstract string `json:"abstract"`
	Content  string `json:"content" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Cover    string `json:"cover"`
}
type ArticleCreateResp struct {
	ID int64 `json:"id"` // 文章id
}
