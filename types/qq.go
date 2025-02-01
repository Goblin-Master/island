package types

type QQLoginReq struct {
	Code string `json:"code"`
}
type QQLoginResp struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Atoken   string `json:"atoken"`
	Rtoken   string `json:"rtoken"`
}
