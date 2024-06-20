package adb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-websocket-server/config"
	"log"
	"xorm.io/xorm"
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
	s.Status = 1
	s.Conn = engine
	return nil
}

func (s *Sql) CloseConn() error {
	err := s.Conn.Close()
	if err != nil {
		return err
	}
	s.Status = 0
	return nil
}

var SqlStruct *Sql

func InitMySQL() {
	fmt.Println("init mysql success")
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
