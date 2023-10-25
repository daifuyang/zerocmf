package server

import (
	"fmt"
	"time"
	v1 "zerocmf/internal/server/api/v1"
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
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	return &Server{
		Router: r,
		SrvCtx: srvCtx,
	}
}

func (s *Server) Start() {
	v1.RegisterHTTPServer(s.Router, s.SrvCtx)
}
