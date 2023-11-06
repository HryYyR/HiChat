package test

import (
	adb "hichat_static_server/ADB"
	"testing"
)

func TestXormsql(t *testing.T) {
	adb.InitMySQL()
	// err := adb.Engine.CreateTables(&models.User{})
	adb.Ssql.Table("users")
}
