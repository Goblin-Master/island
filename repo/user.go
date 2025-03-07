package repo

import (
	"gorm.io/gorm"
	"tgwp/model"
	"tgwp/types"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}
func (r *UserRepo) UserDetail(id int64) (resp types.UserDetailResp, err error) {
	err = r.DB.Model(&model.User{}).Where("id = ?", id).Take(&resp).Error
	return
}
