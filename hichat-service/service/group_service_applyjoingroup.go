package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	GroupScripts "go-websocket-server/ADB/MysqlScripts/GroupsScripts"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
)

type applyjoingroupinfo struct {
	GroupName string `json:"GroupName"`
	GroupID   int    `json:"GroupID"`
	Msg       string `json:"Msg"`
	ApplyWay  int    `json:"ApplyWay"`
}

// ApplyJoinGroup 申请加入群聊
func ApplyJoinGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	var rawdata applyjoingroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		log.Println(err)
		util.H(c, http.StatusBadRequest, "非法访问", nil)
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		log.Println(err)
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	if len(rawdata.Msg) > 50 {
		util.H(c, http.StatusBadRequest, "申请理由超字数上限(50字)", nil)
		return
	}

	groupRepository := c.MustGet("groupRepository").(GroupScripts.GroupRepository)

	applycount, err := groupRepository.GetUserApplyJoinGroupCount(userdata.ID, 0)
	if err != nil {
		log.Println(err)
		util.H(c, http.StatusInternalServerError, "查询关系失败", nil)
		return
	}
	if applycount >= 5 {
		util.H(c, http.StatusBadRequest, "申请已达上限", nil)
		return
	}

	//log.Println(userdata.ID, rawdata.GroupID)
	exitgroup, err := groupRepository.CheckUserIsExistInGroup(userdata.ID, rawdata.GroupID)
	if err != nil {
		log.Println(err)
		util.H(c, http.StatusInternalServerError, "查询关系失败", nil)
		return
	}
	if exitgroup {
		util.H(c, http.StatusBadRequest, "你已在该群聊中", nil)
		return
	}

	//查申请
	applymsgdata := &models.ApplyJoinGroup{
		ApplyUserID: userdata.ID,
		GroupID:     rawdata.GroupID,
	}
	exist, tempapplydata, err2 := applymsgdata.CheckApplyExit()
	if err2 != nil {
		log.Println("check apply error", err)
		util.H(c, http.StatusInternalServerError, "查询申请失败", nil)
		return
	}
	//log.Println("applydata:", exist, tempapplydata.HandleStatus)
	if exist && tempapplydata.HandleStatus == 0 {
		util.H(c, http.StatusBadRequest, "你已经申请过了", nil)
		return
	}

	applygroupdata := &models.Group{
		ID: rawdata.GroupID,
	}
	groupinfo, exist, err := applygroupdata.CheckGroupExit()
	if err != nil {
		log.Println("check group error", err)
		util.H(c, http.StatusInternalServerError, "查询群聊失败", nil)
		return
	}
	if !exist {
		util.H(c, http.StatusBadRequest, "群聊不存在", nil)
		return
	}

	applydata := models.ApplyJoinGroup{
		ApplyUserID:   userdata.ID,
		ApplyUserName: userdata.UserName,
		GroupID:       rawdata.GroupID,
		ApplyMsg:      rawdata.Msg,
		ApplyWay:      rawdata.ApplyWay,
	}
	//log.Printf("%#v\n", applydata)
	err = applydata.InsertApply()
	if err != nil {
		log.Println("insert error", err)
		util.H(c, http.StatusInternalServerError, "申请失败", nil)
		return
	}

	msg := models.Message{
		UserID:   groupinfo.CreaterID, //此id为群主的id,所以消息应该是通知群主
		UserName: userdata.UserName,
		GroupID:  applydata.GroupID,
		MsgType:  config.MsgTypeRefreshGroupNotice,
	}
	msgbyte, _ := json.Marshal(msg)
	//向群主发送验证申请信息
	//models.ServiceCenter.Clients[applygroup.CreaterID].Send <- msgbyte
	models.TransmitMsg(msgbyte, config.MsgTypeRefreshGroupNotice)

	util.H(c, http.StatusOK, "申请成功", nil)
}
