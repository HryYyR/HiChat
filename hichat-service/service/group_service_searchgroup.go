package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type searchgroupinfo struct {
	Searchstr string
}

// SearchGroup 搜索群聊
func SearchGroup(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)
	var rawdata searchgroupinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法访问", nil)
		return
	}
	err = json.Unmarshal(rawbyte, &rawdata)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	lenOfTrimStr := len(strings.TrimSpace(rawdata.Searchstr))

	if lenOfTrimStr == 0 {
		util.H(c, http.StatusBadRequest, "关键词不能为空", nil)
		return
	}

	if lenOfTrimStr > 50 {
		util.H(c, http.StatusBadRequest, "关键词超字数上限(50字)", nil)
		return
	}

	var searchint int
	v, err := strconv.Atoi(rawdata.Searchstr)
	if err == nil {
		searchint = v
	}

	grouplist := make([]models.Group, 0)
	err = adb.SqlStruct.Conn.Table("group").Where("group_name like ? or id=?", rawdata.Searchstr+"%", searchint).Where("creater_id !=?", userdata.ID).Find(&grouplist)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "搜索失败", err)
		return
	}

	responsedata := &[]models.Group{}
	for _, group := range grouplist {
		count, err := group.GetMemberCount()
		if err != nil {
			log.Println(err)
		}
		group.MemberCount = count
		*responsedata = append(*responsedata, group)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "搜索成功",
		"grouplist": responsedata,
	})

}
