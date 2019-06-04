/**
 * \file IOBuffer.go
 * \version
 * \author liaojiansheng
 * \date  2019年01月31日 12:22:43
 * \brief 消息收发缓冲区
 *
 */

package base

import (
	"errors"
	"fmt"
	"io"
)

// 不是线程安全的
type IOBuffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

func NewIOBuffer(reader io.Reader, length int) *IOBuffer {
	buf := make([]byte, length)
	return &IOBuffer{reader, buf, 0, 0}
}

func (b *IOBuffer) Len() int {
	return b.end - b.start
}

// 将有用的字节前移
func (b *IOBuffer) grow() error {
	if b.buf == nil {
		return errors.New("buf is nil")
	}
	if b.start == 0 {
		return nil
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
	return nil
}

// 从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *IOBuffer) ReadFromReader() (int, error) {
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
		return nil, errors.New("buf is nil")
	}
	if b.end-b.start >= n {
		buf := b.buf[b.start : b.start+n]
		return buf, nil
	}
	return nil, errors.New("not enough")
}

// 舍弃offset个字段，读取n个字段
func (b *IOBuffer) Read(offset, n int) ([]byte, error) {
	if b.buf == nil {
		return nil, errors.New("buf is nil")
	}
	if offset < 0 || n < 0 {
		return nil, fmt.Errorf("err value offset:%d n:%d", offset, n)
	}
	if b.start+offset+n > b.end {
		return nil, errors.New("not enough")
	}
	b.start += offset
	buf := b.buf[b.start : b.start+n]
	b.start += n
	return buf, nil
}
