package types

type Message struct {
	ID        int64  `json:"id,string"`
	UserID    int64  `json:"user_id,string"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type SendMessageReq struct {
	UserID   int64  `json:"-"`
	IslandID string `json:"island_id"`
	Message  string `json:"message"`
}

type SendMessageResp struct {
	ID        int64 `json:"id,string"`
	Timestamp int64 `json:"timestamp"`
}

type GetMessagesReq struct {
	IslandID  string `form:"island_id"`
	Timestamp int64  `form:"timestamp"`
}

type GetMessagesResp struct {
	Timestamp int64     `json:"timestamp"`
	Length    int       `json:"length"`
	Messages  []Message `json:"messages"`
}
