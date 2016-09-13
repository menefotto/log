package syslogger

import (
	"log/syslog"
	"testing"
	"time"
)

func TestSysLogger(t *testing.T) {
	log := NewLogger("logger-test:", syslog.LOG_ERR)
	log.MustLog("it works!")
	time.Sleep(time.Second * 1)
	log.Close()
}

func TestMsgToLong(t *testing.T) {
	log := NewLogger("logger-test:", syslog.LOG_ERR)
	log.MustLog("it works! to long not to long to long not to long but it is too long or it isnt't na it isn't that long :)")
	time.Sleep(time.Second * 1)
	log.Log("here we go")
	time.Sleep(time.Microsecond * 500)
	log.Close()
}
