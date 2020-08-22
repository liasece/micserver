package roc

import (
	"reflect"
	"testing"
)

var objType1 ObjType = "o1"
var objType2 ObjType = "o2"

var objID1 = "id1"
var objID2 = "id2"

var testPathString1 = string(objType1) + "[" + objID1 + "].f1"
var testPathString2 = string(objType2) + "[" + objID2 + "].f2.ff2"

func TestO(t *testing.T) {
	type args struct {
		objType ObjType
		objID   string
	}
	tests := []struct {
		name string
		args args
		want *Path
	}{
		{
			name: "OK",
			args: args{
				objType: objType1,
				objID:   objID1,
			},
			want: NewPath(string(objType1) + "[" + objID1 + "]"),
		},
		{
			name: "OK",
			args: args{
				objType: objType1,
				objID:   objID2,
			},
			want: NewPath(string(objType1) + "[" + objID2 + "]"),
		},
		{
			name: "OK",
			args: args{
				objType: "",
				objID:   objID1,
			},
			want: NewPath(string("") + "[" + objID1 + "]"),
		},
		{
			name: "OK",
			args: args{
				objType: objType1,
				objID:   "",
			},
			want: NewPath(string(objType1) + "[" + "" + "]"),
		},
		{
			name: "OK",
			args: args{
				objType: "",
				objID:   "",
			},
			want: NewPath(string("") + "[" + "" + "]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := O(tt.args.objType, tt.args.objID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("O() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathStrDecode(t *testing.T) {
	type args struct {
		pathStr string
	}
	tests := []struct {
		name  string
		args  args
		want  ObjType
		want1 string
	}{
		{
			name: "OK",
			args: args{
				pathStr: testPathString1,
			},
			want:  objType1,
			want1: objID1,
		},
		{
			name: "OK",
			args: args{
				pathStr: testPathString2,
			},
			want:  objType2,
			want1: objID2,
		},
		{
			name: "OK",
			args: args{
				pathStr: "[]",
			},
			want:  "",
			want1: "",
		},
		{
			name: "OK",
			args: args{
				pathStr: "t1[o1]]",
			},
			want:  "t1",
			want1: "o1",
		},
		{
			name: "OK",
			args: args{
				pathStr: "[t1[o1]]",
			},
			want:  "",
			want1: "t1[o1",
		},
		{
			name: "OK",
			args: args{
				pathStr: "t1][t1[o1]]",
			},
			want:  "t1]",
			want1: "t1[o1",
		},
		{
			name: "OK",
			args: args{
				pathStr: "][t1[o1]]",
			},
			want:  "]",
			want1: "t1[o1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := pathStrDecode(tt.args.pathStr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pathStrDecode() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("pathStrDecode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewPath(t *testing.T) {
	type args struct {
		pathStr string
	}
	tests := []struct {
		name string
		args args
		want *Path
	}{
		{
			name: "OK",
			args: args{
				pathStr: testPathString1,
			},
			want: O(objType1, objID1).F("f1"),
		},
		{
			name: "OK",
			args: args{
				pathStr: testPathString2,
			},
			want: O(objType2, objID2).F("f2").F("ff2"),
		},
		{
			name: "OK",
			args: args{
				pathStr: "[]",
			},
			want: O("", ""),
		},
		{
			name: "OK",
			args: args{
				pathStr: "[].",
			},
			want: O("", "").F(""),
		},
		{
			name: "OK",
			args: args{
				pathStr: "t1][t1[o1]].f1",
			},
			want: O("t1]", "t1[o1").F("f1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPath(tt.args.pathStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPath() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestPath_GetObjType(t *testing.T) {
	tests := []struct {
		name string
		path *Path
		want ObjType
	}{
		{
			name: "OK",
			path: NewPath(testPathString1),
			want: objType1,
		},
		{
			name: "OK",
			path: NewPath(testPathString2),
			want: objType2,
		},
		{
			name: "OK",
			path: NewPath(""),
			want: "",
		},
		{
			name: "OK",
			path: NewPath("t1[]"),
			want: "t1",
		},
		{
			name: "OK",
			path: NewPath("[o1]"),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.GetObjType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path.GetObjType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_GetObjID(t *testing.T) {
	tests := []struct {
		name string
		path *Path
		want string
	}{
		{
			name: "OK",
			path: NewPath(testPathString1),
			want: objID1,
		},
		{
			name: "OK",
			path: NewPath(testPathString2),
			want: objID2,
		},
		{
			name: "OK",
			path: NewPath(""),
			want: "",
		},
		{
			name: "OK",
			path: NewPath("t1[]"),
			want: "",
		},
		{
			name: "OK",
			path: NewPath("[o1]"),
			want: "o1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.GetObjID(); got != tt.want {
				t.Errorf("Path.GetObjID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_GetPos(t *testing.T) {
	p21 := NewPath(testPathString2)
	p21.Move()
	p22 := NewPath(testPathString2)
	p22.Move()
	p22.Move()
	p23 := NewPath(testPathString2)
	p23.Move()
	p23.Move()
	p23.Move()
	tests := []struct {
		name string
		path *Path
		want int
	}{
		{
			name: "OK",
			path: NewPath(testPathString1),
			want: 0,
		},
		{
			name: "OK",
			path: NewPath(testPathString2),
			want: 0,
		},
		{
			name: "OK",
			path: p21,
			want: 1,
		},
		{
			name: "OK",
			path: p22,
			want: 2,
		},
		{
			name: "OK",
			path: p23,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.GetPos(); got != tt.want {
				t.Errorf("Path.GetPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Get(t *testing.T) {
	p21 := NewPath(testPathString2)
	p21.Move()

	type args struct {
		pos int
	}
	tests := []struct {
		name string
		path *Path
		args args
		want string
	}{
		{
			name: "OK",
			path: NewPath(testPathString1),
			args: args{
				pos: 0,
			},
			want: "f1",
		},
		{
			name: "OK",
			path: NewPath(testPathString1),
			args: args{
				pos: 1,
			},
			want: "",
		},
		{
			name: "OK",
			path: NewPath(testPathString2),
			args: args{
				pos: 1,
			},
			want: "ff2",
		},
		{
			name: "OK",
			path: p21,
			args: args{
				pos: 0,
			},
			want: "f2",
		},
		{
			name: "OK",
			path: p21,
			args: args{
				pos: 1,
			},
			want: "ff2",
		},
		{
			name: "OK",
			path: p21,
			args: args{
				pos: 2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.Get(tt.args.pos); got != tt.want {
				t.Errorf("Path.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Move(t *testing.T) {
	p2 := NewPath(testPathString2)
	tests := []struct {
		name string
		path *Path
		want string
	}{
		{
			name: "OK",
			path: NewPath(testPathString1),
			want: "f1",
		},
		{
			name: "OK",
			path: NewPath(testPathString2),
			want: "f2",
		},
		{
			name: "OK",
			path: p2,
			want: "f2",
		},
		{
			name: "OK",
			path: p2,
			want: "ff2",
		},
		{
			name: "OK",
			path: p2,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.Move(); got != tt.want {
				t.Errorf("Path.Move() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_F(t *testing.T) {
	p1 := NewPath(testPathString1)
	type args struct {
		funcName string
	}
	tests := []struct {
		name string
		path *Path
		args args
		want *Path
	}{
		{
			name: "OK",
			path: NewPath(testPathString1),
			args: args{
				funcName: "t1",
			},
			want: NewPath(testPathString1 + ".t1"),
		},
		{
			name: "OK",
			path: NewPath(""),
			args: args{
				funcName: "t1",
			},
			want: NewPath("" + ".t1"),
		},
		{
			name: "OK",
			path: NewPath("ads[fs]ds]aas]."),
			args: args{
				funcName: "t1",
			},
			want: NewPath("ads[fs]" + "." + ".t1"),
		},
		{
			name: "OK",
			path: p1,
			args: args{
				funcName: "t1",
			},
			want: NewPath(testPathString1 + ".t1"),
		},
		{
			name: "OK",
			path: p1,
			args: args{
				funcName: "tt1",
			},
			want: NewPath(testPathString1 + ".t1.tt1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.F(tt.args.funcName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path.F() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_String(t *testing.T) {
	tests := []struct {
		name string
		path *Path
		want string
	}{
		{
			name: "OK",
			path: O(objType1, objID1).F("f1"),
			want: testPathString1,
		},
		{
			name: "OK",
			path: O(objType2, objID2).F("f2").F("ff2"),
			want: testPathString2,
		},
		{
			name: "OK",
			path: O("", "").F("f2").F("ff2"),
			want: "[].f2.ff2",
		},
		{
			name: "OK",
			path: O("t1", "").F("f2").F("ff2"),
			want: "t1[].f2.ff2",
		},
		{
			name: "OK",
			path: O("", "o1").F("f2").F(""),
			want: "[o1].f2.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.String(); got != tt.want {
				t.Errorf("Path.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Reset(t *testing.T) {
	p21 := NewPath(testPathString2)
	p21.Move()
	p22 := NewPath(testPathString2)
	p22.Move()
	p22.Move()
	p23 := NewPath(testPathString2)
	p23.Move()
	p23.Move()
	p23.Move()

	tests := []struct {
		name string
		path *Path
	}{
		{
			name: "OK",
			path: p21,
		},
		{
			name: "OK",
			path: p22,
		},
		{
			name: "OK",
			path: p23,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.path.Reset()
		})
	}
}
