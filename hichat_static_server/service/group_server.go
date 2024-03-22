package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	adb "hichat_static_server/ADB"
	"hichat_static_server/common"
	"hichat_static_server/models"
	"hichat_static_server/util"
	"net/http"
	"strconv"
	"strings"
)

func SearchGroup(c *gin.Context) {

}

type GetGroupMessageListRequest struct {
	Groupid    int `json:"groupid,omitempty"`
	Currentnum int `json:"currentnum,omitempty"`
}

func GetGroupMessageList(c *gin.Context) {
	//ud, _ := c.Get("userdata")
	//userdata := ud.(*models.UserClaim)

	databyte, _ := c.GetRawData()
	var data GetGroupMessageListRequest
	err := json.Unmarshal(databyte, &data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	group := models.Group{
		ID: data.Groupid,
	}
	fmt.Println(data.Groupid, data.Currentnum)
	grouplist := make([]models.GroupMessage, 0)
	err = group.GetMessageList(&grouplist, data.Currentnum)
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusInternalServerError, "查询消息失败", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": grouplist,
	})
}

type GetGroupMemberListRequest struct {
	GroupId int
}

func GetGroupMemberList(c *gin.Context) {

	databyte, _ := c.GetRawData()
	var data GetGroupMemberListRequest
	err := json.Unmarshal(databyte, &data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	var memberlist []models.UserShowData
	var uidlist []string

	uidliststr := adb.Rediss.HGet("GroupToUserMap", strconv.Itoa(data.GroupId)).Val()

	if len(uidliststr) == 0 {
		err := adb.Ssql.Table("group_user_relative").Cols("user_id").Where("group_id=?", data.GroupId).Find(&uidlist)
		if err != nil {
			fmt.Println(err)
			util.H(c, http.StatusInternalServerError, "服务器错误", nil)
			return
		}
	}

	uidlist = strings.Split(uidliststr, ",")
	for _, uid := range uidlist {
		var tempu models.UserShowData
		udata := adb.Rediss.HGetAll(uid).Val()
		if len(udata) == 0 {
			has, err2 := adb.Ssql.Table("users").Where("id=?", uid).Get(&tempu)
			if err2 != nil {
				continue
			}
			if !has {
				continue
			}
		} else {
			_ = mapstructure.Decode(udata, &tempu)
			tempu.ID, _ = strconv.Atoi(udata["ID"])
			tempu.Age, _ = strconv.Atoi(udata["Age"])
			tempu.CreatedAt, _ = common.ParseTime(udata["CreatedAt"])
		}
		memberlist = append(memberlist, tempu)

	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "查询成功",
		"data": memberlist,
	})

}
