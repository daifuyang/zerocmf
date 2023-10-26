package middleware

import (
	"zerocmf/internal/service"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Router *gin.Engine
	SrvCtx *service.Context
}

func NewMiddleware(router *gin.Engine, srvCtx *service.Context) *Middleware {
	return &Middleware{Router: router, SrvCtx: srvCtx}
}

func (m *Middleware) UseMiddleware() {
	m.ZapLoggerMiddleware()
}
