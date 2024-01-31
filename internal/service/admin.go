package service

import (
	"strconv"
	"strings"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"
	"zerocmf/pkg/hashed"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type admin struct {
	*Context
}

func NewAdmin(c *Context) *admin {
	return &admin{
		Context: c,
	}
}

// 列表

func (s *admin) List(c *gin.Context) {
	paginate, err := s.useruc.Find(c.Request.Context(), &biz.UserListQuery{Current: 1, PageSize: 10, UserType: 1})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", paginate)
}

// 获取单个管理员

func (s *admin) Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()

	user, err := s.useruc.FindUserByUserID(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 获取当前绑定的岗位信息

	// 获取当前绑定的角色信息

	response.Success(c, "获取成功！", user)
}

// 创建管理员
func (s *admin) Add(c *gin.Context) {
	s.Save(c)
}

// 更新管理员信息
func (s *admin) Edit(c *gin.Context) {
	s.Save(c)
}

// 统一保存
func (s *admin) Save(c *gin.Context) {
	// 定义请求结构体
	var req struct {
		UserName    string    `json:"userName"`
		DeptId      *int64    `json:"deptId"`
		PhoneNumber *string   `json:"phoneNumber"`
		Email       *string   `json:"email"`
		LoginName   string    `json:"loginName"`
		UserType    *int      `json:"userType"`
		Password    string    `json:"password"`
		Gender      *int      `json:"gender"`
		Status      *int      `json:"status"`
		PostCodes   *[]string `json:"postCodes"`
		RoleIds     *[]int64  `json:"roleIds"`
		Remark      *string   `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 先处理必填项
	if strings.TrimSpace(req.UserName) == "" {
		response.Error(c, "用户昵称不能为空！")
		return
	}

	if strings.TrimSpace(req.LoginName) == "" {
		response.Error(c, "登录账号不能为空！")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		response.Error(c, "登录密码不能为空！")
		return
	}

	id := c.Param("id")
	ctx := c.Request.Context()

	var user biz.User
	userId, err := utils.UserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if req.UserType == nil {
		var userType = 1
		req.UserType = &userType
	}

	err = copier.Copy(&user, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	salt := s.Config.Mysql.Salt
	password, err := hashed.Password("123456", salt)
	if err != nil {
		response.Error(c, err)
		return
	}
	user.Password = password

	var msg = ""

	if id == "" {
		// 创建
		user.CreateId = userId
		err = s.useruc.Insert(ctx, &user)
		if err != nil {
			response.Error(c, err)
			return
		}
		msg = "创建成功！"
	} else {
		user.UpdateId = userId
		err = s.useruc.Update(ctx, &user)
		if err != nil {
			response.Error(c, err)
			return
		}
		msg = "更新成功！"
	}
	response.Success(c, msg, user)
}

// 删除单个管理员账号
func (s *admin) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()

	// 先判断是否存在管理员
	user, err := s.useruc.FindUserByUserID(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = s.useruc.DeleteOne(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "删除成功！", user)
}
