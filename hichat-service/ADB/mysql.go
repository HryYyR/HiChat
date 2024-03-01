package adb

import (
	"fmt"
	"go-websocket-server/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var Ssql *xorm.Engine

func InitMySQL() {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8mb4", config.MysqlUserName, config.MysqlPassword, config.MysqlDatabase))
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		panic(err.Error())
	}
	fmt.Println("init mysql success")
	engine.SetMapper(names.GonicMapper{})
	// engine.CreateTables(&models.User{})
	Ssql = engine

}
