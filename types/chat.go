package types

type Message struct {
	ID      int64  `json:"id,string"`
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type SendMessageReq struct {
	UserID   int64  `json:"user_id"`
	IslandID int64  `json:"island_id"`
	Message  string `json:"message"`
}

type SendMessageResp struct {
	ID        int64 `json:"id,string"`
	Timestamp int64 `json:"timestamp"`
}
