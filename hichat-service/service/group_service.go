package service

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var GroupLock sync.Mutex

type groupinfo struct {
	Groupname string
	Avatar    string `json:"Avatar"`

	// Createrid   int
	// Creatername string
}

// 创建群聊
func CreateGroup(c *gin.Context) {
	info, _ := c.Get("userdata")
	userdata := info.(*models.UserClaim)

	rawsbyte, _ := c.GetRawData()
	var rowdata groupinfo
	err := json.Unmarshal(rawsbyte, &rowdata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "解析参数失败!",
		})
		return
	}

	UUID := util.GenerateUUID()
	var group = models.Group{
		UUID:        UUID,
		CreaterID:   userdata.ID,
		Avatar:      rowdata.Avatar,
		CreaterName: userdata.UserName,
		GroupName:   rowdata.Groupname,
	}

	GroupLock.Lock()
	defer GroupLock.Unlock() //解锁
	//判断群是否已存在,存在就禁止创建
	isexit, err := adb.Ssql.Table("group").Where("group_name = ?", rowdata.Groupname).Exist()
	if err != nil {
		fmt.Println(err)
	}
	if isexit {
		c.JSON(http.StatusOK, gin.H{
			"msg": "群聊名称已被使用!",
		})
		return
	}

	_, err = adb.Ssql.Table("group").Insert(&group) //插入群聊
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "insert failed",
		})
		return
	}

	var gid models.Group
	_, err = adb.Ssql.Table("group").Where("uuid=?", UUID).Get(&gid) // 查群聊id
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "qurey group failed",
		})
		return
	}

	var gur = models.GroupUserRelative{
		UserID:    userdata.ID,
		GroupUUID: UUID,
		GroupID:   gid.ID,
	}
	err = gur.Association(userdata, gid) //连接关系
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "created group failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "create success",
	})
}

type joingroupinfo struct {
	GroupName string `json:"GroupName"`
	// GroupId   string `json:"GroupId"`
}

// 加入群聊
func JoinGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	// fmt.Printf("%+v\n", ud)

	var rowdata joingroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(rawbyte, &rowdata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Incorrect format",
		})
		return
	}

	//查群是否存在
	var grouplist models.Group
	has, err := adb.Ssql.Table("group").Where("group_name=?", rowdata.GroupName).Get(&grouplist)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "server Error",
		})
		return
	}

	if !has {
		c.JSON(http.StatusOK, gin.H{
			"msg": "group not exists",
		})
		return
	}

	addggur := models.GroupUserRelative{
		UserID:    userdata.ID,
		GroupID:   grouplist.ID,
		GroupUUID: grouplist.UUID,
	}
	err = addggur.Association(userdata, grouplist) //连接关系
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "join group failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "加入成功!",
	})
}

// 退出群
func ExitGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	var rawdata models.Group //只有group_id
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var groupinfo models.Group
	has, err := adb.Ssql.Table("group").Where("id = ?", rawdata.ID).Get(&groupinfo)
	if !has {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Group not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 相关删除操作
	groupuserlist := models.GroupUserList[groupinfo]
	session := adb.Ssql.NewSession()
	GroupLock.Lock()
	if groupinfo.CreaterID == userdata.ID {
		_, err := adb.Ssql.Table("group_user_relative").Where("group_id = ?", groupinfo.ID).Delete() //说明他是群主,删除所有联系
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		_, err = adb.Ssql.Table("group").Where("id = ?", groupinfo.ID).Delete() //删群
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		for _, userclient := range models.ServiceCenter.Clients { //同步用户列表里相关信息
			delete(userclient.Groups, groupinfo.ID)
		}
		delete(models.GroupUserList, groupinfo) //同步群列表里相关信息
	} else { //只删除该用户对群的联系
		_, err := adb.Ssql.Table("group_user_relative").Where("user_id = ? and group_id=?", userdata.ID, rawdata.ID).Delete()
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		delete(models.ServiceCenter.Clients[userdata.ID].Groups, groupinfo.ID) //同步用户列表里相关信息
		// for group, useridlist := range models.GroupUserList {
		// 	var newlist []int
		// 	for _, userid := range useridlist {
		// 		if userid != userdata.ID {
		// 			newlist = append(newlist, userid)
		// 		}
		// 	}
		// 	models.GroupUserList[group] = newlist //同步群列表里相关信息
		// }
		for index, userid := range models.GroupUserList[groupinfo] {
			if userid == userdata.ID {
				models.GroupUserList[groupinfo] = append(models.GroupUserList[groupinfo][:index], models.GroupUserList[groupinfo][index+1:]...)
			}
		}
	}
	GroupLock.Unlock()

	msg := models.Message{
		MsgType:  201,
		UserID:   userdata.ID,
		UserName: userdata.UserName,
		GroupID:  rawdata.ID,
	}
	msgbyte, _ := json.Marshal(msg)
	for _, userid := range groupuserlist {
		models.ServiceCenter.Clients[userid].Send <- msgbyte
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "退出群聊成功!",
	})

}

type searchgroupinfo struct {
	Searchstr string
}

// 搜索群聊
func SearchGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)
	fmt.Println("ok")
	var rawdata searchgroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	fmt.Println("ok")
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	fmt.Println("ok")

	var grouplist []models.Group
	err = adb.Ssql.Table("group").Where("group_name = ? and creater_id !=?", rawdata.Searchstr, userdata.ID).Find(&grouplist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	fmt.Println("ok")

	c.JSON(http.StatusOK, gin.H{
		"msg":       "search success",
		"grouplist": grouplist,
	})
}
