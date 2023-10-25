package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"time"
)

var GroupUserList = make(map[Group][]int, 0) //群和用户的关系列表 k:group  v:user_id

type GroupUserRelative struct {
	ID        int `xorm:"pk autoincr notnull index"`
	UserID    int `xorm:"notnull"`
	GroupID   int `xorm:"notnull"`
	GroupUUID string
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (GroupUserRelative) TableName() string {
	return "group_user_relative"
}

func (r *GroupUserRelative) Association(group Group) error {
	_, err := adb.Ssql.Table("group_user_relative").Insert(&r) //插入关系
	if err != nil {
		fmt.Println(err)
		return err
	}
	ServiceCenter.Clients[r.UserID].Mutex.Lock()
	ServiceCenter.Clients[r.UserID].Groups[group.ID] = group
	GroupUserList[group] = append(GroupUserList[group], r.UserID)
	ServiceCenter.Clients[r.UserID].Mutex.Unlock()

	return nil
}
