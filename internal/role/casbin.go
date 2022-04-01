package role

import (
	casbin "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"micro-shop/internal/util"
	"net/http"
	"os"
)

type Cb struct {
	enf    *casbin.Enforcer
	logger *log.Logger
}

func NewCb(enf *casbin.Enforcer, logger *log.Logger) *Cb {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("service/casbin")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	return &Cb{
		enf:    enf,
		logger: logger,
	}
}

func (c *Cb) Enforce(role, path, act string) bool {
	var ok bool
	var err error
	if ok, err = c.enf.Enforce(role, path, act); err != nil {
		c.logger.Printf("failed to enforce:%s", err.Error())
		return false
	}
	return ok
}

func (c *Cb) ListPolicies() [][]string {
	return c.enf.GetPolicy()
}

func (c *Cb) UpdatePolicy(old []string, new []string) bool {
	if len(old) != 3 || len(new) != 3 {
		return false
	}
	ok, err := c.enf.UpdatePolicy(old, new)
	if err != nil {
		c.logger.Printf("update policy error:%s", err.Error())
		return false
	}
	return ok
}

func (cb *Cb) ListRole(c *gin.Context) {
	polices := cb.enf.GetPolicy()
	var roleMap = make(map[string]bool)
	var roleList = make([]string, 0, 3)
	for _, item := range polices {
		if len(item) != 3 {
			continue
		}
		if _, ok := roleMap[item[0]]; !ok {
			roleMap[item[0]] = true
			roleList = append(roleList, item[0])
		}
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", roleList))
}
