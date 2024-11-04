package adb

import (
	"HiChat/hichat-mq-service/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var Ssql *xorm.Engine

func InitMySQL() {
	//engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8mb4", config.MysqlUserName, config.MysqlPassword, config.MysqlDatabase))
	mysqlconf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.MysqlUserName, config.MysqlPassword, config.MysqlAddress, config.MysqlDatabase)
	fmt.Println(mysqlconf)
	engine, err := xorm.NewEngine("mysql", mysqlconf)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		panic(err.Error())
	}
	engine.SetMapper(names.GonicMapper{})
	// engine.CreateTables(&models.User{})
	Ssql = engine

}
