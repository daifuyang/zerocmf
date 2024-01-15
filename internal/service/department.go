package service

import (
	"strconv"
	"strings"
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
	biz.Dept
	Children []*DeptTree `json:"children,omitempty"`
}

func buildDeptTree(depts []*biz.Dept, parentID int64) []*DeptTree {
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

	// 解析请求参数
	var query biz.DeptListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()
	Dept, err := s.deptuc.Tree(ctx, &query)
	if err != nil {
		response.Error(c, err)
		return
	}

	if strings.TrimSpace(query.DeptName) != "" || query.Status != nil {
		response.Success(c, "获取成功！", Dept)
		return
	}

	deptTree := buildDeptTree(Dept, 0)
	response.Success(c, "获取成功！", deptTree)
}

// 详情
func (s *department) Show(c *gin.Context) {

	var uri struct {
		Id int64 `uri:"id"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Error(c, err)
		return
	}

	id := uri.Id

	ctx := c.Request.Context()
	dept, err := s.deptuc.Show(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", dept)
}

// 新增
func (s *department) Add(c *gin.Context) {
	s.save(c)
}

// 编辑
func (s *department) Edit(c *gin.Context) {
	s.save(c)
}

// 新增和编辑
func (s *department) save(c *gin.Context) {
	var req struct {
		DeptName  string `json:"deptName"`
		ParentId  int64  `json:"parentId"`
		ListOrder *int   `json:"listOrder"`
		Status    *int   `json:"status"`
		Remark    string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	id := c.Param("id")

	ctx := c.Request.Context()
	userId, err := utils.UserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var saveData biz.Dept

	// 新增
	msg := ""
	if id == "" {
		err = copier.Copy(&saveData, &req)
		if err != nil {
			response.Error(c, err)
			return
		}
		saveData.CreateId = userId
		err = s.deptuc.Insert(ctx, &saveData)
		msg = "添加成功！"
	} else {

		// 获取当前部门是否存在
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.Error(c, err)
			return
		}
		one, err := s.deptuc.Show(ctx, idInt)
		if err != nil {
			response.Error(c, err)
			return
		}

		copier.Copy(&one, &req)
		saveData.UpdateId = userId
		err = s.deptuc.Update(ctx, one)
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	if id == 1 {
		response.Error(c, "顶级部门不能删除！")
		return
	}

	ctx := c.Request.Context()

	count, err := s.deptuc.CountByParentId(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}

	if count > 0 {
		response.Error(c, "请先删除子部门！")
		return
	}

	Dept, err := s.deptuc.Delete(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "删除成功", Dept)
}
