package types

type IslandReq struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	XPoint float64 `json:"xPoint"`
	YPoint float64 `json:"yPoint"`
}
type IslandResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type IslandListResp struct {
	List  []IslandReq `json:"list"`
	Count int         `json:"count"`
}
