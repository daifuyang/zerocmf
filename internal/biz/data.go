package biz

type SysInfo struct {
	CreateId  int64      `gorm:"column:create_id;type:bigint(20);comment:创建者" json:"createId"`
	CreatedAt LocalTime  `gorm:"column:created_at;autoCreateTime;index;comment:创建时间" json:"createdAt"`
	UpdateId  int64      `gorm:"column:update_id;type:bigint(20);comment:更新者" json:"updateId"`
	UpdatedAt LocalTime  `gorm:"column:updated_at;autoUpdateTime;index;comment:更新时间" json:"updatedAt"`
	DeletedAt *LocalTime `gorm:"column:deleted_at;default:null;index;comment:删除时间" json:"deletedAt"`
}
