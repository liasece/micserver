/**
 * \file GBFunctionTime.go
 * \version
 * \author wzy
 * \date  2018年02月01日 14:15:45
 * \brief 统计函数执行时间函数
 *
 */

package functime

import (
	"fmt"
	"github.com/liasece/micserver/log"
	"runtime"
	"sync"
	"time"
)

type FunctionTime struct {
	starttime    uint64
	endtime      uint64
	functionname string
}

func (this *FunctionTime) Start(name string) {
	this.functionname = name
	this.starttime = uint64(time.Now().UnixNano() / 1000000)
}
func (this *FunctionTime) Stop() {
	this.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := this.endtime - this.starttime
	if usetime > 500 {
		log.Error("[消耗时间统计],消耗时间超时,%s,%d 毫秒", this.functionname, usetime)
	}
}
func (this *FunctionTime) StopUseTime() uint64 {
	this.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := this.endtime - this.starttime
	return uint64(usetime)
}

type gbFunctionInfo struct {
	functionname string
	callname     string
	callcount    uint32
	usetime      uint64
}

// 性能统计分析
type GBOptimizeAnalysisM struct {
	functionmaps map[string]gbFunctionInfo
	starttime    uint64
	endtime      uint64
	mutex        sync.Mutex
}

var optimize_analysic_s *GBOptimizeAnalysisM

func init() {
	optimize_analysic_s = &GBOptimizeAnalysisM{}
	optimize_analysic_s.functionmaps = make(map[string]gbFunctionInfo)
}

func GetGBOptimizeAnalysisM() *GBOptimizeAnalysisM {
	return optimize_analysic_s
}
func (this *GBOptimizeAnalysisM) StartCheck() {
	this.starttime = uint64(time.Now().UnixNano() / 1000000)
}
func (this *GBOptimizeAnalysisM) StopCheck() {
	this.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := this.endtime - this.starttime
	if usetime > 100 {
		this.mutex.Lock()
		for _, funcinfo := range this.functionmaps {
			log.Debug("[分时消耗统计],消耗时间,%s,%s,%d 毫秒,调用:%d次", funcinfo.functionname, funcinfo.callname, funcinfo.usetime, funcinfo.callcount)
		}
		log.Debug("[分时消耗统计],消耗时间总计:%d毫秒", usetime)
		this.functionmaps = make(map[string]gbFunctionInfo)
		this.mutex.Unlock()
	} else {
		this.mutex.Lock()
		this.functionmaps = make(map[string]gbFunctionInfo)
		this.mutex.Unlock()
	}
}

// 结束数据会存入 GBOptimizeAnalysisM 中
type FunctionTimeAnalysic struct {
	starttime    uint64
	endtime      uint64
	functionname string
	callname     string
}

func (this *FunctionTimeAnalysic) Start() {
	optimize_analysis := base.GetGBServerConfigM().GetProp("optimize_analysis")
	if optimize_analysis != "true" {
		return
	}
	this.starttime = uint64(time.Now().UnixNano() / 1000000)
	pc, file, line, _ := runtime.Caller(1)
	this.functionname = fmt.Sprintf("%s:%d", file, line)
	f := runtime.FuncForPC(pc)
	this.callname = f.Name()
}
func (this *FunctionTimeAnalysic) Stop() {
	optimize_analysis := base.GetGBServerConfigM().GetProp("optimize_analysis")
	if optimize_analysis != "true" {
		return
	}
	this.endtime = uint64(time.Now().UnixNano() / 1000000)
	usetime := this.endtime - this.starttime
	optimize_analysic_s.mutex.Lock()
	if oldinfo, found := optimize_analysic_s.functionmaps[this.functionname]; found {
		oldinfo.usetime += usetime
		oldinfo.callcount++
		optimize_analysic_s.functionmaps[this.functionname] = oldinfo
	} else {
		newinfo := gbFunctionInfo{}
		newinfo.functionname = this.functionname
		newinfo.callname = this.callname
		newinfo.usetime += usetime
		newinfo.callcount = 1
		optimize_analysic_s.functionmaps[this.functionname] = newinfo
	}
	optimize_analysic_s.mutex.Unlock()
}
