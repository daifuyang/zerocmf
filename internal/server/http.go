package server

import (
	v1 "zerocmf/internal/server/api/v1"
	"zerocmf/internal/server/middleware"
	"zerocmf/internal/service"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer new an HTTP server.

type Server struct {
	Router *gin.Engine
	SrvCtx *service.Context
}

func NewHTTPServer(srvCtx *service.Context) *Server {

	service.NewMenu(srvCtx).ImportMenu()

	r := gin.New()

	return &Server{
		Router: r,
		SrvCtx: srvCtx,
	}
}

func (s *Server) Start() {
	middleware.NewMiddleware(s.Router, s.SrvCtx).UseMiddleware() // 初始化中间件
	v1.RegisterHTTPServer(s.Router, s.SrvCtx)
}
