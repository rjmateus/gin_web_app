package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func authenticate() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar   --> base64: Zm9vOmJhcg==
		"manu": "123", // user:manu password:123  --> base64: bWFudToxMjM=
	})
}

func loadUser(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if userAuth, ok := user_roles[user]; ok {
		c.Set(UserKey, userAuth)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

var BasicAuth = AuthMethod{
	Handlers: []gin.HandlerFunc{authenticate(), loadUser},
}

/* example curl for /admin with basicauth header
	curl -X GET
  	http://localhost:8080/admin \
  	-H 'authorization: Basic Zm9vOmJhcg=='
*/
