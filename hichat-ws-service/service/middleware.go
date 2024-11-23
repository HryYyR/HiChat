package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	GroupScripts "go-websocket-server/ADB/MysqlScripts/GroupsScripts"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"go-websocket-server/Token_packge"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
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

func FlowControl(c *gin.Context) {
	userclaim, _ := c.Get("userdata")
	userdata := userclaim.(*models.UserClaim)
	rkey := fmt.Sprintf("FC%d", userdata.ID)

	ok := adb.Rediss.Exists(rkey).Val() //判断是否存在
	if ok > 0 {
		//存在就+1,判断是否触发限流
		adb.Rediss.Incr(rkey)
		Flow, err := adb.Rediss.Get(rkey).Int()
		if err != nil || Flow > config.FlowControlNum {
			adb.Rediss.Expire(rkey, config.FlowControlTime)
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

// DependencyInjection 依赖注入
func DependencyInjection(userRepository UsersScripts.UserRepository, groupRepository GroupScripts.GroupRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userRepository", userRepository)
		c.Set("groupRepository", groupRepository)
		c.Next()
	}
}
