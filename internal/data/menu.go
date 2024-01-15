package data

import (
	"context"
	"fmt"
	"time"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/v8"
)

type menuRepo struct {
	data *Data
}

var (
	menuCachePrefix = "cache:menu:id:"
)

// 查询全部菜单
func (repo *menuRepo) Find(ctx context.Context) (menus []*biz.Menu, err error) {
	tx := repo.data.db.Where("deleted_at is null").Find(&menus)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	return
}

// 查询单个菜单
func (repo *menuRepo) FindOne(ctx context.Context, id int64) (*biz.Menu, error) {

	var Menu *biz.Menu
	key := fmt.Sprintf("%s%v", userCachePrefix, id)

	err := repo.data.RGet(ctx, &Menu, key)
	if err != redis.Nil {
		return Menu, nil
	}

	tx := repo.data.db.Where("menu_id = ? AND deleted_at is null", id).First(&Menu)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, Menu)
	if err != nil {
		return nil, err
	}

	return Menu, nil
}

// 根据名称查询单个菜单
func (r *menuRepo) FindOneByMenuName(ctx context.Context, menuName string) (*biz.Menu, error) {
	var Menu *biz.Menu
	err := r.data.db.Where("menu_name = ? AND deleted_at is null", menuName).First(&Menu).Error
	if err != nil {
		return nil, err
	}
	return Menu, nil
}

// 跟新单个菜单
func (r *menuRepo) Update(ctx context.Context, menu *biz.Menu) error {
	// redis重置
	r.data.rdb.Del(ctx, fmt.Sprintf("%s%v", menuCachePrefix, menu.MenuID))
	return r.data.db.Save(&menu).Error
}

// 插入菜单
func (r *menuRepo) Insert(ctx context.Context, menu *biz.Menu) error {
	tx := r.data.db.Create(&menu)
	return tx.Error
}

// 删除单个菜单
func (repo *menuRepo) Delete(ctx context.Context, id int64) error {
	repo.data.rdb.Del(ctx, fmt.Sprintf("%s%v", menuCachePrefix, id))
	return repo.data.db.Model(&biz.Menu{}).Where("menu_id = ?", id).Update("deleted_at", time.Now()).Error
}

func NewMenuRepo(data *Data) biz.MenuRepo {
	return &menuRepo{
		data: data,
	}
}
