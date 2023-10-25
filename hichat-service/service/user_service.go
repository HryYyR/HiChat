package service

import (
	"encoding/json"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 申请添加好友
func ApplyAddUser(c *gin.Context) {
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "getrawdata failed",
		})
		return
	}
	var data models.ApplyAddUser
	err = json.Unmarshal(rawbyte, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "json Unmarshal failed",
		})
		return
	}

	// fmt.Printf("%+v", data)
	exit, err := adb.Ssql.Table("apply_add_user").
		Where("pre_apply_user_id=?  and apply_user_id=?  and handle_status=0",
			data.PreApplyUserID, data.ApplyUserID).Exist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "查询申请信息失败!",
			"error": err.Error(),
		})
		return
	}
	if exit {
		c.JSON(http.StatusOK, gin.H{
			"msg": "申请已存在!",
		})
		return
	}
	exit, err = adb.Ssql.Table("apply_add_user").
		Where("pre_apply_user_id=?  and apply_user_id=?  and handle_status=0",
			data.ApplyUserID, data.PreApplyUserID).Exist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "查询申请信息失败!",
			"error": err.Error(),
		})
		return
	}
	if exit {
		c.JSON(http.StatusOK, gin.H{
			"msg": "申请已存在!",
		})
		return
	}

	if _, err = adb.Ssql.Table("apply_add_user").Insert(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "申请添加好友失败!",
		})
		return
	}
	msg := models.Message{
		UserID:   data.ApplyUserID,
		UserName: data.ApplyUserName,
		GroupID:  0,
		MsgType:  config.MsgTypeRefreshFriendNotice,
	}
	msgbyte, _ := json.Marshal(msg)
	models.ServiceCenter.Clients[data.PreApplyUserID].Send <- msgbyte
	models.ServiceCenter.Clients[data.ApplyUserID].Send <- msgbyte

	c.JSON(http.StatusOK, gin.H{
		"msg": "申请成功!",
	})
}

type handleadduserinfo struct {
	ApplyID      int `json:"ApplyID"`
	HandleStatus int `json:"HandleStatus"`
}

// 处理添加好友
func HandleAddUser(c *gin.Context) {
	data := new(handleadduserinfo)
	if err := util.HandleJsonArgument(c, data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "参数有误",
			"error": err.Error(),
		})
		return
	}
	var applyadduserdata models.ApplyAddUser
	exit, err := adb.Ssql.Table("apply_add_user").ID(data.ApplyID).Get(&applyadduserdata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "查询申请失败",
			"error": err.Error(),
		})
		return
	}
	if !exit {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "申请不存在",
		})
		return
	}

	session := adb.Ssql.NewSession()
	// 更新申请
	if _, err = session.Table("apply_add_user").ID(data.ApplyID).Update(&models.ApplyAddUser{HandleStatus: data.HandleStatus}); err != nil {
		session.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "更新申请失败",
			"error": err.Error(),
		})
		return
	}

	if data.HandleStatus == -1 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "拒绝成功",
		})
		return
	} else {
		insertdata := &models.UserUserRelative{
			PreUserID:    applyadduserdata.PreApplyUserID,
			PreUserName:  applyadduserdata.PreApplyUserName,
			BackUserID:   applyadduserdata.ApplyUserID,
			BackUserName: applyadduserdata.ApplyUserName,
		}
		// 插入关系
		if _, err = session.Table("user_user_relative").Insert(&insertdata); err != nil {
			session.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "处理好友请求失败",
				"error": err.Error(),
			})
			return
		}
		session.Commit()

		refreshmsg := models.Message{
			UserID:   applyadduserdata.ApplyUserID,
			UserName: applyadduserdata.ApplyUserName,
			GroupID:  0,
			MsgType:  config.MsgTypeRefreshFriend,
		}
		msgbyte, _ := json.Marshal(refreshmsg)
		models.ServiceCenter.Clients[applyadduserdata.PreApplyUserID].Send <- msgbyte
		models.ServiceCenter.Clients[applyadduserdata.ApplyUserID].Send <- msgbyte

		c.JSON(http.StatusOK, gin.H{
			"msg": "同意成功",
		})
		return
	}

}
