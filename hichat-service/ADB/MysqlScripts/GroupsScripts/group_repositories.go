package UsersScripts

import (
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
)

type GroupRepository interface {
	SelectUserGroupList(userid int) error //example
}

type groupRepository struct {
	db *adb.Sql
}

func NewGroupRepository(db *adb.Sql) GroupRepository {
	return &groupRepository{db: db}
}

// SelectUserGroupList 查询用户加入的群列表(没有详情)
func (r *groupRepository) SelectUserGroupList(userid int) error {
	var gur []models.GroupUserRelative
	return r.db.Conn.Table("group_user_relative").Where("user_id=?", userid).Find(&gur)
}
