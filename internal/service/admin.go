package service

import (
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
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
	current, pageSize := utils.ParsePagination(c)
	paginate, err := s.uc.Find(c.Request.Context(), &biz.UserListQuery{Current: current, PageSize: pageSize, UserType: 1})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", paginate)
}

// 获取单个管理员

func (s *admin) Show(c *gin.Context) {

}

// 创建管理员
func (s *admin) Add(c *gin.Context) {

}

// 更新管理员信息
func (s *admin) Update(c *gin.Context) {

}

// 统一保存
func (s *admin) Save(c *gin.Context) {

}

// 删除单个管理员账号
func (s *admin) Delete(c *gin.Context) {

}
