package adb

import (
	"fmt"
	"hichat_static_server/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var Ssql *xorm.Engine

func InitMySQL() {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.MysqlUserName, config.MysqlPassword, config.MysqlAddress, config.MysqlDatabase))
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println("init mysql success")
	//engine.ShowSQL(true)
	engine.SetMapper(names.GonicMapper{})
	// engine.CreateTables(&models.User{})
	Ssql = engine

}
