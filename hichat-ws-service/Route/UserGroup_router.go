package Route

import (
	"github.com/gin-gonic/gin"
	"go-websocket-server/service"
)

func InItUserGroupRouter(usergroup *gin.RouterGroup) {
	usergroup.POST("/creategroup", service.CreateGroup)         //创建群聊
	usergroup.POST("/handlejoingroup", service.HandleJoinGroup) //处理加入群聊
	usergroup.POST("/applyjoingroup", service.ApplyJoinGroup)   //申请加入群聊
	usergroup.POST("/exitgroup", service.ExitGroup)             //退出群聊
	//usergroup.POST("/searchGroup", service.SearchGroup)                           //搜索群聊
	usergroup.POST("/applyadduser", service.ApplyAddUser)                         //申请添加好友
	usergroup.POST("/handleadduser", service.HandleAddUser)                       //处理添加好友
	usergroup.POST("/deletefriend", service.DeleteUser)                           //删除好友
	usergroup.POST("/startusertouservideocall", service.StartUserToUserVideoCall) //开始视频通话
	usergroup.POST("/test", service.TestUserService)                              //测试
}
