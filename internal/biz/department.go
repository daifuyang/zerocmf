package biz

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SysDept struct {
	DeptID    int64     `gorm:"type:bigint;primary_key;autoIncrement;comment:'部门id'" uri:"id" json:"deptId"`
	ParentID  int64     `gorm:"type:bigint;default:0;comment:'父部门id'" json:"parentId"`
	Ancestors string    `gorm:"type:varchar(50);default:'';comment:'祖级列表'" json:"ancestors"`
	DeptName  string    `gorm:"type:varchar(30);default:'';comment:'部门名称'" json:"deptName"`
	ListOrder int       `gorm:"column:list_order;default:0;comment:显示顺序" json:"listOrder"`
	Leader    string    `gorm:"type:varchar(20);default:null;comment:'负责人'" json:"leader"`
	Phone     string    `gorm:"type:varchar(11);default:null;comment:'联系电话'" json:"phone"`
	Email     string    `gorm:"type:varchar(50);default:null;comment:'邮箱'" json:"email"`
	Status    string    `gorm:"type:char(1);default:'0';comment:'部门状态（0正常 1停用）'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	CreateId  int64     `gorm:"type:varchar(64);default:0;comment:'创建者'" json:"createdId"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;index" json:"updatedAt"`
	UpdateId  int64     `gorm:"type:varchar(64);default:0;comment:'更新者'" json:"updatedId"`
	DeletedAt time.Time `gorm:"default:null;comment:删除时间;index" json:"deletedAt"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

// 数据库迁移
func (biz *SysDept) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}

	db.Where("")

	return nil
}

type DepartmentRepo interface {
	GetOneById(ctx context.Context, id int64) (*SysDept, error)
	Index(ctx context.Context) ([]*SysDept, error)
	Add(ctx context.Context, sysDept *SysDept) error    // 添加部门
	Update(ctx context.Context, sysDept *SysDept) error // 更新部门
}

type Depatmentusecase struct {
	repo DepartmentRepo
}

func NewDepatmentusecase(repo DepartmentRepo) *Depatmentusecase {
	return &Depatmentusecase{repo: repo}
}

// 查看部门列表
func (uc *Depatmentusecase) Tree(ctx context.Context) ([]*SysDept, error) {
	return uc.repo.Index(ctx)
}

// 添加部门
func (uc *Depatmentusecase) Add(ctx context.Context, sysDept *SysDept) error {

	// 查询parentId是否合法
	parentId := sysDept.ParentID

	// 并入参数ancestors
	ancestors := "0"
	if parentId > 0 {
		parent, err := uc.repo.GetOneById(ctx, parentId)
		if err != nil {
			return err
		}

		ancestors = parent.Ancestors + "," + strconv.FormatInt(parent.DeptID, 10)
	}
	sysDept.Ancestors = ancestors

	return uc.repo.Add(ctx, sysDept)
}

// 查看部门
func (uc *Depatmentusecase) Show(ctx context.Context, id int64) (*SysDept, error) {
	return uc.repo.GetOneById(ctx, id)
}

// 编辑部门
func (uc *Depatmentusecase) Update(ctx context.Context, sysDept *SysDept) error {
	return uc.repo.Update(ctx, sysDept)
}
