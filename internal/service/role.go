package service

import (
	"strconv"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type role struct {
	*Context
}

func NewRole(c *Context) *role {
	return &role{
		Context: c,
	}
}

// 获取所有角色
func (s *role) List(c *gin.Context) {

	var query biz.SysRoleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, err)
		return
	}

	data, err := s.rc.Find(c.Request.Context(), &query)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", data)
}

// 获取单个角色
func (s *role) Show(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	data, err := s.rc.FindOne(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	menuIds, err := s.rc.FindPermissions(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	data.MenuIds = menuIds

	response.Success(c, "获取成功！", data)
}

// 创建新角色
func (s *role) Add(c *gin.Context) {
	s.Save(c)
}

// 更新角色信息
func (s *role) Update(c *gin.Context) {
	s.Save(c)
}

// 统一保存角色
func (s *role) Save(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		RoleName          string `json:"roleName"`
		ListOrder         *int   `json:"listOrder"`
		DataScope         int    `json:"dataScope"`
		MenuCheckStrictly *bool  `json:"menuCheckStrictly"`
		DeptCheckStrictly *bool  `json:"deptCheckStrictly"`
		Status            *int   `json:"status"`
		Remark            string `json:"remark"`
		MenuIds           []int  `json:"menuIds"` // 角色拥有的菜单权限id
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	var saveData biz.SysRole

	userId, err := utils.UserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	msg := "添加成功！"

	ctx := c.Request.Context()

	if id == "" {
		err := copier.Copy(&saveData, &req)
		if err != nil {
			response.Error(c, err)
			return
		}

		saveData.CreateId = userId

		err = s.rc.Insert(ctx, &saveData)
		if err != nil {
			response.Error(c, err)
			return
		}
	} else {

		idInt, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			response.Error(c, err)
			return
		}
		role, err := s.rc.FindOne(ctx, idInt)
		if err != nil {
			response.Error(c, err)
			return
		}

		err = copier.Copy(&role, &req)
		if err != nil {
			response.Error(c, err)
			return
		}

		saveData.UpdateId = userId
		err = s.rc.Update(ctx, role)
		if err != nil {
			response.Error(c, err)
			return
		}
		msg = "修改成功！"
	}
	response.Success(c, msg, saveData)
}

// 删除角色信息
func (s *role) Delete(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	sysRole, err := s.rc.Delete(c.Request.Context(), id)

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "删除成功", sysRole)

}
