package biz

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID        uint64    `gorm:"primaryKey;autoIncrement;size:20" json:"user_id"`
	DeptID        *uint64   `gorm:"comment:'部门ID';size:20;type:bigint(20)" json:"dept_id"`
	LoginName     string    `gorm:"not null;comment:'登录账号';size:30;type:varchar(30)" json:"login_name"`
	UserName      string    `gorm:"default:'';comment:'用户昵称';size:30;type:varchar(30)" json:"user_name"`
	UserType      string    `gorm:"default:'00';comment:'用户类型（00系统用户 01注册用户）';size:2;type:varchar(2)" json:"user_type"`
	Email         string    `gorm:"default:'';comment:'用户邮箱';size:50;type:varchar(50)" json:"email"`
	PhoneNumber   string    `gorm:"default:'';comment:'手机号码';size:11;type:varchar(11)" json:"phonenumber"`
	Sex           string    `gorm:"default:'0';comment:'用户性别（0男 1女 2未知）';size:1;type:char(1)" json:"sex"`
	Avatar        string    `gorm:"default:'';comment:'头像路径';size:100;type:varchar(100)" json:"avatar"`
	Password      string    `gorm:"default:'';comment:'密码';size:50;type:varchar(50)" json:"-"`
	Salt          string    `gorm:"default:'';comment:'盐加密';size:20;type:varchar(20)" json:"-"`
	Status        string    `gorm:"default:'0';comment:'帐号状态（0正常 1停用）';size:1;type:char(1)" json:"status"`
	DelFlag       string    `gorm:"default:'0';comment:'删除标志（0代表存在 2代表删除）';size:1;type:char(1)" json:"-"`
	LoginIP       string    `gorm:"default:'';comment:'最后登录IP';size:128;type:varchar(128)" json:"login_ip"`
	LoginDate     time.Time `gorm:"comment:'最后登录时间'" json:"login_date"`
	PwdUpdateDate time.Time `gorm:"comment:'密码最后更新时间'" json:"pwd_update_date"`
	CreateBy      string    `gorm:"default:'';comment:'创建者';size:64;type:varchar(64)" json:"create_by"`
	CreateTime    time.Time `gorm:"comment:'创建时间'" json:"create_time"`
	UpdateBy      string    `gorm:"default:'';comment:'更新者';size:64;type:varchar(64)" json:"update_by"`
	UpdateTime    time.Time `gorm:"comment:'更新时间'" json:"update_time"`
	Remark        string    `gorm:"comment:'备注';size:500;type:varchar(500)" json:"remark"`
}

// TableName 指定表名
func (*User) TableName() string {
	return "sys_user"
}

// 数据库迁移
func (bz *User) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&bz)
	return err
}

// 定义data层接口
type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
}

type Userusecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *Userusecase {
	return &Userusecase{repo: repo}
}

type Register struct {
	Account    string `form:"account" json:"account"` // 账号 邮箱 手机号
	Code       uint   `form:"code" json:"code"`       // 手机验证码
	Password   string `json:"password"`               // 密码
	RePassword string `json:"repassword"`             //
	Type       uint   `json:"type"`                   // 注册类型 0：邮箱手机号注册，1：账号密码注册
}

// 业务方组装层，调用repo数据层
func (uc *Userusecase) Register(ctx context.Context) (ps []*User, err error) {
	return
}
