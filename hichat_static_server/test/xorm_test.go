package test

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"testing"
)

func TestXormsql(t *testing.T) {
	adb.InitMySQL()

	var u models.Users
	u.ID = 1008

	//var grouplist []models.GroupDetail
	//
	//// 查询用户加入的群列表(没有详情)
	//var gur []models.GroupUserRelative
	//if err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).Find(&gur); err != nil {
	//	fmt.Println("查询用户加入的群列表error:", err)
	//	t.Fatal(err)
	//}
	//
	//// 查询用户的所有消息
	//usermessagelist := make([]models.GroupMessage, 0)
	//if err := adb.Ssql.Table("group_message").Desc("id").Find(&usermessagelist); err != nil {
	//	fmt.Println("查询所有消息error:", err)
	//	t.Fatal(err)
	//}
	////MessageTimeSort(usermessagelist, "desc")
	//
	//var unreadmsglist []models.GroupUnreadMessage
	//if err := adb.Ssql.Table("group_unread_message").Where("user_id = ?", u.ID).Find(&unreadmsglist); err != nil {
	//	t.Fatal(err)
	//}
	//unreadmsgmap := make(map[int]int, len(unreadmsglist)+1) //k:groupid   v:unreadnum
	//for _, UnreadMessage := range unreadmsglist {
	//	unreadmsgmap[UnreadMessage.GroupID] = UnreadMessage.UnreadNumber
	//}
	//
	//for _, g := range gur {
	//	var group models.Group                //群详情
	//	var messagelist []models.GroupMessage //群消息列表
	//
	//	//  根据群id查询群的详细信息
	//	exit, err := adb.Ssql.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
	//	if !exit {
	//		continue
	//	}
	//	if err != nil {
	//		fmt.Println("根据群id查询群的详细信息error:", err)
	//		t.Fatal(err)
	//	}
	//
	//	group.UnreadMessage = unreadmsgmap[g.GroupID] //放入未读消息数量
	//
	//	// 将该群聊的消息放入消息列表
	//	for _, m := range usermessagelist {
	//		// fmt.Printf("%+v-----%+v\n", m.GroupID, g.ID)
	//		if len(messagelist) >= 10 {
	//			break
	//		}
	//		if m.GroupID == g.GroupID {
	//			messagelist = append(messagelist, m)
	//		}
	//	}
	//	sort.Slice(messagelist, func(i, j int) bool { return messagelist[i].ID < (messagelist[j].ID) })
	//	//MessageTimeSort(messagelist, "asc")
	//
	//	var groupitem = models.GroupDetail{
	//		GroupInfo:   group,
	//		MessageList: messagelist,
	//	}
	//	grouplist = append(grouplist, groupitem)
	//}
	//
	//for _, v := range grouplist {
	//	fmt.Printf("%+v\n", v.GroupInfo)
	//	for _, message := range v.MessageList {
	//		fmt.Printf("%+v\n", message)
	//	}
	//}

	// err := adb.Engine.CreateTables(&models.User{})

	var grouplist []models.GroupDetail
	var gur models.GroupUserRelative

	err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).
		Iterate(&gur, func(i int, bean interface{}) error {
			gur := bean.(*models.GroupUserRelative)

			//群的信息
			var group models.Group
			_, err := adb.Ssql.Table("group").Where("id=?", gur.GroupID).Get(&group)
			if err != nil {
				fmt.Println("根据群id查询群的详细信息error:", err)
				return err
			}

			//群的消息
			tempdata := make([]models.GroupMessage, 0)
			if err := adb.Ssql.Table("group_message").Where("group_id=?", gur.GroupID).
				Asc("id").Limit(20).Find(&tempdata); err != nil {
				t.Fatal(err)
			}

			var groupitem = models.GroupDetail{
				GroupInfo:   group,
				MessageList: tempdata,
			}
			grouplist = append(grouplist, groupitem)

			return nil
		})
	if err != nil {
		t.Fatal(err)
	}

	//未读数量
	var unreadmsglist []models.GroupUnreadMessage
	if err := adb.Ssql.Table("group_unread_message").Where("user_id = ?", u.ID).Find(&unreadmsglist); err != nil {
		t.Fatal(err)
	}
	unreadmsgmap := make(map[int]int, len(unreadmsglist)+1) //k:groupid   v:unreadnum
	for _, UnreadMessage := range unreadmsglist {
		unreadmsgmap[UnreadMessage.GroupID] = UnreadMessage.UnreadNumber
	}

	for i, v := range grouplist {
		grouplist[i].GroupInfo.UnreadMessage = unreadmsgmap[v.GroupInfo.ID]
	}

}
