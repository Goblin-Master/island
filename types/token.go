package types

type TokenReq struct {
	Token string `json:"token"`
}
type TokenResp struct {
	Atoken string `json:"atoken"`
	Rtoken string `json:"rtoken"`
}
