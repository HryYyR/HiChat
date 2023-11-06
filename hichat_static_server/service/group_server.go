package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"hichat_static_server/models"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "信息输入有误!",
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询消息失败!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": grouplist,
	})

}
