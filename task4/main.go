package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique; not null; type:varchar(20)"`
	Password string `gorm:"not null; type:varchar(255)"`
	Email    string `gorm:"unique; not null; type:varchar(20)"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null; type:varchar(20)"`
	Content string `gorm:"not null; type:varchar(255)"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null; type:varchar(255)"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

func main() {
	// mysql
	var dsn string = "root:a123456@tcp(127.0.0.1:3306)/renren_fast?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql 连接成功 !!!")

	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // 密码
		DB:       0,                // 数据库
		PoolSize: 20,               // 连接池大小
	})
	fmt.Println("redis 连接成功 !!!")

	//db.AutoMigrate(&User{}, &Post{}, &Comment{})

	r := gin.Default()

	// 登录
	r.POST("/login", func(c *gin.Context) {
		var param map[string]interface{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user User
		result := db.Debug().Raw(" SELECT * FROM users WHERE username = ? ", param["username"]).Scan(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}

		if user.ID == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "用户名不存在！",
				"data":    nil,
			})
			return
		}

		// 密码验证
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param["password"].(string))); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "密码错误！",
				"data":    nil,
			})
			return
		}

		// 生成 jwt token
		redisJwtValue, err := GetJwtToken(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// redis 存储 jwt token
		var redisJwtKey string = "redis_jwt_key:" + strconv.FormatUint(uint64(user.ID), 10) + ":" + user.Username
		err = rdb.Set(context.Background(), redisJwtKey, redisJwtValue, 7200*time.Second).Err()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "登录成功！",
			"data":    map[string]string{"token": redisJwtValue},
		})
		return
	})

	// 注册
	r.POST("/register", func(c *gin.Context) {
		var param map[string]interface{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 密码加密
		bytes, err := bcrypt.GenerateFromPassword([]byte(param["password"].(string)), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 注册
		var user User = User{}
		user.Username = param["username"].(string)
		user.Password = string(bytes)
		user.Email = param["email"].(string)
		result := db.Debug().Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}

		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "注册成功！",
			"data":    user.ID,
		})
		return
	})

	// 文章，创建
	r.PUT("post", JWTAuth(), func(c *gin.Context) {
		var param map[string]interface{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user User
		db.Debug().Raw(" SELECT * FROM users WHERE id = ? AND username = ? ", c.GetFloat64("id"), c.GetString("username")).Scan(&user)
		fmt.Println(user)

		var post Post = Post{}
		post.Title = param["title"].(string)
		post.Content = param["content"].(string)
		post.UserID = user.ID
		db.Debug().Create(&post)

		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "文章创建成功！",
			"data":    "",
		})
		return
	})

	// 文章，查询
	r.GET("post", JWTAuth(), func(c *gin.Context) {
		var post_id = c.Query("post_id")

		var post []Post
		if post_id != "" {
			db.Debug().Raw(" SELECT * FROM posts WHERE user_id = ? AND id = ? ", c.GetFloat64("id"), post_id).Scan(&post)
		} else {
			db.Debug().Raw(" SELECT * FROM posts WHERE user_id = ? ", c.GetFloat64("id")).Scan(&post)
		}
		fmt.Println(post)

		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "文章创建成功！",
			"data":    post,
		})
		return
	})

	// 文章，更新
	r.POST("post", JWTAuth(), func(c *gin.Context) {
		var param map[string]interface{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var post Post = Post{}
		db.Debug().Raw(" SELECT * FROM posts WHERE id = ? ", param["id"]).Scan(&post)
		if post.UserID != uint(c.GetFloat64("id")) {
			// response
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "只能更新自己的文章",
				"data":    nil,
			})
			return
		}

		db.Debug().Exec("UPDATE posts SET title = ?, content = ? WHERE id = ? ", param["title"], param["content"], param["id"])
		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "文章更新成功！",
			"data":    nil,
		})
		return
	})

	// 文章，删除
	r.DELETE("post", JWTAuth(), func(c *gin.Context) {
		var post_id = c.Query("post_id")

		var post Post = Post{}
		db.Debug().Raw(" SELECT * FROM posts WHERE id = ? ", post_id).Scan(&post)
		if post.UserID != uint(c.GetFloat64("id")) {
			// response
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "只能删除自己的文章",
				"data":    nil,
			})
			return
		}

		db.Debug().Delete(&post)
		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "文章删除成功！",
			"data":    nil,
		})
		return
	})

	// 评论，创建
	r.PUT("comment", JWTAuth(), func(c *gin.Context) {
		var param map[string]interface{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var comment Comment = Comment{}
		comment.Content = param["content"].(string)
		comment.UserID = uint(c.GetFloat64("id"))
		v, _ := strconv.ParseFloat(param["post_id"].(string), 64)
		comment.PostID = uint(v)
		db.Debug().Create(&comment)

		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "文章评论成功！",
			"data":    nil,
		})
		return
	})

	// 评论，查询
	r.GET("comment", JWTAuth(), func(c *gin.Context) {
		var post_id = c.Query("post_id")

		var comment []Comment
		db.Debug().Raw(" SELECT * FROM comments WHERE post_id = ? ", post_id).Scan(&comment)

		// response
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "评论查询成功！",
			"data":    comment,
		})
		return
	})

	r.Run(":8080")
}
