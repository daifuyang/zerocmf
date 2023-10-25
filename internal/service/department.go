package service

import (
	"strconv"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"
	"zerocmf/pkg/response"

	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
)

type department struct {
	*Context
}

func NewDeparment(c *Context) *department {
	return &department{
		Context: c,
	}
}

type DeptTree struct {
	biz.SysDept
	Children []*DeptTree `json:"children,omitempty"`
}

func buildDeptTree(depts []*biz.SysDept, parentID int64) []*DeptTree {
	tree := make([]*DeptTree, 0)
	for _, dept := range depts {
		if dept.ParentID == parentID {
			child := &DeptTree{}
			copier.Copy(child, dept)
			child.Children = buildDeptTree(depts, dept.DeptID)
			tree = append(tree, child)
		}
	}
	return tree
}

// 查询列表树
func (s *department) Tree(c *gin.Context) {
	ctx := c.Request.Context()
	sysDept, err := s.dc.Tree(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	deptTree := buildDeptTree(sysDept, 0)
	response.Success(c, "获取成功！", deptTree)
}

// 新增

func (s *department) Add(c *gin.Context) {
	s.edit(c)
}

// 详情
func (s *department) Show(c *gin.Context) {

	var uri struct {
		Id int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Error(c, err)
		return
	}

	id := uri.Id

	ctx := c.Request.Context()
	dept, err := s.Context.dc.Show(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", dept)
}

// 编辑
func (s *department) Update(c *gin.Context) {
	s.edit(c)
}

// 新增和编辑
func (s *department) edit(c *gin.Context) {
	var req biz.SysDept
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBind)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBind)
		return
	}

	ctx := c.Request.Context()
	userId, err := utils.UserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 新增
	msg := ""
	if id == 0 {
		req.CreateId = userId
		err = s.Context.dc.Add(ctx, &req)
		msg = "添加成功！"
	} else {

		// 获取当前部门是否存在
		one, err := s.Context.dc.Show(ctx, id)
		if err != nil {
			response.Error(c, err)
			return
		}

		copier.Copy(&one, &req)

		req.DeptID = id
		req.UpdateId = userId
		err = s.Context.dc.Update(ctx, one)
		if err != nil {
			response.Error(c, err)
			return
		}
		msg = "更新成功！"
	}

	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, msg, req)
}

// 删除
func (s *department) Delete(c *gin.Context) {
	response.Success(c, "hello delete", nil)
}
