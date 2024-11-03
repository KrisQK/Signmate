package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var r *gin.Engine
var db *gorm.DB

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(64);not null"`
	Password string `gorm:"type:varchar(64);not null"`
	Token    string `gorm:"type:varchar(64)"`
}

type image struct {
	Url      string
	Category string
}

var testimonial = map[string]string{
	"Bill Danson":   "Jackie has been absolutely great. Very accommodating and a pleasure to deal with. Will be back here the next time for sure. I want to have another car wrapping for my new Ford. Great job guys!",
	"Blank Richard": "Thanks Jackie for the awesome LED lightbox. Your great communication and attention to service are truely hard to come by these days. My partner really happy with that. Thanks sooo much I will definately be comming back for other signages and car graphics.",
	"Shannon Wood":  "Exceptional service. Provided great help and assistance in getting sizes correct for three of my vehicles. Also ensured a date deadline was met and provided excellent communication. Outstanding designer in every respect. Thank you.",
	"Carl Madden":   "I was very pleased with the level of professional service I received from signmate limited, the quality of the graphics for my logo was creative and perfect for my company in attracting new customers. ",
	"Tomma":         "Excellent Service! Definitely Highly Recommended... They did amzing vehicle graphic on my Holden. These guys are really frienly and professional. I had a good service, thank you Signmate!",
	"Mr Dan":        " Signmate Limited helped me created a ideal image for my company store, through this change I have attracted new customers and in turn created more profit.",
}

func main() {
	r = gin.Default()

	Template()
	Database()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/index.html")
	})

	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"testimonial": testimonial,
		})
	})

	r.GET("/about.html", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"testimonial": testimonial,
		})
	})

	r.GET("/services.html", func(c *gin.Context) {
		c.HTML(200, "services.html", gin.H{
			"testimonial": testimonial,
		})
	})

	r.GET("/services-car.html", func(c *gin.Context) {
		c.HTML(200, "services-car.html", gin.H{})
	})

	r.GET("/services-shop.html", func(c *gin.Context) {
		c.HTML(200, "services-shop.html", gin.H{})
	})

	r.GET("/services-signage.html", func(c *gin.Context) {
		c.HTML(200, "services-signage.html", gin.H{})
	})

	r.GET("/services-printing.html", func(c *gin.Context) {
		c.HTML(200, "services-printing.html", gin.H{})
	})

	r.GET("/gallery.html", func(c *gin.Context) {

		category := make([]string, 0)
		var images []image

		dirs, err := os.ReadDir("./front/assets/gallery")
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range dirs {
			if file.IsDir() {
				category = append(category, file.Name())
			}
		}

		for _, c := range category {
			entry, err := os.ReadDir(filepath.Join("./front/assets/gallery", c))
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, e := range entry {
				if !e.IsDir() {

					images = append(images, image{
						Url:      "assets/gallery/" + c + "/" + e.Name(),
						Category: c,
					})
				}
			}
		}

		c.HTML(200, "gallery.html", gin.H{
			"category": category,
			"images":   images,
		})
	})

	r.GET("/contact.html", func(c *gin.Context) {
		c.HTML(200, "contact.html", gin.H{})
	})

	r.POST("/api/contact", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		subject := c.PostForm("subject")
		message := c.PostForm("message")
		fmt.Println(name, email, phone, subject, message)
		c.String(200, "Submit Success! We will contact u soon!")
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	r.Run("0.0.0.0:80")
}

func Template() {
	var files []string
	filepath.Walk("front", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "html") {
			files = append(files, path)
		}
		return nil
	})

	r.LoadHTMLFiles(files...)

	r.Static("/assets", "./front/assets")
}

func Database() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to migrate database")
	}

	if res := db.First(&User{}); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		db.Create(&User{Username: "admin", Password: "123456"})
		fmt.Println("Created admin account")
	}
}
