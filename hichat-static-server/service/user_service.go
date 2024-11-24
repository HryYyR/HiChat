package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	adb "hichat_static_server/ADB"
	"hichat_static_server/common"
	"hichat_static_server/config"
	"hichat_static_server/models"
	"hichat_static_server/util"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type edituserdataform struct {
	City      string `json:"city"`
	Age       string `json:"age"`
	Introduce string `json:"introduce"`
	Avatar    string `json:"avatar"`
}

func EditUserData(c *gin.Context) {
	var mlock sync.Mutex
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	databyte, _ := c.GetRawData()
	var data edituserdataform
	err := json.Unmarshal(databyte, &data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}
	// fmt.Printf("%+v\n", data)
	age, err := strconv.Atoi(data.Age)
	if err != nil {
		util.H(c, http.StatusBadRequest, "信息输入有误", nil)
		return
	}

	if age < 0 || age > 200 {
		util.H(c, http.StatusBadRequest, "信息输入有误", nil)
		return
	}

	mlock.Lock()
	defer mlock.Unlock()

	session := adb.Ssql.NewSession()
	defer session.Close()

	if _, err := session.Table("users").Where("id=?", userdata.ID).Update(&models.Users{
		City:      data.City,
		Age:       age,
		Introduce: data.Introduce,
		Avatar:    data.Avatar,
	}); err != nil {
		log.Println(err)
		session.Rollback()
		util.H(c, http.StatusInternalServerError, "修改失败", nil)
		return
	}

	//TODO 无法保证redis和nebula同时修改成功
	var completeuserdata models.Users
	tempuserdata := models.Users{ID: userdata.ID}
	err = tempuserdata.GetUserData(&completeuserdata)
	fmt.Printf("%+v\n", completeuserdata)
	if err != nil {
		session.Rollback()
		log.Println(err.Error())
		util.H(c, http.StatusInternalServerError, "修改失败", nil)
		return
	}

	completeuserdata.Age = age
	completeuserdata.Introduce = data.Introduce
	completeuserdata.City = data.City
	completeuserdata.Avatar = data.Avatar

	//todo nebula的头像没修改
	if config.IsStartNebula {
		err = completeuserdata.UpdateUser2Nebula()
		if err != nil {
			session.Rollback()
			util.H(c, http.StatusInternalServerError, "修改失败", nil)
			return
		}
	}

	//直接删除缓存
	err = adb.Rediss.Del(strconv.Itoa(userdata.ID)).Err()
	if err != nil {
		session.Rollback()
		log.Println(err.Error())
		util.H(c, http.StatusInternalServerError, "修改失败", nil)
		return
	}
	session.Commit()

	util.H(c, http.StatusOK, "修改成功", nil)
}

func GetUserData(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	var targetuserdata models.UserShowData

	result, err := adb.Rediss.HGetAll(strconv.Itoa(data.ID)).Result()
	if err == nil && len(result) != 0 {
		fmt.Println("redis")
		_ = mapstructure.Decode(result, &targetuserdata)
		targetuserdata.ID, _ = strconv.Atoi(result["ID"])
		targetuserdata.Age, _ = strconv.Atoi(result["Age"])
		targetuserdata.CreatedAt, _ = common.ParseTime(result["CreatedAt"])
		c.JSON(http.StatusOK, gin.H{
			"data": targetuserdata,
		})
		return
	}

	var userdata models.Users
	exit, err := adb.Ssql.Omit("Password,Salt,Grade,IP").Table("users").ID(data.ID).Get(&userdata)
	if !exit {
		util.H(c, http.StatusBadRequest, "用户不存在!", nil)
		return
	}
	if err != nil {
		util.H(c, http.StatusInternalServerError, "查询用户信息失败!", nil)
		return
	}
	targetuserdata = models.UserShowData{
		ID:        userdata.ID,
		UserName:  userdata.UserName,
		NikeName:  userdata.NikeName,
		Email:     userdata.Email,
		Avatar:    userdata.Avatar,
		City:      userdata.City,
		Age:       userdata.Age,
		Introduce: userdata.Introduce,
		CreatedAt: userdata.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": targetuserdata,
	})
}

// GetUserGroupList 获取用户的群聊列表
func GetUserGroupList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "JSON格式不正确!", err)
		return
	}
	grouplist := make([]models.GroupDetail, 0)
	err = data.GetUserGroupList(&grouplist)
	// fmt.Println("消息长度为:", len(grouplist[0].MessageList))
	if err != nil {
		util.H(c, http.StatusInternalServerError, "获取失败!", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功!",
		"data": grouplist,
	})
}

// GetUserFriendList 获取用户的好友列表
func GetUserFriendList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "JSON格式不正确!", err)
		return
	}

	friendlist := make([]models.FriendResponse, 0)
	if err = data.GetFriendListAndMEssage(&friendlist); err != nil {
		util.H(c, http.StatusInternalServerError, "获取用户的好友列表失败!", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的好友列表成功!",
		"data": friendlist,
	})
}

// GetUserApplyJoinGroupList //获取用户的群聊通知列表
func GetUserApplyJoinGroupList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "JSON格式不正确!", err)
		return
	}

	applyjoingrouplist := make([]models.ApplyJoinGroupResponse, 0)
	if err = data.GetApplyMsgList(&applyjoingrouplist); err != nil {
		util.H(c, http.StatusInternalServerError, "获取用户的群聊通知列表失败!", err)
	}

	fmt.Println(applyjoingrouplist)
	util.GroupTimeSort(applyjoingrouplist, "desc")

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的群聊通知列表成功!",
		"data": applyjoingrouplist,
	})
}

// GetUserApplyAddFriendList //获取用户的好友申请列表
func GetUserApplyAddFriendList(c *gin.Context) {
	data := new(models.Users)
	err := util.HandleJsonArgument(c, data)
	if err != nil {
		util.H(c, http.StatusBadRequest, "JSON格式不正确!", err)
		return
	}

	applyaddfriendlist := make([]models.ApplyAddUser, 0)
	if err = data.GetApplyAddUserList(&applyaddfriendlist); err != nil {
		util.H(c, http.StatusInternalServerError, "获取用户的好友申请列表失败!", err)
	}
	util.UserTimeSort(applyaddfriendlist, "desc")

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取用户的好友申请列表成功!",
		"data": applyaddfriendlist,
	})
}

type searchfriendinfo struct {
	Searchstr string
}

// SearchUser 搜索用户
func SearchUser(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	uudata := models.Users{
		ID:       userdata.ID,
		UserName: userdata.UserName,
	}

	var uufriendlist []models.Friend //发起人的好友列表
	if err := uudata.GetFriendList(&uufriendlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查询好友列表失败!",
		})
		return
	}
	uufriendmap := make(map[int]int)
	for _, f := range uufriendlist {
		uufriendmap[int(f.Id)] = int(f.Id)
	}

	var data searchfriendinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = json.Unmarshal(rawbyte, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if len(data.Searchstr) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "搜索成功!",
			"data": []models.Users{},
		})
		return
	}

	//todo  前端其实只需要id列表，点击申请才会去获取详细数据
	friendlist := make([]models.Users, 0)
	err = adb.Ssql.Table("users").Omit("ip,password,salt,grade,uuid").Where("user_name LIKE ?  and user_name !=?",
		data.Searchstr+"%", userdata.UserName).Find(&friendlist)
	if err != nil {
		fmt.Println(err)
		util.H(c, http.StatusInternalServerError, "查询好友列表失败!", err)
		return
	}

	resultfriendlist := make([]models.Users, 0)
	// 筛除已经成为的好友
	for _, ff := range friendlist {
		if _, ok := uufriendmap[ff.ID]; !ok {
			resultfriendlist = append(resultfriendlist, ff)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "搜索成功!",
		"data": resultfriendlist,
	})
}

type getusermessagelistinfo struct {
	Targetuserid int
	Currentnum   int
}

// GetUserMessageList 获取用户间消息列表
func GetUserMessageList(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)

	var requestdata getusermessagelistinfo
	rawbyte, err := c.GetRawData()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	err = json.Unmarshal(rawbyte, &requestdata)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	user := &models.UserUserRelative{
		PreUserID:  userdata.ID,
		BackUserID: requestdata.Targetuserid,
	}

	fmt.Println(requestdata.Currentnum)

	grouplist := make([]models.UserMessageItem, 0)
	err = user.GetUserMessageList(&grouplist, requestdata.Currentnum)
	if err != nil {
		util.H(c, http.StatusInternalServerError, "获取失败", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": grouplist,
	})

}

type AiMessageRequest struct {
	Msg     string `json:"msg"`
	MsgType int    `json:"msgtype"`
}
type AiMessageList struct {
	List map[int]*chains.LLMChain
	Llm  *ollama.LLM
}

var AiMsgList = &AiMessageList{List: make(map[int]*chains.LLMChain), Llm: nil}

// AiMessage Ai问答
func AiMessage(c *gin.Context) {
	ud, _ := c.Get("userdata")
	userdata := ud.(*models.UserClaim)
	var requestdata AiMessageRequest
	rawbyte, err := c.GetRawData()
	if err != nil {
		util.H(c, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	err = json.Unmarshal(rawbyte, &requestdata)
	if err != nil {
		util.H(c, http.StatusBadRequest, "非法格式", nil)
		return
	}

	var mu sync.Mutex

	if AiMsgList.Llm == nil {
		llm, err := ollama.New(ollama.WithModel("llama3.1"))
		if err != nil {
			util.H(c, http.StatusInternalServerError, "获取失败", err)
			return
		}
		AiMsgList.Llm = llm
	}

	var response string
	if userchains, ok := AiMsgList.List[userdata.ID]; !ok {
		conversationWindowBuffer := memory.NewConversationWindowBuffer(30)
		llmChain := chains.NewConversation(AiMsgList.Llm, conversationWindowBuffer)
		mu.Lock()
		AiMsgList.List[userdata.ID] = &llmChain
		mu.Unlock()
		response, err = chains.Run(context.Background(), llmChain, requestdata.Msg)
		if err != nil {
			util.H(c, http.StatusInternalServerError, "获取失败", err)
			return
		}
	} else {
		response, err = chains.Run(context.Background(), userchains, requestdata.Msg)
		if err != nil {
			util.H(c, http.StatusInternalServerError, "获取失败", err)
			return
		}
	}

	util.H(c, http.StatusOK, response, nil)

}
