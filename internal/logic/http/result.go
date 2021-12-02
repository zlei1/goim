package http

import (
	"github.com/gin-gonic/gin"
)

const (
	OK = 0
	RequestErr = 400
	ServerErr = 500
)

type resp struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func result(c *gin.Context, data interface{}, code int) {
	c.JSON(200, resp{
		Code: code,
		Data: data,
	})
}

func errors(c *gin.Context, code int, msg string) {
	c.JSON(200, resp{
		Code:    code,
		Message: msg,
	})
}