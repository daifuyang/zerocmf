package biz

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SysMenu struct {
	MenuID    int64      `gorm:"column:menu_id;primaryKey;comment:菜单ID" json:"menuId"`
	MenuName  string     `gorm:"column:menu_name;not null;comment:菜单名称" json:"menuName" binding:"required"`
	ParentID  int64      `gorm:"column:parent_id;default:0;comment:父菜单ID" json:"parentId"`
	ListOrder float64    `gorm:"column:list_order;default:10000;comment:显示顺序" json:"listOrder"` // 修改字段名
	Path      string     `gorm:"column:path;default:'';comment:路由地址" json:"path"`
	IsFrame   int        `gorm:"column:is_frame;default:0;comment:是否为外链（0：否 1：是）" json:"isFrame"`
	MenuType  int        `gorm:"column:menu_type;default:0;comment:菜单类型（0目录 1菜单 2按钮）" json:"menuType"`
	Visible   string     `gorm:"column:visible;default:1;comment:菜单状态（0：隐藏 1：显示）" json:"visible"`
	Status    int        `gorm:"column:status;default:1;comment:菜单状态（0：停用 1：显示）" json:"status"`
	Perms     string     `gorm:"column:perms;default:null;comment:权限标识" json:"perms"`
	Icon      string     `gorm:"column:icon;default:'';comment:菜单图标" json:"icon"`
	CreateId  int64      `gorm:"column:create_id;default:0;comment:创建者" json:"createId"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;index;comment:创建时间" json:"createdAt"`
	UpdateId  int64      `gorm:"column:update_id;default:0;comment:更新者" json:"updateId"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime;index;comment:更新时间" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;default:null;index;comment:删除时间" json:"deletedAt"`
	Remark    string     `gorm:"column:remark;default:'';comment:备注" json:"remark"`
	Children  []*SysMenu `gorm:"-" json:"children,omitempty"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

// 数据库迁移
func (biz *SysMenu) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}
	return nil
}

type MenuRepo interface {
	Find(ctx context.Context) (menus []*SysMenu, err error)                   // 查看全部
	FindOne(ctx context.Context, id int64) (*SysMenu, error)                  // 查询一条
	FindOneByMenuName(ctx context.Context, menuName string) (*SysMenu, error) // 根据菜单名称查找
	Insert(ctx context.Context, menu *SysMenu) (err error)                    // 插入一条
	Update(ctx context.Context, menu *SysMenu) (err error)                    // 更新一条
	Delete(ctx context.Context, id int64) error                               // 删除一条
}

type Menusecase struct {
	repo MenuRepo
}

func NewMenusecase(repo MenuRepo) *Menusecase {
	return &Menusecase{repo: repo}
}

// 根据菜单名称获取菜单
func (biz *Menusecase) FindOneByMenuName(ctx context.Context, menuName string) (sysMenu *SysMenu, err error) {
	return biz.repo.FindOneByMenuName(ctx, menuName)
}

// 获取全部数据
func (biz *Menusecase) Find(ctx context.Context) (sysMenus []*SysMenu, err error) {
	return biz.repo.Find(ctx)
}

// 查看一条数据
func (biz *Menusecase) FindOne(ctx context.Context, id int64) (*SysMenu, error) {
	return biz.repo.FindOne(ctx, id)
}

// 插入一条数据
func (biz *Menusecase) Insert(ctx context.Context, menu *SysMenu) (err error) {
	return biz.repo.Insert(ctx, menu)
}

// 更新一条数据
func (biz *Menusecase) Update(ctx context.Context, menu *SysMenu) (err error) {
	return biz.repo.Update(ctx, menu)
}

// 软删除一条数据
func (biz *Menusecase) Delete(ctx context.Context, id int64) (*SysMenu, error) {
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
