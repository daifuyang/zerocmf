package biz

import (
	"context"

	"gorm.io/gorm"
)

type Role struct {
	RoleID            int64  `gorm:"column:role_id;type:int(11);primaryKey;comment:角色ID" json:"roleId"`
	RoleName          string `gorm:"column:role_name;not null;comment:角色名称" json:"roleName"`
	ListOrder         int    `gorm:"column:list_order;type:int(8);default:0;comment:显示顺序" json:"listOrder"`
	DataScope         int    `gorm:"column:data_scope;type:tinyint(2);default:0;comment:数据范围（0：全部数据权限 1：自定数据权限 2：本部门数据权限 3：本部门及以下数据权限）" json:"dataScope"`
	MenuCheckStrictly bool   `gorm:"column:menu_check_strictly;type:tinyint(2);default:1;comment:菜单树选择项是否关联显示" json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `gorm:"column:dept_check_strictly;type:tinyint(2);default:1;comment:部门树选择项是否关联显示" json:"deptCheckStrictly"`
	Status            int    `gorm:"column:status;not null;type:tinyint(2);default:1;comment:角色状态（1：正常 0：停用）" json:"status"`
	Remark            string `gorm:"comment:备注;size:500;type:varchar(500)" json:"remark"`
	MenuIds           []*int `gorm:"-" json:"menuIds"` // 角色拥有的菜单权限id
	SysInfo
}

type RoleListQuery struct {
	RoleName string `form:"roleName"`
	Status   *int   `form:"status"`
	PaginateQuery
}

// 设置表名
func (*Role) TableName() string {
	return "cmf_role"
}

// 数据库迁移
func (biz *Role) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}

	// 创建演示数据

	if err := db.Where("role_name = '超级管理员'").FirstOrCreate(&Role{
		RoleName: "超级管理员",
		SysInfo: SysInfo{
			CreateId: 1,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

type RoleRepo interface {
	Find(ctx context.Context, listQuery *RoleListQuery) (interface{}, error) // 查看全部
	FindOne(ctx context.Context, id int64) (*Role, error)                    // 查询一条
	FindPermissions(ctx context.Context, id int64) ([]*int, error)           // 查询角色授权通过的权限
	Insert(ctx context.Context, role *Role) (err error)                      // 插入一条
	Update(ctx context.Context, role *Role) (err error)                      // 更新一条
	Delete(ctx context.Context, id int64) error                              // 删除一条
}

type Roleusecase struct {
	repo RoleRepo
}

func NewRoleusecase(repo RoleRepo) *Roleusecase {
	return &Roleusecase{
		repo: repo,
	}
}

// 获取全部数据
func (biz *Roleusecase) Find(ctx context.Context, listQuery *RoleListQuery) (interface{}, error) {
	return biz.repo.Find(ctx, listQuery)
}

// 查看一条数据
func (biz *Roleusecase) FindOne(ctx context.Context, id int64) (*Role, error) {
	return biz.repo.FindOne(ctx, id)
}

// 查询角色授权通过的权限
func (biz *Roleusecase) FindPermissions(ctx context.Context, id int64) ([]*int, error) {
	return biz.repo.FindPermissions(ctx, id)
}

// 新增一条数据
func (biz *Roleusecase) Insert(ctx context.Context, role *Role) error {
	return biz.repo.Insert(ctx, role)
}

// 更新一条数据
func (biz *Roleusecase) Update(ctx context.Context, role *Role) error {

	_, err := biz.repo.FindOne(ctx, role.RoleID)
	if err != nil {
		return err
	}

	return biz.repo.Update(ctx, role)
}

// 删除一条数据
func (biz *Roleusecase) Delete(ctx context.Context, id int64) (*Role, error) {
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
