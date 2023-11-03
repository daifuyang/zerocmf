package data

import (
	"context"
	"fmt"
	"time"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/v8"
)

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) biz.RoleRepo {
	return &roleRepo{data: data}
}

var (
	roleCachePrefix = "cache:role:id:"
)

// 查询全部
func (repo *roleRepo) Find(ctx context.Context) (roles []*biz.SysRole, err error) {
	tx := repo.data.db.Where("deleted_at is null").Find(&roles)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	return
}

// 根据角色id获取单个角色
func (repo *roleRepo) FindOne(ctx context.Context, id int64) (*biz.SysRole, error) {
	var sysRole *biz.SysRole
	key := fmt.Sprintf("%s%v", roleCachePrefix, id)

	err := repo.data.RGet(ctx, &sysRole, key)
	if err != redis.Nil {
		return sysRole, nil
	}

	tx := repo.data.db.Where("role_id = ? AND deleted_at is null", id).First(&sysRole)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, sysRole)
	if err != nil {
		return nil, err
	}

	return sysRole, nil
}

// 插入一条数据
func (repo *roleRepo) Insert(ctx context.Context, role *biz.SysRole) (err error) {
	tx := repo.data.db.Create(&role)
	return tx.Error
}

// 更新一条数据
func (repo *roleRepo) Update(ctx context.Context, role *biz.SysRole) (err error) {
	// redis重置
	repo.data.rdb.Del(ctx, fmt.Sprintf("%s%v", roleCachePrefix, role.RoleID))
	return repo.data.db.Save(&role).Error
}

// 删除一条数据
func (repo *roleRepo) Delete(ctx context.Context, id int64) error {
	repo.data.rdb.Del(ctx, fmt.Sprintf("%s%v", roleCachePrefix, id))
	return repo.data.db.Model(&biz.SysRole{}).Where("role_id = ?", id).Update("deleted_at", time.Now()).Error
}
