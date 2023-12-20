package service

import "github.com/gin-gonic/gin"

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

}

// 获取当个岗位
func (s *post) Show(c *gin.Context) {

}

// 新增单个岗位
func (s *post) Add(c *gin.Context) {

}

// 更新单个岗位
func (s *post) Update(c *gin.Context) {

}

// 【解耦】统一保存新增和更新逻辑
func (s *post) Save(c *gin.Context) {

}

// 删除单个岗位
func (s *post) Delete(c *gin.Context) {

}
