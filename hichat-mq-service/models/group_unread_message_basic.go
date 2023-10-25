package models

import "time"

type GroupUnreadMessage struct {
	ID           int       `xorm:"pk autoincr notnull index"`
	UserName     string    `xorm:"notnull"`
	UserID       int       `xorm:"notnull"`
	GroupID      int       `xorm:"notnull"`
	UnreadNumber int       `xorm:"notnull default(0)"`
	CreatedAt    time.Time `xorm:"created"`
	DeletedAt    time.Time `xorm:"deleted"`
	UpdatedAt    time.Time `xorm:"updated"`
}

func (GroupUnreadMessage) TableName() string {
	return "group_unread_message"
}

func (g GroupUnreadMessage) GetAllUnreadMsg() {

}
