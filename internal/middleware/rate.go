package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
)

type Rule struct {
	Key   string
	Limit float64
	Per   int
}

func Limiter(rl *RouterLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		idx := strings.Index(c.Request.URL.Path, "/v1")
		if idx < 1 {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
			c.Abort()
			return
		}
		if limiter, ok := rl.limiters[c.Request.URL.Path[:idx]]; ok {
			if limiter.Allow() {
				c.Next()
			}
			c.JSON(http.StatusOK, gin.H{"code": 502, "msg": "busy", "data": ""})
			c.Abort()
			return
		}
		c.Next()
	}
}

func NewLimiter(r *Rule) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(r.Limit), r.Per)
	return limiter
}

type RouterLimiter struct {
	limiters map[string]*rate.Limiter
}

func NewRouterLimiter() *RouterLimiter {
	return &RouterLimiter{
		limiters: map[string]*rate.Limiter{},
	}
}

func (rl *RouterLimiter) AddLimiter(rules ...Rule) {
	for _, item := range rules {
		if _, ok := rl.limiters[item.Key]; !ok {
			rl.limiters[item.Key] = NewLimiter(&item)
		}
	}
}
