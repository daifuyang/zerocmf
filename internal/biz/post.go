package biz

import (
	"context"

	"gorm.io/gorm"
)

// SysPost 表示 sys_post 表的数据模型
type SysPost struct {
	PostID    int64  `gorm:"column:post_id;type:int(11);primaryKey;comment:岗位ID" json:"postId"`
	PostCode  string `gorm:"column:post_code;type:varchar(64);not null;comment:岗位编码" json:"postCode"`
	PostName  string `gorm:"column:post_name;type:varchar(50);not null;comment:岗位名称" json:"postName"`
	ListOrder int    `gorm:"column:list_order;not null;type:int(8);comment:显示顺序" json:"listOrder"`
	Status    int    `gorm:"column:status;type:tinyint(2);not null;default 1;comment:状态（1正常;0停用）" json:"status"`
	Remark    string `gorm:"comment:备注;type:varchar(500)" json:"remark"`
	SysInfo
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

	// 创建演示数据

	if err := db.Where("post_code = 'CEO'").FirstOrCreate(&SysPost{
		PostID:    0,
		PostCode:  "CEO",
		PostName:  "董事长",
		ListOrder: 0,
		Status:    1,
		Remark:    "",
		SysInfo: SysInfo{
			CreateId: 1,
		},
	}).Error; err != nil {
		return err
	}

	if err := db.Where("post_code = 'SE'").FirstOrCreate(&SysPost{
		PostCode: "SE",
		PostName: "项目经理",
		Status:   1,
		SysInfo: SysInfo{
			CreateId: 1,
		},
	}).Error; err != nil {
		return err
	}

	if err := db.Where("post_code = 'HR'").FirstOrCreate(&SysPost{
		PostCode: "HR",
		PostName: "人力资源",
		Status:   1,
		SysInfo: SysInfo{
			CreateId: 1,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

// 定义repo实体接口 （依赖倒置原则）
type PostRepo interface {
	Find(ctx context.Context, listQuery *SysPostListQuery) (interface{}, error) // 查看全部
	First(query interface{}, args ...interface{}) (*SysPost, error)             // 根据条件查询一条
	FindOne(ctx context.Context, id int64) (*SysPost, error)                    // 查询一条
	Insert(ctx context.Context, post *SysPost) error                            // 插入一条
	Update(ctx context.Context, post *SysPost) error                            // 更新一条
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

// 根据条件查询一条
func (biz *Postusecase) First(query interface{}, args ...interface{}) (*SysPost, error) {
	return biz.repo.First(query, args...)
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

	_, err := biz.repo.FindOne(ctx, post.PostID)
	if err != nil {
		return err
	}

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
