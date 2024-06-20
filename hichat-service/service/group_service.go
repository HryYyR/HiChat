package service

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var GroupLock sync.Mutex

type groupinfo struct {
	Groupname string
	Avatar    string `json:"Avatar"`
}

// GroupDetail 用于返回结果的结构体
type GroupDetail struct {
	GroupInfo   models.Group
	MessageList []models.GroupMessage
}

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	info, _ := c.Get("userdata")
	userdata := info.(*models.UserClaim)

	rawsbyte, _ := c.GetRawData()
	var rowdata groupinfo
	err := json.Unmarshal(rawsbyte, &rowdata)
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusBadRequest, "解析参数失败", err)
		return
	}

	UUID := util.GenerateUUID()
	var group = models.Group{
		UUID:        UUID,
		CreaterID:   userdata.ID,
		Avatar:      rowdata.Avatar,
		CreaterName: userdata.UserName,
		GroupName:   rowdata.Groupname,
		MemberCount: 1,
	}

	//判断群是否已存在,存在就禁止创建
	isexit, err := adb.SqlStruct.Conn.Table("group").Where("group_name = ?", rowdata.Groupname).Exist()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "发生了未知的错误", nil)
		fmt.Println(err)
		log.Println(err)
		return
	}
	if isexit {
		util.H(c, http.StatusBadRequest, "群聊名称已被使用", nil)
		return
	}

	//_, err = adb.SqlStruct.Conn.Table("group").Insert(&group) //插入群聊
	//if err != nil {
	//	fmt.Println(err)
	//	util.H(c, http.StatusInternalServerError, "插入群聊失败", nil)
	//	return
	//}

	GroupLock.Lock()
	defer GroupLock.Unlock() //解锁

	session := adb.SqlStruct.Conn.NewSession()
	defer session.Close()
	session.Begin()

	var groupdata models.Group

	fullgroupdata, err := group.InsertGroup(session)
	if err != nil {
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "群聊创建失败", nil)
		fmt.Println(err)
		return
	}
	groupdata = fullgroupdata

	//连接关系
	var gur = models.GroupUserRelative{
		UserID:    userdata.ID,
		GroupUUID: UUID,
		GroupID:   groupdata.ID,
	}
	err = gur.Association(groupdata, session)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "群聊创建失败", nil)
		log.Println(err)
		fmt.Println(err)
		session.Rollback()
		return
	}

	session.Commit()

	responsedata := GroupDetail{
		GroupInfo:   groupdata,
		MessageList: []models.GroupMessage{},
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建成功",
		"data": responsedata,
	})
}

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
		//fmt.Println(err)
		util.H(c, http.StatusBadRequest, "非法访问", nil)
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		//fmt.Println(err)
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	if len(rawdata.Msg) > 50 {
		util.H(c, http.StatusBadRequest, "申请理由超字数上限(50字)", nil)
		return
	}

	applycount, err := adb.SqlStruct.Conn.Table("apply_join_group").Where("apply_user_id = ? and handle_status=?", userdata.ID, 0).Count()
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusInternalServerError, "查询关系失败", nil)
		return
	}
	if applycount >= 5 {
		util.H(c, http.StatusBadRequest, "申请已达上限", nil)
		return
	}

	//fmt.Println(userdata.ID, rawdata.GroupID)
	exitgroup, err := adb.SqlStruct.Conn.Table("group_user_relative").Where("user_id = ? and group_id=?", userdata.ID, rawdata.GroupID).Exist()
	if err != nil {
		fmt.Println(err)
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
	exit, tempapplydata, err2 := applymsgdata.CheckApplyExit()
	if err2 != nil {
		fmt.Println("check apply error", err)
		util.H(c, http.StatusInternalServerError, "查询申请失败", nil)
		return
	}
	//fmt.Println("applydata:", exit, tempapplydata.HandleStatus)
	if exit && tempapplydata.HandleStatus == 0 {
		util.H(c, http.StatusBadRequest, "你已经申请过了", nil)
		return
	}

	applygroupdata := &models.Group{
		ID: rawdata.GroupID,
	}
	applygroup, exit, err := applygroupdata.CheckGroupExit()
	if err != nil {
		fmt.Println("check group error", err)
		util.H(c, http.StatusInternalServerError, "查询群聊失败", nil)
		return
	}
	if !exit {
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
	fmt.Printf("%#v\n", applydata)
	err = applydata.InsertApply()
	if err != nil {
		fmt.Println("insert error", err)
		util.H(c, http.StatusInternalServerError, "申请失败", nil)
		return
	}
	//if _, err = adb.SqlStruct.Conn.Table("apply_join_group").Insert(&applydata); err != nil {
	//	util.H(c, http.StatusInternalServerError, "申请失败", nil)
	//	return
	//}

	msg := models.Message{
		UserID:   userdata.ID,
		UserName: userdata.UserName,
		GroupID:  applydata.GroupID,
		MsgType:  config.MsgTypeRefreshGroupNotice,
	}
	msgbyte, _ := json.Marshal(msg)

	//向群主发送验证申请信息
	models.ServiceCenter.Clients[applygroup.CreaterID].Send <- msgbyte
	models.TransmitMsg(msgbyte, config.MsgTypeRefreshGroupNotice)

	util.H(c, http.StatusOK, "申请成功", nil)
}

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

	// 查询申请是否存在
	tempapplyjoindata := &models.ApplyJoinGroup{
		ID: rawdata.ApplyID,
	}
	applyExit, applyjoindata, err := tempapplyjoindata.CheckApplyExit()
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusInternalServerError, "查询申请失败", nil)
		return
	}
	if !applyExit {
		fmt.Println(err)
		util.H(c, http.StatusBadRequest, "申请不存在", nil)
		return
	}

	//has, err := adb.SqlStruct.Conn.Table("apply_join_group").Where("id = ?", rawdata.ApplyID).Get(&applyjoindata)
	//if err != nil {
	//	fmt.Println(err)
	//	util.H(c, http.StatusInternalServerError, "查询申请失败", nil)
	//	return
	//}
	//if !has {
	//	fmt.Println(err)
	//	util.H(c, http.StatusBadRequest, "申请不存在", nil)
	//	return
	//}

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

	//has, err := adb.SqlStruct.Conn.Table("users").Where("id=?", applyjoindata.ApplyUserID).Get(&applyuserdata)
	//if err != nil {
	//	fmt.Println(err)
	//	util.H(c, http.StatusInternalServerError, "查询用户信息失败", nil)
	//	return
	//}
	//if !has {
	//	util.H(c, http.StatusBadRequest, "用户不存在", nil)
	//	return
	//}

	//上锁,以免数据错误
	var mute sync.Mutex
	mute.Lock()
	defer mute.Unlock()

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

	//
	//has, err := adb.SqlStruct.Conn.Table("group").Where("id=?", applyjoindata.GroupID).Get(&grouplist)
	//if err != nil {
	//	fmt.Println(err)
	//	util.H(c, http.StatusInternalServerError, "查询群聊失败", nil)
	//	return
	//}
	//if !has {
	//	util.H(c, http.StatusBadRequest, "群聊不存在", nil)
	//	return
	//}

	// 拒绝申请
	if rawdata.HandleStatus == -1 {
		if _, err = adb.SqlStruct.Conn.Table("apply_join_group").ID(applyjoindata.ID).Update(models.ApplyJoinGroup{HandleStatus: rawdata.HandleStatus}); err != nil {
			util.H(c, http.StatusInternalServerError, "拒绝失败", nil)
			return
		}

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
		if _, err = session.Table("apply_join_group").ID(applyjoindata.ID).Update(models.ApplyJoinGroup{HandleStatus: rawdata.HandleStatus}); err != nil {
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

		//修改redis相关
		//result, _ := adb.Rediss.HExists("GroupToUserMap", strconv.Itoa(grouplist.ID)).Result()
		//if result {
		//	datastr, _ := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(grouplist.ID)).Result()
		//	splitarr := strings.Split(datastr, ",")
		//	splitarr = append(splitarr, strconv.Itoa(applyuserdata.ID))
		//	resstr := strings.Join(splitarr, ",")
		//	adb.Rediss.HSet("GroupToUserMap", strconv.Itoa(grouplist.ID), resstr)
		//}

		groupmsg := models.GroupMessage{
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
		adb.Rediss.HSet(fmt.Sprintf("group%d", grouplist.ID), "MemberCount", grouplist.MemberCount+1)

		// 通知群里的其他成员有用户加入
		for _, uid := range models.GroupUserList[grouplist.ID] {
			models.ServiceCenter.Clients[uid].Send <- msgbyte
			break
		}
		models.TransmitMsg(msgbyte, config.MsgTypeJoinGroup)

		//通知申请用户刷新群聊列表
		selfmsg := models.Message{
			UserID:   applyuserdata.ID,
			UserName: applyuserdata.UserName,
			GroupID:  grouplist.ID,
			MsgType:  config.MsgTypeRefreshGroup,
		}
		selfmsgbyte, _ := json.Marshal(selfmsg)
		models.ServiceCenter.Clients[applyuserdata.ID].Send <- selfmsgbyte
		models.TransmitMsg(selfmsgbyte, config.MsgTypeRefreshGroup)

		util.H(c, http.StatusOK, "用户已加入", nil)
	}

	//通知申请人,申请已被处理(刷新通知列表)
	groupmsg := models.GroupMessage{
		UserID:   applyuserdata.ID,
		UserName: applyuserdata.UserName,
		GroupID:  grouplist.ID,
		MsgType:  config.MsgTypeRefreshGroupNotice,
	}
	bytemsg, err := json.Marshal(groupmsg)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "处理失败", nil)
		return
	}
	models.ServiceCenter.Clients[applyuserdata.ID].Send <- bytemsg
	models.TransmitMsg(bytemsg, config.MsgTypeRefreshGroupNotice)

}

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

	//查用户是否存在
	var handleuserdata models.Users
	has, err := adb.SqlStruct.Conn.Table("users").Where("id = ?", userdata.ID).Get(&handleuserdata)
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
	groupexist, groupuserlist, err := groupinfo.GetGroupUserIdLIst()
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
				session.Rollback()
				util.H(c, http.StatusInternalServerError, "操作失败", err)
				return
			}
			strArr := util.IntArrToStrArr(willbedeleteuseridlist)
			adb.Rediss.HSet("GroupToUserMap", strconv.Itoa(groupinfo.ID), strings.Join(strArr, ","))
		} else { //redis存在,保存int[]和str[]
			willbedeleteuseridlist = util.StrArrToIntArr(strings.Split(result, ","))
		}

		//redis 事务
		redisSession := adb.Rediss.Pipeline()
		defer redisSession.Close()
		//删redis群聊映射用户和用户映射群聊 数据
		uidarrstr := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(groupinfo.ID)).Val()

		if len(uidarrstr) == 0 {
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
		fmt.Println("ok")

		//删redis群聊消息
		rkey := fmt.Sprintf("gm%s", strconv.Itoa(groupinfo.ID))
		redisSession.Del(rkey)

		rrkey := fmt.Sprintf("group%s", strconv.Itoa(groupinfo.ID))
		redisSession.Del(rrkey)

		_, err := redisSession.Exec()
		if err != nil {
			session.Rollback()
			redisSession.Discard()
			util.H(c, http.StatusInternalServerError, "操作失败", nil)
			return
		}

		//唯一性约束自动删除:关系,消息,未读消息
		_, err = session.Table("group").Where("id = ?", groupinfo.ID).Delete() //删群
		if err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "操作失败", err)
			return
		}

		//删用户的群聊列表里对应群聊
		for _, userid := range willbedeleteuseridlist {
			delete(models.ServiceCenter.Clients[userid].Groups, groupinfo.ID)
		}

		//同步群列表里相关信息
		delete(models.GroupUserList, groupinfo.ID)

		// 通知群里的其他成员群聊已解散
		groupmsg := models.GroupMessage{
			MsgType: config.MsgTypeDissolveGroup,
			GroupID: groupinfo.ID,
		}
		msgbyte, _ := json.Marshal(groupmsg)
		for _, userid := range groupuserlist {
			models.ServiceCenter.Clients[userid].Send <- msgbyte
		}
		models.TransmitMsg(msgbyte, config.MsgTypeDissolveGroup)

	} else {
		//只断开该用户对群的联系
		usergrouprealtive := &models.GroupUserRelative{
			UserID:  userdata.ID,
			GroupID: groupinfo.ID,
		}
		err := usergrouprealtive.DisAssociation(session, groupinfo)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
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
		delete(models.ServiceCenter.Clients[userdata.ID].Groups, groupinfo.ID)

		msgbyte, _ := json.Marshal(groupmsg)
		adb.Rediss.RPush(fmt.Sprintf("gm%s", strconv.Itoa(groupmsg.GroupID)), string(msgbyte))
		// 通知群里的其他成员有用户退出
		for _, userid := range groupuserlist {
			models.ServiceCenter.Clients[userid].Send <- msgbyte
		}
		models.TransmitMsg(msgbyte, config.MsgTypeQuitGroup)

	}
	session.Commit()

	util.H(c, http.StatusOK, "退出群聊成功", nil)
}

type searchgroupinfo struct {
	Searchstr string
}

// SearchGroup 搜索群聊
func SearchGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)
	var rawdata searchgroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法访问", nil)
		return
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	if len(strings.TrimSpace(rawdata.Searchstr)) == 0 {
		util.H(c, http.StatusBadRequest, "关键词不能为空", nil)
		return
	}

	if len(strings.TrimSpace(rawdata.Searchstr)) > 50 {
		util.H(c, http.StatusBadRequest, "关键词超字数上限(50字)", nil)
		return
	}

	var searchint int
	v, err := strconv.Atoi(rawdata.Searchstr)
	if err == nil {
		searchint = v
	}

	grouplist := make([]models.Group, 0)
	err = adb.SqlStruct.Conn.Table("group").Where("group_name like ? or id=?", rawdata.Searchstr+"%", searchint).Where("creater_id !=?", userdata.ID).Find(&grouplist)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "搜索失败", err)
		return
	}

	responsedata := &[]models.Group{}
	for _, group := range grouplist {
		count, err := group.GetMemberCount()
		if err != nil {
			log.Println(err)
		}
		group.MemberCount = count
		*responsedata = append(*responsedata, group)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "搜索成功",
		"grouplist": responsedata,
	})

}
