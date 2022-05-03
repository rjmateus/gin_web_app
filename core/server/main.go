package server

import (
	"gin_web_app/core/auth"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	var secret = []byte("secret")
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore(secret)))
	
	auth := auth.InMemoryAuth
	r.POST("/login", auth.Login)
	r.GET("/logout", auth.Logout)

	return r
}
