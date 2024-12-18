package service

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/config"
	"hichat_static_server/models"
	"hichat_static_server/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

var NoVerificationList = []string{
	"login",
	"register",
	"emailcode",
}

func IdentityCheck(c *gin.Context) {
	for _, s := range NoVerificationList {
		if strings.Contains(c.Request.URL.Path, s) {
			c.Next()
			return
		}
	}

	token := c.GetHeader("Authorization")
	userclaim, err := util.DecryptToken(token)

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

func FlowControl(c *gin.Context) {
	for _, s := range NoVerificationList {
		if strings.Contains(c.Request.URL.Path, s) {
			c.Next()
			return
		}
	}

	userclaim, _ := c.Get("userdata")
	userdata := userclaim.(*models.UserClaim)
	rkey := fmt.Sprintf("FC%d", userdata.ID)

	ok := adb.Rediss.Exists(rkey).Val() //判断是否存在
	if ok > 0 {
		//存在就+1,判断是否触发限流
		adb.Rediss.Incr(rkey)
		Flow, err := adb.Rediss.Get(rkey).Int()
		if err != nil || Flow > config.FlowControlNum {
			util.H(c, http.StatusForbidden, "Access Denied", nil)
			c.Abort()
			return
		}
	} else {
		//不存在就创建并设置过期时间
		adb.Rediss.Incr(rkey)
		adb.Rediss.Expire(rkey, config.FlowControlTime)
	}

	c.Next()
}
