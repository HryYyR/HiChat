package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	GroupScripts "go-websocket-server/ADB/MysqlScripts/GroupsScripts"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
)

type joingroupinfo struct {
	ApplyID      int `json:"ApplyID"`
	HandleStatus int `json:"HandleStatus"`
}

// HandleJoinGroup 处理加入群聊
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
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	groupRepository := c.MustGet("groupRepository").(GroupScripts.GroupRepository)
	// 查询申请是否存在
	tempapplyjoindata := &models.ApplyJoinGroup{
		ID: rawdata.ApplyID,
	}
	applyExit, applyjoindata, err := tempapplyjoindata.CheckApplyExit()
	if err != nil {
		log.Println(err)
		util.H(c, http.StatusInternalServerError, "查询申请失败", nil)
		return
	}
	if !applyExit {
		util.H(c, http.StatusBadRequest, "申请不存在", nil)
		return
	}
	if applyjoindata.HandleStatus != 0 {
		util.H(c, http.StatusBadRequest, "申请已处理", nil)
		return
	}

	// 查用户是否存在
	userdata := &models.Users{
		ID: applyjoindata.ApplyUserID,
	}

	applyuserdata, userExit, err := userdata.CheckUserExit()
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusInternalServerError, "查询用户信息失败", nil)
		return
	}
	if !userExit {
		util.H(c, http.StatusBadRequest, "用户不存在", nil)
		return
	}

	//查群是否存在
	tempgroupdata := &models.Group{
		ID: applyjoindata.GroupID,
	}
	grouplist, groupExit, err := tempgroupdata.CheckGroupExit()
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusInternalServerError, "查询群聊失败", nil)
		return
	}
	if !groupExit {
		util.H(c, http.StatusBadRequest, "群聊不存在", nil)
		return
	}

	// 拒绝申请
	if rawdata.HandleStatus == -1 {
		if _, err := groupRepository.UpdateApplyJoinGroupStatus(applyjoindata.ID, rawdata.HandleStatus); err != nil {
			log.Println(err)
			util.H(c, http.StatusInternalServerError, "拒绝失败", nil)
			return
		}
		//if _, err = adb.SqlStruct.Conn.Table("apply_join_group").ID(applyjoindata.ID).Update(models.ApplyJoinGroup{HandleStatus: rawdata.HandleStatus}); err != nil {
		//	log.Println(err)
		//	util.H(c, http.StatusInternalServerError, "拒绝失败", nil)
		//	return
		//}

		util.H(c, http.StatusOK, "拒绝成功", nil)

	} else if rawdata.HandleStatus == 1 { //同意申请
		// 同意申请
		session := adb.SqlStruct.Conn.NewSession()
		session.Begin()
		defer session.Close()

		addggur := models.GroupUserRelative{
			UserID:    applyjoindata.ApplyUserID,
			GroupID:   grouplist.ID,
			GroupUUID: grouplist.UUID,
		}
		err = addggur.Association(grouplist, session) //连接关系
		if err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "同意加入群聊失败", nil)
			return
		}
		// 修改申请状态
		if _, err := groupRepository.UpdateApplyJoinGroupStatus(applyjoindata.ID, rawdata.HandleStatus); err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "同意失败", nil)
			return
		}
		// 修改群聊总人数
		if _, err = session.Table("group").ID(applyjoindata.GroupID).Update(models.Group{MemberCount: grouplist.MemberCount + 1}); err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "更新失败", nil)
			return
		}

		groupmsg := models.Message{
			UserID:     applyuserdata.ID,
			UserName:   applyuserdata.UserName,
			UserAvatar: applyuserdata.Avatar,
			UserAge:    applyuserdata.Age,
			UserCity:   applyuserdata.City,
			GroupID:    grouplist.ID,
			Msg:        fmt.Sprintf("%s加入了群聊", applyuserdata.UserName),
			MsgType:    config.MsgTypeJoinGroup,
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

		msgbyte, _ := json.Marshal(groupmsg)
		adb.Rediss.RPush(fmt.Sprintf("gm%d", grouplist.ID), string(msgbyte))
		err := adb.Rediss.HIncrBy(fmt.Sprintf("group%d", grouplist.ID), "MemberCount", 1).Err() //redis成员人数+1
		if err != nil {
			log.Println(err)
		}
		// 通知群里的其他成员有用户加入
		models.TransmitMsg(msgbyte, config.MsgTypeJoinGroup)

		//通知申请用户刷新群聊列表
		selfmsg := models.Message{
			UserID:   applyjoindata.ApplyUserID,
			UserName: applyjoindata.ApplyUserName,
			GroupID:  grouplist.ID,
			MsgType:  config.MsgTypeRefreshGroupAndNotice,
		}
		selfmsgbyte, _ := json.Marshal(selfmsg)
		models.TransmitMsg(selfmsgbyte, config.MsgTypeRefreshGroupAndNotice)

		util.H(c, http.StatusOK, "用户已加入", nil)
		return
	}

	//通知申请人,申请已被处理(刷新通知列表)
	groupmsg := models.Message{
		UserID:   applyuserdata.ID,
		UserName: applyuserdata.UserName,
		GroupID:  grouplist.ID,
		MsgType:  config.MsgTypeRefreshGroupNotice,
	}
	bytemsg, err := json.Marshal(groupmsg)
	if err != nil {
		log.Println(err)
		util.H(c, http.StatusInternalServerError, "处理失败", nil)
		return
	}
	//models.ServiceCenter.Clients[applyuserdata.ID].Send <- bytemsg
	models.TransmitMsg(bytemsg, config.MsgTypeRefreshGroupNotice)

}
