package test

import (
	"fmt"
	"github.com/goinggo/mapstructure"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"hichat_static_server/tool"
	"strconv"
	"testing"
	"time"
)

func TestXormsql(t *testing.T) {
	adb.InitMySQL()
	tstr := "2006-01-02 15:04:05"
	fmt.Println(time.Parse(time.DateTime, time.Now().Format(tstr)))
}

func TestSave(t *testing.T) {
	adb.InitMySQL()
	adb.InitRedis()

	var g models.Group
	get, err := adb.Ssql.Table("group").Where("id=?", 1).Get(&g)
	if err != nil {
		t.Fatal(err)
		return
	}
	if !get {
		t.Fatalf("select error")
		return
	}

	key := fmt.Sprintf("group%d", g.ID)
	_, err = adb.Rediss.HMSet(key, map[string]interface{}{
		"ID":          g.ID,
		"GroupName":   g.GroupName,
		"Avatar":      g.Avatar,
		"CreaterID":   g.CreaterID,
		"CreaterName": g.CreaterName,
		"Grade":       g.Grade,
		"MemberCount": g.MemberCount,
		"CreatedAt":   tool.FormatTime(g.CreatedAt),
		"DeletedAt":   tool.FormatTime(g.DeletedAt),
	}).Result()
	if err != nil {
		t.Fatal(err)
		return
	}
	adb.Rediss.Expire(key, time.Hour*360)

}

func TestGetInfo(t *testing.T) {
	adb.InitMySQL()
	adb.InitRedis()

	g := &models.Group{
		ID: 1,
	}

	var groupinfo models.Group
	key := fmt.Sprintf("group%d", g.ID)
	//从redis获取数据
	var gdata = adb.Rediss.HGetAll(key).Val()
	if len(gdata) != 0 {
		fmt.Println("走redis")
		_ = mapstructure.Decode(gdata, &groupinfo)
		groupinfo.ID, _ = strconv.Atoi(gdata["ID"])
		groupinfo.CreaterID, _ = strconv.Atoi(gdata["CreaterID"])
		groupinfo.Grade, _ = strconv.Atoi(gdata["Grade"])
		groupinfo.MemberCount, _ = strconv.Atoi(gdata["MemberCount"])
	}
	fmt.Println("走mysql")
	exit, err := adb.Ssql.Table("group").Where("id =?", g.ID).Get(&groupinfo)
	if !exit {
	}
	if err != nil {
		fmt.Println("mysql查询失败", err)
		t.Fatal(err)
	}

	err = groupinfo.SaveToRedis()
	if err != nil {
		fmt.Println("保存到redis失败", err)
		t.Fatal(err)
	} else {
		fmt.Println("保存到redis成功", err)
	}

	fmt.Printf("%#v", groupinfo)

}
