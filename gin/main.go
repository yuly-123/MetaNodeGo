package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	user := r.Group("/user", gin.BasicAuth(gin.Accounts{
		"yuly": "a123456",
	}))

	user.GET("/name", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		fmt.Println(user)
		c.String(http.StatusOK, user)
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}

	//r.RunTLS("8080", "./cert/tls.crt", "./cert/tls.key")
}
