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

var SqlStruct = &Sql{
	Conn:   nil,
	Status: 0,
}

// CreateConn 创建mysql连接
func (s *Sql) CreateConn() error {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8", config.MysqlUserName, config.MysqlPassword, config.MysqlAddress)
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log.Println("连接数据库失败: ", err)
		return err
	}
	_, _ = engine.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4", config.MysqlDatabase))

	mysqlconf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.MysqlUserName, config.MysqlPassword, config.MysqlAddress, config.MysqlDatabase)
	engine, err = xorm.NewEngine("mysql", mysqlconf)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("初始化 MySQL :", mysqlconf)
	engine.SetMaxIdleConns(config.MysqlMaxIdleConns)
	engine.SetMaxOpenConns(config.MysqlMaxOpenConns)
	engine.SetMapper(names.GonicMapper{})
	engine.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 1000)) //开启缓存,缓存struct的记录数为1000条

	s.Status = 1
	s.Conn = engine
	return nil
}

// CloseConn 关闭mysql连接
func (s *Sql) CloseConn() {
	s.Conn = nil
	s.Status = 0
}

// GetMySQLConn 获取mysql连接,如果连接不可用就创建连接
func GetMySQLConn() *xorm.Engine {
	if SqlStruct.Conn == nil || SqlStruct.Status == 0 {
		err := SqlStruct.CreateConn()
		if err != nil {
			log.Panic(err)
		}
	}
	return SqlStruct.Conn
}
