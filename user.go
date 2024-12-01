package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Auth")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		fmt.Println(token)

		var user User
		res := db.First(&user, "token = ?", token)
		if res.Error != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		c.Set("User", user)
		c.Next()
	}
}
