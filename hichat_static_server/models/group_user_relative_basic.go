package models

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"time"
)

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

// Association
func (r *GroupUserRelative) Association() error {
	_, err := adb.Ssql.Table("group_user_relative").Insert(&r) //插入关系
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
