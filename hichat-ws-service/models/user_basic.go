package models

import (
	"github.com/golang-jwt/jwt/v4"
	adb "go-websocket-server/ADB"
	"time"
)

type Users struct {
	ID        int    `xorm:"pk autoincr notnull index"`
	UUID      string `xorm:"notnull unique"`
	UserName  string `xorm:"notnull"`
	NikeName  string `xorm:"notnull"`
	Password  string `xorm:"notnull"`
	Email     string `xorm:"notnull"`
	Salt      string `xorm:"notnull"`
	IP        string
	Avatar    string
	City      string `xorm:"default('')"`
	Age       int    `xorm:"default(1)"`
	Introduce string
	Grade     int       `xorm:"default(1)"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
	UpdatedAt time.Time `xorm:"updated"`
	LoginTime string    `xorm:"updated"`
	LoginOut  time.Time
}

func (u *Users) TableName() string {
	return "users"
}

// CheckUserExit 检查用户是否存在
func (u *Users) CheckUserExit() (Users, bool, error) {
	var applyuserdata Users
	has, err := adb.SqlStruct.Conn.Table("users").Where("id=?", u.ID).Get(&applyuserdata)
	if err != nil {
		return applyuserdata, false, err
	}
	if !has {
		return applyuserdata, false, nil
	}
	return applyuserdata, true, nil
}

type UserClaim struct {
	ID       int
	UUID     string
	UserName string
	jwt.RegisteredClaims
}
