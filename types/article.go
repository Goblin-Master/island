package types

import "time"

type ArticleCreateReq struct {
	Island   string `json:"island" binding:"required"`
	UserID   int64  `json:"userid"`
	Abstract string `json:"abstract"`
	Content  string `json:"content" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Cover    string `json:"cover"`
}
type ArticleCreateResp struct {
	ID int64 `json:"id,string"` // 文章id
}

type Article struct {
	Island       string    `json:"island"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       int64     `json:"userid,string"`
	Abstract     string    `json:"abstract"`
	Content      string    `json:"content" `
	Title        string    `json:"title" `
	Cover        string    `json:"cover"`
	Username     string    `json:"username"`
	Avatar       string    `json:"avatar"`
	DiggCount    int       `json:"digg_count"`
	CollectCount int       `json:"collect_count"`
	ID           int64     `json:"id,string"`
}

type ArticleList struct {
	List  []Article `json:"list"`
	Count int       `json:"count"`
}

type ArticleDiggReq struct {
	UserID    int64 `form:"userid"`
	ArticleID int64 `form:"article_id"`
}
type ArticleCollectReq struct {
	UserID    int64 `form:"userid"`
	ArticleID int64 `form:"article_id"`
}
