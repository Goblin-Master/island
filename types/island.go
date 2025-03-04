package types

type IslandReq struct {
	ID     int64   `json:"id"` // 岛屿id
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	XPoint float64 `json:"xPoint"`
	YPoint float64 `json:"yPoint"`
}
type IslandResp struct {
	UserID int64  `json:"userid,string"`
	ID     int64  `json:"id,string"`
	Name   string `json:"name"`
}
type ListReq struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	XPoint float64 `json:"xPoint"`
	YPoint float64 `json:"yPoint"`
	Path   string  `json:"path"`
	UserID int64   `json:"userid"`
}
type ListResp struct {
	ID     int64   `json:"id,string"`
	Name   string  `json:"name"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	XPoint float64 `json:"xPoint"`
	YPoint float64 `json:"yPoint"`
	Path   string  `json:"path"`
	UserID int64   `json:"userid,string"`
}
type IslandListResp struct {
	List  []ListResp `json:"list"`
	Count int        `json:"count"`
}
type IslandDeleteReq struct {
	ID int64 `form:"id"`
}
