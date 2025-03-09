package types

type CreateQuestionReq struct {
	Type      int      `json:"type"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Options   []string `json:"options"`
	Answers   []string `json:"answers"`
	Difficult int      `json:"difficult"`
	// 可选
	QuestionBankID string `json:"question_bank_id"`
}

type CreateQuestionResp struct {
	QuestionID int64 `json:"question_id,string"`
}

type CreateQuestionBankReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	// 可选
	IslandID string `json:"island_id"`
}

type CreateQuestionBankResp struct {
	QuestionBankID int64 `json:"question_bank_id,string"`
}

type AddQuestionReq struct {
	QuestionBankID string `json:"question_bank_id"`
	QuestionID     string `json:"question_id"`
}

type AddQuestionResp struct {
}

type AddQuestionBankReq struct {
	QuestionBankID string `json:"question_bank_id"`
	IslandID       string `json:"island_id"`
}

type AddQuestionBankResp struct {
}

type GetQuestionReq struct {
	QuestionID string `form:"question_id"`
}

type GetQuestionResp struct {
	Type      int      `json:"type"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Options   []string `json:"options"`
	Answers   []string `json:"answers"`
	Difficult int      `json:"difficult"`
}

type LiteQuestionBank struct {
	QuestionBankID int64  `json:"question_bank_id,string"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Count          int    `json:"count"`
}

type GetQuestionBankReq struct {
	QuestionBankID string `form:"question_bank_id"`
}

type GetQuestionBankResp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

type LiteQuestion struct {
	QuestionID int64  `json:"question_id,string"`
	Title      string `json:"title"`
	Type       int    `json:"type"`
}

type GetQuestionListReq struct {
	QuestionBankID string `form:"question_bank_id"`
	Page           int    `form:"page"`
	PageSize       int    `form:"page_size"`
}

type GetQuestionListResp struct {
	Length    int            `json:"length"`
	Questions []LiteQuestion `json:"questions"`
}

type GetQuestionBankListReq struct {
	IslandID string `form:"island_id"`
}

type GetQuestionBankListResp struct {
	Length        int                `json:"length"`
	QuestionBanks []LiteQuestionBank `json:"question_banks"`
}

type RunCodeReq struct {
	Language    string `json:"language"`
	Code        string `json:"code"`
	Stdin       string `json:"stdin"`
	TimeLimit   int    `json:"time_limit"`
	MemoryLimit int    `json:"memory_limit"`
	Ans         string `json:"ans"`
}

type RunCodeResp struct {
	Output     string `json:"output"`
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
