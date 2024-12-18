package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"log"
	"strconv"
	"strings"
	"time"
	"xorm.io/xorm"
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

func (*GroupUserRelative) TableName() string {
	return "group_user_relative"
}

// Association 连接关系
func (r *GroupUserRelative) Association(group Group, session *xorm.Session) error {
	_, err := session.Table("group_user_relative").Insert(r) //插入关系
	if err != nil {
		return err
	}

	//redis插入到GroupToUserMap的关系
	var insertGroupToUserMapStr string
	strres := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(r.GroupID)).Val()
	if len(strres) == 0 {
		insertGroupToUserMapStr = strconv.Itoa(r.UserID)
	} else {
		insertGroupToUserMapStr = fmt.Sprintf("%s,%d", strres, r.UserID)
	}
	adb.Rediss.HSet("GroupToUserMap", strconv.Itoa(r.GroupID), insertGroupToUserMapStr)

	//redis插入到UserToGroupMap的关系
	var insertUserToGroupMapStr string
	str2res := adb.Rediss.HGet("UserToGroupMap", strconv.Itoa(r.UserID)).Val()
	if len(str2res) == 0 {
		insertUserToGroupMapStr = strconv.Itoa(r.GroupID)
	} else {
		insertUserToGroupMapStr = fmt.Sprintf("%s,%d", str2res, r.GroupID)
	}

	//fmt.Println("reactive groupidarr:", insertUserToGroupMapStr)
	adb.Rediss.HSet("UserToGroupMap", strconv.Itoa(r.UserID), insertUserToGroupMapStr)
	//fmt.Printf("groupid%d,userid%d\n", group.ID, r.UserID)

	return nil
}

// DisAssociation  断开连接关系
func (r *GroupUserRelative) DisAssociation(session *xorm.Session, groupdata Group) error {

	_, err := session.Table("group_user_relative").Where("user_id = ? and group_id=?", r.UserID, r.GroupID).Delete()
	if err != nil {
		return err
	}
	//更新人数
	_, err = session.Table("group").Where("id = ?", r.GroupID).Update(Group{MemberCount: groupdata.MemberCount - 1})
	if err != nil {
		return err
	}

	//todo  如果缓存不一致就忽略缓存，继续执行接下来的操作
	struid := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(r.GroupID)).Val()
	if len(struid) == 0 {
		log.Println("查询GroupToUserMap失败，群聊:", r.GroupID, " 不存在")
		return nil
	}

	strgid := adb.Rediss.HGet("UserToGroupMap", strconv.Itoa(r.UserID)).Val()
	if len(strgid) == 0 {
		log.Println("查询GroupToUserMap失败，用户:", r.UserID, " 不存在")
		return nil
	}

	struidarr := strings.Split(struid, ",")
	strgidarr := strings.Split(strgid, ",")

	for i, s := range struidarr {
		if s == strconv.Itoa(r.UserID) {
			if i+1 == len(struidarr) {
				struidarr = struidarr[:i]
			} else {
				struidarr = append(struidarr[:i], struidarr[i+1:]...)

			}
		}
	}

	for i, s := range strgidarr {
		if s == strconv.Itoa(r.GroupID) {
			if i+1 == len(strgidarr) {
				strgidarr = strgidarr[:i]
			} else {
				struidarr = append(struidarr[:i], struidarr[i+1:]...)
			}
		}
	}
	ansuid := strings.Join(struidarr, ",")
	ansgid := strings.Join(strgidarr, ",")

	fmt.Println("uid: ", ansuid, " gid: ", ansgid)
	adb.Rediss.HSet("GroupToUserMap", strconv.Itoa(r.GroupID), ansuid)
	adb.Rediss.HSet("UserToGroupMap", strconv.Itoa(r.UserID), ansgid)

	return nil
}

// DisAssociationAll  断开所有用户与此群的链接（解散群聊）
func (r *GroupUserRelative) DisAssociationAll(session *xorm.Session) error {
	_, err := session.Table(r.TableName()).Where("group_id=?", r.GroupID).Delete()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
