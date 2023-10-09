package server

import (
	"zerocmf/configs"
	"zerocmf/internal/service"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer new an HTTP server.

func NewHTTPServer(c *configs.Config, test *service.TestService) *gin.Engine {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		test.CreateArticle(c.Request.Context())
		c.JSON(200, "pong")
	})
	return r
}
