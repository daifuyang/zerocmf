package biz

import (
	"context"
	"net/http"
	"time"
	"zerocmf/pkg/hashed"

	v4Oauth2 "github.com/go-oauth2/oauth2/v4"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type User struct {
	UserID      uint64    `gorm:"primaryKey;autoIncrement;size:20" json:"userId"`
	DeptID      *uint64   `gorm:"index;comment:部门ID;size:20;type:bigint(20)" json:"deptId"`
	LoginName   string    `gorm:"comment:登录账号;size:30;type:varchar(30)" json:"loginName"`
	UserName    string    `gorm:"comment:用户昵称;size:30;type:varchar(30)" json:"userName"`
	ListOrder   int       `gorm:"column:list_order;default:0;comment:显示顺序" json:"listOrder"`
	UserType    uint      `gorm:"default:1;comment:用户类型（0:系统用户 1:注册用户）;size:2;type:tinyint(2)" json:"userType"`
	Email       string    `gorm:"default:null;comment:用户邮箱;size:50;type:varchar(50)" json:"email"`
	PhoneNumber string    `gorm:"default:null;comment:手机号码;size:11;type:varchar(11)" json:"phoneNumber"`
	Sex         uint      `gorm:"default:0;comment:用户性别（0男 1女 2未知）;size:2;type:tinyint(2)" json:"sex"`
	Avatar      string    `gorm:"comment:头像路径;size:100;type:varchar(100)" json:"avatar"`
	Password    string    `gorm:"not null;comment:密码;size:100;type:varchar(100)" json:"-"`
	Salt        string    `gorm:"comment:盐加密;size:20;type:varchar(20)" json:"salt"`
	Status      uint      `gorm:"default:1;comment:帐号状态（0：停用 ,1：启用）;size:2;type:tinyint(2)" json:"status"`
	LoginIP     string    `gorm:"default:'';comment:最后登录IP;size:128;type:varchar(128)" json:"loginIP"`
	LoginedAt   time.Time `gorm:"autoCreateTime" json:"loginedAt"`
	PwdUpdateAt time.Time `gorm:"autoUpdateTime" json:"pwdUpdatedAt"`
	OperateId   uint64    `gorm:"comment:操作人;index" json:"operateId"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;index" json:"updatedAt"`
	DeletedAt   time.Time `gorm:"default:null;comment:删除时间;index" json:"deletedAt"`
	Remark      string    `gorm:"comment:备注;size:500;type:varchar(500)" json:"remark"`
}

// TableName 指定表名
func (*User) TableName() string {
	return "sys_user"
}

// 数据库迁移
func (biz *User) AutoMigrate(db *gorm.DB, salt string) error {
	err := db.AutoMigrate(&biz)
	if err != nil {
		return err
	}
	// 创建管理员admin

	password, err := hashed.Password("123456", salt)
	if err != nil {
		return err
	}

	tx := db.Where("login_name", "admin").FirstOrCreate(&User{
		LoginName: "admin",
		Password:  password,
		Salt:      salt,
	})
	return tx.Error
}

// 定义data层接口
type UserRepo interface {
	FindOne(ctx context.Context, id int64) (*User, error)
	Token(ctx context.Context, user *User) (*oauth2.Token, error)        // 获取token
	ValidationBearerToken(req *http.Request) (v4Oauth2.TokenInfo, error) //验证token
	CreateUser(ctx context.Context, user *User) error
	FindUserByAccount(ctx context.Context, account string) (*User, error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
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

type SmsCode struct {
	Account string `form:"account"` // 账号 邮箱 手机号
}

type Login struct {
	Account  string `json:"account"`  // 账号 邮箱 手机号
	Code     uint   ` json:"code"`    // 手机验证码
	Password string `json:"password"` // 密码
}

// 业务方组装层，调用repo数据层

// 登录
func (uc *Userusecase) Token(ctx context.Context, user *User) (*oauth2.Token, error) {
	return uc.repo.Token(ctx, user)
}

// 验证token
func (uc *Userusecase) ValidationBearerToken(req *http.Request) (v4Oauth2.TokenInfo, error) {
	return uc.repo.ValidationBearerToken(req)
}

// 注册
func (uc *Userusecase) Register(ctx context.Context, user *User) error {
	return uc.repo.CreateUser(ctx, user)
}

// 根据id查询单个用户
func (uc *Userusecase) FindUserByUserID(ctx context.Context, userID int64) (*User, error) {
	return uc.repo.FindOne(ctx, userID)
}

// 根据手机号查询单个用户
func (uc *Userusecase) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
	return uc.repo.FindUserByPhoneNumber(ctx, phoneNumber)
}

// 根据账号推断手机号，邮箱，登录账号
func (uc *Userusecase) FindUserByAccount(ctx context.Context, account string) (*User, error) {
	return uc.repo.FindUserByAccount(ctx, account)
}
