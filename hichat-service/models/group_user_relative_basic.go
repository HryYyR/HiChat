package models

import (
	"fmt"
	"time"
	"xorm.io/xorm"
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

func (r *GroupUserRelative) Association(group Group, session *xorm.Session) error {
	_, err := session.Table("group_user_relative").Insert(&r) //插入关系
	if err != nil {
		fmt.Println(err)
		return err
	}
	ServiceCenter.Clients[r.UserID].Mutex.Lock()
	defer ServiceCenter.Clients[r.UserID].Mutex.Unlock()
	ServiceCenter.Clients[r.UserID].Groups[group.ID] = group
	GroupUserList[group] = append(GroupUserList[group], r.UserID)
	return nil
}
