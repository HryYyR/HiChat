package UsersScripts

import (
	"go-websocket-server/models"
	"xorm.io/xorm"
)

type UserRepository interface {
	SelectUserGroupList(userid int) (gur []models.GroupUserRelative, err error) //查询用户加入的群列表(没有详情)

	GetUserByUsername(username string) (user models.Users, exist bool, err error) //通过用户名获取一条用户信息数据
	GetUserByUserID(userID int) (userdata models.Users, exist bool, err error)    //通过用户ID获取一条用户信息数据

	CheckUserIsExist(userid int) (bool, models.Users, error)                                                      //检查用户是否存在
	CheckUserIsFriend(userid int, targetuserid int) (userrelative models.UserUserRelative, exist bool, err error) //检查是否为好友
	DeleteFriendRelative(userid int, targetUserid int, session *xorm.Session) (bool, error)                       //删除好友关系
	ConnectFriendRelative(relative *models.UserUserRelative, session *xorm.Session) (bool, error)                 //连接好友关系
}

type userRepository struct {
	mysqldb *xorm.Engine
}

func NewUserRepository(db *xorm.Engine) UserRepository {
	return &userRepository{mysqldb: db}
}
