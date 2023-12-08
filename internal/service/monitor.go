package service

import (
	"github.com/gin-gonic/gin"
)

type monitor struct {
	*Context
}

func NewMonitor(c *Context) *monitor {
	return &monitor{c}
}

func (s *monitor) Index(c *gin.Context) {

}
