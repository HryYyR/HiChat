package service

import (
	"encoding/json"
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type emailcode struct {
	Email string
}

func EmailCode(c *gin.Context) {
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "getrawdata failed",
		})
		return
	}

	var emaildata emailcode
	err = json.Unmarshal(rawbyte, &emaildata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "json Unmarshal failed",
		})
		return
	}
	// 查询邮箱是否已被注册
	hasemail, err := adb.Ssql.Table("users").Where("email = ?  ", emaildata.Email).Exist()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询邮箱失败!",
		})
		return
	}
	if hasemail {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "邮箱已被注册!",
		})
		return
	}

	mail := adb.Rediss.Get(emaildata.Email).Val()
	if len(mail) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请勿重复发送!",
		})
		return
	}
	code := util.RandCode()
	adb.Rediss.Set(emaildata.Email, code, 1*time.Minute) //验证码存redis
	go util.MailSendCode(emaildata.Email, code)          //发送验证码
	// if err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"msg": "发送验证码失败,请稍后再试!",
	// 	})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送验证码成功!",
	})

}
