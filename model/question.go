package model

import "time"

type Question struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// Type 1.选择题 2.多选题 3.填空题 4.主观题
	Type       int    `json:"type" gorm:"type:int;not null;default:1;comment:'题目类型'"`
	ID         int64  `json:"id" gorm:"type:bigint;not null;primary_key;comment:'题目ID'"`
	Title      string `json:"title" gorm:"type:text;not null;comment:'题目标题'"`
	Content    string `json:"content" gorm:"type:text;not null;comment:'题目内容'"`
	Options    string `json:"options" gorm:"type:text;not null;comment:'选项列表'"`
	Answer     string `json:"answer" gorm:"type:text;not null;comment:'答案'"`
	Difficulty int    `json:"difficulty" gorm:"type:int;not null;default:1;comment:'题目难度'"`
}

func (i *Question) TableName() string {
	return "questions"
}

type QuestionBank struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ID          int64  `json:"id" gorm:"type:bigint;not null;primary_key;uniqueIndex;comment:'题库ID'"`
	Title       string `json:"title" gorm:"type:text;not null;comment:'题库名称'"`
	Description string `json:"description" gorm:"type:text;not null;comment:'题库描述'"`
}

func (i *QuestionBank) TableName() string {
	return "question_banks"
}

type QuestionBankQuestion struct {
	CommonModel
	CreatedAt      time.Time `gorm:"index:idx_created_at"`
	QuestionBankID int64     `json:"question_bank_id" gorm:"type:bigint;not null;uniqueIndex:idx_bank_question;index:idx_question_bank_id;comment:'题库ID'"`
	QuestionID     int64     `json:"question_id" gorm:"type:bigint;not null;uniqueIndex:idx_bank_question;comment:'题目ID'"`

	QuestionBank QuestionBank `gorm:"foreignKey:QuestionBankID;references:ID"`
	Question     Question     `gorm:"foreignKey:QuestionID;references:ID"`
}

func (i *QuestionBankQuestion) TableName() string {
	return "question_bank_questions"
}

type IslandQuestionBank struct {
	CommonModel
	CreatedAt      time.Time `gorm:"index:idx_created_at"`
	IslandID       int64     `json:"island_id" gorm:"type:bigint;not null;uniqueIndex:idx_island_bank;comment:'岛屿ID'"`
	QuestionBankID int64     `json:"question_bank_id" gorm:"type:bigint;not null;uniqueIndex:idx_island_bank;comment:'题库ID'"`

	Island       Island       `gorm:"foreignKey:IslandID;references:ID"`
	QuestionBank QuestionBank `gorm:"foreignKey:QuestionBankID;references:ID"`
}

func (i *IslandQuestionBank) TableName() string {
	return "island_question_banks"
}
