// main.go (auth-server)

package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// This map will store the username/token key value pairs
var users = make(map[string]string)

// We are using a list of predefined users. In a production app,
// users will most likely be authenticated directly against a database
var seedUsers = []user{
	user{
		Username: "user1",
		Password: "pass1",
	},
	user{
		Username: "user2",
		Password: "pass2",
	},
	user{
		Username: "user3",
		Password: "pass3",
	},
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	s := gin.Default()

	s.POST("/login", login)
	s.POST("/authenticate", authenticate)
	s.POST("/logout", logout)

	s.Run(":8001")
}

func login(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	if token := validateUser(username, password); token == "" {
		// Respond with an HTTP error if authentication fails
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		// If authentication succeeds, add the username and its token
		// to the users map for later reference
		// Respond with an HTTP success status and include the token
		// in the response
		users[username] = token
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func authenticate(c *gin.Context) {
	// Obtain the POSTed username and token values
	username := c.PostForm("username")
	token := c.PostForm("token")

	if v, ok := users[username]; ok && v == token {
		// If the username/token pair is found in the users map,
		// respond with an HTTP success status
		c.JSON(http.StatusOK, nil)
	} else {
		// If the username/token pair is not found in the users map,
		// respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func logout(c *gin.Context) {
	// Obtain the POSTed username and token values
	username := c.PostForm("username")
	token := c.PostForm("token")

	if v, ok := users[username]; ok && v == token {
		// If the username/token pair is found in the users map,
		// remove this username from the users map
		// and respond with an HTTP success status
		delete(users, username)
		c.JSON(http.StatusOK, nil)
	} else {
		// If the username/token pair is not found in the users map,
		// respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func generateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}

func validateUser(username, password string) string {
	// Validate the username/password against the seed values defined above
	// In a production app,
	// users will most likely be authenticated directly against a database
	for _, u := range seedUsers {
		if username == u.Username {
			if u.Password == password {
				return generateSessionToken()
			}
			return ""
		}
	}
	return ""
}
