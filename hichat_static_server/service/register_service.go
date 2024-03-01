package service

import (
	"encoding/json"
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"hichat_static_server/util"
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
	databyte, _ := c.GetRawData()

	var data registerpostform
	err := json.Unmarshal(databyte, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
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
		Avatar:   "static/icon.png",
	}
	_, err = adb.Ssql.Table(&models.Users{}).Insert(&user)
	if err != nil {
		fmt.Println(err)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "注册失败!",
		})
		return
	}

	var UserData models.Users
	has, _ := adb.Ssql.Table("users").Where("user_name=? and email=?", data.Username, data.Email).Get(&UserData)
	if has {
		associateGroup := &models.GroupUserRelative{
			UserID:    UserData.ID,
			GroupID:   1,
			GroupUUID: "1",
		}
		_ = associateGroup.Association()
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "注册成功!",
	})
}
