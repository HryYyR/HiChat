package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	adb "hichat_static_server/ADB"
	"time"
)

type UserUserRelative struct {
	ID           int    `xorm:"pk autoincr notnull index"`
	PreUserID    int    `xorm:"notnull unique(BackUserID)"` //被申请人
	PreUserName  string //被申请人
	BackUserID   int    `xorm:"notnull"`
	BackUserName string
	CreatedAt    time.Time `xorm:"created"`
	DeletedAt    time.Time `xorm:"deleted"`
	UpdatedAt    time.Time `xorm:"updated"`
}

func (r *UserUserRelative) GetUserMessageList(messagelist *[]UserMessageItem, currentnum int) error {
	var key string
	if r.PreUserID > r.BackUserID {
		key = fmt.Sprintf("%d%d", r.PreUserID, r.BackUserID)
	} else {
		key = fmt.Sprintf("%d%d", r.BackUserID, r.PreUserID)
	}

	lLen := adb.Rediss.LLen(key).Val()
	if lLen == 0 {
		return nil
	}
	if (int(lLen) - currentnum) <= 0 {
		return nil
	}
	start := lLen - int64(currentnum) - 20
	if start <= 0 {
		start = 0
	}
	stop := lLen - int64(currentnum) - 1
	if stop <= 0 {
		stop = 0
	}

	result, err := adb.Rediss.LRange(key, start, stop).Result()
	if err != nil {
		return err
	}

	for _, msgstr := range result {
		var msgstruct UserMessageItem
		bufferString := bytes.NewBufferString(msgstr).Bytes()
		err := json.Unmarshal(bufferString, &msgstruct)
		if err != nil {
			continue
		}
		//fmt.Println(msgstruct)
		*messagelist = append(*messagelist, msgstruct)
	}
	return nil
}
