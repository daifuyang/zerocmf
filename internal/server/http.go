package server

import (
	v1 "zerocmf/internal/server/api/v1"
	"zerocmf/internal/service"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer new an HTTP server.

func NewHTTPServer(svcCtx *service.Context) *gin.Engine {
	r := gin.New()
	v1.RegisterHTTPServer(r, svcCtx)
	return r
}
