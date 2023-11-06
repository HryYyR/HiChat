package service

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
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
		MemberCount: 1,
	}

	GroupLock.Lock()
	defer GroupLock.Unlock() //解锁
	session := adb.Ssql.NewSession()
	//判断群是否已存在,存在就禁止创建
	isexit, err := adb.Ssql.Table("group").Where("group_name = ?", rowdata.Groupname).Exist()
	if err != nil {
		fmt.Println(err)
	}
	if isexit {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "群聊名称已被使用!",
		})
		return
	}

	_, err = session.Table("group").Insert(&group) //插入群聊
	if err != nil {
		session.Rollback()
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "插入群聊失败!",
		})
		return
	}

	var groupdata models.Group
	_, err = session.Table("group").Where("uuid=?", UUID).Get(&groupdata) // 查群聊完整信息
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询群聊失败",
		})
		return
	}

	var gur = models.GroupUserRelative{
		UserID:    userdata.ID,
		GroupUUID: UUID,
		GroupID:   groupdata.ID,
	}
	err = gur.Association(groupdata) //连接关系
	if err != nil {
		session.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "创建群聊失败",
		})
		return
	}
	session.Commit()

	responsedata := models.GroupDetail{
		GroupInfo:   groupdata,
		MessageList: []models.GroupMessage{},
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建成功",
		"data": responsedata,
	})
}

type applyjoingroupinfo struct {
	GroupName string `json:"GroupName"`
	GroupID   int    `json:"GroupID"`
	Msg       string `json:"Msg"`
	ApplyWay  int    `json:"ApplyWay"`
}

// 申请加入群聊
func ApplyJoinGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	var rawdata applyjoingroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "格式化json失败!",
		})
		return
	}
	fmt.Println(userdata.ID, rawdata.GroupID)
	exitgroup, err := adb.Ssql.Table("group_user_relative").Where("user_id = ? and group_id=?", userdata.ID, rawdata.GroupID).Exist()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询关系失败!",
		})
		return
	}
	if exitgroup {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "你已在该群聊中!",
		})
		return
	}

	var applymsgdata models.ApplyJoinGroup
	exitapply, err := adb.Ssql.Table("apply_join_group").Where("apply_user_id = ? and group_id=?", userdata.ID, rawdata.GroupID).OrderBy("created_at DESC").Get(&applymsgdata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询关系失败!",
		})
		return
	}
	if exitapply && applymsgdata.HandleStatus == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "你已经申请过了!",
		})
		return
	}

	var applygroup models.Group
	exitgroup, err = adb.Ssql.Table("group").Where("id = ?", rawdata.GroupID).Get(&applygroup)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询群聊失败!",
		})
		return
	}
	if !exitgroup {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "群聊不存在!",
		})
		return
	}

	applydata := models.ApplyJoinGroup{
		ApplyUserID:   userdata.ID,
		ApplyUserName: userdata.UserName,
		GroupID:       rawdata.GroupID,
		ApplyMsg:      rawdata.Msg,
		ApplyWay:      rawdata.ApplyWay,
	}
	if _, err = adb.Ssql.Table("apply_join_group").Insert(&applydata); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "申请失败!",
		})
		return
	}

	msg := models.Message{
		UserID:   userdata.ID,
		UserName: userdata.UserName,
		GroupID:  applydata.GroupID,
		MsgType:  config.MsgTypeJoinGroup,
	}
	msgbyte, _ := json.Marshal(msg)
	//向群主发送验证申请信息
	models.ServiceCenter.Clients[applygroup.CreaterID].Send <- msgbyte

	c.JSON(http.StatusOK, gin.H{
		"msg": "申请成功!",
	})
}

type joingroupinfo struct {
	ApplyID      int `json:"ApplyID"`
	HandleStatus int `json:"HandleStatus"`
}

// 处理加入群聊
func HandleJoinGroup(c *gin.Context) {
	// ud, _ := c.Get("userdata")
	// userdata := ud.(*models.UserClaim)

	var rawdata joingroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Incorrect format",
		})
		return
	}

	// 查询申请是否存在
	var applyjoindata models.ApplyJoinGroup
	has, err := adb.Ssql.Table("apply_join_group").Where("id = ?", rawdata.ApplyID).Get(&applyjoindata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询申请失败!",
		})
		return
	}
	if !has {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "申请不存在!",
		})
		return
	}

	//查群是否存在
	var grouplist models.Group
	has, err = adb.Ssql.Table("group").Where("id=?", applyjoindata.GroupID).Get(&grouplist)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询群聊失败!",
		})
		return
	}

	if !has {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "群聊不存在!",
		})
		return
	}

	// 查用户是否存在
	var applyuserdata models.Users
	has, err = adb.Ssql.Table("users").Where("id=?", applyjoindata.ApplyUserID).Get(&applyuserdata)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询用户信息失败!",
		})
		return
	}
	if !has {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户不存在!",
		})
		return
	}

	// 拒绝申请
	if rawdata.HandleStatus == -1 {
		if _, err = adb.Ssql.Table("apply_join_group").ID(applyjoindata.ID).Update(models.ApplyJoinGroup{HandleStatus: rawdata.HandleStatus}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "拒绝失败!",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "拒绝成功!",
		})
		return
	}

	// 同意申请
	session := adb.Ssql.NewSession()
	defer session.Close()
	session.Begin()
	addggur := models.GroupUserRelative{
		UserID:    applyjoindata.ApplyUserID,
		GroupID:   grouplist.ID,
		GroupUUID: grouplist.UUID,
	}
	err = addggur.Association(grouplist) //连接关系
	if err != nil {
		session.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "同意加入群聊失败!",
		})
		return
	}
	// 修改申请状态
	if _, err = session.Table("apply_join_group").ID(applyjoindata.ID).Update(models.ApplyJoinGroup{HandleStatus: rawdata.HandleStatus}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "同意失败!",
		})
		return
	}
	groupmsg := models.GroupMessage{
		UserID:   applyuserdata.ID,
		UserName: applyuserdata.UserName,
		GroupID:  grouplist.ID,
		Msg:      fmt.Sprintf("%s加入了群聊", applyuserdata.UserName),
		MsgType:  config.MsgTypeJoinGroup,
	}
	_, err = session.Table("group_message").Insert(&groupmsg)
	if err != nil {
		session.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	session.Commit()

	groupuserlist := models.GroupUserList[grouplist]
	// 通知群里的其他成员有用户加入
	msg := models.Message{
		UserID:   applyuserdata.ID,
		UserName: applyuserdata.UserName,
		GroupID:  grouplist.ID,
		MsgType:  config.MsgTypeJoinGroup,
	}
	msgbyte, _ := json.Marshal(msg)
	for _, userid := range groupuserlist {
		models.ServiceCenter.Clients[userid].Send <- msgbyte
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "用户已加入!",
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
		fmt.Println("json")
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
	defer session.Close()
	session.Begin()
	GroupLock.Lock()
	if groupinfo.CreaterID == userdata.ID { //说明他是群主,删除所有联系
		var willbedeleteuseridlist []int //将被删除的用户群聊关系的用户id
		err := session.Table("group_user_relative").Cols("user_id").Where("group_id = ?", groupinfo.ID).Find(&willbedeleteuseridlist)
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		_, err = session.Table("group_user_relative").Where("group_id = ?", groupinfo.ID).Delete()
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		_, err = session.Table("group").Where("id = ?", groupinfo.ID).Delete() //删群
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// for _, userclient := range models.ServiceCenter.Clients { //同步用户列表里相关信息
		// 	delete(userclient.Groups, groupinfo.ID)
		// }

		for _, userid := range willbedeleteuseridlist {
			delete(models.ServiceCenter.Clients[userid].Groups, groupinfo.ID)
		}
		delete(models.GroupUserList, groupinfo) //同步群列表里相关信息
	} else { //只删除该用户对群的联系
		_, err := session.Table("group_user_relative").Where("user_id = ? and group_id=?", userdata.ID, rawdata.ID).Delete()
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// 保存退出消息
		groupmsg := models.GroupMessage{
			UserID:   userdata.ID,
			UserName: userdata.UserName,
			GroupID:  groupinfo.ID,
			Msg:      fmt.Sprintf("%s退出了群聊", userdata.UserName),
			MsgType:  config.MsgTypeQuitGroup,
		}
		_, err = session.Table("group_message").Insert(&groupmsg)
		if err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		delete(models.ServiceCenter.Clients[userdata.ID].Groups, groupinfo.ID) //同步用户列表里相关信息
		for index, userid := range models.GroupUserList[groupinfo] {
			if userid == userdata.ID {
				models.GroupUserList[groupinfo] = append(models.GroupUserList[groupinfo][:index], models.GroupUserList[groupinfo][index+1:]...)
			}
		}
	}
	session.Commit()
	GroupLock.Unlock()

	// 通知群里的其他成员有用户退出
	msg := models.Message{
		MsgType:  config.MsgTypeRefreshGroup,
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
	var rawdata searchgroupinfo
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
	grouplist := make([]models.Group, 0)
	err = adb.Ssql.Table("group").Where("group_name LIKE ? and creater_id !=?", "%"+rawdata.Searchstr+"%", userdata.ID).Find(&grouplist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var responsedata []models.Group
	for _, group := range grouplist {
		group.MemberCount = group.GetMemberCount()
		responsedata = append(responsedata, group)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":       "search success",
		"grouplist": responsedata,
	})
}
