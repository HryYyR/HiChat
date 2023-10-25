package models

import "time"

type ApplyAddUser struct {
	ID               int    `xorm:"pk autoincr notnull index"`
	PreApplyUserID   int    `xorm:"notnull"` //被申请人
	PreApplyUserName string `xorm:"notnull"`
	ApplyUserID      int    `xorm:"notnull"` //申请人
	ApplyUserName    string `xorm:"notnull"`
	ApplyMsg         string
	ApplyWay         int       `xorm:"default(1)"`
	HandleStatus     int       `xorm:"default(0)"`
	CreatedAt        time.Time `xorm:"created"`
	DeletedAt        time.Time `xorm:"deleted"`
	UpdatedAt        time.Time `xorm:"updated"`
}
