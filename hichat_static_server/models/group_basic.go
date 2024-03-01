package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	adb "hichat_static_server/ADB"
	"sort"
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

//
//type GroupUserRelative struct {
//	ID        int `xorm:"pk autoincr notnull index"`
//	UserID    int
//	GroupID   int
//	GroupUUID string
//	CreatedAt time.Time `xorm:"created"`
//	DeletedAt time.Time `xorm:"deleted"`
//	UpdatedAt time.Time `xorm:"updated"`
//}

func (g *Group) GetMessageList(grouplist *[]GroupMessage, currentnum int) error {
	msglist := make([]GroupMessage, 0)

	count, err := adb.Ssql.Table("group_message").Where("group_id = ?", g.ID).Count()
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

	sort.Slice(msglist, func(i, j int) bool {
		return msglist[i].ID < (msglist[j].ID)
	})
	*grouplist = msglist
	return nil
}
func getMsgListFromDatabase(g *Group, currentnum int, msglist *[]GroupMessage) error {
	var msgdata []GroupMessage
	err := adb.Ssql.Table("group_message").Where("group_id = ?", g.ID).Desc("id").Limit(10, currentnum).Find(&msgdata)
	if err != nil {
		return err
	}
	*msglist = msgdata
	return nil
}
func getMsgListFromCache(g *Group, currentnum int, msglist *[]GroupMessage) error {
	var msgdata []GroupMessage
	key := fmt.Sprintf("gm%s", strconv.Itoa(g.ID))
	//判断redis缓存是否存在
	lLen, err := adb.Rediss.LLen(key).Result()
	fmt.Println("len", lLen)
	if err != nil {
		return err
	}
	if lLen == 0 {
		return nil
	}
	start := lLen - int64(currentnum) - 10
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
			continue
		}
		//fmt.Println(msgstruct)
		msgdata = append(msgdata, msgstruct)
	}
	*msglist = msgdata
	return nil
}
