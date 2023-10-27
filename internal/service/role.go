package service

import (
	"zerocmf/internal/utils"

	"github.com/gin-gonic/gin"
)

type role struct {
	*Context
}

func NewRole(c *Context) *role {
	return &role{
		Context: c,
	}
}

// 获取所有角色(分页)
func (s *role) List(c *gin.Context) {
	current, pageSize := utils.ParsePagination(c)
}

// 获取单个角色
func (s *role) Show(c *gin.Context) {

}

// 创建新角色
func (s *role) Add(c *gin.Context) {

}

// 跟新角色信息
func (s *role) Update(c *gin.Context) {

}

// 统一保存角色
func (s *role) Save(c *gin.Context) {

}

// 删除角色信息
func (s *role) Delete(c *gin.Context) {

}
