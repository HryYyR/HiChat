package models

import "time"

type UserUserRelative struct {
	ID           int    `xorm:"pk autoincr notnull index"`
	PreUserID    int    `xorm:"notnull unique(BackUserID)"` //被申请人
	PreUserName  string //被申请人
	BackUserID   int    `xorm:"notnull"`
	BackUserName string
	CreatedAt    time.Time `xorm:"created"`
	DeletedAt    time.Time `xorm:"deleted"`
	UpdatedAt    time.Time `xorm:"updated"`
}
