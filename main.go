package main

import (
	"fmt"
	"gin_web_app/core/auth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func withRole(role string) func(*gin.Context) {

	return func(c *gin.Context) {
		user := c.MustGet(auth.UserKey).(auth.AuthUser)
		allow := false
		for _, val := range user.Roles {
			if val == role {
				allow = true
			}
		}
		if !allow {
			// Abort the request with the appropriate error code
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized role"})
			return
		}
		// Continue down the chain to handler etc
		c.Next()
	}
}

var secretMain = []byte("secret")

func setupRouterMain() *gin.Engine {
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore(secretMain)))

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong\n")
	})

	auth := auth.GetAuthentionApp("mem")
	if auth.Login != nil {
		r.POST("/login", auth.Login)

	}
	if auth.Logout != nil {
		r.GET("/logout", auth.Logout)
	}

	authorized := r.Group("/")
	authorized.Use(auth.Handlers...)

	// Get user value
	authorized.GET("all", all)
	authorized.GET("user", withRole("user"), user)
	authorized.GET("admin", withRole("admin"), admin)

	return r
}

func user(c *gin.Context) {
	user := c.MustGet(auth.UserKey).(auth.AuthUser)
	c.String(http.StatusOK, fmt.Sprintf("user: %s\n", user))
}
func admin(c *gin.Context) {
	user := c.MustGet(auth.UserKey).(auth.AuthUser)
	c.String(http.StatusOK, fmt.Sprintf("admin: %s\n", user))
}

func all(c *gin.Context) {
	user := c.MustGet(auth.UserKey).(auth.AuthUser)
	c.String(http.StatusOK, fmt.Sprintf("all: %s\n", user))
}

func main() {
	r := setupRouterMain()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
