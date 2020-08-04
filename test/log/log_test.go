package main

import (
	"path/filepath"
	"testing"

	"github.com/liasece/micserver/log"
)

func Test_LogFileFlush(t *testing.T) {
	logPath := filepath.Join("log.log")
	log.SetDefaultLogger(log.NewLogger(nil, log.Options().FilePaths(logPath).RotateTimeLayout("060102").Topic("Test_LogFileFlush")))

	log.Syslog("test %s", "Syslog", log.String("field", "test"), log.Bool("field1", true))
	log.Debug("test %s", "Debug", log.String("field", "test"), log.Bool("field1", true))
	log.Info("test %s", "Info", log.String("field", "test"), log.Bool("field1", true))
	log.Warn("test %s", "Warn", log.String("field", "test"), log.Bool("field1", true))
	log.Error("test %s", "Error", log.String("field", "test"), log.Bool("field1", true))
	log.DPanic("test %s", "DPanic", log.String("field", "test"), log.Bool("field1", true))
	func() {
		defer func() {
			if err := recover(); err != nil {
				// t.Error("panic", err)
			} else {
				t.Error("want panic")
			}
		}()
		log.Panic("test %s", "Panic", log.String("field", "test"), log.Bool("field1", true))
	}()
	t.Skip("Fatal")
	log.Fatal("test %s", "Fatal", log.String("field", "test"), log.Bool("field1", true))
}
