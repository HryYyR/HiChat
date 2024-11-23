package models

import (
	adb "go-websocket-server/ADB"
	"time"
)

type ApplyJoinGroup struct {
	ID            int    `xorm:"pk autoincr notnull index"`
	ApplyUserID   int    `xorm:"notnull"`
	ApplyUserName string `xorm:"notnull"`
	GroupID       int    `xorm:"notnull"`
	ApplyMsg      string
	ApplyWay      int       `xorm:"default(1)"`
	HandleStatus  int       `xorm:"default(0)"`
	CreatedAt     time.Time `xorm:"created"`
	DeletedAt     time.Time `xorm:"deleted"`
	UpdatedAt     time.Time `xorm:"updated"`
}

func (*ApplyJoinGroup) TableName() string {
	return "apply_join_group"
}

// CheckApplyExit todo  +redis
// 检查申请的状态
func (a *ApplyJoinGroup) CheckApplyExit() (bool, ApplyJoinGroup, error) {
	var applydata ApplyJoinGroup
	var exit bool
	var err error

	if a.ID == 0 {
		exit, err = adb.SqlStruct.Conn.Table("apply_join_group").Where("apply_user_id = ? and group_id=?", a.ApplyUserID, a.GroupID).OrderBy("created_at DESC").Get(&applydata)

	} else {
		exit, err = adb.SqlStruct.Conn.Table("apply_join_group").Where("id=?", a.ID).OrderBy("created_at DESC").Get(&applydata)
	}

	if err != nil {
		return false, applydata, err
	}
	if !exit {
		return false, applydata, nil
	}
	return true, applydata, nil
}

// InsertApply 插入申请
func (a *ApplyJoinGroup) InsertApply() error {
	applydata := *a
	if _, err := adb.SqlStruct.Conn.Table("apply_join_group").Insert(&applydata); err != nil {
		return err
	}
	return nil
}
