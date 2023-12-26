package biz

import (
	"context"
	"strconv"

	"gorm.io/gorm"
)

type SysDept struct {
	DeptID    int64      `gorm:"column:deptId;type:int(11);primaryKey;comment:部门id" uri:"id" json:"deptId"`
	ParentID  int64      `gorm:"column:parentId;type:bigint;default:0;comment:父部门id" json:"parentId"`
	Ancestors string     `gorm:"column:ancestors;type:varchar(50);default:'';comment:祖级列表" json:"ancestors"`
	DeptName  string     `gorm:"column:deptName;type:varchar(30);default:'';comment:部门名称" json:"deptName"`
	ListOrder int        `gorm:"column:list_order;type:int(11);default:0;comment:显示顺序" json:"listOrder"`
	Leader    string     `gorm:"column:leader;type:varchar(20);default:null;comment:负责人" json:"leader"`
	Phone     string     `gorm:"column:phone;type:varchar(11);default:null;comment:联系电话" json:"phone"`
	Email     string     `gorm:"column:email;type:varchar(50);default:null;comment:'邮箱'" json:"email"`
	Status    int64      `gorm:"column:status;type:tinyint(2);default:1;comment:菜单状态（0：停用 1：显示）" json:"status"`
	DeletedAt *LocalTime `gorm:"column:deleted_at;default:null;comment:删除时间" json:"deletedAt"`
	SysInfo
}

func (SysDept) TableName() string {
	return "sys_dept"
}

// 列表筛选条件
type SysDeptListQuery struct {
	DeptName string `form:"deptName"`
	Status   *int   `form:"status"`
}

// 数据库迁移
func (biz *SysDept) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}

	tx := db.First(&biz)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		tx = db.Create(&SysDept{
			DeptName:  "zerocmf",
			ListOrder: 1,
			Status:    1,
		})
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

type DepartmentRepo interface {
	Find(ctx context.Context, listQuery *SysDeptListQuery) ([]*SysDept, error) // 查看全部
	FindOne(ctx context.Context, id int64) (*SysDept, error)                   // 查看单条
	Insert(ctx context.Context, sysDept *SysDept) error                        // 添加部门
	Update(ctx context.Context, sysDept *SysDept) error                        // 更新部门
	Delete(ctx context.Context, id int64) error                                // 删除部门
}

type Depatmentusecase struct {
	repo DepartmentRepo
}

func NewDepatmentusecase(repo DepartmentRepo) *Depatmentusecase {
	return &Depatmentusecase{repo: repo}
}

// 查看部门列表
func (biz *Depatmentusecase) Tree(ctx context.Context, listQuery *SysDeptListQuery) ([]*SysDept, error) {
	return biz.repo.Find(ctx, listQuery)
}

// 添加部门
func (biz *Depatmentusecase) Insert(ctx context.Context, sysDept *SysDept) error {

	// 查询parentId是否合法
	parentId := sysDept.ParentID

	// 并入参数ancestors
	ancestors := "0"
	if parentId > 0 {
		parent, err := biz.repo.FindOne(ctx, parentId)
		if err != nil {
			return err
		}

		ancestors = parent.Ancestors + "," + strconv.FormatInt(parent.DeptID, 10)
	}
	sysDept.Ancestors = ancestors

	return biz.repo.Insert(ctx, sysDept)
}

// 查看部门
func (biz *Depatmentusecase) Show(ctx context.Context, id int64) (*SysDept, error) {
	return biz.repo.FindOne(ctx, id)
}

// 编辑部门
func (biz *Depatmentusecase) Update(ctx context.Context, sysDept *SysDept) error {

	_, err := biz.repo.FindOne(ctx, sysDept.DeptID)
	if err != nil {
		return err
	}

	return biz.repo.Update(ctx, sysDept)
}

// 删除部门
func (biz *Depatmentusecase) Delete(ctx context.Context, id int64) (*SysDept, error) {
	one, err := biz.repo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	err = biz.repo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return one, nil
}
