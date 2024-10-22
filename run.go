package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"io/fs"
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

func main() {
	r = gin.Default()

	Template()
	Database()

	r.GET("/:any", func(c *gin.Context) {
		if strings.HasSuffix(c.Param("any"), ".html") {
			c.HTML(200, c.Param("any"), gin.H{})
			return
		}

		c.HTML(404, "404.html", gin.H{})
	})

	r.Run("0.0.0.0:80")
}

func Template() {
	var files []string
	filepath.Walk("./front", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
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
