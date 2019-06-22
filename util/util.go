package util

import (
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/url"
	// "os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"crypto/hmac"
	"crypto/sha1"
)

func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func GetStringHash(str string) uint32 {
	return GetHash(StringToBytes(str))
}

// GetHash returns a murmur32 hash for the data slice.
func GetHash(data []byte) uint32 {
	// Seed is set to 37, same as C# version of emitter
	var h1 uint32 = 37

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}

	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= 0x1b873593

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // rotl32(h1, 13)
		h1 = h1*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= 0x1b873593
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return (h1 << 24) | (((h1 >> 8) << 16) & 0xFF0000) | (((h1 >> 16) << 8) & 0xFF00) | (h1 >> 24)
}

func HmacSha1(content []byte, key []byte) string {
	//hmac ,use sha1
	mac := hmac.New(sha1.New, key)
	// mac := hmac.New(md5.New, key)
	_, err := mac.Write(content)
	if err != nil {
	}
	return string(mac.Sum(nil))
}

func UrlEncodeSortByKeys(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := url.QueryEscape(k + "=")
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteString(url.QueryEscape("&"))
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// func GetIPv4Addrs() {
// 	ifaces, _ := net.Interfaces()
// 	for _, ifi := range ifaces {
// 		fmt.Print("\nip addrs key,ifname:%s,\n", ifi.Name)
// 	}

// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	for _, address := range addrs {
// 		// 检查ip地址判断是否回环地址
// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				fmt.Println(ipnet.IP.String())
// 			}

// 		}
// 	}
// }

// GetIPv4ByInterface return IPv4 address from a specific interface IPv4 addresses
func GetIPv4ByInterface(name string) string {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}

	return ""
}

func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

func MysqlRealEscapeStringBack(value string) string {
	replace := map[string]string{"\\\\": "\\", `\'`: "'", "\\\\0": "\\0", "\\n": "\n", "\\r": "\r", `\"`: `"`, "\\Z": "\x1a"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

// 乱序切片
func SliceOutOfOrder(in []string) []string {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}
func SliceOutOfOrderByInt(in []uint64) []uint64 {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

func RandBetween(min, max int) int {
	if min >= max || max == 0 {
		return max
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := r.Intn(max-min) + min
	return random
}

func RandString(count int) string {
	var randomstr string
	for r := 0; r < count; r++ {
		i := RandBetween(65, 90)
		a := rune(i)
		randomstr += string(a)
	}
	return randomstr
}

type valueWeightItem struct {
	weight uint32
	value  uint64
}

// 权值对，根据权重随机一个值出来
type GBValueWeightPair struct {
	allweight uint32
	valuelist []valueWeightItem
}

func NewValueWeightPair() *GBValueWeightPair {
	vwp := new(GBValueWeightPair)
	vwp.valuelist = make([]valueWeightItem, 0)
	return vwp
}

func (this *GBValueWeightPair) Add(weight uint32, value uint64) {
	valueinfo := valueWeightItem{}
	valueinfo.weight = weight
	valueinfo.value = value
	this.valuelist = append(this.valuelist, valueinfo)
	this.allweight += weight
}
func (this *GBValueWeightPair) Random() uint64 {
	if this.allweight > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randvalue := uint32(r.Intn(int(this.allweight)))
		addweight := uint32(0)
		for i := 0; i < len(this.valuelist); i++ {
			addweight += this.valuelist[i].weight
			if randvalue <= addweight {
				return this.valuelist[i].value
			}
		}
	}
	return 0
}

// 根据权重列表随机出一个结果，返回命中下标
func RandWeight(weight []uint32) uint32 {
	total := uint32(0)
	for _, v := range weight {
		total += v
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randvalue := uint32(r.Intn(int(total)))
	tmp := uint32(0)
	for i, v := range weight {
		tmp += v
		if tmp >= randvalue {
			return uint32(i)
		}
	}
	return 0
}

type UtilWeightInterface interface {
	GetWeight() uint32
}

// 根据权重列表随机出一个结果，返回命中下标
func RandWeightStruct(weight []UtilWeightInterface) interface{} {
	if len(weight) == 0 {
		return nil
	}
	total := uint32(0)
	for _, v := range weight {
		total += v.GetWeight()
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randvalue := uint32(r.Intn(int(total)))
	tmp := uint32(0)
	for i, v := range weight {
		tmp += v.GetWeight()
		if tmp <= randvalue {
			return weight[i]
		}
	}
	return weight[0]
}

func LoadJsonFromFile(filename string, v interface{}) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("[LoadJsonFromFile] Load %s failed \n%s ",
			filename, err.Error())
	}
	err = json.Unmarshal([]byte(content), v)
	if err != nil {
		return fmt.Errorf(
			"[LoadJsonFromFile] Load %s failed, Unmarshal failed :\n%s ",
			filename, err.Error())
	}
	return nil
}

func StringSplitToUint32(str string, sli string) []uint32 {
	strlist := strings.Split(str, sli)
	res := make([]uint32, len(strlist))
	for i, str := range strlist {
		tmpint, err := strconv.Atoi(str)
		if err == nil {
			res[i] = uint32(tmpint)
		}
	}
	return res
}

func StringSplitToInt32(str string, sli string) []int32 {
	strlist := strings.Split(str, sli)
	res := make([]int32, len(strlist))
	for i, str := range strlist {
		tmpint, err := strconv.Atoi(str)
		if err == nil {
			res[i] = int32(tmpint)
		}
	}
	return res
}

// 获取系统毫秒时间
func GetTimeMs() uint64 {
	return uint64(time.Now().UnixNano()) / 1000000
}

func Abs(n int32) uint32 {
	if n < 0 {
		return uint32(-n)
	}
	return uint32(n)
}
