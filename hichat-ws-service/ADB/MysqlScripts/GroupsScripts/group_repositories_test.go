package GroupScripts

import (
	adb "go-websocket-server/ADB"
	"testing"
	"xorm.io/xorm"
)

func Test_groupRepository_ByGroupNameCheckGroupIsExist(t *testing.T) {
	type fields struct {
		db *xorm.Engine
	}
	type args struct {
		groupName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantBeUse bool
		wantErr   bool
	}{
		{
			name: "test1",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				groupName: "test",
			},
			wantBeUse: true,
			wantErr:   false,
		},
		{
			name: "test2",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				groupName: "none",
			},
			wantBeUse: false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &groupRepository{
				db: tt.fields.db,
			}
			gotBeUse, err := r.ByGroupNameCheckGroupIsExist(tt.args.groupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ByGroupNameCheckGroupIsExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBeUse != tt.wantBeUse {
				t.Errorf("ByGroupNameCheckGroupIsExist() gotBeUse = %v, want %v", gotBeUse, tt.wantBeUse)
			}
		})
	}
}

func Test_groupRepository_GetUserApplyJoinGroupCount(t *testing.T) {
	type fields struct {
		db *xorm.Engine
	}
	type args struct {
		userid int
		Status int
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantApplycount int64
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid: 1015,
				Status: 1,
			},
			wantApplycount: 0,
			wantErr:        false,
		},
		{
			name: "test2",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid: 1016,
				Status: 1,
			},
			wantApplycount: 1,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &groupRepository{
				db: tt.fields.db,
			}
			gotApplycount, err := r.GetUserApplyJoinGroupCount(tt.args.userid, tt.args.Status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserApplyJoinGroupCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotApplycount != tt.wantApplycount {
				t.Errorf("GetUserApplyJoinGroupCount() gotApplycount = %v, want %v", gotApplycount, tt.wantApplycount)
			}
		})
	}
}

func Test_groupRepository_GetUserApplyJoinGroupCount1(t *testing.T) {
	type fields struct {
		db *xorm.Engine
	}
	type args struct {
		userid int
		Status int
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantApplycount int64
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid: 1015,
				Status: 1,
			},
			wantApplycount: 0,
			wantErr:        false,
		},
		{
			name: "test2",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid: 1016,
				Status: 1,
			},
			wantApplycount: 1,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &groupRepository{
				db: tt.fields.db,
			}
			gotApplycount, err := r.GetUserApplyJoinGroupCount(tt.args.userid, tt.args.Status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserApplyJoinGroupCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotApplycount != tt.wantApplycount {
				t.Errorf("GetUserApplyJoinGroupCount() gotApplycount = %v, want %v", gotApplycount, tt.wantApplycount)
			}
		})
	}
}

func Test_groupRepository_CheckUserIsExistInGroup(t *testing.T) {
	type fields struct {
		db *xorm.Engine
	}
	type args struct {
		userid  int
		groupId int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantExist bool
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid:  1015,
				groupId: 1,
			},
			wantExist: false,
			wantErr:   false,
		},
		{
			name: "test2",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid:  1016,
				groupId: 238,
			},
			wantExist: true,
			wantErr:   false,
		},
		{
			name: "test2",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				userid:  1015,
				groupId: 238,
			},
			wantExist: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &groupRepository{
				db: tt.fields.db,
			}
			gotExist, err := r.CheckUserIsExistInGroup(tt.args.userid, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserIsExistInGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExist != tt.wantExist {
				t.Errorf("CheckUserIsExistInGroup() gotExist = %v, want %v", gotExist, tt.wantExist)
			}
		})
	}
}

func Test_groupRepository_UpdateApplyJoinGroupStatus(t *testing.T) {
	type fields struct {
		db *xorm.Engine
	}
	type args struct {
		applyId int
		Status  int
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantUpdateCount int64
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				applyId: 448,
				Status:  0,
			},
			wantUpdateCount: 1,
			wantErr:         false,
		},
		{
			name: "test2",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				applyId: 448,
				Status:  -1,
			},
			wantUpdateCount: 1,
			wantErr:         false,
		},
		{
			name: "test3",
			fields: fields{
				db: adb.GetMySQLConn(),
			},
			args: args{
				applyId: 448,
				Status:  1,
			},
			wantUpdateCount: 1,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &groupRepository{
				db: tt.fields.db,
			}
			gotUpdateCount, err := r.UpdateApplyJoinGroupStatus(tt.args.applyId, tt.args.Status)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateApplyJoinGroupStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUpdateCount != tt.wantUpdateCount {
				t.Errorf("UpdateApplyJoinGroupStatus() gotUpdateCount = %v, want %v", gotUpdateCount, tt.wantUpdateCount)
			}
		})
	}
}
