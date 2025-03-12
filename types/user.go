package types

type UserDetailReq struct {
	UserID string `form:"user_id"`
}
type UserDetailResp struct {
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}
type FocusUserReq struct {
	UserID  int64  `form:"user_id"`
	FocusID string `form:"focus_id"`
}
