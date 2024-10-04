package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/util"
	"log"
	"strconv"
	"strings"
	"time"
	"xorm.io/xorm"
)

type Group struct {
	ID            int    `xorm:"pk autoincr notnull index"`
	UUID          string `xorm:"notnull unique"`
	CreaterID     int    `xorm:"notnull"`
	CreaterName   string `xorm:"notnull"`
	GroupName     string `xorm:"notnull unique"`
	Avatar        string
	Grade         int `xorm:"default(1)"`
	MemberCount   int
	UnreadMessage int
	CreatedAt     time.Time `xorm:"created"`
	DeletedAt     time.Time `xorm:"deleted"`
	UpdatedAt     time.Time `xorm:"updated"`
	Status        int       `xorm:"notnull default(0)"`
}

func (g *Group) TableName() string {
	return "group"
}

// ByGroupIDSetGroupStatus 通过群聊id修改群聊状态（启用 | 已删除）
func (g *Group) ByGroupIDSetGroupStatus(session *xorm.Session, status int) error {
	g.Status = status
	_, err := session.ID(g.ID).Cols("status").Update(g)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetUserGroupList todo 没用redis
// 获取指定用户的群聊列表
func GetUserGroupList(ID int) (map[int]Group, error) {
	grouplist := make(map[int]Group)
	var groupidlist []int

	if err := adb.SqlStruct.Conn.Cols("group_id").Table("group_user_relative").Where("user_id=?", ID).Find(&groupidlist); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, groupid := range groupidlist {
		var groupitem Group
		if _, err := adb.SqlStruct.Conn.Table("group").Where("id=?", groupid).Get(&groupitem); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		grouplist[groupid] = groupitem
	}

	return grouplist, nil
}

// GetMemberCount 获取成员数量
func (g *Group) GetMemberCount() (int, error) {
	result, err := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(g.ID)).Result()
	if err != nil || len(result) == 0 {
		sqlres, err := adb.SqlStruct.Conn.Table("group_user_relative").Where("group_id=?", g.ID).Count()
		if err != nil {
			fmt.Println("获取用户成员数量失败!", err)
			return 0, err
		}
		return int(sqlres), nil
	}
	redisres := strings.Split(result, ",")
	return len(redisres), nil

}

// InsertGroup 插入群聊
func (g *Group) InsertGroup(session *xorm.Session) (Group, error) {
	var groupdata Group

	//插入mysql
	_, err := session.Table("group").Insert(g)
	if err != nil {
		return groupdata, err
	}

	//查询完整
	has, err := session.Table("group").Where("id=?", g.ID).Get(&groupdata)
	if !has {
		return groupdata, fmt.Errorf(`group %d not found`, g.ID)
	}
	if err != nil {
		return groupdata, err
	}

	return groupdata, nil
}

// CheckGroupExit 检查群聊是否存在
func (g *Group) CheckGroupExit() (Group, bool, error) {
	var groupdata Group
	exitgroup, err := adb.SqlStruct.Conn.Table("group").Where("id = ?", g.ID).Get(&groupdata)
	if err != nil {
		return groupdata, false, err
	}
	if !exitgroup {
		return groupdata, false, nil
	}
	return groupdata, true, nil
}

// GetGroupUserIdLIst 获取群聊内用户的id数组
func (g *Group) GetGroupUserIdLIst() (bool, []int, error) {
	var useridarr []int
	result := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(g.ID)).Val()
	if len(result) == 0 {
		if err := adb.SqlStruct.Conn.Cols("user_id").Table("group_user_relative").Where("group_id=?", g.ID).Find(&useridarr); err != nil {
			return false, useridarr, err
		}
		if len(useridarr) == 0 {
			return false, useridarr, nil
		}
	} else {
		strarr := strings.Split(result, ",")
		useridarr = util.StrArrToIntArr(strarr)
	}

	return true, useridarr, nil
}
