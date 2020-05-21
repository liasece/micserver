package httpconn

import (
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
)

// WriterReturnHTTPStrs 返回 HTTP 消息
func WriterReturnHTTPStrs(writer http.ResponseWriter, strs []string) {
	str := ""
	for _, v := range strs {
		str += v
	}
	WriterReturnHTTPStr(writer, str)
}

// WriterReturnHTTPStr 返回 HTTP 消息
func WriterReturnHTTPStr(writer http.ResponseWriter, str string) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("content-type", "application/json")
	writer.Header().Add("cache-control", "no-cache")
	writer.Header().Add("Accept-Encoding", "gzip, deflate")
	writer.Header().Add("Pragma", "no-cache")
	writer.Header().Set("connection", "keep-alive")
	writer.WriteHeader(http.StatusOK)
	log.Debug("%s", str)

	if writer.Header().Get("Use-Encrypt") == "Yes" {
		aesstr, _ := util.AesEncrypt([]byte(str))
		encodeString := base64.StdEncoding.EncodeToString(aesstr)
		n, err := io.WriteString(writer, encodeString)
		if err != nil {
			log.Error("[WriterReturnHTTPStr] io.WriteString Err[%s] N[%d]",
				err.Error(), n)
		}
	} else {
		n, err := io.WriteString(writer, str)
		if err != nil {
			log.Error("[WriterReturnHTTPStr] io.WriteString Err[%s] N[%d]",
				err.Error(), n)
		}
	}
}

// HTTPConn HTTP 连接对象
type HTTPConn struct {
	writer http.ResponseWriter
	Tempid uint64 // 唯一编号
	Openid string
	BufStr []string
	Holder chan int
}

// SetWriter 设置 HTTP 返回
func (httpConn *HTTPConn) SetWriter(w http.ResponseWriter) {
	httpConn.writer = w
}

// AppendBufStr 增加一个 HTTP 返回值到缓冲区中
func (httpConn *HTTPConn) AppendBufStr(str string) {
	if httpConn.BufStr == nil {
		httpConn.BufStr = make([]string, 0)
	}
	httpConn.BufStr = append(httpConn.BufStr, str)
}

// ReturnHTTPStr 直接返回一个 HTTP 请求返回值
func (httpConn *HTTPConn) ReturnHTTPStr(str string) {
	WriterReturnHTTPStr(httpConn.writer, str)
}

// Hold 保持住一个 HTTP 请求
func (httpConn *HTTPConn) Hold() {
	if httpConn.Holder == nil {
		httpConn.Holder = make(chan int)
	}
}

// Wait 等待一个 HTTP 请求返回完成
func (httpConn *HTTPConn) Wait() int {
	if httpConn.Holder == nil {
		return -1
	}

	// 超时时间为3秒
	select {
	case <-httpConn.Holder:
		return 0
	case <-time.After(time.Second * 3):
		return 1
	}
}

// Release HTTP 请求返回完成，释放该HTTP请求的保持
func (httpConn *HTTPConn) Release() {
	if httpConn.Holder == nil {
		return
	}
	httpConn.Holder <- 1
}
