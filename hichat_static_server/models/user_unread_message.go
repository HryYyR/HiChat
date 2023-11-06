package models

import "time"

type UserUnreadMessage struct {
	ID           int       `xorm:"pk autoincr notnull"`
	UserName     string    `xorm:"notnull"`
	UserID       int       `xorm:"notnull index(uid_fid_unique)"`
	FriendID     int       `xorm:"notnull index(uid_fid_unique)"`
	UnreadNumber int       `xorm:"notnull default(0)"`
	CreatedAt    time.Time `xorm:"created"`
	DeletedAt    time.Time `xorm:"deleted"`
	UpdatedAt    time.Time `xorm:"updated"`
}

func (UserUnreadMessage) TableName() string {
	return "user_unread_message"
}
