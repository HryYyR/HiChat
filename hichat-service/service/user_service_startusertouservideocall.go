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
	"log"
	"net/http"
	"strconv"
)

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

	userRepository := c.MustGet("userRepository").(UsersScripts.UserRepository)

	receiveuserinfo, has, err := userRepository.GetUserByUserID(data.Userid)
	//has, err := adb.SqlStruct.Conn.Table("users").Where("id=?", data.Userid).Get(&receiveuserinfo)
	if !has {
		util.H(c, http.StatusBadRequest, "用户不存在", nil)
		return
	}
	if err != nil {
		log.Println(err)
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
