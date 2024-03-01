package test

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"testing"
	"time"
)

func TestXormsql(t *testing.T) {
	adb.InitMySQL()
	tstr := "2006-01-02 15:04:05"
	fmt.Println(time.Parse(time.DateTime, time.Now().Format(tstr)))
}
