package auth

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const UserKey = "userContextKey"

type AuthMethod struct {
	Login  gin.HandlerFunc
	Logout gin.HandlerFunc

	Handlers []gin.HandlerFunc
}

type AuthUser struct {
	User     string
	password string
	Roles    []string
}

var user_roles = map[string]AuthUser{
	"foo":  {User: "foo", password: "bar", Roles: []string{"user", "admin"}},
	"manu": {User: "manu", password: "123", Roles: []string{"user"}},
}

func GetAuthentionApp(auth string) AuthMethod {
	if auth == "basic" { // suitable for API usage and single login
		return BasicAuth
	} else if auth == "mem" {
		return InMemoryAuth
	}
	log.Fatal("Unable to start:", auth)
	os.Exit(1)
	return InMemoryAuth
}
