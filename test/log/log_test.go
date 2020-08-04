package main

import (
	"path/filepath"
	"testing"

	"github.com/liasece/micserver/log"
)

func Test_LogFileFlush(t *testing.T) {
	logPath := filepath.Join("log.log")
	log.SetDefaultLogger(log.NewLogger(nil, log.Options().FilePaths(logPath).RotateTimeLayout("060102")))

	log.Syslog("test Syslog", log.String("field", "test"), log.Bool("field1", true))
	log.Debug("test Debug", log.String("field", "test"), log.Bool("field1", true))
	log.Info("test Info", log.String("field", "test"), log.Bool("field1", true))
	log.Warn("test Warn", log.String("field", "test"), log.Bool("field1", true))
	log.Error("test Error %s", "test", log.String("field", "test"), log.Bool("field1", true))
	func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		log.DPanic("test DPanic", log.String("field", "test"), log.Bool("field1", true))
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		log.Panic("test Panic", log.String("field", "test"), log.Bool("field1", true))
	}()
	// log.Fatal("test Fatal", log.String("field", "test"), log.Bool("field1", true))
}
