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
	"strings"
	"time"
)

// ExitGroup 退出群
func ExitGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	var rawdata models.Group //只有group_id
	rawbyte, err := c.GetRawData()
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法操作", err)
		return
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", err)
		return
	}

	userRepository := c.MustGet("userRepository").(UsersScripts.UserRepository)

	//查用户是否存在
	//var handleuserdata models.Users
	//has, err := adb.SqlStruct.Conn.Table("users").Where("id = ?", userdata.ID).Get(&handleuserdata)
	has, handleuserdata, err := userRepository.CheckUserIsExist(userdata.ID)
	if !has {
		util.H(c, http.StatusBadRequest, "用户不存在", nil)
		return
	}
	if err != nil {
		log.Println(err.Error())
		util.H(c, http.StatusInternalServerError, "操作失败", err)
		return
	}

	//查群聊是否存在
	tempgroupinfo := &models.Group{
		ID: rawdata.ID,
	}
	groupinfo, groupExist, err := tempgroupinfo.CheckGroupExit()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "操作失败", err)
		return
	}
	if !groupExist {
		util.H(c, http.StatusBadRequest, "群聊不存在", nil)
		return
	}

	// 相关删除操作

	//在服务里找到群聊对应用户id
	groupexist, _, err := groupinfo.GetGroupUserIdLIst()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "操作失败", err)
		return
	}
	if !groupexist {
		util.H(c, http.StatusBadRequest, "群聊不存在", err)
		return
	}

	//fmt.Printf("%+v\n", groupuserlist)
	session := adb.SqlStruct.Conn.NewSession()
	defer session.Close()
	session.Begin()

	var willbedeleteuseridlist []int        //将被删除的用户群聊关系的用户id
	if groupinfo.CreaterID == userdata.ID { //说明他是群主,删除所有联系
		result := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(groupinfo.ID)).Val()
		//查群里的user的id 如果redis不存在,查mysql,存redis
		if len(result) == 0 {
			err := session.Table("group_user_relative").Cols("user_id").Where("group_id = ?", groupinfo.ID).Find(&willbedeleteuseridlist)
			if err != nil {
				log.Println(err)
				session.Rollback()
				util.H(c, http.StatusInternalServerError, "操作失败", err)
				return
			}
			strArr := util.IntArrToStrArr(willbedeleteuseridlist)
			adb.Rediss.HSet("GroupToUserMap", strconv.Itoa(groupinfo.ID), strings.Join(strArr, ","))
		} else { //redis存在,保存int[]和str[]
			willbedeleteuseridlist = util.StrArrToIntArr(strings.Split(result, ","))
		}

		//唯一性约束自动删除:关系,消息,未读消息
		//todo 群不用删,但是关系也没删
		gur := models.GroupUserRelative{GroupID: groupinfo.ID}
		err := gur.DisAssociationAll(session)
		if err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "操作失败", err)
			return
		}
		err = groupinfo.ByGroupIDSetGroupStatus(session, -1)
		if err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "操作失败", err)
			return
		}
		//_, err = session.Table("group").Where("id = ?", groupinfo.ID).Delete() //删群
		//if err != nil {
		//	session.Rollback()
		//	util.H(c, http.StatusInternalServerError, "操作失败", err)
		//	return
		//}

		//redis 事务
		redisSession := adb.Rediss.Pipeline()
		defer redisSession.Close()
		//删redis群聊映射用户和用户映射群聊 数据
		uidarrstr := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(groupinfo.ID)).Val()

		if len(uidarrstr) == 0 {
			log.Println(err)
			session.Rollback()
			redisSession.Discard()
			util.H(c, http.StatusInternalServerError, "操作失败", nil)
			return
		}
		redisSession.HDel("GroupToUserMap", strconv.Itoa(groupinfo.ID))
		uidarr := strings.Split(uidarrstr, ",")

		for _, uid := range uidarr {
			gidarrstr := adb.Rediss.HGet("UserToGroupMap", uid).Val()
			gidarr := strings.Split(gidarrstr, ",")
			strSlice := util.DeleteStrSlice(gidarr, strconv.Itoa(groupinfo.ID))
			redisSession.HSet("UserToGroupMap", uid, strings.Join(strSlice, ","))
		}

		//删redis群聊消息
		rkey := fmt.Sprintf("gm%s", strconv.Itoa(groupinfo.ID))
		redisSession.Del(rkey)
		rrkey := fmt.Sprintf("group%s", strconv.Itoa(groupinfo.ID))
		redisSession.Del(rrkey)
		redisSession.HDel("GroupToUserMap", strconv.Itoa(groupinfo.ID))

		//将通知名单转为[]byte
		willBeDeleteUserIdListBytes, err := util.IntsToBytes(willbedeleteuseridlist)
		if err != nil {
			log.Println(err)
			session.Rollback()
			redisSession.Discard()
			util.H(c, http.StatusInternalServerError, "操作失败", nil)
			return
		}

		_, err = redisSession.Exec()
		if err != nil {
			log.Println(err)
			session.Rollback()
			redisSession.Discard()
			util.H(c, http.StatusInternalServerError, "操作失败", nil)
			return
		}

		// 通知群里的其他成员群聊已解散
		groupmsg := models.Message{
			MsgType: config.MsgTypeDissolveGroup,
			GroupID: groupinfo.ID,
			Context: willBeDeleteUserIdListBytes, //通知名单
		}
		msgbyte, _ := json.Marshal(groupmsg)
		models.TransmitMsg(msgbyte, config.MsgTypeDissolveGroup)

		//删用户的群聊列表里对应群聊
		for _, userid := range willbedeleteuseridlist {
			models.ServiceCenter.Mutex.Lock()
			delete(models.ServiceCenter.Clients[userid].Groups, groupinfo.ID)
			models.ServiceCenter.Mutex.Unlock()
		}
	} else {
		//只断开该用户对群的联系
		usergrouprealtive := &models.GroupUserRelative{
			UserID:  userdata.ID,
			GroupID: groupinfo.ID,
		}
		err := usergrouprealtive.DisAssociation(session, groupinfo)
		if err != nil {
			log.Println(err)
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "退出失败", err)
			return
		}

		// 保存退出消息
		groupmsg := models.GroupMessage{
			UserID:     userdata.ID,
			UserName:   userdata.UserName,
			UserCity:   handleuserdata.City,
			UserAge:    handleuserdata.Age,
			UserAvatar: handleuserdata.Avatar,
			GroupID:    groupinfo.ID,
			Msg:        fmt.Sprintf("%s退出了群聊", userdata.UserName),
			MsgType:    config.MsgTypeQuitGroup,
			CreatedAt:  time.Now().Local(),
		}
		_, err = session.Table("group_message").Insert(&groupmsg)
		if err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "操作失败", err)
			return
		}

		//同步用户列表里相关信息
		models.ServiceCenter.Mutex.Lock()
		delete(models.ServiceCenter.Clients[userdata.ID].Groups, groupinfo.ID)
		models.ServiceCenter.Mutex.Unlock()

		bytes, _ := json.Marshal(groupmsg)
		adb.Rediss.RPush(fmt.Sprintf("gm%s", strconv.Itoa(groupmsg.GroupID)), bytes)
		err = adb.Rediss.HIncrBy(fmt.Sprintf("group%s", strconv.Itoa(groupmsg.GroupID)), "MemberCount", -1).Err()
		if err != nil {
			log.Println(err)
		}

		//用于传输的消息体
		transmitmsg := models.Message{
			UserID:     userdata.ID,
			UserName:   userdata.UserName,
			UserCity:   handleuserdata.City,
			UserAge:    strconv.Itoa(handleuserdata.Age),
			UserAvatar: handleuserdata.Avatar,
			GroupID:    groupinfo.ID,
			Msg:        fmt.Sprintf("%s退出了群聊", userdata.UserName),
			MsgType:    config.MsgTypeQuitGroup,
			CreatedAt:  time.Now().Local(),
		}
		msgbyte, _ := json.Marshal(transmitmsg)
		// 通知群里的其他成员有用户退出
		models.TransmitMsg(msgbyte, config.MsgTypeQuitGroup)
	}

	session.Commit()

	util.H(c, http.StatusOK, "退出群聊成功", nil)
}
