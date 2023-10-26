package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *Middleware) ZapLoggerMiddleware() {
	m.Router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next() // 执行后续中间件

		// Since returns the time elapsed since t.
		// It is shorthand for time.Now().Sub(t).
		cost := time.Since(start)
		logger := m.SrvCtx.Logger
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost), // 运行时间
		)

	})

	m.Router.Use(gin.Recovery())
}
