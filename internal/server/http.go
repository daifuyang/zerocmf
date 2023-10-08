package server

import "github.com/gin-gonic/gin"

// NewHTTPServer new an HTTP server.
func NewHTTPServer(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
