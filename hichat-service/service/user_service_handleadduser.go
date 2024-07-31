package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"net/http"
	"strconv"
)

type handleadduserinfo struct {
	ApplyID      int `json:"ApplyID"`
	HandleStatus int `json:"HandleStatus"`
}

// HandleAddUser 处理添加好友
func HandleAddUser(c *gin.Context) {
	data := new(handleadduserinfo)
	if err := util.HandleJsonArgument(c, data); err != nil {
		util.H(c, http.StatusBadRequest, "参数有误", err)
		return
	}
	var applyadduserdata models.ApplyAddUser
	exit, err := adb.SqlStruct.Conn.Table("apply_add_user").ID(data.ApplyID).Get(&applyadduserdata)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询申请失败", err)
		return
	}
	if !exit {
		util.H(c, http.StatusBadRequest, "申请不存在", nil)
		return
	}

	session := adb.SqlStruct.Conn.NewSession()
	// 更新申请
	if _, err = session.Table("apply_add_user").ID(data.ApplyID).Update(&models.ApplyAddUser{HandleStatus: data.HandleStatus}); err != nil {
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "更新申请失败", err)
		return
	}

	if data.HandleStatus == -1 {
		util.H(c, http.StatusOK, "拒绝成功", nil)
	} else {
		insertdata := models.UserUserRelative{
			PreUserID:    applyadduserdata.PreApplyUserID,
			PreUserName:  applyadduserdata.PreApplyUserName,
			BackUserID:   applyadduserdata.ApplyUserID,
			BackUserName: applyadduserdata.ApplyUserName,
		}
		// 插入关系
		if _, err = session.Table("user_user_relative").Insert(&insertdata); err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}

		RelatStr := strconv.Itoa(applyadduserdata.ApplyUserID)
		ulist := adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(applyadduserdata.PreApplyUserID)).Val()
		if len(ulist) != 0 {
			RelatStr = fmt.Sprintf("%s,%d", ulist, applyadduserdata.ApplyUserID)
		}
		adb.Rediss.HSet("UserToUserRelative", strconv.Itoa(applyadduserdata.PreApplyUserID), RelatStr)

		preRelatStr := strconv.Itoa(applyadduserdata.PreApplyUserID)
		preulist := adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(applyadduserdata.ApplyUserID)).Val()
		if len(preulist) != 0 {
			preRelatStr = fmt.Sprintf("%s,%d", preulist, applyadduserdata.ApplyUserID)
		}
		adb.Rediss.HSet("UserToUserRelative", strconv.Itoa(applyadduserdata.ApplyUserID), preRelatStr)

		session.Commit()

		util.H(c, http.StatusOK, "同意成功", nil)
	}

	//通知相关用户
	refreshmsg := models.UserMessage{
		UserID:        applyadduserdata.ApplyUserID,
		ReceiveUserID: applyadduserdata.PreApplyUserID,
		MsgType:       config.MsgTypeRefreshFriendAndNotice,
	}
	msgbyte, _ := json.Marshal(refreshmsg)
	models.TransmitMsg(msgbyte, config.MsgTypeRefreshFriendAndNotice)

}
