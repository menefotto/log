// Package log provides a New function that given a prefix string and a file or
// lack of empty string "" builds either a syslog or a filename named local log.
// It provides 2 functions and one of them its called Log, which taken a strings
// logs it, if the string is longer than 80 characters is going to be cut at 80
// characters.
// The other one is Close which must be called when the logger isn't needed any
// longer, IT'S IMPORTANT otherwise message may get lost and never reach the log
// whether is a syslog or local file.
// The main difference between other log package or the standard library one is
// to provide a simpler to use interface, and "ELIMINATE THE COST usually
// associated with the standard logging methods". All of this is achieved by letting
// a background cooroutine to the actual logging so the every call to Log is as
// much expensive as a message sent to a go channel.

package log

import (
	"log"
	"log/syslog"
	"os"
	"time"
)

// New creates a new logger takes a prefix with is going to be added to every line
// logged and a filename, if the file is omited and is given an empty string like
// so "" the local syslog will be used, if everything fails a nil logger will be
// returned else a valid logger instance is returned.
func New(prefix string, filename string) *Logger {
	var (
		stdlog  *log.Logger
		logfile *os.File
		err     error
		flag    int = log.Lshortfile
	)

	switch {
	case filename == "":
		priority := syslog.LOG_EMERG | syslog.LOG_USER
		stdlog, err = syslog.NewLogger(priority, flag)
		if err != nil {
			return nil
		}
	default:
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
		logfile, err = os.OpenFile(filename, flag, 0660)
		if err != nil {
			return nil
		}
		stdlog = log.New(logfile, prefix, flag)
	}

	stdlog.SetPrefix(prefix)

	log := &Logger{
		logger:   stdlog,
		file:     logfile,
		Messages: make(chan string, 1024),
		Done:     make(chan bool, 1),
		Sync:     make(chan bool, 1),
	}

	go func(log *Logger) {
		for {
			select {
			case <-log.Done:

				close(log.Messages)
				for msg := range log.Messages {
					log.log(msg)
				}

				time.Sleep(time.Millisecond * 250)

				defer func() {
					if log.file != nil {
						log.file.Close()
					}
				}()

				return
			case <-log.Sync:
				if log.file != nil {
					_ = log.file.Sync()
				}
			case msg := <-log.Messages:
				log.log(msg)
			}
		}
	}(log)

	return log

}

// Logger struct
type Logger struct {
	logger   *log.Logger
	file     *os.File
	Messages chan string
	Done     chan bool
	Sync     chan bool
}

// Log logs the message but doesn't assure that the message it being sent
func (s *Logger) Log(msg string) {
	s.toLog(msg)
}

// SyncLog logs the message and in case of a file logger syncs it to disk immidiately
// assuring logs are always written to disk.
func (s *Logger) SyncLog(msg string) {
	s.toLog(msg)
	s.Sync <- true
}

// Typically called in defer log.Close() fashion for short lived objects otherwise
// must be called when the logger isn't needed any longer, once called if any
// messages are still left on the Messages channel, they will be sent to syslog before
// actually exting, any message sent after the close method is called will a result
// in a sent operation to a close channel.
func (s *Logger) Close() {
	s.Done <- true
}

const (
	MsgMaxLen = 79
)

func (s *Logger) log(msg string) {
	s.logger.Output(2, msg)
}

func (s *Logger) toLog(msg string) {
	var truncated string

	if len(msg) > MsgMaxLen {
		log.Println(len(msg))
		truncated = msg[:79]
	} else {
		truncated = msg
	}

	s.Messages <- string(truncated)
}
