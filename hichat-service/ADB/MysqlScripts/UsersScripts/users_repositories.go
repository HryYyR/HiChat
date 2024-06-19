package UsersScripts

import (
	"go-websocket-server/models"
	"xorm.io/xorm"
)

type UserRepository interface {
	SelectUserGroupList(userid int) error
	GetUserByUsername(username string) (*models.Users, error)
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

func (r *userRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users
	_, err := r.db.Where("username = ?", username).Get(&user)
	return &user, err
}
