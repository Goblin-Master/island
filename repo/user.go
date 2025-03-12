package repo

import (
	"errors"
	"gorm.io/gorm"
	"tgwp/log/zlog"
	"tgwp/model"
	"tgwp/response"
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
func (r *UserRepo) IsExist(focusID int64) (err error) {
	return r.DB.Model(&model.User{}).Where("id = ?", focusID).Take(&model.User{}).Error
}
func (r *UserRepo) FocusUser(userID int64, focusID int64) (resp string, err error) {
	var focus model.Focus
	err = r.DB.Model(&model.Focus{}).Where("user_id = ? and focus_id = ?", userID, focusID).Take(&focus).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.DB.Model(&model.Focus{}).Create(&model.Focus{
			UserID:  userID,
			FocusID: focusID,
		}).Error
		if err != nil {
			return "", response.ErrResp(err, response.USER_FOCUS_ERROR)
		}
		return "关注成功", nil
	} else if err != nil {
		zlog.Errorf("关注失败: %v", err)
		return "", response.ErrResp(err, response.COMMON_FAIL)
	} else {
		err = r.DB.Delete(&focus).Error
		if err != nil {
			return "", response.ErrResp(err, response.USER_CANCEL_FOCUS)
		}
		return "取消关注成功", nil
	}
}
