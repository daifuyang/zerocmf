package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

func Success(c *gin.Context, msg string, data interface{}) {
	res := response{
		Code: 1,
		Data: data,
		Msg:  msg,
	}
	c.JSON(http.StatusOK, res)
}

func Error(c *gin.Context, msg string, data interface{}) {
	res := response{
		Code: 0,
		Data: data,
		Msg:  msg,
	}
	c.JSON(http.StatusOK, res)
}
