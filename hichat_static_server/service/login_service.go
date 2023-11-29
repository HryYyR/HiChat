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

type loginpostform struct {
	Username string
	Password string
}

func Login(c *gin.Context) {
	databyte, _ := c.GetRawData()

	var data loginpostform
	err := json.Unmarshal(databyte, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "内容格式有误,请检查后重试!",
		})
		return
	}
	//fmt.Printf("%+v\n", data)

	if data.Username == "" || data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "内容格式有误,请检查后重试!",
		})
		return
	}

	// 查询用户是否存在
	var userdata models.Users
	hasuser, err := adb.Ssql.Table(&models.Users{}).Where("user_name = ?", data.Username).Get(&userdata)
	if err != nil {
		log.Println(err)
		// fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "failed to query userdata",
		})
		return
	}
	if !hasuser {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "user not exit",
		})
		return
	}
	if util.Md5(data.Password+userdata.Salt) != userdata.Password {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "password mismatch",
		})
		return
	}
	token, err := util.GenerateToken(userdata.ID, userdata.UUID, userdata.UserName, 24*time.Hour)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Generate Token failed",
		})
		return
	}

	// 获取用户的信息列表
	ResponseUserData := new(models.ResponseUserData)
	err = userdata.Login(ResponseUserData)
	if err != nil {
		fmt.Println(err)
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"err": "登陆失败!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":      "登陆成功!",
		"token":    token,
		"userdata": ResponseUserData,
	})

}

func Test(c *gin.Context) {
	userdata := models.Users{
		ID: 1008,
	}

	friendmsg := new(models.ResponseUserData)
	err := userdata.Login(friendmsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": friendmsg,
	})

}
