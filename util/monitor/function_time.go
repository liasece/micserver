/**
 * \file GBFunctionTime.go
 * \version
 * \author wzy
 * \date  2018年02月01日 14:15:45
 * \brief 统计函数执行时间函数
 *
 */

package monitor

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/liasece/micserver/log"
)

// FunctionTime function time
type FunctionTime struct {
	starttime    uint64
	endtime      uint64
	functionname string
}

// Start func
func (ftime *FunctionTime) Start(name string) {
	ftime.functionname = name
	ftime.starttime = uint64(time.Now().UnixNano() / 1000000)
}

// Stop func
func (ftime *FunctionTime) Stop() {
	ftime.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := ftime.endtime - ftime.starttime
	if usetime > 500 {
		log.Error("[消耗时间统计],消耗时间超时,%s,%d 毫秒", ftime.functionname, usetime)
	}
}

// StopUseTime func
func (ftime *FunctionTime) StopUseTime() uint64 {
	ftime.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := ftime.endtime - ftime.starttime
	return uint64(usetime)
}

type gbFunctionInfo struct {
	functionname string
	callname     string
	callcount    uint32
	usetime      uint64
}

// GBOptimizeAnalysisM 性能统计分析
type GBOptimizeAnalysisM struct {
	functionmaps map[string]gbFunctionInfo
	starttime    uint64
	endtime      uint64
	mutex        sync.Mutex
}

var optimizeAnalysic *GBOptimizeAnalysisM

func init() {
	optimizeAnalysic = &GBOptimizeAnalysisM{}
	optimizeAnalysic.functionmaps = make(map[string]gbFunctionInfo)
}

// GetGBOptimizeAnalysisM func
func GetGBOptimizeAnalysisM() *GBOptimizeAnalysisM {
	return optimizeAnalysic
}

// StartCheck func
func (ftime *GBOptimizeAnalysisM) StartCheck() {
	ftime.starttime = uint64(time.Now().UnixNano() / 1000000)
}

// StopCheck func
func (ftime *GBOptimizeAnalysisM) StopCheck() {
	ftime.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := ftime.endtime - ftime.starttime
	if usetime > 100 {
		ftime.mutex.Lock()
		for _, funcinfo := range ftime.functionmaps {
			log.Debug("[分时消耗统计],消耗时间,%s,%s,%d 毫秒,调用:%d次", funcinfo.functionname, funcinfo.callname, funcinfo.usetime, funcinfo.callcount)
		}
		log.Debug("[分时消耗统计],消耗时间总计:%d毫秒", usetime)
		ftime.functionmaps = make(map[string]gbFunctionInfo)
		ftime.mutex.Unlock()
	} else {
		ftime.mutex.Lock()
		ftime.functionmaps = make(map[string]gbFunctionInfo)
		ftime.mutex.Unlock()
	}
}

// FunctionTimeAnalysic 结束数据会存入 GBOptimizeAnalysisM 中
type FunctionTimeAnalysic struct {
	starttime    uint64
	endtime      uint64
	functionname string
	callname     string
}

// Start func
func (ftime *FunctionTimeAnalysic) Start() {
	ftime.starttime = uint64(time.Now().UnixNano() / 1000000)
	pc, file, line, _ := runtime.Caller(1)
	ftime.functionname = fmt.Sprintf("%s:%d", file, line)
	f := runtime.FuncForPC(pc)
	ftime.callname = f.Name()
}

// Stop func
func (ftime *FunctionTimeAnalysic) Stop() {
	ftime.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := ftime.endtime - ftime.starttime
	optimizeAnalysic.mutex.Lock()
	if oldinfo, found := optimizeAnalysic.functionmaps[ftime.functionname]; found {
		oldinfo.usetime += usetime
		oldinfo.callcount++
		optimizeAnalysic.functionmaps[ftime.functionname] = oldinfo
	} else {
		newinfo := gbFunctionInfo{}
		newinfo.functionname = ftime.functionname
		newinfo.callname = ftime.callname
		newinfo.usetime += usetime
		newinfo.callcount = 1
		optimizeAnalysic.functionmaps[ftime.functionname] = newinfo
	}
	optimizeAnalysic.mutex.Unlock()
}
