package test

import (
	adb "hichat-file-service/ADB"
	"hichat-file-service/models"
	"testing"
)

func TestXormTest(t *testing.T) {
	adb.InitMySQL()
	// err := adb.Engine.CreateTables(&models.User{})
	err := adb.Ssql.Sync2(new(models.UsersFile))
	if err != nil {
		t.Fatal(err)
	}
}
