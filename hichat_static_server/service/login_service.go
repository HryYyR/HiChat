package service

import (
	"encoding/json"
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"hichat_static_server/rpcserver"
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
	json.Unmarshal(databyte, &data)
	fmt.Printf("%+v\n", data)

	if data.Username == "" || data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid username or password",
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
			"msg": "failed to query userdata  ",
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
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "password mismatch",
		})
		return
	}
	// 获取用户的群列表
	// usergrouplist, err := userdata.GetUserGroupList()
	// jsonData, _ := json.MarshalIndent(usergrouplist, "", "  ")
	// fmt.Println(string(jsonData))
	// fmt.Println("--------------------------------")
	rpcusergrouplist, _ := rpcserver.GetUserGroupList(userdata.ID)
	// rpcjsonData, _ := json.MarshalIndent(rpcusergrouplist, "", "  ")
	// fmt.Println(string(rpcjsonData))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "qurey grouplist failed",
		})
		return
	}
	token, err := util.GenerateToken(userdata.ID, userdata.UUID, userdata.UserName, 24*time.Hour)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Generate Token failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "login success",
		"token": token,
		"userdata": models.ResponseUserData{
			ID:          userdata.ID,
			UserName:    userdata.UserName,
			NikeName:    userdata.NikeName,
			Email:       userdata.Email,
			CreatedTime: userdata.CreatedAt,
			LoginTime:   userdata.LoginTime,
			GroupList:   rpcusergrouplist,
		},
	})

}
