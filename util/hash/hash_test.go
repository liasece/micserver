package hash

import "testing"

func TestGetHash(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "OK",
			args: args{
				data: []byte("TestHash1"),
			},
			want: 1062913394,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHash(tt.args.data); got != tt.want {
				t.Errorf("GetHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
