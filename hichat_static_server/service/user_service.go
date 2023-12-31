package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"hichat_static_server/util"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type edituserdataform struct {
	City      string `json:"city"`
	Age       string `json:"age"`
	Introduce string `json:"introduce"`
}

func EditUserData(c *gin.Context) {
	var mlock sync.Mutex
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	databyte, _ := c.GetRawData()
	var data edituserdataform
	err := json.Unmarshal(databyte, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "信息输入有误!",
		})
		return
	}
	// fmt.Printf("%+v\n", data)
	age, err := strconv.Atoi(data.Age)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "信息输入有误!",
		})
		return
	}

	if age < 0 || age > 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "信息输入有误!",
		})
		return
	}

	mlock.Lock()
	defer mlock.Unlock()

	if _, err := adb.Ssql.Table("users").Where("id=?", userdata.ID).Update(&models.Users{
		City:      data.City,
		Age:       age,
		Introduce: data.Introduce,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "修改失败!",
		})
		return
	}
	err = adb.Rediss.Del(strconv.Itoa(userdata.ID)).Err()
	if err != nil {
		log.Println(err.Error())
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

// GetUserGroupList 获取用户的群聊列表
func GetUserGroupList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}
	grouplist := make([]models.GroupDetail, 0)
	err = data.GetUserGroupList(&grouplist)
	// fmt.Println("消息长度为:", len(grouplist[0].MessageList))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取用户的群聊列表!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的群聊列表!",
		"data": grouplist,
	})
}

// GetUserFriendList 获取用户的好友列表
func GetUserFriendList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	friendlist := make([]models.FriendResponse, 0)
	if err = data.GetFriendListAndMEssage(&friendlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取用户的好友列表失败!",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的好友列表成功!",
		"data": friendlist,
	})
}

// GetUserApplyJoinGroupList //获取用户的群聊通知列表
func GetUserApplyJoinGroupList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	applyjoingrouplist := make([]models.ApplyJoinGroupResponse, 0)
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

// GetUserApplyAddFriendList //获取用户的好友申请列表
func GetUserApplyAddFriendList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式不正确!",
		})
		return
	}

	applyaddfriendlist := make([]models.ApplyAddUser, 0)
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

type searchfriendingo struct {
	Searchstr string
}

// SearchUser 搜索用户
func SearchUser(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	uudata := models.Users{
		ID:       userdata.ID,
		UserName: userdata.UserName,
	}

	var uufriendlist []models.Friend //发起人的好友列表
	if err := uudata.GetFriendList(&uufriendlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询好友列表失败!",
		})
		return
	}
	uufriendmap := make(map[int]int, 0)
	for _, f := range uufriendlist {
		uufriendmap[int(f.Id)] = int(f.Id)
	}

	var data searchfriendingo
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = json.Unmarshal(rawbyte, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if len(data.Searchstr) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "搜索成功!",
			"data": []models.Users{},
		})
		return
	}

	friendlist := make([]models.Users, 0)
	err = adb.Ssql.Table("users").Omit("ip,password,salt,grade,uuid").Where("user_name LIKE ?  and user_name !=?",
		"%"+data.Searchstr+"%", userdata.UserName).Find(&friendlist)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "搜索失败!",
		})
		return
	}

	resultfriendlist := make([]models.Users, 0)
	// 筛除已经成为的好友
	for _, ff := range friendlist {
		if _, ok := uufriendmap[ff.ID]; !ok {
			resultfriendlist = append(resultfriendlist, ff)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "搜索成功!",
		"data": resultfriendlist,
	})
}
