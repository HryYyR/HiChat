package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/rpcserver"
	"go-websocket-server/util"
	"net/http"
)

// ApplyAddUser 申请添加好友
func ApplyAddUser(c *gin.Context) {
	rawbyte, err := c.GetRawData()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "获取数据失败", err)
		return
	}
	var data models.ApplyAddUser
	err = json.Unmarshal(rawbyte, &data)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "数据解析失败", err)
		return
	}

	// fmt.Printf("%+v", data)
	exit, err := adb.SqlStruct.Conn.Table("apply_add_user").
		Where("pre_apply_user_id=?  and apply_user_id=?  and handle_status=0",
			data.PreApplyUserID, data.ApplyUserID).Exist()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询申请信息失败", err)
		return
	}
	if exit {
		util.H(c, http.StatusOK, "申请已存在", nil)
		return
	}
	exit, err = adb.SqlStruct.Conn.Table("apply_add_user").
		Where("pre_apply_user_id=?  and apply_user_id=?  and handle_status=0",
			data.ApplyUserID, data.PreApplyUserID).Exist()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询申请信息失败", err)
		return
	}
	if exit {
		util.H(c, http.StatusOK, "申请已存在", nil)
		return
	}

	if _, err = adb.SqlStruct.Conn.Table("apply_add_user").Insert(&data); err != nil {
		util.H(c, http.StatusInternalServerError, "申请添加好友失败", err)
		return
	}
	msg := models.Message{
		UserID:   data.ApplyUserID,
		UserName: data.ApplyUserName,
		GroupID:  0,
		MsgType:  config.MsgTypeRefreshFriendNotice,
	}
	msgbyte, _ := json.Marshal(msg)
	models.ServiceCenter.Clients[data.PreApplyUserID].Send <- msgbyte
	models.ServiceCenter.Clients[data.ApplyUserID].Send <- msgbyte
	util.H(c, http.StatusOK, "申请成功", nil)

}

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
		insertdata := &models.UserUserRelative{
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
		session.Commit()

		util.H(c, http.StatusOK, "同意成功", nil)
	}

	//通知相关用户
	refreshmsg := models.Message{
		UserID:   applyadduserdata.ApplyUserID,
		UserName: applyadduserdata.ApplyUserName,
		GroupID:  0,
		MsgType:  config.MsgTypeRefreshFriend,
	}
	msgbyte, _ := json.Marshal(refreshmsg)
	models.ServiceCenter.Clients[applyadduserdata.PreApplyUserID].Send <- msgbyte
	models.ServiceCenter.Clients[applyadduserdata.ApplyUserID].Send <- msgbyte

}

type CheckLoginStatusReq struct {
	Userid int `json:"userid"`
}

// StartUserToUserVideoCall 发起用户与用户的远程视频
func StartUserToUserVideoCall(c *gin.Context) {
	info, _ := c.Get("userdata")
	userdata := info.(*models.UserClaim)

	data := new(CheckLoginStatusReq)
	if err := util.HandleJsonArgument(c, data); err != nil {
		util.H(c, http.StatusBadRequest, "参数有误", err)
		return
	}

	var userinfo models.Users
	has, err := adb.SqlStruct.Conn.Table("users").Where("id=?", data.Userid).Get(&userinfo)
	if !has {
		util.H(c, http.StatusBadRequest, "用户不存在", nil)
		return
	}
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	if models.ServiceCenter.Clients[userinfo.ID].Status == false {
		util.H(c, http.StatusBadRequest, "对方不在线", nil)
		return
	} else {
		callmsg := models.UserMessage{
			UserID:          userdata.ID,
			UserName:        userdata.UserName,
			ReceiveUserID:   userinfo.ID,
			ReceiveUserName: userinfo.UserName,
			Msg:             fmt.Sprintf("%s发起了视频通话", userdata.UserName),
			MsgType:         config.MsgTypeStartUserToUserVideoCall,
		}

		fmt.Println("开始调用远程,创建房间")
		_, err := rpcserver.CallNoticeVideoStreamServer(callmsg)
		if err != nil {
			util.H(c, http.StatusInternalServerError, "创建房间失败", err)
			return
		}

		callmsgbyte, _ := json.Marshal(callmsg)
		models.ServiceCenter.Clients[userinfo.ID].Send <- callmsgbyte

		fmt.Println("创建房间成功")

		util.H(c, http.StatusOK, "ok", nil)

	}

}
