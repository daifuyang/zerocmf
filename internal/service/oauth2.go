package service

import (
	"context"

	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
)

func (s *Context) Callback(c *gin.Context) {

	oauthConfig := s.oauthConfig
	code := c.Query("code")
	if code == "" {
		response.Error(c, response.ErrBind)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(200, token)
}
