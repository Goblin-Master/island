package types

type ImageResp struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
	WebPath  string `json:"webPath"`
}
type ImageListResp struct {
	List  []ImageResp
	Count int `json:"count"`
}
