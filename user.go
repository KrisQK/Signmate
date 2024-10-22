package main

//func AuthMid() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token, err := c.Cookie("Auth")
//		if err != nil {
//			c.HTML(401, "admin/login.html", gin.H{"error": err.Error()})
//			return
//		}
//
//		fmt.Println(token)
//
//		var user User
//		res := db.First(&user, "token = ?", token)
//		if res.Error != nil {
//			c.HTML(401, "admin/login.html", gin.H{"error": err.Error()})
//			return
//		}
//
//		c.Set("User", user)
//		c.Next()
//	}
//}
//
//func init() {
//
//	r.GET("/login", func(c *gin.Context) {
//		c.HTML(200, "admin/login.html", gin.H{})
//	})
//
//	r.POST("/login", func(c *gin.Context) {
//		username := c.PostForm("username")
//		password := c.PostForm("password")
//
//		var user User
//		db.First(&user, "username = ?", username)
//		if user.ID == 0 {
//			c.HTML(400, "admin/login.html", gin.H{"error": "用户不存在！"})
//		} else if password == user.Password {
//			token := uuid.NewString()
//
//			res := db.Model(user).Update("token", token)
//
//			if res.Error == nil {
//				c.SetCookie("Auth", token, 3600, "/", "", false, false)
//				c.HTML(200, "admin/login.html", gin.H{"success": "登录成功！"})
//			} else {
//				c.HTML(500, "admin/login.html", gin.H{"error": res.Error.Error()})
//			}
//		} else {
//			c.HTML(400, "admin/login.html", gin.H{"error": "用户名或密码错误"})
//		}
//	})
//}
