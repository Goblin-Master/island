package types

import (
	"tgwp/repo/list"
	"time"
)

type UserDetailReq struct {
	UserID string `form:"user_id"`
}
type UserDetailResp struct {
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
type FocusUserReq struct {
	UserID  int64  `form:"user_id"`
	FocusID string `form:"focus_id"`
}
type FocusListReq struct {
	list.PageInfo
	UserID string `form:"user_id"`
}
type FocusList struct {
	UserID    int64     `json:"user_id,string"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
type FocusListResp struct {
	Count int         `json:"count"`
	List  []FocusList `json:"list"`
}
type CountFansReq struct {
	UserID string `form:"user_id"`
}
