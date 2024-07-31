package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"go-websocket-server/config"
	"go-websocket-server/models"
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

	userRepository := c.MustGet("userRepository").(UsersScripts.UserRepository)

	if data.ApplyUserID == data.PreApplyUserID {
		util.H(c, http.StatusBadRequest, "不能添加自己为好友", nil)
		return
	}

	if len(data.ApplyMsg) > 50 {
		util.H(c, http.StatusBadRequest, "申请理由超字数上限(50字)", nil)
		return
	}

	//检查目标用户是否为好友
	_, exist, err := userRepository.CheckUserIsFriend(data.ApplyUserID, data.PreApplyUserID)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询申请信息失败", err)
		return
	}
	if exist {
		util.H(c, http.StatusBadRequest, "不能添加好友", nil)
		return
	}

	// fmt.Printf("%+v", data)
	exist, err = adb.SqlStruct.Conn.Table("apply_add_user").
		Where("pre_apply_user_id=?  and apply_user_id=?  and handle_status=0",
			data.PreApplyUserID, data.ApplyUserID).Exist()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询申请信息失败", err)
		return
	}
	if exist {
		util.H(c, http.StatusOK, "申请已存在", nil)
		return
	}

	exist, err = adb.SqlStruct.Conn.Table("apply_add_user").
		Where("pre_apply_user_id=?  and apply_user_id=?  and handle_status=0",
			data.ApplyUserID, data.PreApplyUserID).Exist()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询申请信息失败", err)
		return
	}
	if exist {
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
