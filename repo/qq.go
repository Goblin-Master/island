package repo

import (
	"gorm.io/gorm"
	"tgwp/model"
)

type QQLoginRepo struct {
	DB *gorm.DB
}

func NewQQLoginRepo(db *gorm.DB) *QQLoginRepo {
	return &QQLoginRepo{
		DB: db,
	}
}

func (r *QQLoginRepo) IsExist(openid string) (int64, error) {
	var user model.User
	err := r.DB.Where("open_id = ?", openid).Take(&user).Error
	return user.ID, err
}
func (r *QQLoginRepo) CreateUser(user model.User) (err error) {
	err = r.DB.Create(&user).Error
	return
}
