package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"strconv"
	"time"
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

	userRepository := c.MustGet("userRepository").(UsersScripts.UserRepository)
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
	session.Begin()
	defer session.Close()
	// 更新申请
	if _, err = session.Table("apply_add_user").ID(data.ApplyID).Update(&models.ApplyAddUser{HandleStatus: data.HandleStatus}); err != nil {
		log.Println(err)
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "更新申请失败", err)
		return
	}

	if data.HandleStatus == -1 {
		util.H(c, http.StatusOK, "拒绝成功", nil)
	} else {
		insertdata := models.UserUserRelative{
			PreUserID:    applyadduserdata.ApplyUserID,
			PreUserName:  applyadduserdata.ApplyUserName,
			BackUserID:   applyadduserdata.PreApplyUserID,
			BackUserName: applyadduserdata.PreApplyUserName,
		}

		// 插入关系
		ok, err := userRepository.ConnectFriendRelative(&insertdata, session)
		if err != nil || !ok {
			log.Println(err)
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}

		//即将被添加到映射关系的id
		RelatStr := strconv.Itoa(applyadduserdata.ApplyUserID)
		preRelatStr := strconv.Itoa(applyadduserdata.PreApplyUserID)

		//preapply:被申请 , apply:主动申请
		//获取UserToUserRelative映射关系
		ulist := adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(applyadduserdata.PreApplyUserID)).Val()
		if len(ulist) != 0 {
			RelatStr = fmt.Sprintf("%s,%d", ulist, applyadduserdata.ApplyUserID)
		}
		preulist := adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(applyadduserdata.ApplyUserID)).Val()
		if len(preulist) != 0 {
			preRelatStr = fmt.Sprintf("%s,%d", preulist, applyadduserdata.PreApplyUserID)
		}

		//redis 事务
		redisSession := adb.Rediss.Pipeline()
		defer redisSession.Close()

		//修改UserToUserRelative映射关系
		err = redisSession.HSet("UserToUserRelative", strconv.Itoa(applyadduserdata.PreApplyUserID), RelatStr).Err()
		if err != nil {
			log.Println(err)
			redisSession.Discard()
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}

		err = redisSession.HSet("UserToUserRelative", strconv.Itoa(applyadduserdata.ApplyUserID), preRelatStr).Err()
		if err != nil {
			log.Println(err)
			redisSession.Discard()
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}

		//查询uuid
		srcuserdata, exist, err := userRepository.GetUserByUserID(applyadduserdata.ApplyUserID)
		if err != nil || !exist {
			log.Println(err)
			redisSession.Discard()
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}
		dstuserdata, exist, err := userRepository.GetUserByUserID(applyadduserdata.PreApplyUserID)
		if err != nil || !exist {
			log.Println(err)
			redisSession.Discard()
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}

		//nebula用户与用户关系联系
		err = adb.NebulaInstance.InsertEdge("UserAddUser", srcuserdata.UUID, dstuserdata.UUID,
			[]string{"start_id", "end_id", "created_at"}, []any{srcuserdata.ID, dstuserdata.ID, time.Now()})
		if err != nil {
			log.Println(err)
			redisSession.Discard()
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "处理好友请求失败", err)
			return
		}

		redisSession.Exec()
		util.H(c, http.StatusOK, "同意成功", nil)

	}
	session.Commit()

	//通知相关用户
	refreshmsg := models.UserMessage{
		UserID:        applyadduserdata.ApplyUserID,
		ReceiveUserID: applyadduserdata.PreApplyUserID,
		MsgType:       config.MsgTypeRefreshFriendAndNotice,
	}
	msgbyte, _ := json.Marshal(refreshmsg)
	models.TransmitMsg(msgbyte, config.MsgTypeRefreshFriendAndNotice)

}
