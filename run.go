package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed front/*
var HtmlDir embed.FS

//go:embed front/assets/**/*
var AssetDir embed.FS

func main() {

	r := gin.Default()

	// embed files
	tmpl := template.New("")
	tmpl = template.Must(tmpl.ParseFS(HtmlDir, "front/*.html"))
	r.SetHTMLTemplate(tmpl)

	assfs, _ := fs.Sub(AssetDir, "front/assets")
	r.StaticFS("assets", http.FS(assfs))

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	r.GET("/about.html", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{})
	})

	r.GET("/service.html", func(c *gin.Context) {
		c.HTML(200, "service.html", gin.H{})
	})

	r.GET("/review.html", func(c *gin.Context) {
		c.HTML(200, "review.html", gin.H{})
	})
	r.GET("/project.html", func(c *gin.Context) {
		c.HTML(200, "project.html", gin.H{})
	})
	r.GET("/contact.html", func(c *gin.Context) {
		c.HTML(200, "contact.html", gin.H{})
	})
	r.GET("/blog.html", func(c *gin.Context) {
		c.HTML(200, "blog.html", gin.H{})
	})
	r.GET("/faq.html", func(c *gin.Context) {
		c.HTML(200, "faq.html", gin.H{})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(200, "404.html", gin.H{})
	})

	r.Run("0.0.0.0:80")
}
