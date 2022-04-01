package middleware

import (
	"log"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func Jwt(logger *log.Logger, enf *casbin.Enforcer) gin.HandlerFunc {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("token:")
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if len(tokenStr) < 8 {
			c.JSON(http.StatusOK, util.GetResponse(401, "token is nil", nil))
			c.Abort()
			return
		}
		claim, err := util.ParseToken(tokenStr)
		if err != nil {
			logger.Printf("parse token :%s", err.Error())
			c.JSON(http.StatusOK, util.GetResponse(501, "parse token err", nil))
			c.Abort()
			return
		}
		if claim.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusOK, util.GetResponse(402, "token is out of date", nil))
			return
		}
		idx := strings.Index(c.Request.URL.Path, "/v1")
		if idx < 1 {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found", "data": "{}"})
			c.Abort()
			return
		}
		path := c.Request.URL.Path[:idx]
		method := c.Request.Method
		if ok, err := enf.Enforce(claim.RoleName, path, method); err != nil || !ok {
			if err != nil {
				logger.Printf("enforce error:%s", err.Error())
			}
			c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
			c.Abort()
			return
		}

		c.Set("userUuid", claim.UserUuid)
		c.Set("userName", claim.Username)
		c.Set("roleName", claim.RoleName)
		c.Next()
	}
}
