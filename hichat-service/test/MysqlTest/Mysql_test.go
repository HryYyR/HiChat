package MysqlTest

import (
	"fmt"
	adb "go-websocket-server/ADB"
	GroupScripts "go-websocket-server/ADB/MysqlScripts/GroupsScripts"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"testing"
)

func TestSelectUserGroupList(t *testing.T) {
	userRepository := UsersScripts.NewUserRepository(adb.GetMySQLConn())
	grouplist, err := userRepository.SelectUserGroupList(1015)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(grouplist, err)

	grouplist, err = userRepository.SelectUserGroupList(1014)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(grouplist, err)
}

// TestCheckGroupIsExist 检查群聊是否存在
func TestCheckGroupIsExist(t *testing.T) {
	groupRepository := GroupScripts.NewGroupRepository(adb.GetMySQLConn())
	group, exist, err := groupRepository.CheckGroupIsExist(238)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(group, exist, err)

	group, exist, err = groupRepository.CheckGroupIsExist(237)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(group, exist, err)
}

func TestCheckUserIsExist(t *testing.T) {
	groupRepository := UsersScripts.NewUserRepository(adb.GetMySQLConn())
	exist, user, err := groupRepository.CheckUserIsExist(1015)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user, exist, err)

	exist, user, err = groupRepository.CheckUserIsExist(1014)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user, exist, err)
}

func TestCheckUserIsFriend(t *testing.T) {
	groupRepository := UsersScripts.NewUserRepository(adb.GetMySQLConn())
	userrelative, exist, err := groupRepository.CheckUserIsFriend(1015, 1016)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userrelative, exist, err)

	userrelative, exist, err = groupRepository.CheckUserIsFriend(1014, 1015)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userrelative, exist, err)
}

func TestGetUserByUserIDAndName(t *testing.T) {
	groupRepository := UsersScripts.NewUserRepository(adb.GetMySQLConn())
	userrelative, exist, err := groupRepository.GetUserByUserID(1015)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userrelative, exist, err)

	userrelative, exist, err = groupRepository.GetUserByUserID(1014)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userrelative, exist, err)

	userrelative, exist, err = groupRepository.GetUserByUsername("nekoni")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userrelative, exist, err)

	userrelative, exist, err = groupRepository.GetUserByUsername("none")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userrelative, exist, err)
}
