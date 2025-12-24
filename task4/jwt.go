package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

const JwtHmacKey string = "ujlZosSLQrexDjtZevQRJNtmjV2kjwGKQ+OyiUeCf/A="

// 创建 token
func GetJwtToken(user User) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Second * 7200).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwtHmacKey))
}

// jwt验证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token 的三段格式验证：Header.Payload.Signature
		// Header：包含算法类型和token类型
		// Payload：包含声明（claims）数据
		// Signature：签名部分
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		jwtBody := strings.Split(tokenString, ".")
		header, _ := base64.URLEncoding.DecodeString(jwtBody[0])
		payload, _ := base64.URLEncoding.DecodeString(jwtBody[1])
		signature := jwtBody[2]
		fmt.Println("header: ", string(header))   // {"alg":"HS256","typ":"JWT"}
		fmt.Println("payload: ", string(payload)) // {"exp":1766558655,"id":4,"username":"zhangshan"}
		fmt.Println("signature: ", signature)     // GVemASPYhswkMeMfierHJcEXuPyDyyiRtZRboz2Ytwc

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(JwtHmacKey), nil
		})

		if err != nil {
			// response
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok == false || token.Valid == false {
			// response
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "invalid token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		if ok && token.Valid {
			c.Set("id", claims["id"])
			c.Set("username", claims["username"])
			c.Next()
		}
	}
}
