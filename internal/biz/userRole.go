package biz

import (
	"context"

	"gorm.io/gorm"
)

// 用户和角色的关联表
type UserRole struct {
	UserID int64 `gorm:"column:user_id;type:int(11);primaryKey;comment:用户ID" json:"userId"`
	RoleID int64 `gorm:"column:role_id;type:int(11);primaryKey;comment:角色ID" json:"roleId"`
}

// 设置表名
func (*UserRole) TableName() string {
	return "cmf_user_role"
}

// 数据库迁移
func (biz *UserRole) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}
	return nil
}

// 定义data层接口
type UserRoleRepo interface {
	Save(ctx context.Context, userRole *UserRole) error // 新增
}

type UserRoleusecase struct {
	repo UserRoleRepo
}

func NewUserRoleUsecase(repo UserRoleRepo) *UserRoleusecase {
	return &UserRoleusecase{repo: repo}
}

// todo 获取用户绑定的角色列表
