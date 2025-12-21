package main

import (
	"gorm.io/gorm"
)

// 进阶gorm
// 题目1：模型定义
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
//
// 题目2：关联查询
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
//
// 题目3：钩子函数
// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func Run_advanced(db *gorm.DB) {
	//db.Debug().AutoMigrate(&User{}, &Post{}, &Comment{})
	// CREATE TABLE `users` (`id` bigint AUTO_INCREMENT,`name` varchar(20),PRIMARY KEY (`id`))
	// CREATE TABLE `posts` (`id` bigint AUTO_INCREMENT,`content` varchar(255),`user_id` bigint,PRIMARY KEY (`id`),CONSTRAINT `fk_posts_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`))
	// CREATE TABLE `comments` (`id` bigint AUTO_INCREMENT,`chat` varchar(255),`post_id` bigint,PRIMARY KEY (`id`),CONSTRAINT `fk_comments_post` FOREIGN KEY (`post_id`) REFERENCES `posts`(`id`))

	//var user1 User = User{Name: "张三"}
	//var user2 User = User{Name: "李四"}
	//db.Debug().Create(&user1)
	//db.Debug().Create(&user2)
	//fmt.Println(user1)
	//fmt.Println(user2)
	//
	//var post1 Post = Post{Content: "c++从入门到放弃", UserId: user1.ID}
	//var post2 Post = Post{Content: "linux之成哥的私房菜", UserId: user1.ID}
	//var post3 Post = Post{Content: "java编程思想中的梦想", UserId: user2.ID}
	//var post4 Post = Post{Content: "asp.net基础到实战应用", UserId: user2.ID}
	//db.Debug().Create(&post1)
	//db.Debug().Create(&post2)
	//db.Debug().Create(&post3)
	//db.Debug().Create(&post4)
	//fmt.Println(post1)
	//fmt.Println(post2)
	//fmt.Println(post3)
	//fmt.Println(post4)
	//
	//var comment1 Comment = Comment{Chat: "说的好好", PostId: post1.ID}
	//var comment2 Comment = Comment{Chat: "说的对对对", PostId: post1.ID}
	//var comment3 Comment = Comment{Chat: "说的妙", PostId: post4.ID}
	//var comment4 Comment = Comment{Chat: "加油", PostId: post4.ID}
	//var comment5 Comment = Comment{Chat: "看这里么的", PostId: post4.ID}
	//var comment6 Comment = Comment{Chat: "哼哈嘿", PostId: post4.ID}
	//db.Debug().Create(&comment1)
	//db.Debug().Create(&comment2)
	//db.Debug().Create(&comment3)
	//db.Debug().Create(&comment4)
	//db.Debug().Create(&comment5)
	//db.Debug().Create(&comment6)
	//fmt.Println(comment1)
	//fmt.Println(comment2)
	//fmt.Println(comment3)
	//fmt.Println(comment4)
	//fmt.Println(comment5)
	//fmt.Println(comment6)

	//var user User = User{ID: 2}
	//db.Debug().Preload("Posts.Comments").First(&user)
	//fmt.Println("user:", user)

	//var sql string = " SELECT c.post_id, count(0) AS cc, p.content " +
	//	" FROM comments AS c " +
	//	" LEFT JOIN posts AS p ON p.id=c.post_id " +
	//	" GROUP BY c.post_id " +
	//	" ORDER BY cc DESC " +
	//	" LIMIT 1 "
	//var m map[string]interface{}
	//db.Debug().Raw(sql).Scan(&m)
	//fmt.Println("文章编号:", m["post_id"], ", 评论数量:", m["cc"], ", 文章:", m["content"])

	//var p Post = Post{Content: "docker框架深入剖析", UserId: 1}
	//db.Debug().Create(&p)

	var c Comment = Comment{ID: 1}
	db.Debug().First(&c)
	db.Debug().Delete(&c)
}

func (p *Post) AfterCreate(db *gorm.DB) (err error) {
	var u User = User{ID: p.UserId}
	db.Debug().First(&u, p.UserId)
	u.PostCount++
	db.Debug().Save(&u)
	return
}

func (c *Comment) AfterCreate(db *gorm.DB) (err error) {
	var p Post = Post{ID: c.PostId}
	db.Debug().First(&p, c.PostId)
	p.CommentCount++
	p.CommentStatus = "有评论"
	db.Debug().Save(&p)
	return
}

func (c *Comment) AfterDelete(db *gorm.DB) (err error) {
	var p Post = Post{ID: c.PostId}
	db.Debug().First(&p, c.PostId)
	if p.CommentCount == 1 {
		p.CommentCount = 0
		p.CommentStatus = "无评论"
	} else {
		p.CommentCount--
	}
	db.Debug().Save(&p)
	return
}

type User struct {
	ID        int64
	Name      string `gorm:"type:varchar(20)"`
	PostCount int32
	Posts     []Post
}
type Post struct {
	ID            int64
	Content       string `gorm:"type:varchar(255)"`
	CommentCount  int32
	CommentStatus string `gorm:"type:varchar(20);default:'无评论'"`
	Comments      []Comment
	UserId        int64
	User          User
}
type Comment struct {
	ID     int64
	Chat   string `gorm:"type:varchar(255)"`
	PostId int64
	Post   Post
}
