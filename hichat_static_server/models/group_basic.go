package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	adb "hichat_static_server/ADB"
	"hichat_static_server/common"
	"hichat_static_server/tool"
	"log"
	"strconv"
	"time"
)

type UserGroupList struct {
	GroupDetail           []GroupDetail
	ApplyJoinGroupMessage []ApplyJoinGroup
	ApplyAddUserMessage   []ApplyAddUser
	FriendList            []FriendResponse
}

type GroupDetail struct {
	GroupInfo   Group
	MessageList []GroupMessage
}

type Group struct {
	ID            int    `xorm:"pk autoincr notnull index"`
	UUID          string `xorm:"notnull unique"`
	CreaterID     int    `xorm:"notnull"`
	CreaterName   string `xorm:"notnull"`
	GroupName     string `xorm:"notnull"`
	Avatar        string
	Grade         int `xorm:"default(1)"`
	MemberCount   int
	UnreadMessage int
	CreatedAt     time.Time `xorm:"created"`
	DeletedAt     time.Time `xorm:"deleted"`
	UpdatedAt     time.Time `xorm:"updated"`
}

// GroupMessage 群聊消息
type GroupMessage struct {
	ID          int `xorm:"pk autoincr"`
	UserID      int `xorm:"notnull"`
	UserUUID    string
	UserName    string
	UserAvatar  string
	UserCity    string
	UserAge     string
	GroupID     int    `xorm:"notnull"`
	Msg         string `xorm:"notnull"`
	MsgType     int
	IsReply     bool //是否是回复消息
	ReplyUserID int  //如果是,被回复的用户id
	Context     []byte
	CreatedAt   time.Time `xorm:"created"`
	DeletedAt   time.Time `xorm:"deleted"`
	UpdatedAt   time.Time `xorm:"updated"`
}

//	type GroupUserRelative struct {
//		ID        int `xorm:"pk autoincr notnull index"`
//		UserID    int
//		GroupID   int
//		GroupUUID string
//		CreatedAt time.Time `xorm:"created"`
//		DeletedAt time.Time `xorm:"deleted"`
//		UpdatedAt time.Time `xorm:"updated"`
//	}

// SaveToRedis 保存群聊信息到redis
func (g *Group) SaveToRedis() error {
	key := fmt.Sprintf("group%d", g.ID)
	_, err := adb.Rediss.HMSet(key, map[string]interface{}{
		"ID":          g.ID,
		"GroupName":   g.GroupName,
		"Avatar":      g.Avatar,
		"CreaterID":   g.CreaterID,
		"CreaterName": g.CreaterName,
		"Grade":       g.Grade,
		"MemberCount": g.MemberCount,
		"CreatedAt":   tool.FormatTime(g.CreatedAt),
		"DeletedAt":   tool.FormatTime(g.DeletedAt),
	}).Result()
	if err != nil {
		log.Printf("Save to Redis failed for group %d: %v", g.ID, err)
		return err
	}
	adb.Rediss.Expire(key, time.Hour*360)
	return nil
}

// GetGroupInfo 获取群聊信息
func (g *Group) GetGroupInfo() (Group, error) {
	var groupinfo Group
	key := fmt.Sprintf("group%d", g.ID)
	//从redis获取数据
	var gdata = adb.Rediss.HGetAll(key).Val()
	if len(gdata) > 3 {
		log.Println("走redis")
		_ = mapstructure.Decode(gdata, &groupinfo)
		groupinfo.ID, _ = strconv.Atoi(gdata["ID"])
		groupinfo.CreaterID, _ = strconv.Atoi(gdata["CreaterID"])
		groupinfo.Grade, _ = strconv.Atoi(gdata["Grade"])
		groupinfo.MemberCount, _ = strconv.Atoi(gdata["MemberCount"])
		groupinfo.CreatedAt, _ = common.ParseTime(gdata["CreatedAt"])
		log.Printf("%+v", gdata)
		return groupinfo, nil
	}
	log.Println("走mysql")
	exit, err := adb.Ssql.Table("group").Where("id =?", g.ID).Get(&groupinfo)
	if !exit {
		return Group{}, fmt.Errorf("用户不存在")
	}
	if err != nil {
		log.Println("mysql查询失败", err)
		return Group{}, err
	}

	err = groupinfo.SaveToRedis()
	if err != nil {
		log.Println("保存到redis失败", err)
		return groupinfo, nil
	}

	return groupinfo, nil
}

// GetMessageList 获取消息列表
func (g *Group) GetMessageList(grouplist *[]GroupMessage, currentnum int) error {
	msglist := make([]GroupMessage, 0)

	count, err := adb.Ssql.Table("group_message").Where("group_id = ? ", g.ID).Count()
	if err != nil {
		return err
	}
	if currentnum >= int(count) {
		return nil
	}

	err = getMsgListFromCache(g, currentnum, &msglist)
	if err != nil {
		err = getMsgListFromDatabase(g, currentnum, &msglist)
		if err != nil {
			return err
		}
	}

	//sort.Slice(msglist, func(i, j int) bool {
	//	return msglist[i].ID < (msglist[j].ID)
	//})
	*grouplist = msglist
	return nil
}

// 从数据库获取消息列表
func getMsgListFromDatabase(g *Group, currentnum int, msglist *[]GroupMessage) error {
	var msgdata []GroupMessage
	err := adb.Ssql.Table("group_message").Where("group_id = ?", g.ID).Desc("id").Limit(20, currentnum).Find(&msgdata)
	if err != nil {
		return err
	}
	*msglist = msgdata
	return nil
}

// 从redis获取消息列表
func getMsgListFromCache(g *Group, currentnum int, msglist *[]GroupMessage) error {
	var msgdata []GroupMessage
	key := fmt.Sprintf("gm%s", strconv.Itoa(g.ID))
	//判断redis缓存是否存在
	lLen := adb.Rediss.LLen(key).Val()
	log.Println("len", lLen)
	if lLen == 0 {
		return nil
	}
	start := lLen - int64(currentnum) - 20
	if start <= 0 {
		start = 0
	}
	stop := lLen - int64(currentnum) - 1
	if stop <= 0 {
		stop = 0
	}
	result, err := adb.Rediss.LRange(key, start, stop).Result()
	if err != nil {
		return err
	}
	for _, msgstr := range result {
		var msgstruct GroupMessage
		bufferString := bytes.NewBufferString(msgstr).Bytes()
		err := json.Unmarshal(bufferString, &msgstruct)
		if err != nil {
			//todo: 解析有错误的不应该放入聊天记录,但是 因为只是类型转换错误导致的失败 而放弃这条记录,
			//todo: 会导致总记录数量不正确,最终导致拉取记录出问题
			//log.Println(err)
			//continue
		}
		//log.Println(msgstruct)
		msgdata = append(msgdata, msgstruct)
	}
	*msglist = msgdata
	return nil
}
