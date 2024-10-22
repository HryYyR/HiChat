package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	GroupScripts "go-websocket-server/ADB/MysqlScripts/GroupsScripts"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
)

type groupinfo struct {
	Groupname string
	Avatar    string `json:"Avatar"`
}

// GroupDetail 用于返回结果的结构体
type GroupDetail struct {
	GroupInfo   models.Group
	MessageList []models.Message
}

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	info, _ := c.Get("userdata")
	userdata := info.(*models.UserClaim)

	rawsbyte, _ := c.GetRawData()
	var rowdata groupinfo
	err := json.Unmarshal(rawsbyte, &rowdata)
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusBadRequest, "解析参数失败", err)
		return
	}

	groupRepository := c.MustGet("groupRepository").(GroupScripts.GroupRepository)

	UUID := util.GenerateUUID()
	var group = models.Group{
		UUID:        UUID,
		CreaterID:   userdata.ID,
		Avatar:      rowdata.Avatar,
		CreaterName: userdata.UserName,
		GroupName:   rowdata.Groupname,
		MemberCount: 1,
	}

	//判断群名称是否占用,是则禁止创建
	isexit, err := groupRepository.ByGroupNameCheckGroupIsExist(rowdata.Groupname)
	//isexit, err := adb.SqlStruct.Conn.Table("group").Where("group_name = ?", rowdata.Groupname).Exist()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "发生了未知的错误", nil)
		fmt.Println(err)
		log.Println(err)
		return
	}
	if isexit {
		util.H(c, http.StatusBadRequest, "群聊名称已被使用", nil)
		return
	}

	GroupLock.Lock()
	defer GroupLock.Unlock() //解锁

	session := adb.SqlStruct.Conn.NewSession()
	defer session.Close()
	session.Begin()

	var groupdata models.Group

	fullgroupdata, err := group.InsertGroup(session)
	if err != nil {
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "群聊创建失败", nil)
		fmt.Println(err)
		return
	}
	groupdata = fullgroupdata

	//连接关系
	var gur = models.GroupUserRelative{
		UserID:    userdata.ID,
		GroupUUID: UUID,
		GroupID:   groupdata.ID,
	}
	err = gur.Association(groupdata, session)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "群聊创建失败", nil)
		log.Println(err)
		fmt.Println(err)
		session.Rollback()
		return
	}

	session.Commit()

	responsedata := GroupDetail{
		GroupInfo:   groupdata,
		MessageList: []models.Message{},
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建成功",
		"data": responsedata,
	})
}
