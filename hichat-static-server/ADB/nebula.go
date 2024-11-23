package adb

import (
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"hichat_static_server/config"
	"log"
)

const (
	useHTTP2 = false
)

type nebulaStruct struct {
	Conn   *nebula.SessionPool
	Status int
}

var NebulaInstance = &nebulaStruct{
	Conn:   nil,
	Status: 0,
}

func (n *nebulaStruct) createNebula() {

	hostAddress := nebula.HostAddress{Host: config.NebulaAddress, Port: config.NebulaPort}

	// Create configs for session pool
	config, err := nebula.NewSessionPoolConf(
		config.NebulaUserName,
		config.NebulaPassWord,
		[]nebula.HostAddress{hostAddress},
		"HiChat",
		nebula.WithHTTP2(useHTTP2),
	)
	if err != nil {
		log.Panic(fmt.Sprintf("创建 Nebula 配置文件失败, %s\n", err.Error()))
	}

	sessionPool, err := nebula.NewSessionPool(*config, nebula.DefaultLogger{})
	if err != nil {
		log.Panic(fmt.Sprintf("初始化 Nebula失败, %s\n", err.Error()))
	}

	n.Conn = sessionPool

}

func (n *nebulaStruct) GetNebulaSession() *nebula.SessionPool {
	if n.Conn == nil {
		n.createNebula()
	}
	return n.Conn
}

func (n *nebulaStruct) CloseNebula() {
	n.Conn.Close()
	n.Conn = nil
}
