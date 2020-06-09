/**
 * \file IOBuffer.go
 * \version
 * \author Jansen
 * \date  2019年01月31日 12:22:43
 * \brief 消息收发缓冲区
 *
 */

/*
Package buffer 无拷贝IO缓冲区实现
*/
package buffer

import (
	"io"
	"reflect"
	"testing"

	"github.com/liasece/micserver/log"
)

func TestNewIOBuffer(t *testing.T) {
	type args struct {
		reader io.Reader
		length int
	}
	tests := []struct {
		name    string
		args    args
		want    *IOBuffer
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				reader: nil,
				length: 1,
			},
			want: &IOBuffer{
				buf:           make([]byte, 1),
				maxLength:     1,
				defaultLength: 1,
			},
		},
		{
			name: "neglength",
			args: args{
				reader: nil,
				length: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIOBuffer(tt.args.reader, tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIOBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIOBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_SetBanAutoResize(t *testing.T) {
	buf1, _ := NewIOBuffer(nil, 10)

	type args struct {
		value bool
	}
	tests := []struct {
		name string
		b    *IOBuffer
		args args
	}{
		{
			name: "sec",
			b:    buf1,
			args: args{
				value: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetBanAutoResize(tt.args.value)
		})
	}
}

func TestIOBuffer_Len(t *testing.T) {
	buf1, _ := NewIOBuffer(nil, 1)
	buf2, _ := NewIOBuffer(nil, 1)
	buf3, _ := NewIOBuffer(nil, 1)
	buf2.Write([]byte{0x10})
	buf3.Write([]byte{0x10, 0x20})

	tests := []struct {
		name string
		b    *IOBuffer
		want int
	}{
		{
			name: "sec",
			b:    buf1,
			want: 0,
		},
		{
			name: "sec",
			b:    buf2,
			want: 1,
		},
		{
			name: "sec",
			b:    buf3,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Len(); got != tt.want {
				t.Errorf("IOBuffer.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_RemainSize(t *testing.T) {
	buf1, _ := NewIOBuffer(nil, 1)
	buf2, _ := NewIOBuffer(nil, 1)
	buf3, _ := NewIOBuffer(nil, 1)
	buf2.Write([]byte{0x10})
	buf3.Write([]byte{0x10, 0x20})

	tests := []struct {
		name string
		b    *IOBuffer
		want int
	}{
		{
			name: "sec",
			b:    buf1,
			want: 1,
		},
		{
			name: "sec_full",
			b:    buf2,
			want: 0,
		},
		{
			name: "sec_oversize",
			b:    buf3,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.RemainSize(); got != tt.want {
				t.Errorf("IOBuffer.RemainSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_TotalSize(t *testing.T) {
	buf1, _ := NewIOBuffer(nil, 1)
	buf2, _ := NewIOBuffer(nil, 1)
	buf3, _ := NewIOBuffer(nil, 1)
	buf2.Write([]byte{0x10})
	buf3.Write([]byte{0x10, 0x20})

	tests := []struct {
		name string
		b    *IOBuffer
		want int
	}{
		{
			name: "sec",
			b:    buf1,
			want: 1,
		},
		{
			name: "sec_full",
			b:    buf2,
			want: 1,
		},
		{
			name: "sec_oversize",
			b:    buf3,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.TotalSize(); got != tt.want {
				t.Errorf("IOBuffer.TotalSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_SetLogger(t *testing.T) {
	buf1, _ := NewIOBuffer(nil, 1)

	type args struct {
		l *log.Logger
	}
	tests := []struct {
		name string
		b    *IOBuffer
		args args
	}{
		{
			name: "sec",
			b:    buf1,
			args: args{
				l: log.GetDefaultLogger(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetLogger(tt.args.l)
		})
	}
}

func TestIOBuffer_DefaultSize(t *testing.T) {
	buf1, _ := NewIOBuffer(nil, 1)
	buf2, _ := NewIOBuffer(nil, 1)
	buf3, _ := NewIOBuffer(nil, 1)
	buf2.Write([]byte{0x10})
	buf3.Write([]byte{0x10, 0x20})

	tests := []struct {
		name string
		b    *IOBuffer
		want int
	}{
		{
			name: "sec",
			b:    buf1,
			want: 1,
		},
		{
			name: "sec_full",
			b:    buf2,
			want: 1,
		},
		{
			name: "sec_oversize",
			b:    buf3,
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.DefaultSize(); got != tt.want {
				t.Errorf("IOBuffer.DefaultSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

type treader struct{}

func (*treader) Read(p []byte) (n int, err error) {
	if len(p) <= 0 {
		return 0, nil
	}
	p[0] = 0x10
	return 1, nil
}

func TestIOBuffer_ReadFromReader(t *testing.T) {
	buf1, _ := NewIOBuffer(&treader{}, 1)
	buf2, _ := NewIOBuffer(&treader{}, 1)
	buf2.Write([]byte{0x10})

	tests := []struct {
		name    string
		b       *IOBuffer
		want    int
		wantErr bool
	}{
		{
			name: "sec",
			b:    buf1,
			want: 1,
		},
		{
			name: "sec",
			b:    buf2,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.ReadFromReader()
			if (err != nil) != tt.wantErr {
				t.Errorf("IOBuffer.ReadFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IOBuffer.ReadFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_Seek(t *testing.T) {
	buf1, _ := NewIOBuffer(&treader{}, 1)
	buf2, _ := NewIOBuffer(&treader{}, 1)
	buf2.Write([]byte{0x10})

	type args struct {
		n int
	}
	tests := []struct {
		name    string
		b       *IOBuffer
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "0",
			b:    buf1,
			args: args{
				n: 0,
			},
			want: []byte{},
		},
		{
			name: "notenough",
			b:    buf1,
			args: args{
				n: 1,
			},
			wantErr: true,
		},
		{
			name: "sec",
			b:    buf2,
			args: args{
				n: 1,
			},
			want: []byte{0x10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Seek(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("IOBuffer.Seek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IOBuffer.Seek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_SeekAll(t *testing.T) {
	buf1, _ := NewIOBuffer(&treader{}, 1)
	buf2, _ := NewIOBuffer(&treader{}, 1)
	buf2.Write([]byte{0x10})

	tests := []struct {
		name    string
		b       *IOBuffer
		want    []byte
		wantErr bool
	}{
		{
			name: "0",
			b:    buf1,
			want: []byte{},
		},
		{
			name: "sec",
			b:    buf2,
			want: []byte{0x10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.SeekAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("IOBuffer.SeekAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IOBuffer.SeekAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_Read(t *testing.T) {
	buf1, _ := NewIOBuffer(&treader{}, 1)
	buf2, _ := NewIOBuffer(&treader{}, 1)
	buf2.Write([]byte{0x10})
	buf3, _ := NewIOBuffer(&treader{}, 1)
	buf3.Write([]byte{0x10, 0x20})

	type args struct {
		offset int
		n      int
	}
	tests := []struct {
		name    string
		b       *IOBuffer
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "0",
			b:    buf1,
			args: args{
				offset: 0,
				n:      1,
			},
			wantErr: true,
		},
		{
			name: "sec",
			b:    buf2,
			args: args{
				offset: 0,
				n:      1,
			},
			want: []byte{0x10},
		},
		{
			name: "sec",
			b:    buf3,
			args: args{
				offset: 1,
				n:      1,
			},
			want: []byte{0x20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Read(tt.args.offset, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("IOBuffer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IOBuffer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOBuffer_Write(t *testing.T) {
	buf1, _ := NewIOBuffer(&treader{}, 1)
	buf2, _ := NewIOBuffer(&treader{}, 1)
	buf3, _ := NewIOBuffer(&treader{}, 1)
	buf3.SetBanAutoResize(true)

	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		b       *IOBuffer
		args    args
		wantErr bool
		seekAll []byte
	}{
		{
			name: "sec",
			b:    buf1,
			args: args{
				src: []byte{0x10},
			},
			seekAll: []byte{0x10},
		},
		{
			name: "sec",
			b:    buf2,
			args: args{
				src: []byte{0x10, 0x20},
			},
			seekAll: []byte{0x10, 0x20},
		},
		{
			name: "sec",
			b:    buf3,
			args: args{
				src: []byte{0x10, 0x20},
			},
			wantErr: true,
			seekAll: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Write(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("IOBuffer.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if b, err := tt.b.SeekAll(); err != nil || !reflect.DeepEqual(b, tt.seekAll) {
				t.Errorf("IOBuffer.SeekAll() error = %v, wantErr nil, bytes = %v, wantBytes %v", err, b, tt.seekAll)
			}
		})
	}
}

func TestIOBuffer_MoveStartPtr(t *testing.T) {
	buf1, _ := NewIOBuffer(&treader{}, 1)
	buf2, _ := NewIOBuffer(&treader{}, 1)
	buf2.Write([]byte{0x10})
	buf3, _ := NewIOBuffer(&treader{}, 1)
	buf3.Write([]byte{0x10, 0x20})
	buf4, _ := NewIOBuffer(&treader{}, 1)
	buf4.Write([]byte{0x10, 0x20})

	type args struct {
		n int
	}
	tests := []struct {
		name    string
		b       *IOBuffer
		args    args
		wantErr bool
	}{
		{
			name: "0",
			b:    buf1,
			args: args{
				n: 0,
			},
		},
		{
			name: "sec",
			b:    buf2,
			args: args{
				n: 1,
			},
		},
		{
			name: "sec",
			b:    buf3,
			args: args{
				n: 2,
			},
		},
		{
			name: "sec",
			b:    buf4,
			args: args{
				n: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.MoveStartPtr(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("IOBuffer.MoveStartPtr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
