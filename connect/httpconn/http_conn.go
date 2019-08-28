package httpconn

import (
	"encoding/base64"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"io"
	"net/http"
	"time"
)

func WriterReturnHttpStrs(writer http.ResponseWriter, strs []string) {
	str := ""
	for _, v := range strs {
		str += v
	}
	WriterReturnHttpStr(writer, str)
}

func WriterReturnHttpStr(writer http.ResponseWriter, str string) {
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
			log.Error("[WriterReturnHttpStr] io.WriteString Err[%s] N[%d]",
				err.Error(), n)
		}
	} else {
		n, err := io.WriteString(writer, str)
		if err != nil {
			log.Error("[WriterReturnHttpStr] io.WriteString Err[%s] N[%d]",
				err.Error(), n)
		}
	}
}

type HttpConn struct {
	writer http.ResponseWriter
	Tempid uint64 // 唯一编号
	Openid string
	BufStr []string
	Holder chan int
}

func (this *HttpConn) SetWriter(w http.ResponseWriter) {
	this.writer = w
}

func (this *HttpConn) AppendBufStr(str string) {
	if this.BufStr == nil {
		this.BufStr = make([]string, 0)
	}
	this.BufStr = append(this.BufStr, str)
}

func (this *HttpConn) ReturnHttpStr(str string) {
	WriterReturnHttpStr(this.writer, str)
}

func (this *HttpConn) Hold() {
	if this.Holder == nil {
		this.Holder = make(chan int)
	}
}
func (this *HttpConn) Wait() int {
	if this.Holder == nil {
		return -1
	}

	// 超时时间为3秒
	select {
	case <-this.Holder:
		return 0
	case <-time.After(time.Second * 3):
		return 1
	}
}
func (this *HttpConn) Release() {
	if this.Holder == nil {
		return
	}
	this.Holder <- 1
}
