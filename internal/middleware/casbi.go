package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"micro-shop/internal/util"
	"net/http"
)

func Casbin(enf *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName := c.GetString("roleName")
		if len(roleName) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(403, "request error", nil))
			c.Abort()
		}
		path := c.Request.URL.Path
		method := c.Request.Method
		if ok, err := enf.Enforce(roleName, path, method); err != nil || !ok {
			c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
			c.Abort()
		}
		c.Next()
	}
}
