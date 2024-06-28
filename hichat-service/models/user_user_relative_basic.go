package models

import "time"

type UserUserRelative struct {
	ID           int    `xorm:"pk autoincr notnull index"`
	PreUserID    int    `xorm:"notnull"` //被申请人
	PreUserName  string //被申请人
	BackUserID   int    `xorm:"notnull"`
	BackUserName string
	CreatedAt    time.Time `xorm:"created"`
	DeletedAt    time.Time `xorm:"deleted"`
	UpdatedAt    time.Time `xorm:"updated"`
}

func (u *UserUserRelative) TableName() string {
	return "user_user_relative"
}

//type Friend struct {
//	Id        int32
//	UserName  string
//	NikeName  string
//	Email     string
//	Avatar    string
//	City      string
//	Age       string
//	CreatedAt time.Time
//	DeletedAt time.Time
//	UpdatedAt time.Time
//}
