package httpconn

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
)

// HTTPDecode HTTP 连接消息编解码器
func HTTPDecode(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Use-Encrypt") != "Yes" {
		return
	}

	writer.Header().Set("Use-Encrypt", "Yes")

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(request.Body)
	if err != nil {
		log.Error("[HTTPDecode] buf.ReadFrom Err[%s] N[%d]",
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
	splitMsg := strings.Split(message, ",")
	for _, msg := range splitMsg {
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

// ParseFromHTTP 获取请求中指定键值的字符串值
func ParseFromHTTP(request *http.Request, keyname string) string {
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

// ParseIntFromHTTP 获取请求中指定键值的整数值
func ParseIntFromHTTP(request *http.Request, keyname string) int {
	valuestr := ParseFromHTTP(request, keyname)
	valuenum, _ := strconv.Atoi(valuestr)
	return valuenum
}

// ParseUInt64FromHTTP 获取请求中指定键值的uint64值
func ParseUInt64FromHTTP(request *http.Request, keyname string) uint64 {
	valuestr := ParseFromHTTP(request, keyname)
	valuenum, _ := strconv.ParseUint(valuestr, 10, 64)
	return valuenum
}

// ClientPool http客户端连接管理器
type ClientPool struct {
	alllink     map[uint64]*HTTPConn // 所有连接
	mutex       sync.Mutex
	starttempid uint64
}

var httptaskmanager *ClientPool

func init() {
	httptaskmanager = &ClientPool{}
	httptaskmanager.alllink = make(map[uint64]*HTTPConn)
	httptaskmanager.starttempid = 1000000000
}

// AddHTTPTask 增加一个 HTTP 客户端连接
func (cp *ClientPool) AddHTTPTask(
	writer http.ResponseWriter) *HTTPConn {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.starttempid++
	wstask := new(HTTPConn)
	wstask.SetWriter(writer)
	curtime := uint64(time.Now().Unix())
	wstask.Tempid = (curtime << 32) + cp.starttempid
	cp.alllink[wstask.Tempid] = wstask
	log.Debug("[ClientPool.AddHTTPTask] 添加新的HTTP连接,%d,%d",
		wstask.Tempid,
		len(cp.alllink))

	notify := writer.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		cp.mutex.Lock()
		defer cp.mutex.Unlock()
		delete(cp.alllink, wstask.Tempid) // 这里应该需要加锁
		log.Debug("[ClientPool.AddHTTPTask] 删除关闭的的HTTP连接,%d,%d",
			wstask.Tempid, len(cp.alllink))
	}()
	return wstask
}

// RemoveHTTPTask 移除一个 HTTP 客户端连接
func (cp *ClientPool) RemoveHTTPTask(tempid uint64) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if _, found := cp.alllink[tempid]; found {
		delete(cp.alllink, tempid) // 这里应该需要加锁
		log.Debug("[ClientPool.RemoveHTTPTask] 删除HTTP连接,%d,%d",
			tempid, len(cp.alllink))
	}
}

// GetHTTPTask 获取一个 HTTP 客户端连接
func (cp *ClientPool) GetHTTPTask(tempid uint64) *HTTPConn {
	if value, found := cp.alllink[tempid]; found {
		return value
	}
	return nil
}
