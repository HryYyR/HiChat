package service

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type registerpostform struct {
	Username string
	Password string
	Email    string
	Code     string
}

func Register(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")
	// email := c.PostForm("email")
	// fmt.Println(username, password, email)

	databyte, _ := c.GetRawData()

	var data registerpostform
	json.Unmarshal(databyte, &data)
	fmt.Printf("%+v\n", data)

	if data.Username == "" || data.Password == "" || !util.EmailValid(data.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid username or password or email",
		})
		return
	}
	emailcode := adb.Rediss.Get(data.Email).Val()
	if emailcode != data.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效验证码!",
		})
		return
	}

	// 查询邮箱是否已被注册
	hasemail, err := adb.Ssql.Table("users").Where("email = ? or user_name=? ", data.Email, data.Username).Exist()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询邮箱失败!",
		})
		return
	}
	if hasemail {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "邮箱或者用户名已被注册!",
		})
		return
	}

	salt := fmt.Sprintln(time.Now().Unix())
	encodepwd := util.Md5(data.Password + salt)
	user := &models.Users{
		UUID:     util.GenerateUUID(),
		UserName: data.Username,
		NikeName: data.Username,
		Password: encodepwd,
		Email:    data.Email,
		Salt:     salt,
	}
	_, err = adb.Ssql.Table(&models.Users{}).InsertOne(&user)
	if err != nil {
		fmt.Println(err)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "register failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "registered successfully",
	})
}

// /user/RefreshGroupList
// 获取指定id用户的数据(刷新数据)
func RefreshGroupList(c *gin.Context) {
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "getrawdata failed",
		})
		return
	}
	var user models.Users
	err = json.Unmarshal(rawbyte, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "json Unmarshal failed",
		})
		return
	}

	usergrouplist, err := user.GetUserGroupList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "get usergrouplist failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":           "query userdata success",
		"usergrouplist": usergrouplist,
	})
}

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
	err = util.MailSendCode(emaildata.Email, code)       //发送验证码
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "发送验证码失败,请稍后再试!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送验证码成功!",
	})

}
