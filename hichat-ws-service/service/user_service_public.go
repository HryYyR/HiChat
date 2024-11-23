package service

import (
	"github.com/gin-gonic/gin"
	"go-websocket-server/util"
	"net/http"
)

// TestUserService 发起用户与用户的远程视频
func TestUserService(c *gin.Context) {
	util.H(c, http.StatusOK, "ok", nil)
}
