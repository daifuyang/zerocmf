package data

import (
	"context"
	"fmt"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/v8"
)

type MenuRepo struct {
	data *Data
}

var (
	menuCachePrefix = "cache:user:userId:"
)

// Delete implements biz.MenuRepo.
func (*MenuRepo) Delete(ctx context.Context, id int64) (*biz.SysMenu, error) {
	panic("unimplemented")
}

// Find implements biz.MenuRepo.
func (repo *MenuRepo) Find(ctx context.Context) (menus []*biz.SysMenu, err error) {
	tx := repo.data.db.Find(&menus)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	return
}

// FindOne implements biz.MenuRepo.
func (repo *MenuRepo) FindOne(ctx context.Context, id int64) (*biz.SysMenu, error) {

	var sysMenu *biz.SysMenu
	key := fmt.Sprintf("%s%v", userCachePrefix, id)

	err := repo.data.RGet(ctx, &sysMenu, key)
	if err != redis.Nil {
		return sysMenu, nil
	}

	tx := repo.data.db.Where("menu_id = ? AND status = 1", id).First(&sysMenu)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, sysMenu)
	if err != nil {
		return nil, err
	}

	return sysMenu, nil
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
