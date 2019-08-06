package httpconn

import (
	"bytes"
	"sync"
	//	"compress/gzip"
	"encoding/base64"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"io/ioutil"
	"math"
	"net/http"
	// "net/rpc"
	"strconv"
	"strings"
	"time"
)

func HttpDecode(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Use-Encrypt") != "Yes" {
		return
	}

	writer.Header().Set("Use-Encrypt", "Yes")

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(request.Body)
	if err != nil {
		log.Error("[HttpDecode] buf.ReadFrom Err[%s] N[%d]",
			err.Error(), n)
		return
	}
	decodeBytes, _ := base64.StdEncoding.DecodeString(buf.String())
	decode, _ := util.AesDecrypt([]byte(decodeBytes))
	if decode == nil {
		return
	}

	message := string(decode)
	message = strings.Replace(message, "{", "", -1)
	message = strings.Replace(message, "}", "", -1)
	message = strings.Replace(message, "\"", "", -1)
	log.Debug("aes message:%s", message)

	var newstr string
	split_msg := strings.Split(message, ",")
	for _, msg := range split_msg {
		slt := strings.Split(msg, ":")
		if len(slt) == 2 {
			newstr += slt[0] + "=" + slt[1] + "&"
		} else if len(slt) == 3 {
			newstr += slt[0] + "=" + slt[1] + ":" + slt[2] + "&"
		}
	}
	newlen := int64(math.Dim(float64(len(newstr)), 1))
	newstr = newstr[:newlen]
	log.Debug("aes newstr:%s len=%d", newstr, newlen)

	s := strings.NewReader(string(newstr))
	request.Body = ioutil.NopCloser(s)
	request.ContentLength = newlen
}

// 获取请求中指定键值的字符串值
func ParseFromHttp(request *http.Request, keyname string) string {
	keyvalues := request.PostFormValue(keyname)
	if len(keyvalues) > 0 {
		return keyvalues
	}
	keyvalues = request.FormValue(keyname)
	if len(keyvalues) > 0 {
		return keyvalues
	}
	return ""
}

// 获取请求中指定键值的整数值
func ParseIntFromHttp(request *http.Request, keyname string) int {
	valuestr := ParseFromHttp(request, keyname)
	valuenum, _ := strconv.Atoi(valuestr)
	return valuenum
}

// 获取请求中指定键值的uint64值
func ParseUInt64FromHttp(request *http.Request, keyname string) uint64 {
	valuestr := ParseFromHttp(request, keyname)
	valuenum, _ := strconv.ParseUint(valuestr, 10, 64)
	return valuenum
}

// http客户端连接管理器
type ClientConnPool struct {
	alllink     map[uint64]*HttpConn // 所有连接
	mutex       sync.Mutex
	starttempid uint64
}

var httptaskmanager_s *ClientConnPool

func init() {
	httptaskmanager_s = &ClientConnPool{}
	httptaskmanager_s.alllink = make(map[uint64]*HttpConn)
	httptaskmanager_s.starttempid = 1000000000
}

func (this *ClientConnPool) AddHttpTask(
	writer http.ResponseWriter) *HttpConn {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.starttempid++
	wstask := new(HttpConn)
	wstask.SetWriter(writer)
	curtime := uint64(time.Now().Unix())
	wstask.Tempid = (curtime << 32) + this.starttempid
	this.alllink[wstask.Tempid] = wstask
	log.Debug("[ClientConnPool.AddHttpTask] 添加新的HTTP连接,%d,%d",
		wstask.Tempid,
		len(this.alllink))

	notify := writer.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		this.mutex.Lock()
		defer this.mutex.Unlock()
		delete(this.alllink, wstask.Tempid) // 这里应该需要加锁
		log.Debug("[ClientConnPool.AddHttpTask] 删除关闭的的HTTP连接,%d,%d",
			wstask.Tempid, len(this.alllink))
	}()
	return wstask
}

func (this *ClientConnPool) RemoveHttpTask(tempid uint64) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, found := this.alllink[tempid]; found {
		delete(this.alllink, tempid) // 这里应该需要加锁
		log.Debug("[ClientConnPool.RemoveHttpTask] 删除HTTP连接,%d,%d",
			tempid, len(this.alllink))
	}
}

func (this *ClientConnPool) GetHttpTask(tempid uint64) *HttpConn {
	if value, found := this.alllink[tempid]; found {
		return value
	}
	return nil
}
