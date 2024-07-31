package UsersScripts

import (
	"go-websocket-server/models"
)

// SelectUserGroupList 查询用户加入的群列表(没有详情)
func (r *userRepository) SelectUserGroupList(userid int) (gur []models.GroupUserRelative, err error) {
	err = r.mysqldb.Table("group_user_relative").Where("user_id=?", userid).Find(&gur)
	return
}

// GetUserByUsername 通过用户名获取一条数据
func (r *userRepository) GetUserByUsername(username string) (user models.Users, exist bool, err error) {
	exist, err = r.mysqldb.Where("user_name = ?", username).Get(&user)
	return
}

// GetUserByUserID 通过用户ID获取一条用户信息数据
func (r *userRepository) GetUserByUserID(userid int) (userdata models.Users, exist bool, err error) {
	exist, err = r.mysqldb.Where("id = ?", userid).Get(&userdata)
	return
}

// CheckUserIsExist 检查用户是否存在
func (r *userRepository) CheckUserIsExist(userid int) (bool, models.Users, error) {
	var user models.Users
	exist, err := r.mysqldb.Where("id = ?", userid).Get(&user)
	return exist, user, err
}

// CheckUserIsFriend 检查是否为好友
func (r *userRepository) CheckUserIsFriend(userid int, targetuserid int) (userrelative models.UserUserRelative, exist bool, err error) {
	exist, err = r.mysqldb.Table(userrelative.TableName()).Where("pre_user_id=? and back_user_id=?", userid, targetuserid).Or("pre_user_id=? and back_user_id=?", targetuserid, userid).Get(&userrelative)
	return userrelative, exist, err
}

// DeleteFriendRelative 删除好友关系
func (r *userRepository) DeleteFriendRelative(userid int, targetUserid int) (bool, error) {
	_, err := r.mysqldb.Table("user_user_relative").Where("pre_user_id=? and back_user_id=?", userid, targetUserid).Delete()
	if err != nil {
		return false, err
	}
	_, err = r.mysqldb.Table("user_user_relative").Where("pre_user_id=? and back_user_id=?", targetUserid, userid).Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}
