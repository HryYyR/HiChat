package UsersScripts

import (
	"go-websocket-server/models"
)

type UserService struct {
	userRepo UserRepository
}

type UserInterface interface {
	SelectUserGroupList(userid int) error
	GetUserByUsername(username string) (*models.Users, error)
}

func (s *UserService) SelectUserGroupList(userid int) error {
	return s.userRepo.SelectUserGroupList(userid)
}

func (s *UserService) GetUserByUsername(username string) (*models.Users, error) {
	return s.userRepo.GetUserByUsername(username)
}

func NewUserService(userRepo UserRepository) UserInterface {
	return &UserService{userRepo: userRepo}
}
