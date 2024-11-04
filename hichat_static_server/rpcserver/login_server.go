package rpcserver

import (
	"context"
	"errors"
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/models"
	"hichat_static_server/proto"
	"hichat_static_server/util"
	"time"
)

type Server struct {
	proto.UnimplementedLoginServer
}

func NewServer() *Server {
	return &Server{}
}

// UserLogin todo 废弃
func (s *Server) UserLogin(ctx context.Context, in *proto.UserData) (*proto.LoginResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, errors.New("内容格式有误,请检查后重试")
	}

	// 查询用户是否存在
	var userdata models.Users
	hasuser, err := adb.Ssql.Table(&models.Users{}).Where("user_name = ?", in.Username).Get(&userdata)
	if err != nil {
		return nil, errors.New("查询用户信息失败")
	}
	if !hasuser {
		return nil, errors.New("用户不存在")
	}

	if util.Md5(in.Password+userdata.Salt) != userdata.Password {
		return nil, errors.New("密码错误")
	}

	token, err := util.GenerateToken(userdata.ID, userdata.UUID, userdata.UserName, 24*time.Hour)
	if err != nil {
		return nil, errors.New("生成签名失败")
	}

	//获取登录信息
	ResponseUserData := new(models.ResponseUserData)
	err = userdata.Login(ResponseUserData)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("登录失败")
	}

	ProtoStruct := ResponseUserData.ResponseUserDataToProto()

	res := new(proto.LoginResponse)
	res.Token = token
	res.Msg = "登陆成功"
	res.Userdata = ProtoStruct

	return res, nil
}
