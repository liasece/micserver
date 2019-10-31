/**
 * \file IOBuffer.go
 * \version
 * \author liaojiansheng
 * \date  2019年01月31日 12:22:43
 * \brief 消息收发缓冲区
 *
 */

package buffer

import (
	"errors"
	"io"
)

var (
	ErrBufNil    = errors.New("buf is nil")
	ErrReaderNil = errors.New("reader is nil")
	ErrNotEnough = errors.New("not enough")
	ErrOverSize  = errors.New("oversize")
	ErrLess0     = errors.New("less than 0")
)

// 不是线程安全的
type IOBuffer struct {
	reader    io.Reader
	buf       []byte
	start     int
	end       int
	maxLength int
}

func NewIOBuffer(reader io.Reader, length int) *IOBuffer {
	buf := make([]byte, length)
	return &IOBuffer{reader, buf, 0, 0, length}
}

func (b *IOBuffer) Len() int {
	return b.end - b.start
}

// 将有用的字节前移
func (b *IOBuffer) grow() error {
	if b.buf == nil {
		return ErrBufNil
	}
	if b.start == 0 {
		return nil
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
	return nil
}

// 当前剩余大小
func (b *IOBuffer) RemainSize() int {
	return b.maxLength - (b.end - b.start)
}

// 总大小
func (b *IOBuffer) TotalSize() int {
	return b.maxLength
}

// 从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *IOBuffer) ReadFromReader() (int, error) {
	if b.reader == nil {
		return 0, ErrReaderNil
	}
	gerr := b.grow()
	if gerr != nil {
		return 0, gerr
	}
	n, err := b.reader.Read(b.buf[b.end:])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

// 返回n个字节，而不产生移位
func (b *IOBuffer) Seek(n int) ([]byte, error) {
	if b.buf == nil {
		return nil, ErrBufNil
	}
	if b.end-b.start >= n {
		buf := b.buf[b.start : b.start+n]
		return buf, nil
	}
	return nil, ErrNotEnough
}

// 返回所有字节，而不产生移位
func (b *IOBuffer) SeekAll() ([]byte, error) {
	if b.buf == nil {
		return nil, ErrBufNil
	}
	return b.buf[b.start:b.end], nil
}

// 舍弃offset个字段，读取n个字段
func (b *IOBuffer) Read(offset, n int) ([]byte, error) {
	if b.buf == nil {
		return nil, ErrBufNil
	}
	if b.maxLength < n {
		return nil, ErrOverSize
	}
	if offset < 0 || n < 0 {
		return nil, ErrLess0
	}
	if b.start+offset+n > b.end {
		return nil, ErrNotEnough
	}
	b.start += offset
	buf := b.buf[b.start : b.start+n]
	b.start += n
	return buf, nil
}

// 写入一段数据，要么全部成功，要么全部不成功
func (b *IOBuffer) Write(src []byte) error {
	gerr := b.grow()
	if gerr != nil {
		return gerr
	}
	size := len(src)
	if size > b.RemainSize() {
		return ErrOverSize
	}
	b.end += copy(b.buf[b.end:], src)
	return nil
}

func (b *IOBuffer) MoveStart(n int) error {
	tmpn := b.start + n
	if tmpn < 0 {
		return ErrLess0
	}
	if tmpn < b.end {
		return ErrOverSize
	}
	b.start = tmpn
	return nil
}
