package biz

import "time"

type SysRole struct {
	RoleID            int64      `gorm:"column:role_id;primaryKey;comment:角色ID" json:"roleId"`
	RoleName          string     `gorm:"column:role_name;not null;comment:角色名称" json:"roleName"`
	RoleSort          int        `gorm:"column:role_sort;not null;comment:显示顺序" json:"roleSort"`
	DataScope         string     `gorm:"column:data_scope;default:1;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）" json:"dataScope"`
	MenuCheckStrictly bool       `gorm:"column:menu_check_strictly;default:1;comment:菜单树选择项是否关联显示" json:"menuCheckStrictly"`
	DeptCheckStrictly bool       `gorm:"column:dept_check_strictly;default:1;comment:部门树选择项是否关联显示" json:"deptCheckStrictly"`
	Status            string     `gorm:"column:status;not null;comment:角色状态（0正常 1停用）" json:"status"`
	DelFlag           string     `gorm:"column:del_flag;default:0;comment:删除标志（0代表存在 2代表删除）" json:"delFlag"`
	CreateId          int64      `gorm:"column:create_id;default:0;comment:创建者" json:"createId"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime;index;comment:创建时间" json:"createdAt"`
	UpdateId          int64      `gorm:"column:update_id;default:0;comment:更新者" json:"updateId"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime;index;comment:更新时间" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;default:null;index;comment:删除时间" json:"deletedAt"`
}

// 设置表名
func (SysRole) TableName() string {
	return "sys_role"
}
