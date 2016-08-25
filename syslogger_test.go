package syslogger

import (
	"log/syslog"
	"testing"
)

func TestSysLogger(t *testing.T) {
	log := NewSysLogger("logger-test", syslog.LOG_ERR)
	log.MustLog("it works!")
	//	time.Sleep(time.Second * 1)
	//log.Close()
}
