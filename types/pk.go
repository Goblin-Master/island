package types

type QuestionsInfo struct {
	QuestionID  int64 `json:"question_id,string"`
	User1Submit bool  `json:"user1_submit"`
	User2Submit bool  `json:"user2_submit"`
	User1Score  int64 `json:"user1_score"`
	User2Score  int64 `json:"user2_score"`
	YourSubmit  bool  `json:"your_submit"`
	YourScore   int64 `json:"your_score"`
}

type RedisRoomInfo struct {
	RoomID                    int64           `json:"room_id"`
	User1ID                   int64           `json:"user1_id"`
	User2ID                   int64           `json:"user2_id"`
	Questions                 []QuestionsInfo `json:"questions"`
	User1ScoreTotal           int64           `json:"user1_score_total"`
	User2ScoreTotal           int64           `json:"user2_score_total"`
	User1FinalSubmitTimestamp int64           `json:"user1_final_submit_timestamp"` // 最后一次提交时间戳 (用于同分比较输赢)
	User2FinalSubmitTimestamp int64           `json:"user2_final_submit_timestamp"` // 最后一次提交时间戳 (用于同分比较输赢)
	StartTimestamp            int64           `json:"start_timestamp"`
	EndTimestamp              int64           `json:"end_timestamp"`
	WinnerID                  int64           `json:"winner_id"`
}

type PKMatchingReq struct {
	UserID         int64  `form:"-"`
	QuestionBankID string `json:"question_bank_id"`
}

type PKMatchingResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	RoomID int64  `json:"room_id,string"`
}

type GetRoomInfoReq struct {
	UserID int64  `form:"-"`
	RoomID string `form:"room_id"`
}

type GetRoomInfoResp struct {
	User1ID                   int64           `json:"user1_id,string"`
	User2ID                   int64           `json:"user2_id,string"`
	Questions                 []QuestionsInfo `json:"questions"`
	User1ScoreTotal           int64           `json:"user1_score_total"`
	User2ScoreTotal           int64           `json:"user2_score_total"`
	User1FinalSubmitTimestamp int64           `json:"user1_final_submit_timestamp"` // 最后一次提交时间戳 (用于同分比较输赢)
	User2FinalSubmitTimestamp int64           `json:"user2_final_submit_timestamp"` // 最后一次提交时间戳 (用于同分比较输赢)
	StartTimestamp            int64           `json:"start_timestamp"`
	EndTimestamp              int64           `json:"end_timestamp"`
	WinnerID                  int64           `json:"winner_id,string"`
}

type SubmitQuestionReq struct {
	UserID     int64    `json:"-"`
	RoomID     string   `json:"room_id"`
	QuestionID string   `json:"question_id"`
	Answer     []string `json:"answer"`
}

type SubmitQuestionResp struct {
	Score int64 `json:"score"`
}
