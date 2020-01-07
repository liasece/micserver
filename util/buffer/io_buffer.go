/**
 * \file IOBuffer.go
 * \version
 * \author Jansen
 * \date  2019年01月31日 12:22:43
 * \brief 消息收发缓冲区
 *
 */

/*
无拷贝IO缓冲区实现
*/
package buffer

import (
	"errors"
	"github.com/liasece/micserver/log"
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
	*log.Logger

	reader        io.Reader
	buf           []byte
	start         int
	end           int
	maxLength     int
	defaultLength int
	banAutoResize bool
}

// 构造一个缓冲区
func NewIOBuffer(reader io.Reader, length int) *IOBuffer {
	buf := make([]byte, length)
	return &IOBuffer{
		reader:        reader,
		buf:           buf,
		start:         0,
		end:           0,
		maxLength:     length,
		defaultLength: length,
	}
}

// 设置缓冲区是否可以根据需求自动调整大小
func (b *IOBuffer) SetBanAutoResize(value bool) {
	b.banAutoResize = value
}

// 当前缓冲区内容的长度
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

func (b *IOBuffer) resize(length int) error {
	newbuf := make([]byte, length)
	if b.end != 0 || b.start != 0 {
		// 向新缓冲区中 grow
		copy(newbuf, b.buf[b.start:b.end])
		b.end -= b.start
		b.start = 0
	}
	b.buf = newbuf
	b.maxLength = length
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
	// 如果缓冲区空了，需要将扩容的内存还回去
	if b.end == 0 && b.maxLength >= b.defaultLength*2+1 {
		b.Syslog("缓冲区扩容恢复 %d->%d", b.maxLength, b.defaultLength)
		b.resize(b.defaultLength)
	}

	leftSize := b.maxLength - b.end
	n, err := b.reader.Read(b.buf[b.end:])
	if err != nil {
		return n, err
	}
	b.end += n
	if n == leftSize && !b.banAutoResize {
		// 缓冲区满，扩容一次，最大容忍超过默认值的16倍
		targetLength := b.maxLength * 2
		if targetLength <= b.defaultLength*16 {
			b.Syslog("缓冲区满，扩容 %d->%d", b.maxLength, targetLength)
			b.resize(targetLength)
		} else {
			b.Error("缓冲区满，扩容失败！ now[%d] default[%d]",
				b.maxLength, b.defaultLength)
		}
	}
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
	// 如果缓冲区空了，需要将扩容的内存还回去
	if b.end == 0 && b.maxLength >= b.defaultLength*2 {
		b.Syslog("缓冲区扩容恢复 %d->%d", b.maxLength, b.defaultLength)
		b.resize(b.defaultLength)
	}

	size := len(src)
	if size > b.RemainSize() && !b.banAutoResize {
		// 缓冲区满，扩容一次，最大容忍超过默认值的16倍
		targetLength := b.end + size
		if targetLength <= b.defaultLength*16 {
			b.Syslog("缓冲区满，扩容 %d->%d", b.maxLength, targetLength)
			b.resize(targetLength)
		} else {
			b.Error("缓冲区满，扩容失败！ now[%d] default[%d]",
				b.maxLength, b.defaultLength)
		}
		// return ErrOverSize
	}

	size = len(src)
	if size > b.RemainSize() {
		return ErrOverSize
	}

	b.end += copy(b.buf[b.end:], src)
	return nil
}

// 修改缓冲区内容起始指针
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
