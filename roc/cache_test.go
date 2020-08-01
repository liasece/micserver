// Package roc 每个 micserver 进程会持有一份 ROC 缓存单例，维护了所有已知的ROC对象所处的位置
package roc

import (
	"reflect"
	"testing"
)

const TestObjType1 ObjType = "TestObjType1"
const TestObjType2 ObjType = "TestObjType2"

const TestObjID1 string = "TestObjID1"
const TestObjID2 string = "TestObjID2"

const TestModuleID1 string = "TestModuleID1"
const TestModuleID2 string = "TestModuleID2"

func TestCache_catchGetTypeMust(t *testing.T) {
	type args struct {
		objType ObjType
	}
	tests := []struct {
		name string
		c    *Cache
		args args
	}{
		{
			name: "OK",
			c:    &Cache{},
			args: args{
				objType: TestObjType1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.catchGetTypeMust(tt.args.objType); got == nil {
				t.Errorf("Cache.catchGetTypeMust() = %v, want !nil", got)
			}
		})
	}
}

func TestCache_catchGetServerMust(t *testing.T) {
	c1 := &Cache{}
	c1.Set(TestObjType1, TestObjID1, TestModuleID1)

	type args struct {
		moduleID string
	}
	tests := []struct {
		name string
		c    *Cache
		args args
		want *catchServerInfo
	}{
		{
			name: "OK",
			c:    c1,
			args: args{
				moduleID: TestModuleID1,
			},
			want: &catchServerInfo{moduleID: TestModuleID1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.catchGetServerMust(tt.args.moduleID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.catchGetServerMust() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	type args struct {
		objType  ObjType
		objID    string
		moduleID string
	}
	tests := []struct {
		name    string
		c       *Cache
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objID:    TestObjID1,
				moduleID: TestModuleID1,
			},
			wantErr: false,
		},
		{
			name: "error",
			c:    &Cache{},
			args: args{
				objType:  "",
				objID:    TestObjID1,
				moduleID: TestModuleID1,
			},
			wantErr: true,
		},
		{
			name: "error",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objID:    "",
				moduleID: TestModuleID1,
			},
			wantErr: true,
		},
		{
			name: "error",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objID:    TestObjID1,
				moduleID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Set(tt.args.objType, tt.args.objID, tt.args.moduleID); (err != nil) != tt.wantErr {
				t.Errorf("Cache.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_SetM(t *testing.T) {
	type args struct {
		objType  ObjType
		objIDs   []string
		moduleID string
	}
	tests := []struct {
		name    string
		c       *Cache
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{TestObjID1},
				moduleID: TestModuleID1,
			},
			wantErr: false,
		},
		{
			name: "OK",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objIDs:   nil,
				moduleID: TestModuleID1,
			},
			wantErr: false,
		},
		{
			name: "error",
			c:    &Cache{},
			args: args{
				objType:  "",
				objIDs:   []string{TestObjID1},
				moduleID: TestModuleID1,
			},
			wantErr: true,
		},
		{
			name: "error",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{""},
				moduleID: TestModuleID1,
			},
			wantErr: true,
		},
		{
			name: "error",
			c:    &Cache{},
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{TestObjID1},
				moduleID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetM(tt.args.objType, tt.args.objIDs, tt.args.moduleID); (err != nil) != tt.wantErr {
				t.Errorf("Cache.SetM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_Del(t *testing.T) {
	c1 := &Cache{}
	c1.Set(TestObjType1, TestObjID1, TestModuleID1)
	c1.Set(TestObjType2, TestObjID2, TestModuleID2)

	type args struct {
		objType  ObjType
		objID    string
		moduleID string
	}
	tests := []struct {
		name string
		c    *Cache
		args args
		want bool
	}{
		{
			name: "OK",
			c:    c1,
			args: args{
				objType:  TestObjType1,
				objID:    TestObjID1,
				moduleID: TestModuleID1,
			},
			want: true,
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:  TestObjType1,
				objID:    TestObjID1,
				moduleID: TestModuleID1,
			},
			want: false,
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:  "",
				objID:    TestObjID2,
				moduleID: TestModuleID2,
			},
			want: false,
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:  TestObjType2,
				objID:    "",
				moduleID: TestModuleID2,
			},
			want: false,
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:  TestObjType2,
				objID:    TestObjID2,
				moduleID: "",
			},
			want: false,
		},
		{
			name: "OK",
			c:    c1,
			args: args{
				objType:  TestObjType2,
				objID:    TestObjID2,
				moduleID: TestModuleID2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Del(tt.args.objType, tt.args.objID, tt.args.moduleID); got != tt.want {
				t.Errorf("Cache.Del() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_DelM(t *testing.T) {
	c1 := &Cache{}
	c1.Set(TestObjType1, TestObjID1, TestModuleID1)

	c2 := &Cache{}
	c2.Set(TestObjType1, TestObjID1, TestModuleID1)
	c2.Set(TestObjType2, TestObjID2, TestModuleID2)

	c3 := &Cache{}
	c3.Set(TestObjType1, TestObjID1, TestModuleID1)
	c3.Set(TestObjType1, TestObjID2, TestModuleID1)

	type args struct {
		objType  ObjType
		objIDs   []string
		moduleID string
	}
	tests := []struct {
		name string
		c    *Cache
		args args
		want int
	}{
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:  TestObjType2,
				objIDs:   []string{TestObjID2},
				moduleID: TestModuleID2,
			},
			want: 0,
		},
		{
			name: "OK",
			c:    c1,
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{TestObjID1},
				moduleID: TestModuleID1,
			},
			want: 1,
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{TestObjID1},
				moduleID: TestModuleID1,
			},
			want: 0,
		},
		{
			name: "OK",
			c:    c2,
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{TestObjID1, ""},
				moduleID: TestModuleID1,
			},
			want: 1,
		},
		{
			name: "OK",
			c:    c2,
			args: args{
				objType:  TestObjType2,
				objIDs:   []string{TestObjID2, TestObjID2},
				moduleID: TestModuleID2,
			},
			want: 1,
		},
		{
			name: "OK",
			c:    c3,
			args: args{
				objType:  TestObjType1,
				objIDs:   []string{TestObjID1, TestObjID2},
				moduleID: TestModuleID1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.DelM(tt.args.objType, tt.args.objIDs, tt.args.moduleID); got != tt.want {
				t.Errorf("Cache.DelM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	c1 := &Cache{}
	c1.Set(TestObjType1, TestObjID1, TestModuleID1)
	c1.Set(TestObjType2, TestObjID2, TestModuleID2)

	type args struct {
		objType ObjType
		objID   string
	}
	tests := []struct {
		name string
		c    *Cache
		args args
		want string
	}{
		{
			name: "OK",
			c:    c1,
			args: args{
				objType: TestObjType1,
				objID:   TestObjID1,
			},
			want: TestModuleID1,
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType: TestObjType1,
				objID:   TestObjID2,
			},
			want: "",
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType: TestObjType2,
				objID:   TestObjID1,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Get(tt.args.objType, tt.args.objID); got != tt.want {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_RangeByType(t *testing.T) {
	c1 := &Cache{}
	c1.Set(TestObjType1, TestObjID1, TestModuleID1)
	c1.Set(TestObjType1, TestObjID2, TestModuleID2)
	c1.Set(TestObjType2, TestObjID2, TestModuleID2)

	c2 := &Cache{}

	inTimes1 := 0

	type args struct {
		objType        ObjType
		f              func(id string, location string) bool
		limitModuleIDs map[string]bool
	}
	tests := []struct {
		name string
		c    *Cache
		args args
	}{
		{
			name: "break",
			c:    c1,
			args: args{
				objType: TestObjType1,
				f: func(id string, location string) bool {
					inTimes1++
					if inTimes1 > 1 {
						t.Errorf("inTimes1 = %v, want %v", inTimes1, 1)
					}
					return false
				},
			},
		},
		{
			name: "OK",
			c:    c1,
			args: args{
				objType:        TestObjType1,
				limitModuleIDs: map[string]bool{TestModuleID1: true},
				f: func(id string, location string) bool {
					if location != TestModuleID1 {
						t.Errorf("location = %v, want %v", location, TestModuleID1)
					}
					return true
				},
			},
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:        TestObjType2,
				limitModuleIDs: map[string]bool{TestModuleID1: true},
				f: func(id string, location string) bool {
					t.Errorf("location = %v, want %v", location, "")
					return true
				},
			},
		},
		{
			name: "no target",
			c:    c1,
			args: args{
				objType:        TestObjType1,
				limitModuleIDs: map[string]bool{},
				f: func(id string, location string) bool {
					t.Errorf("location = %v, want %v", location, "")
					return true
				},
			},
		},
		{
			name: "no target",
			c:    c2,
			args: args{
				objType:        TestObjType1,
				limitModuleIDs: map[string]bool{TestModuleID1: true},
				f: func(id string, location string) bool {
					t.Errorf("location = %v, want %v", location, "")
					return true
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.RangeByType(tt.args.objType, tt.args.f, tt.args.limitModuleIDs)
		})
	}
	if inTimes1 != 1 {
		t.Errorf("inTimes1 = %v, want %v", inTimes1, 1)
	}
}

func TestCache_RandomObjIDByType(t *testing.T) {
	c1 := &Cache{}
	c1.Set(TestObjType1, TestObjID1, TestModuleID1)
	c1.Set(TestObjType1, TestObjID2, TestModuleID2)
	c1.Set(TestObjType2, TestObjID2, TestModuleID2)

	c2 := &Cache{}

	type args struct {
		objType        ObjType
		limitModuleIDs map[string]bool
	}
	tests := []struct {
		name   string
		c      *Cache
		args   args
		wantOr []string
		times  int
	}{
		{
			name: "OK",
			c:    c1,
			args: args{
				objType: TestObjType1,
			},
			wantOr: []string{TestObjID1, TestObjID2},
			times:  100,
		},
		{
			name: "OK",
			c:    c1,
			args: args{
				objType: TestObjType2,
			},
			wantOr: []string{TestObjID2},
			times:  100,
		},
		{
			name: "OK",
			c:    c2,
			args: args{
				objType: TestObjType1,
			},
			wantOr: []string{""},
			times:  100,
		},
	}
	in := func(ss []string, str string) bool {
		for _, v := range ss {
			if v == str {
				return true
			}
		}
		return false
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i <= tt.times; i++ {
				if got := tt.c.RandomObjIDByType(tt.args.objType, tt.args.limitModuleIDs); !in(tt.wantOr, got) {
					t.Errorf("Cache.RandomObjIDByType() = %v, wantOr %v", got, tt.wantOr)
				}
			}
		})
	}
}
