package models

import (
	adb "hichat_static_server/ADB"
	"sort"
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

type GroupUserRelative struct {
	ID        int `xorm:"pk autoincr notnull index"`
	UserID    int
	GroupID   int
	GroupUUID string
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (g *Group) GetMessageList(grouplist *[]GroupMessage, currentnum int) error {
	msglist := make([]GroupMessage, 0)
	count, err := adb.Ssql.Table("group_message").Where("group_id = ?", g.ID).Count()
	if err != nil {
		return err
	}
	if currentnum >= int(count) {
		return nil
	}

	err = adb.Ssql.Table("group_message").Where("group_id = ?", g.ID).Desc("id").Limit(10, currentnum).Find(&msglist)
	if err != nil {
		return err
	}

	sort.Slice(msglist, func(i, j int) bool {
		return msglist[i].ID < (msglist[j].ID)
	})

	*grouplist = msglist
	return nil
}
