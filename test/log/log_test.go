package main

import (
	"path/filepath"
	"testing"

	"github.com/liasece/micserver/log"
)

func Test_LogFileFlush(t *testing.T) {
	logPath := filepath.Join("log.log")
	log.SetDefaultLogger(log.NewLogger(&log.Options{
		FilePaths: []string{logPath},
		// AsyncWrite: true,
	}))
	log.Syslog("test")
	log.Info("test")
	log.Debug("test")
	log.Warn("test")
	log.Error("test")
}
