package service

import (
	"fmt"
	"strconv"
	"strings"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
)

type authz struct {
	*Context
}

func NewAuthz(c *Context) *authz {
	return &authz{c}
}

// 查询对应的角色包含的权限
func (s *authz) Index(c *gin.Context) {
	roleId := c.Param("id")
	if strings.TrimSpace(roleId) == "" {
		response.Error(c, "参数错误："+roleId)
		return
	}

	roleIdInt, err := strconv.ParseInt(roleId, 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()

	// 判断当前角色是否存在
	_, err = s.rc.FindOne(ctx, roleIdInt)
	if err != nil {
		response.Error(c, err)
		return
	}

	access := s.Context.e.GetPermissionsForUser(roleId)
	response.Success(c, "获取成功！", access)
}

// 保存角色权限
func (s *authz) Save(c *gin.Context) {
	roleId := c.Param("id")
	if strings.TrimSpace(roleId) == "" {
		response.Error(c, "参数错误："+roleId)
		return
	}
	roleIdInt, err := strconv.ParseInt(roleId, 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()

	// 判断当前角色是否存在
	_, err = s.rc.FindOne(ctx, roleIdInt)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 获取提交的权限列表

	var req struct {
		Access []string `json:"access"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	//获取原来的数据
	originalStrings := s.e.GetPermissionsForUser(roleId)
	for _, p := range originalStrings {
		fmt.Println("ppp", p[1])
	}

	// 将原始字符串数组转换为 map 以便快速查找
	// originalMap := make(map[string]struct{})
	// for _, s := range originalStrings {
	// 	originalMap[s] = struct{}{}
	// }

	for i := 0; i < len(req.Access); i++ {
		s.e.AddPermissionForUser(roleId, req.Access[i])
	}

	response.Success(c, "操作成功！", req)
}
