package models

import "time"

type UsersFile struct {
	ID        int    `xorm:"pk autoincr notnull index"`
	Identity  string `xorm:"unique notnull"`
	Hash      string `xorm:"unique notnull"`
	Name      string
	Ext       string
	Size      int64
	Path      string
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
	UpdatedAt time.Time `xorm:"updated"`
}
