package service

import (
	"strconv"
	"strings"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type post struct {
	*Context
}

func NewPost(c *Context) *post {
	return &post{
		Context: c,
	}
}

// 获取所有岗位
func (s *post) List(c *gin.Context) {
	// 解析请求参数
	var query biz.SysPostListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	data, err := s.postuc.Find(ctx, &query)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", data)
}

// 获取当个岗位
func (s *post) Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()
	data, err := s.postuc.FindOne(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取成功！", data)
}

// 新增单个岗位
func (s *post) Add(c *gin.Context) {
	s.Save(c)
}

// 更新单个岗位
func (s *post) Update(c *gin.Context) {
	s.Save(c)
}

// 【解耦】统一保存新增和更新逻辑
func (s *post) Save(c *gin.Context) {

	id := c.Param("id")

	var req struct {
		PostCode  string `json:"postCode"`
		PostName  string `json:"postName"`
		ListOrder *int   `json:"listOrder"`
		Status    *int   `json:"status"`
		Remark    string `json:"remark"`
	}

	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	var saveData biz.SysPost

	userId, err := utils.UserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	msg := "添加成功！"

	query := []string{"deleted_at is null", "post_code = ?"}
	queryArgs := []interface{}{req.PostCode}
	queryStr := strings.Join(query, " AND ")

	first, firstErr := s.postuc.First(queryStr, queryArgs...)
	if firstErr != nil && firstErr != gorm.ErrRecordNotFound {
		response.Error(c, firstErr)
		return
	}

	if id == "" {

		if first != nil {
			response.Error(c, "该岗位编码已存在！")
			return
		}

		err = copier.Copy(&saveData, &req)
		if err != nil {
			response.Error(c, err)
			return
		}

		saveData.CreateId = userId

		err = s.postuc.Insert(ctx, &saveData)
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
		one, err := s.postuc.FindOne(ctx, idInt)
		if err != nil {
			response.Error(c, err)
			return
		}

		if first != nil && first.PostID != one.PostID {
			response.Error(c, "该岗位编码已存在！")
			return
		}

		err = copier.Copy(&one, &req)
		if err != nil {
			response.Error(c, err)
			return
		}

		saveData.UpdateId = userId
		err = s.postuc.Update(ctx, one)
		if err != nil {
			response.Error(c, err)
			return
		}
		saveData = *one
		msg = "修改成功！"
	}
	response.Success(c, msg, saveData)
}

// 删除单个岗位
func (s *post) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	sysPost, err := s.postuc.Delete(c.Request.Context(), id)

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "删除成功", sysPost)
}
