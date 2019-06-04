/**
 * \file GBRedisManager.go
 * \version
 * \author wzy
 * \date  2018年06月05日 14:28:58
 * \brief 管理redis连接相关 把所有用户所在user服务器的信息加载到内存
 * \brief super 连接所有redis负责往里面写数据,其它只随机连接一个只读数据
 *
 */
package subnet

import (
	"base"
	"base/log"
	// "encoding/xml"
	"fmt"
	"github.com/garyburd/redigo/redis"
	// "io/ioutil"
	// "math/rand"
	"reflect"
	"time"
)

// websocket连接管理器
type GBRedisManager struct {
	redispools map[string]*redis.Pool // redis连接池
	allredis   map[string]redis.Conn  // 所有redis连接
}

var redismanager_s *GBRedisManager

func init() {
	redismanager_s = &GBRedisManager{}
	redismanager_s.allredis = make(map[string]redis.Conn)
	redismanager_s.redispools = make(map[string]*redis.Pool)

}

func GetGBRedisManager() *GBRedisManager {
	return redismanager_s
}

// 连接到redis
func (this *GBRedisManager) connectRedisServer(address string) bool {
	if _, founded := this.redispools[address]; founded {
		log.Error("[redis] connect redis failed Addr[%s] 重复连接",
			address)
		return false
	}
	newpool := &redis.Pool{
		MaxIdle:     30,
		MaxActive:   30,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			log.Debug("[redis] Connect redis succeed. "+
				" Addr[%s] 连接到redis服务器成功", address)
			return redis.Dial("tcp", address)
		},
	}
	this.redispools[address] = newpool

	return true
}

// 初始化连接到redis
func (this *GBRedisManager) InitRedisManager() {
	for _, v := range base.GetGBServerConfigM().RedisConfig.RedisList {
		address := fmt.Sprintf("%s:%d", v.IP, v.Port)
		log.Debug("[GBRedisManager.InitRedisManager] "+
			"初始化Redis连接 IPPort[%s]", address)
		this.connectRedisServer(address)
	}
	log.Debug("[GBRedisManager.InitRedisManager] "+
		"启动初始化成功,连接数量:%d", len(this.redispools))
}

func (this *GBRedisManager) GetString(key string) (string, error) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		tmpvalue, err := redis.String(conn.Do("GET", key))
		if err != nil {
			log.Debug("[redis] 读取数据异常 "+
				"Index[%s] Key[%s] Error[%s] Conn[%p] CountConn[%d]",
				index, key, err.Error(), conn, value.ActiveCount())
			return "", err
		} else {
			return tmpvalue, nil
		}
		// if tmpvalue != "" {
		// 	break
		// }
	}
	return "", nil
}

func (this *GBRedisManager) GetInt(key string) (int, error) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		tmpvalue, err := redis.Int(conn.Do("GET", key))
		if err != nil {
			log.Debug("[redis] 读取数据异常 "+
				"Index[%s] Key[%s] Error[%s] Conn[%p] CountConn[%d]",
				index, key, err.Error(), conn, value.ActiveCount())
			return 0, err
		} else {
			return tmpvalue, nil
		}
		// if tmpvalue != 0 {
		// 	break
		// }
	}
	return 0, nil
}

func (this *GBRedisManager) GetUint64(key string) (uint64, error) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		tmpvalue, err := redis.Uint64(conn.Do("GET", key))
		if err != nil {
			log.Debug("[redis] 读取数据异常 "+
				"Index[%s] Key[%s] Error[%s] Conn[%p] CountConn[%d]",
				index, key, err.Error(), conn, value.ActiveCount())
			return 0, err
		} else {
			return tmpvalue, nil
		}
		// if tmpvalue != 0 {
		// 	break
		// }
	}
	return 0, nil
}

func (this *GBRedisManager) Set(key string, args interface{}) error {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()

		_, err := conn.Do("SET", key, args)
		if err != nil {
			log.Error("[redis] 存取数据异常 "+
				"Index[%s] Key[%s] Error[%s] Conn[%p] CountConn[%d]",
				index, key, err.Error(), conn, value.ActiveCount())
			return err
		} else {
			return nil
		}
	}
	return nil
}

func (this *GBRedisManager) Expire(key string, timeoffset uint64) error {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()

		_, err := conn.Do("EXPIRE", key, timeoffset)
		if err != nil {
			log.Error("[redis] 数据设置过期时间异常,%s,%s,%s",
				index, key, err.Error())
			return err
		} else {
			return nil
		}
	}
	return nil
}

func (this *GBRedisManager) MSet(args ...interface{}) error {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()

		_, err := conn.Do("MSET", args...)
		if err != nil {
			log.Error("[redis] 存取数据异常 "+
				"Index[%s] Error[%s] Conn[%p] CountConn[%d]",
				index, err.Error(), conn, value.ActiveCount())
			return err
		} else {
			return nil
		}
	}
	return nil
}

// 获取指定UUID 至 UserServerId 的Key
func (this *GBRedisManager) GetUUIDToUServerIdKey(uuid uint64) string {
	return fmt.Sprintf("UTUS%d", uuid)
}

// 获取指定OpenId 至 UserServerId 的Key
func (this *GBRedisManager) GetOpenIdToUServerIdKey(openid string) string {
	return fmt.Sprintf("OTUS%s", openid)
}

// 获取用户所在UserServer的serverid
func (this *GBRedisManager) GetUserServerIDByUUID(
	uuid uint64) (uint32, error) {
	rediskey := this.GetUUIDToUServerIdKey(uuid)
	serverid, err := this.GetInt(rediskey)
	return uint32(serverid), err
}

// 获取用户所在UserServer的serverid
func (this *GBRedisManager) GetUserServerIDByOpenid(
	openid string) (uint32, error) {
	rediskey := this.GetOpenIdToUServerIdKey(openid)
	serverid, err := this.GetInt(rediskey)
	return uint32(serverid), err
}

// 新增用户所在UserServer的serverid
func (this *GBRedisManager) AddUserServerIDByUUID(
	uuid uint64, serverid uint32) {
	rediskey := this.GetUUIDToUServerIdKey(uuid)
	err := this.Set(rediskey, serverid)
	if err != nil {
		log.Error("[GBRedisManager.AddUserServerIDByUUID] Err[%s]",
			err.Error())
	}
}

// 新增用户所在UserServer的serverid
func (this *GBRedisManager) AddUserServerIDByUUIDMuch(
	datas map[uint64]uint32) {
	arge := make([]interface{}, len(datas)*2)
	i := 0
	for uuid, serverid := range datas {
		arge[i*2+0] = this.GetUUIDToUServerIdKey(uuid)
		arge[i*2+1] = serverid

		i++
	}
	// 按10000个玩家分组设置
	for len(arge) > 10000*2*2 {
		err := this.MSet(arge[:10000*2*2]...)
		if err != nil {
			log.Error("[GBRedisManager.AddUserServerIDByUUIDMuch] "+
				"MSet 1 Err[%s]",
				err.Error())
		}
		arge = arge[10000*2*2:]
	}
	if len(arge) > 0 {
		err := this.MSet(arge...)
		if err != nil {
			log.Error("[GBRedisManager.AddUserServerIDByUUIDMuch] "+
				"MSet 2 Err[%s]",
				err.Error())
		}
	}
}

// 新增用户所在UserServer的serverid
func (this *GBRedisManager) AddUserServerIDByOpenid(
	openid string, serverid uint32) {
	rediskey := this.GetOpenIdToUServerIdKey(openid)
	err := this.Set(rediskey, serverid)
	if err != nil {
		log.Error("[GBRedisManager.AddUserServerIDByOpenid] Err[%s]",
			err.Error())
	}
}

// 新增用户所在UserServer的serverid
func (this *GBRedisManager) AddUserServerIDByOpenidWithDeadline(
	openid string,
	serverid uint32, timeoffset uint64) {
	rediskey := this.GetOpenIdToUServerIdKey(openid)
	err := this.Set(rediskey, serverid)
	if err != nil {
		log.Error("[GBRedisManager.AddUserServerIDByOpenidWithDeadline] "+
			"Set Err[%s]",
			err.Error())
		return
	}
	err = this.Expire(rediskey, timeoffset)
	if err != nil {
		log.Error("[GBRedisManager.AddUserServerIDByOpenidWithDeadline] "+
			"Expire Err[%s]",
			err.Error())
		return
	}
}

// 新增用户所在UserServer的serverid
func (this *GBRedisManager) AddUserServerIDByOpenidMuch(
	datas map[string]uint32) {
	arge := make([]interface{}, len(datas)*2)
	i := 0
	for openid, serverid := range datas {
		arge[i*2+0] = this.GetOpenIdToUServerIdKey(openid)
		arge[i*2+1] = serverid

		i++
	}
	// 按10000个玩家分组设置
	for len(arge) > 10000*2*2 {
		err := this.MSet(arge[:10000*2*2]...)
		if err != nil {
			log.Error("[GBRedisManager.AddUserServerIDByOpenidMuch] "+
				"MSet 1 Err[%s]",
				err.Error())
		}
		arge = arge[10000*2*2:]
	}
	if len(arge) > 0 {
		err := this.MSet(arge...)
		if err != nil {
			log.Error("[GBRedisManager.AddUserServerIDByOpenidMuch] "+
				"MSet 2 Err[%s]",
				err.Error())
		}
	}
}

// 新增用户所在UserServer的serverid
func (this *GBRedisManager) AddUserServerID(
	openid string, uuid uint64, serverid uint32) {
	this.AddUserServerIDByOpenid(openid, serverid)
	this.AddUserServerIDByUUID(uuid, serverid)
}

// 批量添加用户所在UserServer的serverid

// 获取指定UUID 至 OpenId 的Key
func (this *GBRedisManager) GetUUIDToOpenIdKey(uuid uint64) string {
	return fmt.Sprintf("UTO%d", uuid)
}

// 获取指定OpenId 至 UUID 的Key
func (this *GBRedisManager) GetOpenIdToUUIDKey(openid string) string {
	return fmt.Sprintf("OTU%s", openid)
}

// 获取用户的OpenID
func (this *GBRedisManager) GetOpenidByUUID(uuid uint64) string {
	rediskey := this.GetUUIDToOpenIdKey(uuid)
	value, _ := this.GetString(rediskey)
	return value
}

// 获取用户的UUID
func (this *GBRedisManager) GetUUIDByOpenid(openid string) uint64 {
	rediskey := this.GetOpenIdToUUIDKey(openid)
	uuid, _ := this.GetUint64(rediskey)
	return uuid
}

// 新增用户的UUID 及 Openid 关系
func (this *GBRedisManager) AddOpenidAddUUID(
	openid string, uuid uint64) {
	rediskey := this.GetUUIDToOpenIdKey(uuid)
	err := this.Set(rediskey, openid)
	if err != nil {
		log.Error("[GBRedisManager.AddOpenidAddUUID] "+
			"Set 1 Err[%s]",
			err.Error())
	}

	rediskey1 := this.GetOpenIdToUUIDKey(openid)
	err = this.Set(rediskey1, uuid)
	if err != nil {
		log.Error("[GBRedisManager.AddOpenidAddUUID] "+
			"Set 2 Err[%s]",
			err.Error())
	}
}

// 新增用户的UUID 及 Openid 关系
func (this *GBRedisManager) AddOpenidAddUUIDMuch(
	datas map[uint64]string) {
	arge := make([]interface{}, len(datas)*2*2)
	i := 0
	for uuid, openid := range datas {
		arge[i*4+0] = this.GetUUIDToOpenIdKey(uuid)
		arge[i*4+1] = openid

		arge[i*4+2] = this.GetOpenIdToUUIDKey(openid)
		arge[i*4+3] = uuid

		i++
	}
	// 按10000个玩家分组设置
	for len(arge) > 10000*2*2 {
		err := this.MSet(arge[:10000*2*2]...)
		if err != nil {
			log.Error("[GBRedisManager.AddOpenidAddUUIDMuch] "+
				"MSet 1 Err[%s]",
				err.Error())
		}
		arge = arge[10000*2*2:]
	}
	if len(arge) > 0 {
		err := this.MSet(arge...)
		if err != nil {
			log.Error("[GBRedisManager.AddOpenidAddUUIDMuch] "+
				"MSet 2 Err[%s]",
				err.Error())
		}
	}
}

// 判断Openid是否存在
func (this *GBRedisManager) IsExistUserOpenid(
	openid string) (bool, error) {
	for key, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		rediskey := this.GetOpenIdToUUIDKey(openid)
		is_key_exit, err := redis.Bool(conn.Do("EXISTS", rediskey))
		if err != nil {
			log.Error("[redis] 判断数据是否存着异常 "+
				"Key[%s] Error[%s] Conn[%p] CountConn[%d]",
				key, err.Error(), conn, value.ActiveCount())
			return false, err
		} else {
			return is_key_exit, nil
		}
		// break
	}
	return false, nil
}

// 判断UUID是否存在
func (this *GBRedisManager) IsExistUserUUID(uuid uint64) (bool, error) {
	for key, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		rediskey := this.GetUUIDToOpenIdKey(uuid)
		is_key_exit, err := redis.Bool(conn.Do("EXISTS", rediskey))
		if err != nil {
			log.Error("[redis] 判断数据是否存着异常"+
				"Key[%s] Error[%s] Conn[%p] CountConn[%d]",
				key, err.Error(), conn, value.ActiveCount())
			return false, err
		} else {
			return is_key_exit, nil
		}
		// break
	}
	return false, nil
}

// 获取指定用户的Redis数据的Key
func (this *GBRedisManager) GetUserMemInfoRedisKey(uuid uint64) string {
	return fmt.Sprintf("UMI%d", uuid)
}

// 写入用户信息相关
func (this *GBRedisManager) AddUserInfo2MemDB(
	uuid uint64, v interface{}) {
	args := make([]interface{}, 0)
	key := this.GetUserMemInfoRedisKey(uuid)
	args = append(args, key)
	redisvalue := ""
	vtype := reflect.TypeOf(v)
	vvalue := reflect.ValueOf(v)
	fieldnum := vtype.NumField()
	for i := 0; i < fieldnum; i++ {
		field := vtype.Field(i)
		hkey := field.Name
		hvalue := ""
		switch vvalue.Field(i).Kind() {
		case reflect.Int32, reflect.Int64:
			hvalue = fmt.Sprintf("%d", vvalue.Field(i).Int())
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			hvalue = fmt.Sprintf("%d", vvalue.Field(i).Uint())
		case reflect.String:
			hvalue = vvalue.Field(i).String()
		}
		addvalue := fmt.Sprintf("%s %s ", hkey, hvalue)
		redisvalue += addvalue
		args = append(args, hkey, hvalue)
	}
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		_, err := conn.Do("HMSET", args...)
		if err != nil {
			log.Error("[redis] 存取用户数据异常 "+
				"Index[%s] Key[%s] Value[%s] Error[%s]",
				index, key, redisvalue, err.Error())
		} else {
			log.Debug("[redis] 存取用户数据成功 "+
				"Index[%s] Key[%s] Value[%s]",
				index, key, redisvalue)
		}
	}
}

// 读取用户信息相关
func (this *GBRedisManager) GetUserInfoToStruct(
	uuid uint64, v interface{}) {
	key := this.GetUserMemInfoRedisKey(uuid)
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		valueinfos, err := redis.Values(conn.Do("HGETALL", key))
		if err != nil {
			log.Error("[redis] 读取用户数据异常 "+
				"Index[%s] Key[%s] Error[%s]",
				index, key, err.Error())
			return
		}
		err = redis.ScanStruct(valueinfos, v)
		if err != nil {
			log.Error("[redis] 读取用户数据解析异常 "+
				"Index[%s] Key[%s] Error[%s]",
				index, key, err.Error())
		} else {
		}
		return
	}
}

// 获取用户信息相关
func (this *GBRedisManager) GetUserInfoToString(
	uuid uint64, hashkey string) string {
	key := this.GetUserMemInfoRedisKey(uuid)
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		infostr, err := redis.String(conn.Do("HGET", key, hashkey))
		if err != nil {
			log.Error("[redis] 读取用户数据ToString异常异常 "+
				"Index[%s] Key[%s] Error[%s]",
				index, key, err.Error())
			return ""
		} else {
			return infostr
		}
		// break
	}
	return ""
}

// 向有序集添加一项
func (this *GBRedisManager) ZAddToRedis(addkey string,
	addscore uint32, addmember string) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		_, err := conn.Do("ZADD", addkey, int64(addscore), addmember)
		if err != nil {
			log.Error("[redis] 添加用户排行榜数据异常 "+
				"Index[%s] Member[%s] Error[%s]",
				index, addmember, err.Error())
			return
		} else {
		}
	}
}

// 向有序集添加一项
func (this *GBRedisManager) ZRemToRedis(removekey string,
	removemember string) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		_, err := conn.Do("ZREM", removekey, removemember)
		if err != nil {
			log.Error("[redis] 删除用户排行榜数据异常 "+
				"Index[%s] Member[%s] Error[%s]",
				index, removemember, err.Error())
			return
		} else {
		}
	}
}

// 获取有序集中的前N项
func (this *GBRedisManager) ZRevRangeToRedis(
	rangekey string, maxnum uint32) []interface{} {
	return this.ZRevRangeToRedisInterval(rangekey, 0, maxnum)
}

// 获取有序集中 N - M 项
func (this *GBRedisManager) ZRevRangeToRedisInterval(rangekey string,
	minnum uint32, maxnum uint32) []interface{} {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		res, err := redis.Values(conn.Do("ZREVRANGE",
			rangekey, minnum, maxnum))
		if err != nil {
			log.Error("[redis] 加载用户排行榜数据异常 "+
				"Index[%s] Key[%s] Error[%s]",
				index, rangekey, err.Error())
			return nil
		} else {
			return res
		}
	}
	return nil
}

// 获取有序集中 指定成员的排名
func (this *GBRedisManager) ZRevRankToRedis(rangekey string,
	addmember string) int {
	for _, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		res, err := redis.Int(conn.Do("ZREVRANK", rangekey, addmember))
		if err != nil {
			// 没有找到目标成员
			return -1
		} else {
			return res
		}
	}
	return -1
}

// 获取有序集中 成员数量
func (this *GBRedisManager) ZCount(rangekey string) uint32 {
	resnum := uint32(0)
	for _, value := range this.redispools {
		conn := value.Get()
		res, err := redis.Int(conn.Do("ZCARD", rangekey))
		if err == nil {
			resnum += uint32(res)
		} else {
			log.Error("ZCARD ERR:%s", err.Error())
		}
		conn.Close()
	}
	return resnum
}

// 添加列表信息相关
func (this *GBRedisManager) LPushToRedis(
	addkey string, addvalue string) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		_, err := conn.Do("LPUSH", addkey, addvalue)
		if err != nil {
			log.Error("[redis] LPUSH添加榜数据异常 "+
				"Index[%s] Key[%s] Value[%s] Error[%s]",
				index, addkey, addvalue, err.Error())
			return
		} else {
		}
	}
}

// LTRIM key start stop
func (this *GBRedisManager) LTRIMToRedis(addkey string,
	start uint32, end uint32) {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		_, err := conn.Do("LTRIM", addkey, start, end)
		if err != nil {
			log.Error("[redis] LTRIM 更改数据异常 "+
				"Index[%s] Error[%s]",
				index, err.Error())
			return
		} else {
		}
	}
}

// 获取列表信息相关
// LRANGE addkey  start end [start,end] 的元素
func (this *GBRedisManager) LRangeToRedis(addkey string,
	start uint32, end uint32) []interface{} {
	for key, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		res, err := redis.Values(conn.Do("LRANGE", addkey, start, end))
		if err != nil {
			log.Error("[redis] 加载用户排行榜数据异常,%s,%s,%s",
				key, addkey, err.Error())
			return nil
		} else {
			return res
		}
	}
	return nil
}

func (this *GBRedisManager) SetWxcodeInfo(wxcode string, info string) {
	for key, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		rediskey := fmt.Sprintf("WXC_%s", wxcode)
		// wxcode 在redis中存活时间为一天
		_, err := conn.Do("SETEX", rediskey, 24*60*60, info)
		if err != nil {
			log.Error("[redis] SetWxcodeInfo添加数据异常,%s,%s,%s,%s",
				key, rediskey, info, err.Error())
			return
		} else {
		}
	}
}

func (this *GBRedisManager) GetWxcodeInfo(wxcode string) string {
	for index, value := range this.redispools {
		conn := value.Get()
		defer conn.Close()
		rediskey := fmt.Sprintf("WXC_%s", wxcode)
		tmpvalue, err := redis.String(conn.Do("GET", rediskey))
		if err != nil {
			log.Debug("[redis] 读取数据不成功 "+
				"%s,key:%s,%s,conn:%p,countconn:%d",
				index, rediskey, err.Error(), conn, value.ActiveCount())
			return ""
		} else {
			return tmpvalue
		}
		// if tmpvalue != "" {
		// 	break
		// }
	}
	return ""
}
