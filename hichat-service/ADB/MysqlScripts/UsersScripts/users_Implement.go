package UsersScripts

import "go-websocket-server/models"

// SelectUserGroupList 查询用户加入的群列表(没有详情)
func (r *userRepository) SelectUserGroupList(userid int) error {
	var gur []models.GroupUserRelative
	return r.mysqldb.Table("group_user_relative").Where("user_id=?", userid).Find(&gur)
}

// GetUserByUsername 通过用户名获取一条数据
func (r *userRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users
	_, err := r.mysqldb.Where("username = ?", username).Get(&user)
	return &user, err
}

// GetUserByUserID 通过用户ID获取一条用户信息数据
func (r *userRepository) GetUserByUserID(userid int) (userdata *models.Users, exist bool, err error) {
	exist, err = r.mysqldb.Where("userid = ?", userid).Get(&userdata)
	return
}

// CheckUserIsExist 检查用户是否存在
func (r *userRepository) CheckUserIsExist(userid int) (bool, *models.Users, error) {
	var user models.Users
	exist, err := r.mysqldb.Table(user.TableName()).Where("id = ?", userid).Exist(&user)
	return exist, &user, err
}

// CheckUserIsFriend 检查是否为好友
func (r *userRepository) CheckUserIsFriend(userid int, targetuserid int) (userrelative *models.UserUserRelative, exist bool, err error) {
	exist, err = r.mysqldb.Table(userrelative.TableName()).Where("pre_user_id=? and back_user_id=?", userid, targetuserid).Or("pre_user_id=? and back_user_id=?", targetuserid, userid).Exist(&userrelative)
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
