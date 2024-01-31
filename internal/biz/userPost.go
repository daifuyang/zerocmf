package biz

import (
	"context"

	"gorm.io/gorm"
)

// 用户和部门的关联表
type UserPost struct {
	UserID   int64  `gorm:"column:user_id;type:int(11);primaryKey;comment:用户ID" json:"userId"`
	PostCode string `gorm:"column:post_code;type:varchar(64);primaryKey;comment:部门ID" json:"postCode"`
}

// 设置表名
func (*UserPost) TableName() string {
	return "cmf_user_post"
}

// 数据库迁移
func (biz *UserPost) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}
	return nil
}

// 定义data层接口
type UserPostRepo interface {
	Save(ctx context.Context, userRole *UserRole) error // 新增
}

type UserPostusecase struct {
	repo UserRoleRepo
}

func NewUserPostUsecase(repo UserPostRepo) *UserPostusecase {
	return &UserPostusecase{repo: repo}
}

// todo 获取用户绑定的岗位列表
