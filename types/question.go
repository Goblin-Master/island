package types

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
