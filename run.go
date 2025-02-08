package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"os/exec"
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

type Imonial struct {
	ID   uint `gorm:"primarykey"`
	User string
	Word string
}

type KV struct {
	Key   string
	Value string
}

var ViewTime int

func ViewTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ViewTime += 1

		// 处理请求
		c.Next()
	}
}

func main() {

	go func() {
		fmt.Println("启动GoFile")
		err := exec.Command("./goFile.exe", "-path", "./front/assets/gallery", "-port", "8888").Run()
		if err != nil {
			fmt.Println(err)
		}
	}()

	r = gin.Default()
	r.Use(ViewTimeMiddleware())

	Template()
	Database()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/index.html")
	})

	r.GET("/index.html", func(c *gin.Context) {

		var imonials []Imonial
		db.Limit(6).Order("RANDOM()").Find(&imonials)

		c.HTML(200, "index.html", gin.H{
			"imonials": imonials,
		})
	})

	r.GET("/about.html", func(c *gin.Context) {
		var imonials []Imonial
		db.Limit(6).Order("RANDOM()").Find(&imonials)

		c.HTML(200, "about.html", gin.H{
			"imonials": imonials,
		})
	})

	r.GET("/services.html", func(c *gin.Context) {
		var imonials []Imonial
		db.Limit(6).Order("RANDOM()").Find(&imonials)

		c.HTML(200, "services.html", gin.H{
			"imonials": imonials,
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

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "admin/login.html", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user User
		db.First(&user, "username = ?", username)
		if user.ID == 0 {
			c.HTML(400, "admin/login.html", gin.H{"error": "用户不存在！"})
		} else if password == user.Password {
			token := uuid.NewString()

			res := db.Model(user).Update("token", token)

			if res.Error == nil {
				c.SetCookie("Auth", token, 3600, "/", "", false, false)
				c.HTML(200, "admin/login.html", gin.H{"success": "登录成功！"})
			} else {
				c.HTML(500, "admin/login.html", gin.H{"error": res.Error.Error()})
			}
		} else {
			c.HTML(400, "admin/login.html", gin.H{"error": "用户名或密码错误"})
		}
	})

	adminGroup := r.Group("/admin").Use(AuthMid())

	adminGroup.GET("/", func(c *gin.Context) {
		c.HTML(200, "admin/home.html", gin.H{
			"ViewTime": ViewTime,
		})
	})

	adminGroup.GET("/imonial", func(c *gin.Context) {
		c.HTML(200, "admin/imonial.html", gin.H{})
	})

	adminGroup.GET("api/imonial", func(c *gin.Context) {
		var imonials []Imonial

		db.Find(&imonials)
		c.JSON(200, gin.H{
			"code":  0,
			"msg":   "",
			"count": len(imonials),
			"data":  imonials,
		})
	})

	adminGroup.GET("api/imonial/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		db.Delete(&Imonial{}, id)

		c.JSON(200, gin.H{"msg": "删除成功"})
	})

	adminGroup.POST("api/imonial/add", func(c *gin.Context) {

		user := c.PostForm("user")
		word := c.PostForm("word")
		println(user, word)

		if user == "" || word == "" {
			c.JSON(200, gin.H{"msg": "参数错误"})
			return
		}

		db.Create(&Imonial{User: user, Word: word})

		c.JSON(200, gin.H{"msg": "添加成功"})
	})

	adminGroup.GET("/gallery", func(c *gin.Context) {
		c.String(200, "admin.callery")
	})

	adminGroup.GET("/api/gallery", func(c *gin.Context) {
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

		c.JSON(200, gin.H{
			"category": category,
			"images":   images,
		})
	})

	adminGroup.POST("/api/gallery/category", func(c *gin.Context) {
		name := c.Param("name")
		if strings.Contains(name, ".") {
			c.JSON(200, gin.H{"msg": "no way"})
			return
		}
		err := os.MkdirAll("./front/assets/gallery/"+name, 0750)
		c.JSON(200, gin.H{"msg": "ok", "error": err.Error()})
	})

	adminGroup.POST("/delete/gallery/category", func(c *gin.Context) {
		name := c.Param("name")
		if strings.Contains(name, ".") {
			c.JSON(200, gin.H{"msg": "no way"})
			return
		}
		err := os.RemoveAll("./front/assets/gallery/" + name)
		c.JSON(200, gin.H{"msg": "ok", "error": err.Error()})
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

	err = db.AutoMigrate(&Imonial{})
	if err != nil {
		panic("Failed to migrate database")
	}

	if res := db.First(&Imonial{}); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		db.Create(&Imonial{User: "Bill Danson", Word: "Jackie has been absolutely great. Very accommodating and a pleasure to deal with. Will be back here the next time for sure. I want to have another car wrapping for my new Ford. Great job guys!"})
		db.Create(&Imonial{User: "Blank Richard", Word: "Thanks Jackie for the awesome LED lightbox. Your great communication and attention to service are truely hard to come by these days. My partner really happy with that. Thanks sooo much I will definately be comming back for other signages and car graphics."})
		db.Create(&Imonial{User: "Shannon Wood", Word: "Exceptional service. Provided great help and assistance in getting sizes correct for three of my vehicles. Also ensured a date deadline was met and provided excellent communication. Outstanding designer in every respect. Thank you."})
		db.Create(&Imonial{User: "Carl Madden", Word: "I was very pleased with the level of professional service I received from signmate limited, the quality of the graphics for my logo was creative and perfect for my company in attracting new customers."})
		db.Create(&Imonial{User: "Tomma", Word: "Excellent Service! Definitely Highly Recommended... They did amzing vehicle graphic on my Holden. These guys are really frienly and professional. I had a good service, thank you Signmate!"})
		db.Create(&Imonial{User: "Mr Dan", Word: "Signmate Limited helped me created a ideal image for my company store, through this change I have attracted new customers and in turn created more profit."})
		fmt.Println("Created imonial")
	}

	err = db.AutoMigrate(&KV{})
	if err != nil {
		panic("Failed to migrate database KV")
	}

	if res := db.First(&KV{Key: "Guard"}); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		db.Create(&KV{Key: "Guard", Value: ""})
		fmt.Println("Created 内嵌代码")
	}
}
