package biz

import (
	"context"
	"net/http"
	"zerocmf/pkg/hashed"

	v4Oauth2 "github.com/go-oauth2/oauth2/v4"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type User struct {
	UserID      int64     `gorm:"primaryKey;autoIncrement;size:20" json:"userId"`
	DeptID      *int64    `gorm:"index;comment:部门ID;size:20;type:bigint(20)" json:"deptId"`
	LoginName   string    `gorm:"comment:登录账号;size:30;type:varchar(30)" json:"loginName"`
	UserName    string    `gorm:"comment:用户昵称;size:30;type:varchar(30)" json:"userName"`
	ListOrder   int       `gorm:"column:list_order;default:0;comment:显示顺序" json:"listOrder"`
	UserType    int       `gorm:"default:1;comment:用户类型（1:系统用户 0:注册用户）;size:2;type:tinyint(2)" json:"userType"`
	Email       string    `gorm:"default:null;comment:用户邮箱;size:50;type:varchar(50)" json:"email"`
	PhoneNumber string    `gorm:"default:null;comment:手机号码;size:11;type:varchar(11)" json:"phoneNumber"`
	Gender      int       `gorm:"default:0;comment:用户性别（0男 1女 2未知）;size:2;type:tinyint(2)" json:"gender"`
	Avatar      string    `gorm:"comment:头像路径;size:100;type:varchar(100)" json:"avatar"`
	Password    string    `gorm:"not null;comment:密码;size:100;type:varchar(100)" json:"-"`
	Salt        string    `gorm:"comment:盐加密;size:20;type:varchar(20)" json:"salt"`
	Status      int       `gorm:"default:1;comment:帐号状态（0：停用 ,1：启用）;size:2;type:tinyint(2)" json:"status"`
	LoginIP     string    `gorm:"default:'';comment:最后登录IP;size:128;type:varchar(128)" json:"loginIP"`
	LoginedAt   LocalTime `gorm:"autoCreateTime" json:"loginedAt"`
	PwdUpdateAt LocalTime `gorm:"autoUpdateTime" json:"pwdUpdatedAt"`
	SysInfo
	Remark    string    `gorm:"comment:备注;size:500;type:varchar(500)" json:"remark"`
	PostCodes *[]string `gorm:"-" json:"postIds"`
	RoleIds   *[]int64  `gorm:"-" json:"roleIds"`
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
	Token(ctx context.Context, user *User) (*oauth2.Token, error)        // 获取token
	ValidationBearerToken(req *http.Request) (v4Oauth2.TokenInfo, error) //验证token
	FindUserByAccount(ctx context.Context, account string) (*User, error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	Find(ctx context.Context, query *UserListQuery) (*Paginate, error) // 查询列表
	FindOne(ctx context.Context, id int64) (*User, error)
	Insert(ctx context.Context, user *User) error  // 新增用户
	Update(ctx context.Context, user *User) error  // 更新用户
	DeleteOne(ctx context.Context, id int64) error // 删除用户
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
	Code     uint   `json:"code"`     // 手机验证码
	Password string `json:"password"` // 密码
}

type UserListQuery struct {
	Current  int   `form:"current"`
	PageSize int   `form:"pageSize"`
	UserType int64 `json:"userType"`
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

// 新增
func (uc *Userusecase) Insert(ctx context.Context, user *User) error {
	return uc.repo.Insert(ctx, user)
}

// 更新
func (uc *Userusecase) Update(ctx context.Context, user *User) error {
	return uc.repo.Update(ctx, user)
}

// 删除
func (uc *Userusecase) DeleteOne(ctx context.Context, id int64) error {
	return uc.repo.DeleteOne(ctx, id)
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

// 按分页查询列表
func (uc *Userusecase) Find(ctx context.Context, query *UserListQuery) (*Paginate, error) {
	return uc.repo.Find(ctx, query)
}
