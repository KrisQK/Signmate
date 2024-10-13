package main

import "github.com/gin-gonic/gin"

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("front/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	r.Run("0.0.0.0:80")
}
