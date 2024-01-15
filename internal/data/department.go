package data

import (
	"context"
	"fmt"
	"strings"
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
func (repo *departmentRepo) Find(ctx context.Context, listQuery *biz.DeptListQuery) ([]*biz.Dept, error) {

	// 筛选条件
	query := []string{"deleted_at is null"}
	queryArgs := make([]interface{}, 0)

	if strings.TrimSpace(listQuery.DeptName) != "" {
		query = append(query, "dept_name like ?")
		queryArgs = append(queryArgs, "%"+listQuery.DeptName+"%")
	}
	// 状态
	if listQuery.Status != nil {
		query = append(query, "status = ?")
		queryArgs = append(queryArgs, *listQuery.Status)
	}
	queryStr := strings.Join(query, " and ")
	Dept := []*biz.Dept{}
	tx := repo.data.db.Debug().Where(queryStr, queryArgs...).Order("list_order").Find(&Dept)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return Dept, nil
}

// 查看单个部门
func (repo *departmentRepo) FindOne(ctx context.Context, id int64) (*biz.Dept, error) {
	dept := &biz.Dept{
		DeptID: id,
	}
	tx := repo.data.db.First(dept)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dept, nil
}

// 根据父亲id统计子部门数量
func (repo *departmentRepo) CountByParentId(ctx context.Context, parentId int64) (int64, error) {
	var count int64
	tx := repo.data.db.Model(&biz.Dept{}).Where("parent_id = ?", parentId).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

// 添加部门
func (repo *departmentRepo) Insert(ctx context.Context, Dept *biz.Dept) error {
	tx := repo.data.db.Create(&Dept)
	return tx.Error
}

// 更新部门
func (repo *departmentRepo) Update(ctx context.Context, Dept *biz.Dept) error {
	key := fmt.Sprintf("%s%v", deptCachePrefix, Dept.DeptID)
	repo.data.rdb.Del(ctx, key)
	tx := repo.data.db.Where("dept_id", Dept.DeptID).Save(&Dept)
	return tx.Error
}

// 删除部门
func (repo *departmentRepo) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%v", deptCachePrefix, id)
	repo.data.rdb.Del(ctx, key)
	return repo.data.db.Model(&biz.Dept{}).Where("dept_id = ?", id).Update("deleted_at", time.Now()).Error
}

func NewDeparmentRepo(data *Data) biz.DepartmentRepo {
	return &departmentRepo{
		data: data,
	}
}
