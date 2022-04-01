package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetResponse(code int, msg string, data interface{}) gin.H {
	return gin.H{
		"code": strconv.Itoa(code),
		"msg":  msg,
		"data": data,
	}
}
