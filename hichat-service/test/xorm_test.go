package test

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"testing"
)

func TestXormTest(t *testing.T) {
	adb.InitMySQL()
	// err := adb.Engine.CreateTables(&models.User{})
	err := adb.Ssql.Sync2(new(models.Users), new(models.Group), new(models.GroupMessage), new(models.GroupUserRelative))
	if err != nil {
		t.Fatal(err)
	}
}

func TestXormsql(t *testing.T) {
	adb.InitMySQL()
	// err := adb.Engine.CreateTables(&models.User{})
	adb.Ssql.Table("users")
}

func TestGetUserGroupList(t *testing.T) {
	userdata := models.UserClaim{
		ID:       2,
		UUID:     "2028bdd7d36f",
		UserName: "niko",
	}
	adb.InitMySQL()

	grouplist := make(map[int]models.Group, 0)

	var groupidlist []int
	if err := adb.Ssql.Cols("group_id").Table("group_user_relative").Where("user_id=?", userdata.ID).Find(&groupidlist); err != nil {
		fmt.Println(err.Error())
		t.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", groupidlist)
	for _, groupid := range groupidlist {
		var groupitem models.Group
		if _, err := adb.Ssql.Table("group").Where("id=?", groupid).Get(&groupitem); err != nil {
			t.Fatal(err)
		}
		grouplist[groupid] = groupitem
	}

	fmt.Printf("%+v\n", grouplist)

}

func TestUserToClient(t *testing.T) {
	adb.InitMySQL()
	mockservicecenter := make(map[int]models.UserClient, 0)

	var userdatalist []models.Users
	err := adb.Ssql.Table("users").Find(&userdatalist)
	if err != nil {
		t.Fatal(err)
	}
	for _, userdata := range userdatalist {
		uuid := util.GenerateUUID()
		grouplist, err := models.GetUserGroupList(userdata.ID)
		if err != nil {
			t.Fatal(err)
		}
		client := models.UserClient{
			ClientID:        uuid,
			UserID:          userdata.ID,
			UserUUID:        userdata.UUID,
			UserName:        userdata.UserName,
			Send:            make(chan []byte, 256),
			Status:          false,
			Groups:          grouplist,
			CachingMessages: make(map[int]int, 0),
		}
		mockservicecenter[userdata.ID] = client
	}
}
