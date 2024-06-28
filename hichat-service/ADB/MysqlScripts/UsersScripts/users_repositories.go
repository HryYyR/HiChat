package UsersScripts

import (
	"go-websocket-server/models"
	"xorm.io/xorm"
)

type UserRepository interface {
	SelectUserGroupList(userid int) error
	GetUserByUsername(username string) (*models.Users, error)
	CheckUserIsExist(userid int) (bool, *models.Users, error)
	CheckUserIsFriend(userid int, targetuserid int) (bool, error)
	DeleteFriendRelative(userid int, targetUserid int) (bool, error)
}

type userRepository struct {
	db *xorm.Engine
}

func NewUserRepository(db *xorm.Engine) UserRepository {
	return &userRepository{db: db}
}

// SelectUserGroupList 查询用户加入的群列表(没有详情)
func (r *userRepository) SelectUserGroupList(userid int) error {
	var gur []models.GroupUserRelative
	return r.db.Table("group_user_relative").Where("user_id=?", userid).Find(&gur)
}

// GetUserByUsername 通过用户名获取一条数据
func (r *userRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users
	_, err := r.db.Where("username = ?", username).Get(&user)
	return &user, err
}

// CheckUserIsExist 检查用户是否存在
func (r *userRepository) CheckUserIsExist(userid int) (bool, *models.Users, error) {
	var user models.Users
	exist, err := r.db.Table(user.TableName()).Where("id = ?", userid).Exist(&user)
	return exist, &user, err
}

// CheckUserIsFriend 检查是否为好友
func (r *userRepository) CheckUserIsFriend(userid int, targetuserid int) (bool, error) {
	//todo
	return false, nil
}

// DeleteFriendRelative 删除好友关系
func (r *userRepository) DeleteFriendRelative(userid int, targetUserid int) (bool, error) {
	_, err := r.db.Table("user_user_relative").Where("pre_user_id=? and back_user_id=?", userid, targetUserid).Delete()
	if err != nil {
		return false, err
	}
	_, err = r.db.Table("user_user_relative").Where("pre_user_id=? and back_user_id=?", targetUserid, userid).Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}
