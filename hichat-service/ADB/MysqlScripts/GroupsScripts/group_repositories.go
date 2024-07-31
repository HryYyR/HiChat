package GroupScripts

import (
	"go-websocket-server/models"
	"xorm.io/xorm"
)

type GroupRepository interface {
	CheckGroupIsExist(groupId int) (groupData models.Group, exist bool, err error) //检查群聊是否存在
	ByGroupNameCheckGroupIsExist(groupName string) (BeUse bool, err error)         //检查群聊名称是否占用

	GetUserApplyJoinGroupCount(userid int, Status int) (applycount int64, err error) //GetUserApplyJoinGroupCount 获取此用户的特定申请状态的申请数
	CheckUserIsExistInGroup(userid int, groupId int) (exist bool, err error)         // CheckUserIsExistInGroup 检查用户是否已加入此群聊

	UpdateApplyJoinGroupStatus(applyId int, Status int) (updateCount int64, err error) //// UpdateApplyJoinGroupStatus 更新申请状态
}

type groupRepository struct {
	db *xorm.Engine
}

func NewGroupRepository(db *xorm.Engine) GroupRepository {
	return &groupRepository{db: db}
}

// CheckGroupIsExist 检查群聊是否存在
func (r *groupRepository) CheckGroupIsExist(groupId int) (groupData models.Group, exist bool, err error) {
	exist, err = r.db.Where("id=?", groupId).Get(&groupData)
	if err != nil {
		return models.Group{}, false, err
	}
	return
}

// ByGroupNameCheckGroupIsExist 检查群聊名称是否占用
func (r *groupRepository) ByGroupNameCheckGroupIsExist(groupName string) (BeUse bool, err error) {
	BeUse, err = r.db.Table("group").Where("group_name = ?", groupName).Exist()
	return
}

// GetUserApplyJoinGroupCount 获取此用户的特定申请状态的申请数
func (r *groupRepository) GetUserApplyJoinGroupCount(userid int, Status int) (applycount int64, err error) {
	applycount, err = r.db.Table("apply_join_group").Where("apply_user_id = ? and handle_status=?", userid, Status).Count()
	return
}

// CheckUserIsExistInGroup 检查用户是否已加入此群聊
func (r *groupRepository) CheckUserIsExistInGroup(userid int, groupId int) (exist bool, err error) {
	exist, err = r.db.Table("group_user_relative").Where("user_id = ? and group_id=?", userid, groupId).Exist()
	return
}

// UpdateApplyJoinGroupStatus 更新申请状态
func (r *groupRepository) UpdateApplyJoinGroupStatus(applyId int, Status int) (updateCount int64, err error) {
	updateCount, err = r.db.Table("apply_join_group").ID(applyId).Update(models.ApplyJoinGroup{HandleStatus: Status})
	return
}
