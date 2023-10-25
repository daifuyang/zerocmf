package biz

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SysMenu struct {
	MenuID    int64     `gorm:"column:menu_id;primaryKey;comment:菜单ID" json:"menuId"`
	MenuName  string    `gorm:"column:menu_name;not null;comment:菜单名称" json:"menuName"`
	ParentID  int64     `gorm:"column:parent_id;default:0;comment:父菜单ID" json:"parentId"`
	OrderNum  int       `gorm:"column:order_num;default:0;comment:显示顺序" json:"orderNum"`
	Path      string    `gorm:"column:path;default:'';comment:路由地址" json:"path"`
	Component string    `gorm:"column:component;default:null;comment:组件路径" json:"component"`
	IsFrame   int       `gorm:"column:is_frame;default:0;comment:是否为外链（0：否 1：是）" json:"isFrame"`
	MenuType  int       `gorm:"column:menu_type;default:0;comment:菜单类型（0目录 1菜单 2按钮）" json:"menuType"`
	Visible   string    `gorm:"column:visible;default:1;comment:菜单状态（0：隐藏 1：显示）" json:"visible"`
	Status    int       `gorm:"column:status;default:1;comment:菜单状态（0：停用 1：显示）" json:"status"`
	Perms     string    `gorm:"column:perms;default:null;comment:权限标识" json:"perms"`
	Icon      string    `gorm:"column:icon;default:'';comment:菜单图标" json:"icon"`
	CreateId  int64     `gorm:"column:create_id;default:0;comment:创建者" json:"createId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdateId  int64     `gorm:"column:update_id;default:0;comment:更新者" json:"updateId"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	Remark    string    `gorm:"column:remark;default:'';comment:备注" json:"remark"`
	Children  []SysMenu `gorm:"-" json:"children"`
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
	FindOneByMenuName(ctx context.Context, menuName string) (sysMenu *SysMenu, err error)
	Insert(ctx context.Context, menu *SysMenu) (err error)
	Update(ctx context.Context, menu *SysMenu) (err error)
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

// 插入一条数据

func (biz *Menusecase) Insert(ctx context.Context, menu *SysMenu) (err error) {
	return biz.repo.Insert(ctx, menu)
}

// 更新一条数据
func (biz *Menusecase) Update(ctx context.Context, menu *SysMenu) (err error) {
	return biz.repo.Update(ctx, menu)
}