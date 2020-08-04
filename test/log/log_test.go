package main

import (
	"path/filepath"
	"testing"

	"github.com/liasece/micserver/log"
)

func Test_LogFileFlush(t *testing.T) {
	logPath := filepath.Join("log.log")
	log.SetDefaultLogger(log.NewLogger(nil, log.Options().FilePaths(logPath).RotateTimeLayout("060102")))

	log.Syslog("test", log.String("field", "test"), log.Bool("field1", true))
	log.Info("test", log.String("field", "test"), log.Bool("field1", true))
	log.Debug("test", log.String("field", "test"), log.Bool("field1", true))
	log.Warn("test", log.String("field", "test"), log.Bool("field1", true))
	log.Error("test %s", "test", log.String("field", "test"), log.Bool("field1", true))
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	log.Panic("test", log.String("field", "test"), log.Bool("field1", true))
}
