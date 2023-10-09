package user

import "time"

type User struct {
	ID            uint `gorm:"primaryKey;column:user_id"`
	DeptID        uint `gorm:"column:dept_id"`
	LoginName     string
	UserName      string `gorm:"column:user_name"`
	UserType      string `gorm:"column:user_type;default:00"`
	Email         string
	PhoneNumber   string `gorm:"column:phonenumber"`
	Sex           string `gorm:"column:sex;default:0"`
	Avatar        string
	Password      string
	Salt          string
	Status        string    `gorm:"column:status;default:0"`
	DelFlag       string    `gorm:"column:del_flag;default:0"`
	LoginIP       string    `gorm:"column:login_ip"`
	LoginDate     time.Time `gorm:"column:login_date"`
	PwdUpdateDate time.Time `gorm:"column:pwd_update_date"`
	CreateBy      string    `gorm:"column:create_by"`
	CreateTime    time.Time `gorm:"column:create_time"`
	UpdateBy      string    `gorm:"column:update_by"`
	UpdateTime    time.Time `gorm:"column:update_time"`
	Remark        string    `gorm:"column:remark"`
}

// TableName 指定表名
func (User) TableName() string {
	return "sys_user"
}
