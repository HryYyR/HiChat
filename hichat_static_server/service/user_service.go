package service

import (
	"encoding/json"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"hichat_static_server/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type edituserdataform struct {
	City string `json:"city"`
	Age  string `json:"age"`
}

func EditUserData(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	databyte, _ := c.GetRawData()
	var data edituserdataform
	json.Unmarshal(databyte, &data)
	// fmt.Printf("%+v\n", data)
	age, err := strconv.Atoi(data.Age)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "信息输入有误!",
		})
		return
	}

	if _, err := adb.Ssql.Table("users").Where("id=?", userdata.ID).Update(&models.Users{
		City: data.City,
		Age:  age,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "修改失败!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功!",
	})

}

func GetUserData(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	var targetuserdata models.Users
	exit, err := adb.Ssql.Omit("Password,Salt,Grade,IP").Table("users").ID(data.ID).Get(&targetuserdata)
	if !exit {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户不存在!",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询用户信息失败!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": targetuserdata,
	})
}

// 获取用户的群聊列表
func GetUserGroupList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	grouplist, err := data.GetUserGroupList()
	// fmt.Println("消息长度为:", len(grouplist[0].MessageList))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取用户的群聊列表!",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的群聊列表!",
		"data": grouplist,
	})
}

// 获取用户的好友列表
func GetUserFriendList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	var friendlist []models.Friend
	if err = data.GetFriendList(&friendlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取用户的好友列表失败!",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的好友列表成功!",
		"data": friendlist,
	})
}

// //获取用户的群聊通知列表
func GetUserApplyJoinGroupList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	var applyjoingrouplist []models.ApplyJoinGroup
	if err = data.GetApplyMsgList(&applyjoingrouplist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取用户的群聊通知列表失败!",
		})
	}

	util.GroupTimeSort(applyjoingrouplist, "desc")

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的群聊通知列表成功!",
		"data": applyjoingrouplist,
	})
}

// //获取用户的好友申请列表
func GetUserApplyAddFriendList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	var applyaddfriendlist []models.ApplyAddUser
	if err = data.GetApplyAddUserList(&applyaddfriendlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取用户的好友申请列表失败!",
		})
	}
	util.UserTimeSort(applyaddfriendlist, "desc")

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的好友申请列表成功!",
		"data": applyaddfriendlist,
	})
}
