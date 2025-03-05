package repo

import (
	"gorm.io/gorm"
	"tgwp/model"
)

type IslandRepo struct {
	DB *gorm.DB
}

func NewIslandRepo(db *gorm.DB) *IslandRepo {
	return &IslandRepo{
		DB: db,
	}
}
func (r *IslandRepo) CreateIsland(island model.Island) (err error) {
	err = r.DB.Create(&island).Error
	return
}
func (r *IslandRepo) UpdatesIsland(island model.Island) (err error) {
	err = r.DB.Where("id = ?", island.ID).Updates(&island).Error
	return
}

func (r *IslandRepo) IdentifyIslandName(name string) (exist bool) {
	var island model.Island
	err := r.DB.Where("name = ?", name).Take(&island).Error
	if err == nil {
		exist = true
	}
	return
}
func (r *IslandRepo) IdentifyIslandNameAndId(name string, id int64) (exist bool) {
	var island model.Island
	err := r.DB.Where("name = ? and id <> ?", name, id).Take(&island).Error
	if err == nil {
		exist = true
	}
	return
}
func (r *IslandRepo) IdentifyIslandById(id, user_id int64) (exist bool) {
	var island model.Island
	err := r.DB.Where("id = ? and user_id = ?", id, user_id).Take(&island).Error
	if err == nil {
		exist = true
	}
	return
}
func (r *IslandRepo) DeleteIsland(id, user_id int64) (err error) {
	var island model.Island
	err = r.DB.Where("id = ? and user_id = ?", id, user_id).Take(&island).Error
	if err != nil {
		return
	}
	err = r.DB.Delete(&island).Error
	return
}
