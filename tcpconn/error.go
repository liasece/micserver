/**
 * \file TCPConn.go
 * \version
 * \author liaojiansheng
 * \date  2019年01月31日 12:22:43
 * \brief 连接数据管理器
 *
 */

package tcpconn

import (
	"errors"
	// "fmt"
)

var (
	ErrSendNilData = errors.New("send nil data")
	ErrCloseed     = errors.New("conn has been closed")
	ErrBufferFull  = errors.New("buffer full")
)
