package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"time"
)

type Group struct {
	ID            int    `xorm:"pk autoincr notnull index"`
	UUID          string `xorm:"notnull unique"`
	CreaterID     int    `xorm:"notnull"`
	CreaterName   string `xorm:"notnull"`
	GroupName     string `xorm:"notnull unique"`
	Avatar        string
	Grade         int `xorm:"default(1)"`
	MemberCount   int
	UnreadMessage int
	CreatedAt     time.Time `xorm:"created"`
	DeletedAt     time.Time `xorm:"deleted"`
	UpdatedAt     time.Time `xorm:"updated"`
}

func (Group) TableName() string {
	return "group"
}

// 获取指定用户的群聊列表
func GetUserGroupList(ID int) (map[int]Group, error) {
	grouplist := make(map[int]Group, 0)

	var groupidlist []int

	if err := adb.Ssql.Cols("group_id").Table("group_user_relative").Where("user_id=?", ID).Find(&groupidlist); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, groupid := range groupidlist {
		var groupitem Group
		if _, err := adb.Ssql.Table("group").Where("id=?", groupid).Get(&groupitem); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		grouplist[groupid] = groupitem
	}

	return grouplist, nil
}

// 获取成员数量
func (g Group) GetMemberCount() int {
	num, err := adb.Ssql.Table("group_user_relative").Where("group_id=?", g.ID).Count()
	if err != nil {
		fmt.Println("获取用户成员数量失败!", err)
		return 0
	}
	return int(num)
}
