package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	adb "go-websocket-server/ADB"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DeleteUserReq struct {
	UserID int `json:"userid"`
}

// DeleteUser 删除好友
func DeleteUser(c *gin.Context) {
	info, _ := c.Get("userdata")

	userRepository := c.MustGet("userRepository").(UsersScripts.UserRepository)

	userdata := info.(*models.UserClaim)

	req := new(DeleteUserReq)

	err := util.HandleJsonArgument(c, req)
	if err != nil {
		util.H(c, http.StatusBadRequest, "参数有误", err)
		return
	}

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
		_, exist, err := userRepository.CheckUserIsFriend(userdata.ID, req.UserID)
		if err != nil {
			util.H(c, http.StatusInternalServerError, "删除好友失败", err)
			return
		}
		if !exist {
			util.H(c, http.StatusBadRequest, "该用户不是您的好友", err)
			return
		}
	}

	//开启事务
	session := adb.GetMySQLConn().NewSession()
	session.Begin()
	defer session.Close()

	redissession := adb.Rediss.Pipeline()

	//删除关系
	isdel, err := userRepository.DeleteFriendRelative(userdata.ID, req.UserID, session)
	if err != nil || !isdel {
		log.Println(err)
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}

	//	删除双方redis好友映射表
	if err := redissession.HSet("UserToUserRelative", strconv.Itoa(userdata.ID), strings.Join(flist, ",")).Err(); err != nil && !errors.Is(err, redis.Nil) {
		log.Println(err)
		redissession.Discard()
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}

	fstr = adb.Rediss.HGet("UserToUserRelative", strconv.Itoa(req.UserID)).Val()
	flist = strings.Split(fstr, ",")
	strtid = strconv.Itoa(userdata.ID)
	for i, fid := range flist {
		if fid == strtid {
			if len(flist) > 1 {
				flist = append(flist[:i], flist[i+1:]...)
			} else {
				flist = flist[:0]
			}
			break
		}
	}
	if err := adb.Rediss.HSet("UserToUserRelative", strconv.Itoa(req.UserID), strings.Join(flist, ",")).Err(); err != nil && !errors.Is(err, redis.Nil) {
		log.Println(err)
		redissession.Discard()
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}

	srcuserdata, exist, err := userRepository.GetUserByUserID(req.UserID)
	if err != nil || !exist {
		log.Println(err)
		redissession.Discard()
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}
	dstuserdata, exist, err := userRepository.GetUserByUserID(userdata.ID)
	if err != nil || !exist {
		log.Println(err)
		redissession.Discard()
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}
	err = adb.NebulaInstance.DeleteEdge("UserAddUser ", srcuserdata.UUID, dstuserdata.UUID, true)
	if err != nil {
		log.Println(err)
		redissession.Discard()
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "删除好友失败", err)
		return
	}

	redissession.Exec()
	session.Commit()

	umsg := models.UserMessage{
		UserID:        userdata.ID,
		ReceiveUserID: req.UserID,
		MsgType:       config.MsgTypeRefreshFriend,
	}
	umsgbytes, _ := json.Marshal(umsg)
	models.TransmitMsg(umsgbytes, config.MsgTypeRefreshFriend)

	util.H(c, http.StatusOK, "删除好友成功", err)

}
