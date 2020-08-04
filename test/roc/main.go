package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"sync"
	"time"

	"github.com/liasece/micserver/app"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/module"
	"github.com/liasece/micserver/roc"
)

// roc type consts
const (
	ROCType1 roc.ObjType = "ROCType1"
)

// T1 test roc
type T1 struct {
	ID string
}

// GetROCObjID func
func (t *T1) GetROCObjID() string {
	return t.ID
}

// GetROCObjType func
func (t *T1) GetROCObjType() roc.ObjType {
	return ROCType1
}

// OnROCCall func
func (t *T1) OnROCCall(path *roc.Path, data []byte) ([]byte, error) {
	return []byte{}, nil
}

// Model type
type Model struct {
	module.BaseModule
	wg *sync.WaitGroup
}

//TopRunner func
func (m *Model) TopRunner() {
	roc := m.NewROC(ROCType1)
	for i := 0; i < 1000000; i++ {
		id := fmt.Sprintf("testObjID_%d", i)
		roc.RegObj(&T1{ID: id})
		if i%10000 == 0 {
			log.Debug("roc.RegObj id[%s]", id)
		}
	}
	m.wg.Done()
}

func main() {
	logPath := filepath.Join("log.log")
	log.SetDefaultLogger(log.NewLogger(nil, log.Options().FilePaths(logPath).AsyncWrite(true).Level(log.DebugLevel).RotateTimeLayout("060102")))

	go http.ListenAndServe("localhost:6060", nil)

	app := &app.App{}
	app.Setup(nil)

	wg := &sync.WaitGroup{}
	modules := make([]module.IModule, 0)
	for i := 0; i < 256; i++ {
		wg.Add(1)
		modules = append(modules, &Model{wg: wg})
	}

	{
		err := app.Init(modules)
		if err != nil {
			log.Error("app.Init error[%+v]", err)
		}
	}
	go app.RunAndBlock(modules)
	wg.Wait()
	time.Sleep(time.Second * 2000)
	app.Stop()

	time.Sleep(time.Microsecond * 2000)
}
