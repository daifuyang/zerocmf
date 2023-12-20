package biz

import (
	"context"

	"gorm.io/gorm"
)

// SysPost 表示 sys_post 表的数据模型
type SysPost struct {
	PostID    int64      `gorm:"column:post_id;primaryKey;comment:岗位ID" json:"postId"`
	PostCode  string     `gorm:"column:post_code;size:64;unique;not null;comment:岗位编码" json:"postCode"`
	PostName  string     `gorm:"column:post_name;size:50;not null;comment:岗位名称" json:"postName"`
	ListOrder int        `gorm:"column:list_order;not null;comment:显示顺序" json:"listOrder"`
	Status    string     `gorm:"column:status;size:1;not null;default 1;comment:状态（1正常;0停用）" json:"status"`
	CreateId  int64      `gorm:"column:create_id;default:0;comment:创建者" json:"createId"`
	CreatedAt LocalTime  `gorm:"column:created_at;autoCreateTime;index;comment:创建时间" json:"createdAt"`
	UpdateId  int64      `gorm:"column:update_id;default:0;comment:更新者" json:"updateId"`
	UpdatedAt LocalTime  `gorm:"column:updated_at;autoUpdateTime;index;comment:更新时间" json:"updatedAt"`
	DeletedAt *LocalTime `gorm:"column:deleted_at;default:null;index;comment:删除时间" json:"deletedAt"`
	Remark    string     `gorm:"comment:备注;size:500;type:varchar(500)" json:"remark"`
}

// 列表筛选条件

type SysPostListQuery struct {
	PostCode string `form:"postCode"`
	PostName string `form:"postName"`
	Status   *int   `form:"status"`
	PaginateQuery
}

// TableName 指定表名
func (SysPost) TableName() string {
	return "sys_post"
}

// 数据库迁移
func (biz *SysPost) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}
	return nil
}

// 定义repo实体接口 （依赖倒置原则）
type PostRepo interface {
	Find(ctx context.Context, listQuery *SysPostListQuery) (interface{}, error) // 查看全部
	FindOne(ctx context.Context, id int64) (*SysPost, error)                    // 查询一条
	Insert(ctx context.Context, menu *SysPost) (err error)                      // 插入一条
	Update(ctx context.Context, menu *SysPost) (err error)                      // 更新一条
	Delete(ctx context.Context, id int64) error                                 // 删除一条
}

// 定义业务用例
type Postusecase struct {
	repo PostRepo
}

// 使用wire进行依赖注入
func NewPostusecase(repo PostRepo) *Postusecase {
	return &Postusecase{
		repo: repo,
	}
}

// 获取列表
func (biz *Postusecase) Find(ctx context.Context, listQuery *SysPostListQuery) (interface{}, error) {
	return biz.repo.Find(ctx, listQuery)
}

// 获取一条数据
func (biz *Postusecase) FindOne(ctx context.Context, id int64) (*SysPost, error) {
	return biz.repo.FindOne(ctx, id)
}

// 新增一条数据
func (biz *Postusecase) Insert(ctx context.Context, post *SysPost) error {
	return biz.repo.Insert(ctx, post)
}

// 更新一条数据
func (biz *Postusecase) Update(ctx context.Context, post *SysPost) error {
	return biz.repo.Update(ctx, post)
}

// 删除一条数据
func (biz *Postusecase) Delete(ctx context.Context, id int64) (*SysPost, error) {
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
