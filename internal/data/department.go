package data

import (
	"context"
	"fmt"
	"time"
	"zerocmf/internal/biz"
)

type departmentRepo struct {
	data *Data
}

var (
	deptCachePrefix = "cache:dept:id:"
)

// 查询全部部门列表
func (repo *departmentRepo) Find(ctx context.Context, query *biz.SysDeptListQuery) ([]*biz.SysDept, error) {
	sysDept := []*biz.SysDept{}
	tx := repo.data.db.Find(&sysDept)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sysDept, nil
}

// 查看单个部门
func (repo *departmentRepo) FindOne(ctx context.Context, id int64) (*biz.SysDept, error) {
	dept := &biz.SysDept{
		DeptID: id,
	}
	tx := repo.data.db.First(dept)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dept, nil
}

// 添加部门
func (repo *departmentRepo) Insert(ctx context.Context, sysDept *biz.SysDept) error {
	tx := repo.data.db.Create(&sysDept)
	return tx.Error
}

// 更新部门
func (repo *departmentRepo) Update(ctx context.Context, sysDept *biz.SysDept) error {
	key := fmt.Sprintf("%s%v", deptCachePrefix, sysDept.DeptID)
	repo.data.rdb.Del(ctx, key)
	tx := repo.data.db.Where("dept_id", sysDept.DeptID).Save(&sysDept)
	return tx.Error
}

// 删除部门
func (repo *departmentRepo) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%v", deptCachePrefix, id)
	repo.data.rdb.Del(ctx, key)
	return repo.data.db.Model(&biz.SysDept{}).Where("post_id = ?", id).Update("deleted_at", time.Now()).Error
}

func NewDeparmentRepo(data *Data) biz.DepartmentRepo {
	return &departmentRepo{
		data: data,
	}
}
