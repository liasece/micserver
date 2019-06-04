package subnet

import (
	"bytes"
	"sync"
	//	"compress/gzip"
	"encoding/base64"
	"github.com/liasece/micserver/encode"
	"github.com/liasece/micserver/log"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	// "net/rpc"
	"strconv"
	"strings"
	"time"
)

type GBHttpTask struct {
	writer http.ResponseWriter
	Tempid uint64 // 唯一编号
	Openid string
	BufStr []string
	Holder chan int
}

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
		aesstr, _ := encode.AesEncrypt([]byte(str))
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

func (this *GBHttpTask) AppendBufStr(str string) {
	if this.BufStr == nil {
		this.BufStr = make([]string, 0)
	}
	this.BufStr = append(this.BufStr, str)
}

func (this *GBHttpTask) ReturnHttpStr(str string) {
	WriterReturnHttpStr(this.writer, str)
}

func (this *GBHttpTask) Hold() {
	if this.Holder == nil {
		this.Holder = make(chan int)
	}
}
func (this *GBHttpTask) Wait() int {
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
func (this *GBHttpTask) Release() {
	if this.Holder == nil {
		return
	}
	this.Holder <- 1
}

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
	decode, _ := encode.AesDecrypt([]byte(decodeBytes))
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
func ParseIntFromHttp(request *http.Request, keyname string) int {
	valuestr := ParseFromHttp(request, keyname)
	valuenum, _ := strconv.Atoi(valuestr)
	return valuenum
}
func ParseUInt64FromHttp(request *http.Request, keyname string) uint64 {
	valuestr := ParseFromHttp(request, keyname)
	valuenum, _ := strconv.ParseUint(valuestr, 10, 64)
	return valuenum
}

// httptask连接管理器
type GBHttpTaskManger struct {
	allsockets  map[uint64]*GBHttpTask // 所有连接
	mutex       sync.Mutex
	starttempid uint64
}

var httptaskmanager_s *GBHttpTaskManger

func init() {
	httptaskmanager_s = &GBHttpTaskManger{}
	httptaskmanager_s.allsockets = make(map[uint64]*GBHttpTask)
	httptaskmanager_s.starttempid = 1000000000
}

func GetHttpTaskManger() *GBHttpTaskManger {
	return httptaskmanager_s
}

func (this *GBHttpTaskManger) AddHttpTask(
	writer http.ResponseWriter) *GBHttpTask {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.starttempid++
	wstask := new(GBHttpTask)
	wstask.writer = writer
	curtime := uint64(time.Now().Unix())
	wstask.Tempid = (curtime << 32) + this.starttempid
	this.allsockets[wstask.Tempid] = wstask
	log.Debug("添加新的HTTP连接,%d,%d", wstask.Tempid,
		len(this.allsockets))

	notify := writer.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		this.mutex.Lock()
		defer this.mutex.Unlock()
		delete(this.allsockets, wstask.Tempid) // 这里应该需要加锁
		log.Debug("删除关闭的的HTTP连接,%d,%d", wstask.Tempid,
			len(this.allsockets))
	}()
	return wstask
}

func (this *GBHttpTaskManger) RemoveHttpTask(tempid uint64) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, found := this.allsockets[tempid]; found {
		delete(this.allsockets, tempid) // 这里应该需要加锁
		log.Debug("删除HTTP连接,%d,%d", tempid, len(this.allsockets))
	}
}

func (this *GBHttpTaskManger) GetHttpTask(tempid uint64) *GBHttpTask {
	if value, found := this.allsockets[tempid]; found {
		return value
	}
	return nil
}

func HttpRpcStart(serverinfo string, serviceMethod string,
	args interface{}, reply interface{}) error {
	return GetGBRPCManager().TCPRPCStart(serverinfo, serviceMethod,
		args, reply)

	// client, err := rpc.DialHTTP("tcp", serverinfo)
	// if err != nil {
	// 	log.Error("[RPC] 链接rpc服务器失败 [DialHTTP] "+
	// 		"Method[%s] ServerInfo[%s] Error[%s]",
	// 		serviceMethod, serverinfo, err.Error())
	// 	return err
	// }
	// defer client.Close()
	// err = client.Call(serviceMethod, args, reply)
	// if err != nil {
	// 	log.Error("[RPC] 链接rpc服务器失败 "+
	// 		"Method[%s] ServerInfo[%s] Error[%s]",
	// 		serviceMethod, serverinfo, err.Error())
	// } else {
	// 	log.Debug("[RPC] 链接rpc服务器成功 "+
	// 		"Method[%s] ServerInfo[%s]",
	// 		serviceMethod, serverinfo)
	// 	return nil
	// }
	// return nil
}
