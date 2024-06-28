package adb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-websocket-server/config"
	"log"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/names"
)

type Sql struct {
	Conn   *xorm.Engine
	Status int
}

func (s *Sql) CreateConn() error {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8mb4", config.MysqlUserName, config.MysqlPassword, config.MysqlDatabase))
	if err != nil {
		return err
	}
	engine.SetMaxIdleConns(config.MysqlMaxIdleConns)
	engine.SetMaxOpenConns(config.MysqlMaxOpenConns)
	engine.SetMapper(names.GonicMapper{})
	engine.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 1000)) //开启缓存,缓存struct的记录数为1000条
	s.Status = 1
	s.Conn = engine
	return nil
}

func (s *Sql) CloseConn() {
	s.Status = 0
}

var SqlStruct *Sql

func InitMySQL() {
	fmt.Println("初始化 Mysql 成功")
	// engine.CreateTables(&models.User{})
	SqlStruct = &Sql{
		Conn:   nil,
		Status: 0,
	}
	err := SqlStruct.CreateConn()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

}
