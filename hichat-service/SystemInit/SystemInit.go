package systeminit

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"time"
)

func PrintRoomInfo() {
	fmt.Println("init prinroominfo...")
	for {
		time.Sleep(10 * time.Second)
		prinstr := "\n用户列表 :[\n"
		for _, j := range models.ServiceCenter.Clients {
			prinstr += fmt.Sprintf("id:%v name:%v status:%v groups:[", j.UserID, j.UserName, j.Status)
			for _, k := range j.Groups {
				prinstr += fmt.Sprintf(" %v:%v", k.ID, k.GroupName)
				if k.CreaterID == j.UserID {
					prinstr += "(群主)"
				}
			}
			prinstr += " ]\n"
		}
		prinstr += "]\n"

		printgroupstr := "群列表: [\n"
		for gid, useridarr := range models.GroupUserList {
			printgroupstr += fmt.Sprintf("%v : %v\n", gid, useridarr)
		}
		printgroupstr += "]\n"

		fmt.Println(prinstr)
		fmt.Println(printgroupstr)

	}

}

// InitUserToClient 初始化用户到内存中 ServiceCenter.Clients
func InitUserToClient() error {
	var userdatalist []models.Users

	err := adb.SqlStruct.Conn.Table("users").Find(&userdatalist)
	if err != nil {
		return err
	}
	for _, userdata := range userdatalist {
		uuid := util.GenerateUUID()
		grouplist, err := models.GetUserGroupList(userdata.ID)
		if err != nil {
			return err
		}
		client := models.UserClient{
			ClientID: uuid,
			UserID:   userdata.ID,
			UserUUID: userdata.UUID,
			UserName: userdata.UserName,
			Send:     make(chan []byte, 256),
			Status:   false,
			Groups:   grouplist,
			// CachingMessages: make(map[int]int, 0),
		}
		models.ServiceCenter.Clients[userdata.ID] = client
	}

	return nil
}

// InitClientsToGrouplist 初始化群和用户的关系列表 GroupUserList
func InitClientsToGrouplist() error {
	var grouplist []models.Group
	err := adb.SqlStruct.Conn.Table("group").Find(&grouplist) //查询所有的群
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, g := range grouplist {
		var useridlist []int
		err = adb.SqlStruct.Conn.Cols("user_id").Table("group_user_relative").Where("group_id=?", g.ID).Find(&useridlist)
		if err != nil {
			return err
		}
		models.GroupUserList[g.ID] = useridlist
	}

	return nil
}
