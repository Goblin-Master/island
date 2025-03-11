package repo

import (
	"errors"
	"gorm.io/gorm"
	"tgwp/log/zlog"
	"tgwp/model"
)

type QuestionRepo struct {
	DB *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) *QuestionRepo {
	return &QuestionRepo{
		DB: db,
	}
}

// CreateQuestion 创建问题
func (r *QuestionRepo) CreateQuestion(question model.Question) error {
	return r.DB.Create(&question).Error
}

// CreateQuestionBank 创建题库
func (r *QuestionRepo) CreateQuestionBank(questionBank model.QuestionBank) error {
	return r.DB.Create(&questionBank).Error
}

// AddQuestion 添加题目到题库
func (r *QuestionRepo) AddQuestion(model model.QuestionBankQuestion) error {
	return r.DB.Create(&model).Error
}

// AddQuestionBank 添加题库到岛屿
func (r *QuestionRepo) AddQuestionBank(model model.IslandQuestionBank) error {
	return r.DB.Create(&model).Error
}

// GetQuestion 获取指定题目
func (r *QuestionRepo) GetQuestion(id int64) (model.Question, error) {
	var question model.Question
	err := r.DB.First(&question, id).Error
	return question, err
}

// GetQuestionBank 获取指定题库
func (r *QuestionRepo) GetQuestionBank(id int64) (model.QuestionBank, error) {
	var questionBank model.QuestionBank
	err := r.DB.First(&questionBank, id).Error
	return questionBank, err
}

// GetQuestionBankQuestionCount 获取题库题目数量
func (r *QuestionRepo) GetQuestionBankQuestionCount(id int64) (int, error) {
	var count int64
	err := r.DB.Model(&model.QuestionBankQuestion{}).Where("question_bank_id =?", id).Count(&count).Error
	zlog.Debugf("题库 %d 题目数量: %d", id, count)
	return int(count), err
}

// GetQuestionList 获取题库题目列表
func (r *QuestionRepo) GetQuestionList(questionBankID int64, offset, limit int) ([]model.Question, error) {
	var questions []model.Question

	err := r.DB.
		Select("questions.*").
		Joins("INNER JOIN question_bank_questions ON questions.id = question_bank_questions.question_id").
		Where("question_bank_questions.question_bank_id = ?", questionBankID).
		Order("question_bank_questions.created_at").
		Offset(offset).
		Limit(limit).
		Find(&questions).Error

	return questions, err
}

func (r *QuestionRepo) GetRandomQuestions(questionBankID int64, count int) ([]model.Question, error) {
	var questions []model.Question
	err := r.DB.
		Select("questions.*").
		Joins("INNER JOIN question_bank_questions ON questions.id = question_bank_questions.question_id").
		Where("question_bank_questions.question_bank_id = ?", questionBankID).
		Where("questions.type <= 3").
		Limit(count).
		Order("RAND()").
		Find(&questions).Error

	return questions, err
}

// GetQuestionBankList 获取岛屿题库列表
func (r *QuestionRepo) GetQuestionBankList(islandID int64) ([]model.QuestionBank, error) {
	var questionBanks []model.QuestionBank

	err := r.DB.
		Select("question_banks.*").
		Joins("INNER JOIN island_question_banks ON question_banks.id = island_question_banks.question_bank_id").
		Where("island_question_banks.island_id = ?", islandID).
		Order("island_question_banks.created_at DESC").
		Find(&questionBanks).Error

	return questionBanks, err
}

// CheckQuestionExist 检查题目是否存在
func (r *QuestionRepo) CheckQuestionExist(id int64) bool {
	return !errors.Is(r.DB.First(&model.Question{}, id).Error, gorm.ErrRecordNotFound)
}

// CheckQuestionBankExist 检查题库是否存在
func (r *QuestionRepo) CheckQuestionBankExist(id int64) bool {
	return !errors.Is(r.DB.First(&model.QuestionBank{}, id).Error, gorm.ErrRecordNotFound)
}

// CheckIslandExist 检查岛屿是否存在
func (r *QuestionRepo) CheckIslandExist(id int64) bool {
	return !errors.Is(r.DB.First(&model.Island{}, id).Error, gorm.ErrRecordNotFound)
}
