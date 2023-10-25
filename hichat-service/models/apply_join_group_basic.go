package models

import "time"

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
