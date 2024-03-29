package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Content-Security-Policy", "default-src 'self'; connect-src 'self' http://localhost:3004")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func IdentityCheck(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userclaim, err := Token_packge.DecryptToken(token)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid Authorization",
		})
		c.Abort()
		return
	}
	c.Set("userdata", userclaim)
	c.Next()
}
