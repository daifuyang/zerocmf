package data

import (
	"context"
	"zerocmf/internal/biz"
)

type DepartmentRepo struct {
	data *Data
}

// GetOneById implements biz.DepartmentRepo.
func (repo *DepartmentRepo) GetOneById(ctx context.Context, id int64) (*biz.SysDept, error) {
	dept := &biz.SysDept{
		DeptID: id,
	}
	tx := repo.data.db.First(dept)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dept, nil
}

// Add implements biz.DepartmentRepo.
func (repo *DepartmentRepo) Add(ctx context.Context, sysDept *biz.SysDept) error {
	tx := repo.data.db.Create(&sysDept)
	return tx.Error
}

// update implements biz.DepartmentRepo.
func (repo *DepartmentRepo) Update(ctx context.Context, sysDept *biz.SysDept) error {
	tx := repo.data.db.Where("dept_id", sysDept.DeptID).Save(&sysDept)
	return tx.Error
}

// 查询全部部门列表
func (repo *DepartmentRepo) Index(ctx context.Context) ([]*biz.SysDept, error) {
	sysDept := []*biz.SysDept{}
	tx := repo.data.db.Find(&sysDept)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sysDept, nil
}

func NewDeparmentRepo(data *Data) biz.DepartmentRepo {
	return &DepartmentRepo{
		data: data,
	}
}
