package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userkey = "user"

var secret = []byte("secret")

// login is a handler that parses a form and checks for specific data
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}
	userAuth, ok := user_roles[username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
	if userAuth.password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// AuthRequired is a simple middleware to check the session
func AuthRequiredSession(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user nill"})
		return
	}
	userText := fmt.Sprintf("%s", user)
	fmt.Println(userText)
	if userAuth, ok := user_roles[userText]; ok {
		c.Set(UserKey, userAuth)
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user not fould in map"})
	}
	// Continue down the chain to handler etc
	c.Next()
}

var InMemoryAuth = AuthMethod{
	Login:    login,
	Logout:   logout,
	Handlers: []gin.HandlerFunc{AuthRequiredSession},
}
