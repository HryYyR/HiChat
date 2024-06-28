package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/rpcserver"
	"go-websocket-server/util"
	"net/http"
	"strconv"
	"strings"
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

	if data.ApplyUserID == data.PreApplyUserID {
		util.H(c, http.StatusBadRequest, "不能添加自己为好友", nil)
		return
	}

	if len(data.ApplyMsg) > 50 {
		util.H(c, http.StatusBadRequest, "申请理由超字数上限(50字)", nil)
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

	//fmt.Printf("%#v\n", data)
	if _, err = adb.SqlStruct.Conn.Table("apply_add_user").Insert(&data); err != nil {
		util.H(c, http.StatusInternalServerError, "申请添加好友失败", err)
		return
	}
	msg := models.UserMessage{
		UserID:        data.ApplyUserID,
		ReceiveUserID: data.PreApplyUserID,
		MsgType:       config.MsgTypeRefreshFriendNotice,
	}
	msgbyte, _ := json.Marshal(msg)
	models.TransmitMsg(msgbyte, config.MsgTypeRefreshFriendNotice)
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

	var receiveuserinfo models.Users
	has, err := adb.SqlStruct.Conn.Table("users").Where("id=?", data.Userid).Get(&receiveuserinfo)
	if !has {
		util.H(c, http.StatusBadRequest, "用户不存在", nil)
		return
	}
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	receiveuserstatus := adb.Rediss.HGet("UserClient", strconv.Itoa(receiveuserinfo.ID)).Val()
	if receiveuserstatus == "0" {
		util.H(c, http.StatusBadRequest, "对方不在线", nil)
		return
	} else {
		callmsg := models.UserMessage{
			UserID:          userdata.ID,
			UserName:        userdata.UserName,
			ReceiveUserID:   receiveuserinfo.ID,
			ReceiveUserName: receiveuserinfo.UserName,
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
		//models.ServiceCenter.Clients[receiveuserinfo.ID].Send <- callmsgbyte
		models.TransmitMsg(callmsgbyte, config.MsgTypeStartUserToUserVideoCall)

		fmt.Println("创建房间成功")

		util.H(c, http.StatusOK, "ok", nil)
	}
}

type DeleteUserReq struct {
	UserID int `json:"userid"`
}

// DeleteUser 删除好友
func DeleteUser(c *gin.Context) {
	info, _ := c.Get("userdata")
	userdata := info.(*models.UserClaim)

	req := new(DeleteUserReq)

	err := util.HandleJsonArgument(c, req)
	if err != nil {
		util.H(c, http.StatusBadRequest, "参数有误", err)
		return
	}

	userRepository := UsersScripts.NewUserRepository(adb.SqlStruct.Conn)

	//检查目标用户是否存在
	userIsExist := adb.Rediss.Exists(strconv.Itoa(req.UserID)).Val()
	if userIsExist <= 0 {
		userIsExist, _, err := userRepository.CheckUserIsExist(req.UserID)
		if err != nil {
			util.H(c, http.StatusInternalServerError, "查询用户状态失败", err)
			return
		}
		if !userIsExist {
			util.H(c, http.StatusBadRequest, "用户不存在", err)
			return
		}
	}
	//	检查是否为好友
	fstr := adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(userdata.ID)).Val()
	flist := strings.Split(fstr, ",")
	strtid := strconv.Itoa(req.UserID)
	sign := 0
	for i, fid := range flist {
		if fid == strtid {
			sign = 1
			flist = append(flist[:i], flist[i+1:]...)
			break
		}
	}
	if sign == 0 {
		util.H(c, http.StatusBadRequest, "该用户不是您的好友", err)
		return
	}

	//删除关系
	isdel, err := userRepository.DeleteFriendRelative(userdata.ID, req.UserID)
	if err != nil || !isdel {
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}

	//	删除双方redis好友映射表
	adb.Rediss.HSet("UserToUserRelative", strconv.Itoa(userdata.ID), strings.Join(flist, ","))

	fstr = adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(req.UserID)).Val()
	flist = strings.Split(fstr, ",")
	strtid = strconv.Itoa(userdata.ID)
	sign = 0
	for i, fid := range flist {
		if fid == strtid {
			sign = 1
			if len(flist) > 1 {
				flist = append(flist[:i], flist[i+1:]...)
			} else {
				flist = flist[:0]
			}
			break
		}
	}
	adb.Rediss.HSet("UserToUserRelative", strconv.Itoa(req.UserID), strings.Join(flist, ","))

	umsg := models.UserMessage{
		UserID:        userdata.ID,
		ReceiveUserID: req.UserID,
		MsgType:       config.MsgTypeRefreshFriend,
	}
	umsgbytes, _ := json.Marshal(umsg)
	models.TransmitMsg(umsgbytes, config.MsgTypeRefreshFriend)

	util.H(c, http.StatusOK, "删除好友成功", err)

}
