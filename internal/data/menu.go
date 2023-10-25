package data

import (
	"context"
	"zerocmf/internal/biz"
)

type MenuRepo struct {
	data *Data
}

// FindOneByMenuName implements biz.MenuRepo.
func (r *MenuRepo) FindOneByMenuName(ctx context.Context, menuName string) (*biz.SysMenu, error) {
	var sysMenu *biz.SysMenu
	err := r.data.db.Where("menu_name = ?", menuName).First(&sysMenu).Error
	if err != nil {
		return nil, err
	}
	return sysMenu, nil
}

// Update implements biz.MenuRepo.
func (r *MenuRepo) Update(ctx context.Context, menu *biz.SysMenu) error {
	return r.data.db.Save(&menu).Error
}

func (r *MenuRepo) Insert(ctx context.Context, menu *biz.SysMenu) error {
	tx := r.data.db.Create(&menu)
	return tx.Error
}

func NewMenuRepo(data *Data) biz.MenuRepo {
	return &MenuRepo{
		data: data,
	}
}
